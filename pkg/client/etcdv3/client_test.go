// Copyright 2018 Yunify Inc. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package etcdv3

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/yunify/metad/log"
	"github.com/yunify/metad/store"
)

func init() {
	log.SetLevel("debug")
	rand.Seed(int64(time.Now().Nanosecond()))
}

func TestClientSyncStop(t *testing.T) {

	prefix := fmt.Sprintf("/prefix%v", rand.Intn(1000))

	stopChan := make(chan bool)
	log.Info("prefix is %s", prefix)
	nodes := []string{"http://127.0.0.1:2379"}
	storeClient, err := NewEtcdClient("default", prefix, nodes, "", "", "", false, "", "")
	assert.NoError(t, err)
	go func() {
		time.Sleep(3 * time.Second)
		stopChan <- true
	}()

	metastore := store.New()
	// expect internalSync not block after stopChan has signal
	initWG := &sync.WaitGroup{}
	initWG.Add(1)

	doneWG := &sync.WaitGroup{}
	doneWG.Add(1)

	go func() {
		storeClient.internalSync(prefix, stopChan, initWG, storeClient.newInitStoreFunc(prefix, metastore), newProcessSyncChangeFunc(metastore))
		doneWG.Done()
	}()
	initWG.Wait()
	doneWG.Wait()
}

func TestClientSyncStopWhenInitError(t *testing.T) {

	prefix := fmt.Sprintf("/prefix%v", rand.Intn(1000))

	stopChan := make(chan bool)
	log.Info("prefix is %s", prefix)
	nodes := []string{"http://127.0.0.1:2379"}
	storeClient, err := NewEtcdClient("default", prefix, nodes, "", "", "", false, "", "")
	assert.NoError(t, err)
	go func() {
		time.Sleep(3 * time.Second)
		stopChan <- true
	}()

	metastore := store.New()
	// expect internalSync not block after stopChan has signal
	initWG := &sync.WaitGroup{}
	initWG.Add(1)

	doneWG := &sync.WaitGroup{}
	doneWG.Add(1)
	go func() {
		storeClient.internalSync(prefix, stopChan, initWG, func() error {
			return fmt.Errorf("always error")
		}, newProcessSyncChangeFunc(metastore))
		doneWG.Done()
	}()
	initWG.Wait()
	doneWG.Wait()
}
