// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// Package libconfd provides mini confd lib.
package libconfd

import (
	"errors"
	"sync"
	"time"
)

type Call struct {
	Config *Config
	Client BackendClient
	Error  error
	Done   chan *Call
}

func (call *Call) done() {
	select {
	case call.Done <- call:
		// ok
	default:
		// We don't want to block here. It is the caller's responsibility to make
		// sure the channel has enough buffer space. See comment in Go().
		GetLogger().Debugln("libconfd: discarding Call reply due to insufficient Done chan capacity")
	}
}

type Processor struct {
	pendingMutex sync.Mutex
	pending      []*Call

	closeChan chan bool
	wg        sync.WaitGroup
}

func (p *Processor) isClosing() bool {
	if p.closeChan == nil {
		GetLogger().Panic("closeChan is nil")
	}
	select {
	case <-p.closeChan:
		return true
	default:
		return false
	}
}

func (p *Processor) addPendingCall(call *Call) {
	p.pendingMutex.Lock()
	defer p.pendingMutex.Unlock()

	p.pending = append(p.pending, call)
}
func (p *Processor) getPendingCall() *Call {
	p.pendingMutex.Lock()
	defer p.pendingMutex.Unlock()

	if len(p.pending) == 0 {
		return nil
	}

	call := p.pending[0]
	p.pending = p.pending[1:]
	return call
}
func (p *Processor) clearPendingCall() {
	p.pendingMutex.Lock()
	defer p.pendingMutex.Unlock()

	for _, call := range p.pending {
		err := errors.New("libconfd: processor is shut down")
		GetLogger().Error(err)
		call.Error = err
		call.done()
	}

	p.pending = p.pending[:0]
}

func (p *Processor) checkBackendClient(client BackendClient) error {
	_, err := client.GetValues([]string{"/"})
	if err != nil {
		GetLogger().Error(err)
		return err
	}
	return nil
}

func NewProcessor() *Processor {
	p := &Processor{
		closeChan: make(chan bool),
	}

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()

		for {
			if p.isClosing() {
				p.clearPendingCall()
				return
			}

			call := p.getPendingCall()
			if call == nil {
				time.Sleep(time.Second / 10)
				continue
			}

			p.wg.Add(1)
			go func() {
				GetLogger().Debugln("process start")
				defer GetLogger().Debugln("process done")

				defer p.wg.Done()
				defer call.done()

				p.process(call)
			}()
		}
	}()

	return p
}

func (p *Processor) Go(cfg *Config, client BackendClient, opts ...Options) *Call {
	if client == nil {
		GetLogger().Panic("client is nil")
	}

	call := new(Call)

	call.Config = cfg.Clone().applyOptions(opts...)
	call.Client = client
	call.Done = make(chan *Call, 10) // buffered.

	if err := cfg.Valid(); err != nil {
		GetLogger().Error(err)
		call.Error = err
		call.done()
		return call
	}

	// just print the when check failed
	if err := p.checkBackendClient(client); err != nil {
		GetLogger().Warning(err)
		// donot return
	}

	p.addPendingCall(call)
	return call
}

func (p *Processor) Run(cfg *Config, client BackendClient, opts ...Options) error {
	if err := cfg.Valid(); err != nil {
		GetLogger().Error(err)
		return err
	}
	if client == nil {
		GetLogger().Panic("client is nil")
	}

	call := <-p.Go(cfg, client, opts...).Done
	if err := call.Error; err != nil {
		GetLogger().Error(err)
		return err
	}

	return nil
}

func (p *Processor) Close() error {
	close(p.closeChan)
	p.wg.Wait()
	return nil
}

func (p *Processor) process(call *Call) {
	switch {
	case call.Config.Onetime:
		p.runOnce(call)
	case call.Config.Watch:
		p.runInWatchMode(call)
	default:
		p.runInIntervalMode(call)
	}
}

func (p *Processor) runOnce(call *Call) {
	ts, err := MakeAllTemplateResourceProcessor(call.Config, call.Client)
	if err != nil {
		GetLogger().Error(err)
		call.Error = err
		return
	}

	for _, t := range ts {
		if p.isClosing() {
			return
		}

		if err := t.Process(call); err != nil {
			GetLogger().Error(err)
		}
	}

	return
}

func (p *Processor) runInIntervalMode(call *Call) {
	ts, err := MakeAllTemplateResourceProcessor(call.Config, call.Client)
	if err != nil {
		GetLogger().Warning(err)
		call.Error = err
		return
	}

	for {
		if p.isClosing() {
			return
		}

		for _, t := range ts {
			if p.isClosing() {
				return
			}

			if err := t.Process(call); err != nil {
				GetLogger().Error(err)
				continue
			}
		}

		time.Sleep(time.Duration(call.Config.Interval) * time.Second)
	}
}

func (p *Processor) runInWatchMode(call *Call) {
	ts, err := MakeAllTemplateResourceProcessor(call.Config, call.Client)
	if err != nil {
		GetLogger().Warning(err)
		return
	}

	var wg sync.WaitGroup
	var stopChan = make(chan bool)

	for i := 0; i < len(ts); i++ {
		wg.Add(1)
		go func(t *TemplateResourceProcessor) {
			defer wg.Done()
			p.monitorPrefix(t, &wg, stopChan, call)
		}(ts[i])
	}

	for {
		time.Sleep(time.Second / 2)

		if p.isClosing() {
			close(stopChan)
			break
		}
	}

	wg.Wait()
	return
}

func (p *Processor) monitorPrefix(
	t *TemplateResourceProcessor,
	wg *sync.WaitGroup, stopChan chan bool,
	call *Call,
) {
	keys := t.getAbsKeys()

	// hook keys
	if fn := call.Config.HookAbsKeyAdjuster; fn != nil {
		for i, k := range keys {
			keys[i] = fn(k)
		}
	}

	for {
		if p.isClosing() {
			return
		}

		// watch some key changed
		index, err := t.client.WatchPrefix(t.Prefix, keys, t.lastIndex, stopChan)
		if err != nil {
			GetLogger().Error(err)
		}

		t.lastIndex = index
		if err := t.Process(call); err != nil {
			GetLogger().Error(err)
		}
	}
}
