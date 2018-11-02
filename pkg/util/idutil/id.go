// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package idutil

import (
	"crypto/rand"
	"math/big"
	"os"

	"github.com/sony/sonyflake"
	"github.com/speps/go-hashids"

	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

var sf *sonyflake.Sonyflake

func init() {
	var st sonyflake.Settings

	var enableRandomSeed = os.Getenv("OPENPITRIX_ID_RANDOM_SEED")
	if enableRandomSeed == "yes" {
		st.MachineID = getRandomMachineID
	}

	sf = sonyflake.NewSonyflake(st)
}

func getRandomMachineID() (uint16, error) {
	for {
		i, err := rand.Int(rand.Reader, big.NewInt(65536))
		if err != nil {
			return 0, err
		}
		mid := i.Uint64()
		if 0 < mid && mid < 65536 {
			return uint16(mid), nil
		}
	}
}

func GetIntId() uint64 {
	id, err := sf.NextID()
	if err != nil {
		panic(err)
	}
	return id
}

// format likes: B6BZVN3mOPvx
func GetUuid(prefix string) string {
	id := GetIntId()
	hd := hashids.NewData()
	h, err := hashids.NewWithData(hd)
	if err != nil {
		panic(err)
	}
	i, err := h.Encode([]int{int(id)})
	if err != nil {
		panic(err)
	}

	return prefix + stringutil.Reverse(i)
}

const (
	Alphabet62 = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	Alphabet36 = "abcdefghijklmnopqrstuvwxyz1234567890"
)

// format likes: 300m50zn91nwz5
func GetUuid36(prefix string) string {
	id := GetIntId()
	hd := hashids.NewData()
	hd.Alphabet = Alphabet36
	h, err := hashids.NewWithData(hd)
	if err != nil {
		panic(err)
	}
	i, err := h.Encode([]int{int(id)})
	if err != nil {
		panic(err)
	}

	return prefix + stringutil.Reverse(i)
}

func randString(letters string, n int) string {
	output := make([]byte, n)

	// We will take n bytes, one byte for each character of output.
	randomness := make([]byte, n)

	// read all random
	_, err := rand.Read(randomness)
	if err != nil {
		panic(err)
	}

	l := len(letters)
	// fill output
	for pos := range output {
		// get random item
		random := uint8(randomness[pos])

		// random % 64
		randomPos := random % uint8(l)

		// put into output
		output[pos] = letters[randomPos]
	}

	return string(output)
}

func GetSecret() string {
	return randString(Alphabet62, 50)
}

func GetRefreshToken() string {
	return randString(Alphabet62, 50)
}
