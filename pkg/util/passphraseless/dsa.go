// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package passphraseless

import (
	"crypto/dsa"
	"crypto/rand"
	"encoding/asn1"
	"encoding/pem"
	"math/big"
)

func GenerateDsaKey() (*dsa.PrivateKey, *dsa.PublicKey, error) {
	var private dsa.PrivateKey
	params := &private.Parameters
	err := dsa.GenerateParameters(params, rand.Reader, dsa.L1024N160)
	if err != nil {
		return nil, nil, err
	}
	err = dsa.GenerateKey(&private, rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return &private, &private.PublicKey, nil
}

func EncodeDsaPrivateKey(pkey *dsa.PrivateKey) []byte {
	var key struct {
		Version int
		P       *big.Int
		Q       *big.Int
		G       *big.Int
		Priv    *big.Int
		Pub     *big.Int
	}

	key.P = pkey.PublicKey.P
	key.Q = pkey.PublicKey.Q
	key.G = pkey.PublicKey.G
	key.Priv = pkey.Y
	key.Pub = pkey.X
	key.Version = 0

	value, _ := asn1.Marshal(key)
	return pem.EncodeToMemory(&pem.Block{
		Bytes: value,
		Type:  "DSA PRIVATE KEY",
	})
}

func MakeSSHDsaKeyPair() (string, string, error) {
	pkey, pubkey, err := GenerateDsaKey()
	if err != nil {
		return "", "", err
	}

	pub, err := EncodeSSHKey(pubkey)
	if err != nil {
		return "", "", err
	}

	return string(EncodeDsaPrivateKey(pkey)), string(pub), nil
}
