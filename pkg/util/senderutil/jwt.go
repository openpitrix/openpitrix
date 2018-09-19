// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package senderutil

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

var ErrExpired = fmt.Errorf("access token expired")

func trimKey(k string) []byte {
	return []byte(strings.TrimSpace(k))
}

func Validate(k, s string) (*Sender, error) {
	tok, err := jwt.ParseSigned(s)
	if err != nil {
		return nil, err
	}
	c := &jwt.Claims{}
	sender := &Sender{}
	err = tok.Claims(trimKey(k), c, sender)
	if err != nil {
		return nil, err
	}
	if c.Expiry.Time().Unix() < time.Now().Unix() {
		return nil, ErrExpired
	}
	sender.UserId = c.Subject
	return sender, nil
}

func Generate(k string, expire time.Duration, userId, role string) (string, error) {
	// TODO: use RS512 or ES512 to encrypt token
	// https://auth0.com/blog/brute-forcing-hs256-is-possible-the-importance-of-using-strong-keys-to-sign-jwts/

	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS512, Key: trimKey(k)}, nil)
	if err != nil {
		return "", err
	}
	sender := &Sender{
		Role: role,
	}
	now := time.Now()
	c := &jwt.Claims{
		IssuedAt: jwt.NewNumericDate(now),
		Expiry:   jwt.NewNumericDate(now.Add(expire)),
		// TODO: add jti
		Subject: userId,
	}
	return jwt.Signed(signer).Claims(sender).Claims(c).CompactSerialize()
}
