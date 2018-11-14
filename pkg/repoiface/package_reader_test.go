// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repoiface

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"

	"openpitrix.io/openpitrix/pkg/util/archiveutil"
)

const testBase64 = `H4sIAAAAAAAA/+xXTW/jNhDN2b9iwLPqiP6Ivb4VObcNsN1eFoHByGOHa4lkSSpIHPi/FxTtWKIpp22cDQJoLjY4j+S8J85wuFlfKpat2Qr7P4wUF+9haZqmk/G4+k3TNPxN6fjwvxqnA0qvLiB9l2gCK41l+iJ9814huU9iz4QpPn9AbbgUZEYeKEmIYAWSGdmsSUIOrrRP+ylJCFOqNsGPLdBkmivrx34V8IdCccOt5o/AlCLbj+bZWdw268tMiiVfvV/6v5b/o+FwGOQ/nUwmXf7/DOs99wAAiH1SSGZAmNbsiSR+UGmpUFuOhszgu0dWjjU+OXCWl8ai3sErT6MQANmsYQeC2mI1fGzfk3vX96/qVNJ05ewO8zZnEN2f9wgOBnIJ9h5hswaD+oFnRxP3YRqruVgdL7tkZW494dCp8e+Sa1yQGSxZbvDFu00gTqseZRu7E5gIydrQGbm+nakp7wTaVpJfo+6A3/W9lAaBgV8MrIQfkh+pcg5WVpf/gpSWOc7dsaKtxNohATkH/OWOGVzAYU48leBUOsFrKVUnkKkymNqM/vrmW8wffpebb8YdNmTZPQi5CM9ZI14uLK4alSTyaWjEq5lYuRW+0wQGCYwSmCZAr25j0JYPCbGPWdejwELqUM2mJL+1QgJVPPC8ugzS0fSUNJUfpvTLwEkznI4SGA4mV9MERl/oeHBWrTJZijBjm1JdtyECpX4vizvUTiknkoGl1FXh2t8nVkKmkdm36zeMeAv26E5cmsZ8XLSdxv8r24PMywLnhm+idF7E+6vCwdcWXKT0+5XBrVxJWB07Loxl4qj0/3fp6Al9oq69rDGnsajaJp5U9rYX/N3e9ra9j26tPoW5/t8nVPUA6NtC5efe45X+n9I07P+HdNy9/3+K7ft//+QH8vzc358HN7Td7t8CQWmp4WqeA3zX2TWRfvAAqgp78LRwbU5ri0QyKSzjAjWZhd3Xvm4tZLY+KluEF6y6DIlYcfFIal1cfXF3MzVDPoTRr9wvwfsZqjyBV2UTvesjWid4f3OOL9/HbKubonWl2nXSWM6H4XjMleSe7KVFY+cLZlkIXPIczZOxWDgcPtpRTbeuznbWWWedfWr7JwAA//8oRRHCABoAAA==`

func TestLoadPackage(t *testing.T) {
	c, err := ioutil.ReadAll(base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(testBase64)))
	require.NoError(t, err)
	v, err := LoadPackage(nil, Vmbased, c)
	require.NoError(t, err)
	require.Equal(t, "zk", v.GetName())

	files, err := archiveutil.Load(bytes.NewBuffer(c))
	require.NoError(t, err)
	c, err = archiveutil.Save(files, "test-package")
	require.NoError(t, err)
	v, err = LoadPackage(nil, Vmbased, c)
	require.NoError(t, err)
	require.Equal(t, "zk", v.GetName())
}
