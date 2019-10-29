// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// https://coreos.com/blog/etcd
// https://coreos.com/blog/transactional-memory-with-etcd3.html

package frontgate

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/chai2010/jsonmap"
	"github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const DefaultEtcdMaxOpsPerTxn = 128

type EtcdClientManager struct {
	clientMap map[string]*EtcdClient
	mu        sync.Mutex
}

type EtcdClient struct {
	*clientv3.Client
	maxTxnOps int
}

func NewEtcdClientManager() *EtcdClientManager {
	return &EtcdClientManager{
		clientMap: make(map[string]*EtcdClient),
	}
}

func (p *EtcdClientManager) Get() (*EtcdClient, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, c := range p.clientMap {
		return c, nil
	}

	return nil, fmt.Errorf("frontgate: no valid etcd client")
}

func (p *EtcdClientManager) GetClient(endpoints []string, timeout time.Duration, maxTxnOps int) (*EtcdClient, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	key := fmt.Sprintf("%v:%v", endpoints, timeout)
	if c, ok := p.clientMap[key]; ok {
		return c, nil
	}

	c, err := NewEtcdClient(endpoints, timeout, maxTxnOps)
	if err != nil {
		return nil, err
	}

	p.clientMap[key] = c
	return c, nil
}

func (p *EtcdClientManager) ClearClient(client *EtcdClient) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for k, c := range p.clientMap {
		if c == client {
			delete(p.clientMap, k)
			c.Close()
			return
		}
	}
}

func (p *EtcdClientManager) Clear() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, c := range p.clientMap {
		c.Close()
	}

	p.clientMap = make(map[string]*EtcdClient)
}

func NewEtcdClient(endpoints []string, timeout time.Duration, maxTxnOps int) (*EtcdClient, error) {
	if timeout == 0 {
		timeout = time.Second
	}
	if maxTxnOps == 0 {
		maxTxnOps = DefaultEtcdMaxOpsPerTxn
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: timeout,
	})
	if err != nil {
		return nil, err
	}

	p := &EtcdClient{
		Client:    cli,
		maxTxnOps: maxTxnOps,
	}

	return p, nil
}
func NewEtcdClientWithConfig(cfg clientv3.Config, maxTxnOps int) (*EtcdClient, error) {
	if maxTxnOps == 0 {
		maxTxnOps = DefaultEtcdMaxOpsPerTxn
	}

	cli, err := clientv3.New(cfg)
	if err != nil {
		return nil, err
	}

	p := &EtcdClient{
		Client:    cli,
		maxTxnOps: maxTxnOps,
	}

	return p, nil
}

func (p *EtcdClient) Close() error {
	return p.Client.Close()
}

func (p *EtcdClient) Get(key string) (val string, ok bool) {
	kv := clientv3.NewKV(p.Client)

	resp, err := kv.Get(context.Background(), key)
	if err != nil {
		return "", false
	}

	if len(resp.Kvs) > 0 {
		return string(resp.Kvs[0].Value), true
	}

	return "", false
}

func (p *EtcdClient) Set(key, val string) error {
	kv := clientv3.NewKV(p.Client)

	_, err := kv.Put(context.Background(), key, val)
	if err != nil {
		return err
	}

	return nil
}

func (p *EtcdClient) GetValues(keys ...string) (map[string]string, error) {
	kvc := clientv3.NewKV(p.Client)

	var ops []clientv3.Op
	for _, k := range keys {
		ops = append(ops, clientv3.OpGet(k))
	}

	m := make(map[string]string)
	for startIdx := 0; startIdx < len(ops); startIdx += p.maxTxnOps {
		endIdx := startIdx + p.maxTxnOps
		if endIdx > len(ops) {
			endIdx = len(ops)
		}

		resp, err := kvc.Txn(context.Background()).Then(ops[startIdx:endIdx]...).Commit()
		if err != nil {
			return nil, err
		}

		for _, resp_i := range resp.Responses {
			if respRange := resp_i.GetResponseRange(); respRange != nil {
				for _, kv := range respRange.Kvs {
					m[string(kv.Key)] = string(kv.Value)
				}
			}
		}
	}

	return m, nil
}

func (p *EtcdClient) GetValuesByPrefix(keyPrefix string) (map[string]string, error) {
	m := make(map[string]string)
	kv := clientv3.NewKV(p.Client)

	resp, err := kv.Get(context.Background(), keyPrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	for _, v := range resp.Kvs {
		m[string(v.Key)] = string(v.Value)
	}

	return m, nil
}

func (p *EtcdClient) SetValues(m map[string]string) error {
	kvc := clientv3.NewKV(p.Client)

	var ops []clientv3.Op
	for k, v := range m {
		ops = append(ops, clientv3.OpPut(k, v))
	}

	for startIdx := 0; startIdx < len(ops); startIdx += p.maxTxnOps {
		endIdx := startIdx + p.maxTxnOps
		if endIdx > len(ops) {
			endIdx = len(ops)
		}

		_, err := kvc.Txn(context.Background()).Then(ops[startIdx:endIdx]...).Commit()
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *EtcdClient) SetValuesWithPrefix(prefix string, m map[string]string) error {
	kvc := clientv3.NewKV(p.Client)

	var ops []clientv3.Op
	for k, v := range m {
		ops = append(ops, clientv3.OpPut(prefix+k, v))
	}

	for startIdx := 0; startIdx < len(ops); startIdx += p.maxTxnOps {
		endIdx := startIdx + p.maxTxnOps
		if endIdx > len(ops) {
			endIdx = len(ops)
		}

		_, err := kvc.Txn(context.Background()).Then(ops[startIdx:endIdx]...).Commit()
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *EtcdClient) GetStructValue(keyPrefix string, out interface{}) error {
	m := make(map[string]interface{})
	kv := clientv3.NewKV(p.Client)

	resp, err := kv.Get(context.Background(), keyPrefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, v := range resp.Kvs {
		key, val := string(v.Key), string(v.Value)
		key = strings.TrimPrefix(key, keyPrefix)
		m[key] = string(val)
	}

	jsonMap := jsonmap.NewJsonMapFromKV(m, "/")
	return jsonMap.ToStruct(out)
}

func (p *EtcdClient) SetStructValue(keyPrefix string, val interface{}) error {
	kvc := clientv3.NewKV(p.Client)

	m := make(map[string]string)
	for k, v := range jsonmap.NewJsonMapFromStruct(val).ToMapString("/") {
		m[keyPrefix+k] = v
	}

	var ops []clientv3.Op
	for k, v := range m {
		ops = append(ops, clientv3.OpPut(k, v))
	}

	for startIdx := 0; startIdx < len(ops); startIdx += p.maxTxnOps {
		endIdx := startIdx + p.maxTxnOps
		if endIdx > len(ops) {
			endIdx = len(ops)
		}

		_, err := kvc.Txn(context.Background()).Then(ops[startIdx:endIdx]...).Commit()
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *EtcdClient) DelValues(keys ...string) error {
	kvc := clientv3.NewKV(p.Client)

	var ops []clientv3.Op
	for _, k := range keys {
		ops = append(ops, clientv3.OpDelete(k))
	}

	for startIdx := 0; startIdx < len(ops); startIdx += p.maxTxnOps {
		endIdx := startIdx + p.maxTxnOps
		if endIdx > len(ops) {
			endIdx = len(ops)
		}

		_, err := kvc.Txn(context.Background()).Then(ops[startIdx:endIdx]...).Commit()
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *EtcdClient) DelValuesWithPrefix(keyPrefixs ...string) error {
	kvc := clientv3.NewKV(p.Client)

	var ops []clientv3.Op
	for _, k := range keyPrefixs {
		ops = append(ops, clientv3.OpDelete(k, clientv3.WithPrefix()))
	}

	for startIdx := 0; startIdx < len(ops); startIdx += p.maxTxnOps {
		endIdx := startIdx + p.maxTxnOps
		if endIdx > len(ops) {
			endIdx = len(ops)
		}

		_, err := kvc.Txn(context.Background()).Then(ops[startIdx:endIdx]...).Commit()
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *EtcdClient) GetAllValues() (map[string]string, error) {
	kvc := clientv3.NewKV(p.Client)

	resp, err := kvc.Get(context.Background(), "", clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	for _, v := range resp.Kvs {
		m[string(v.Key)] = string(v.Value)
	}

	return m, nil
}
func (p *EtcdClient) Clear() error {
	kvc := clientv3.NewKV(p.Client)

	_, err := kvc.Delete(context.Background(), "", clientv3.WithPrefix())
	if err != nil {
		return err
	}

	return nil
}

func (p *EtcdClient) isHaltErr(err error) bool {
	if err == nil {
		return false
	}

	// Deprecated: this error should not be relied upon by users; use the status
	// code of Canceled instead.
	if err == grpc.ErrClientConnClosing {
		return true
	}
	if ev, ok := status.FromError(err); ok {
		switch ev.Code() {
		case codes.Canceled, codes.Unavailable, codes.Internal:
			return true
		}
	}

	// not the connection error
	return false
}
