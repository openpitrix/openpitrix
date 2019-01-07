// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package jwtutil

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"

	"openpitrix.io/openpitrix/pkg/sender"
)

var ErrExpired = fmt.Errorf("access token expired")

func trimKey(k string) []byte {
	return []byte(strings.TrimSpace(k))
}

func Validate(k, str string) (*sender.Sender, error) {
	tok, err := jwt.ParseSigned(str)
	if err != nil {
		return nil, err
	}
	c := &jwt.Claims{}
	s := &sender.Sender{}
	err = tok.Claims(trimKey(k), c, s)
	if err != nil {
		return nil, err
	}
	if c.Expiry.Time().Unix() < time.Now().Unix() {
		return nil, ErrExpired
	}
	s.UserId = c.Subject
	return s, nil
}

func Generate(k string, expire time.Duration, userId, role string) (string, error) {
	// TODO: use RS512 or ES512 to encrypt token
	// https://auth0.com/blog/brute-forcing-hs256-is-possible-the-importance-of-using-strong-keys-to-sign-jwts/

	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS512, Key: trimKey(k)}, nil)
	if err != nil {
		return "", err
	}
	s := &sender.Sender{
		Role: role,
	}
	now := time.Now()
	c := &jwt.Claims{
		IssuedAt: jwt.NewNumericDate(now),
		Expiry:   jwt.NewNumericDate(now.Add(expire)),
		// TODO: add jti
		Subject: userId,
	}
	return jwt.Signed(signer).Claims(s).Claims(c).CompactSerialize()
}
