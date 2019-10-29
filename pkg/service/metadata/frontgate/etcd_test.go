// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// go test -test-etcd-enabled
// go test -test-etcd-enabled -run ^name$

package frontgate

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	pbtypes "openpitrix.io/openpitrix/pkg/pb/metadata/types"
)

type TImage struct {
	Url    string
	Width  int
	Height int
}

type TPerson struct {
	Name  string
	Age   int
	Photo TImage
}

func tGetEtcdEndpoints() []string {
	return strings.Split(*tEtcdHost, ",")
}

func _TestDefaultEtcdMaxOpsPerTxn(t *testing.T) {
	// import "github.com/coreos/etcd/embed"
	// Assert(t, embed.DefaultMaxTxnOps == DefaultEtcdMaxOpsPerTxn)
}

func TestEtcdClient(t *testing.T) {
	if !*tEtcdEnabled {
		t.Skip("etcd disabled")
	}

	c, err := NewEtcdClient(tGetEtcdEndpoints(), time.Second/10, 0)
	Assert(t, err == nil, err)
	Assert(t, c != nil)

	err = c.Clear()
	Assert(t, err == nil, err)

	err = c.Set("abc", "abc-value")
	Assert(t, err == nil, err)

	err = c.Set("abc/sub", "abc-sub-value")
	Assert(t, err == nil, err)

	m, err := c.GetAllValues()
	Assert(t, err == nil, err)
	Assert(t, len(m) == 2)
	Assert(t, m["abc"] == "abc-value")
	Assert(t, m["abc/sub"] == "abc-sub-value")

	s, ok := c.Get("abc")
	Assert(t, ok)
	Assert(t, s == "abc-value")

	m0, err := c.GetValues("abc", "abc/sub", "11")
	Assert(t, err == nil, err)
	Assert(t, len(m0) == 2)
	Assert(t, m0["abc"] == "abc-value")
	Assert(t, m0["abc/sub"] == "abc-sub-value")
	Assert(t, m0["11"] == "")

	c.DelValues("abc", "abc/sub")

	s, ok = c.Get("abc")
	Assert(t, !ok)
	Assert(t, s == "")

	m0, err = c.GetValues("abc", "abc/sub", "11")
	Assert(t, err == nil, err)
	Assert(t, len(m0) == 0)

	m, err = c.GetAllValues()
	Assert(t, err == nil, err)
	Assert(t, len(m) == 0)
}

func TestEtcdClient_SetStructValue(t *testing.T) {
	if !*tEtcdEnabled {
		t.Skip("etcd disabled")
	}

	c, err := NewEtcdClient(tGetEtcdEndpoints(), time.Second/10, 0)
	Assert(t, err == nil, err)
	Assert(t, c != nil)

	err = c.Clear()
	Assert(t, err == nil, err)

	person1 := TPerson{
		Age:  2018,
		Name: "openpitrix",
		Photo: TImage{
			Url:    "https://openpitrix.io",
			Width:  100,
			Height: 200,
		},
	}

	err = c.SetStructValue("/person", &person1)
	Assert(t, err == nil, err)

	var person2 TPerson
	c.GetStructValue("/person", &person2)
	Assert(t, reflect.DeepEqual(person1, person2), person1, person2)
}

func TestEtcdClient_GetValuesByPrefix(t *testing.T) {
	if !*tEtcdEnabled {
		t.Skip("etcd disabled")
	}

	c, err := NewEtcdClient(tGetEtcdEndpoints(), time.Second/10, 0)
	Assert(t, err == nil, err)
	Assert(t, c != nil)

	err = c.Clear()
	Assert(t, err == nil, err)

	person1 := TPerson{
		Age:  2018,
		Name: "openpitrix",
		Photo: TImage{
			Url:    "https://openpitrix.io",
			Width:  100,
			Height: 200,
		},
	}

	err = c.SetStructValue("/person", &person1)
	Assert(t, err == nil, err)

	m0, err := c.GetValuesByPrefix("/person")
	Assert(t, err == nil, err)
	Assert(t, len(m0) == 5)

	Assert(t, m0["/person/Age"] == strconv.Itoa(person1.Age))
	Assert(t, m0["/person/Name"] == person1.Name)
	Assert(t, m0["/person/Photo/Url"] == person1.Photo.Url)
	Assert(t, m0["/person/Photo/Width"] == strconv.Itoa(person1.Photo.Width))
	Assert(t, m0["/person/Photo/Height"] == strconv.Itoa(person1.Photo.Height))
}
func TestEtcdClient_metadata(t *testing.T) {
	if !*tEtcdEnabled {
		t.Skip("etcd disabled")
	}

	c, err := NewEtcdClient(tGetEtcdEndpoints(), time.Second/10, 0)
	Assert(t, err == nil, err)
	Assert(t, c != nil)

	err = c.Clear()
	Assert(t, err == nil, err)

	err = c.RegisterMetadata(&pbtypes.SubTask_RegisterMetadata{
		Cnodes: `{
			"/abc":   "abc-value",
			"/123/a": "123-a-value",
			"/123/b": "123-b-value"
		}`,
	})
	Assert(t, err == nil, err)

	m0, err := c.GetValuesByPrefix("/")
	Assert(t, err == nil, err)
	Assert(t, len(m0) == 3)

	err = c.DeregisterMetadata(&pbtypes.SubTask_DeregisterMetadata{
		Cnodes: `{
			"/abc":  "",
			"/123/": ""
		}`, // remove /123/a,/123/b
	})
	Assert(t, err == nil, err)

	m0, err = c.GetValuesByPrefix("/")
	Assert(t, err == nil, err)
	Assert(t, len(m0) == 0)
}

func TestEtcdClient_cmd(t *testing.T) {
	if !*tEtcdEnabled {
		t.Skip("etcd disabled")
	}

	c, err := NewEtcdClient(tGetEtcdEndpoints(), time.Second/10, 0)
	Assert(t, err == nil, err)
	Assert(t, c != nil)

	err = c.Clear()
	Assert(t, err == nil, err)

	err = c.RegisterMetadata(&pbtypes.SubTask_RegisterMetadata{
		Cnodes: `{
			"/cmd/abc":   "abc-value",
			"/cmd/123/a": "123-a-value",
			"/cmd/123/b": "123-b-value"
		}`,
	})
	Assert(t, err == nil, err)

	m0, err := c.GetValuesByPrefix("/")
	Assert(t, err == nil, err)
	Assert(t, len(m0) == 3)

	err = c.DeregisterCmd(&pbtypes.SubTask_DeregisterCmd{
		Cnodes: `{
			"/cmd": ""
		}`,
	})
	Assert(t, err == nil, err)

	m0, err = c.GetValuesByPrefix("/")
	Assert(t, err == nil, err)
	Assert(t, len(m0) == 0)
}
