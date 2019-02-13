// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

module "openpitrix.io/libconfd"

require (
	"github.com/BurntSushi/toml" v0.3.0
	"go.etcd.io/etcd/clientv3" v3.3.0
	"github.com/urfave/cli" v1.20.0
	"golang.org/x/crypto" v0.0.0-20180219163459-432090b8f568
)
