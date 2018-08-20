// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package jwtutil

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/nats-io/nuid"
)

const (
	issuer      = "sts"
	idTokenType = "ID"
	party       = "iam"
	notBefore   = 0
	acr         = "1"
)

var (
	ErrUnauthorizedAccess = errors.New("missing or invalid credentials provided")
)

func MakeToken(id, secret string, d time.Duration, m map[string]interface{}) (tokenStr string, err error) {
	now := time.Now().UTC()
	exp := now.Add(d)

	m["sub"] = id
	m["iss"] = issuer
	m["nbf"] = notBefore
	m["aud"] = party
	m["azp"] = party
	m["typ"] = idTokenType
	m["acr"] = acr
	m["iat"] = now.Unix()
	m["exp"] = exp.Unix()
	m["jti"] = nuid.Next()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(m))
	tokenStr, err = token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func ValidateToken(tokenStr, secret string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnauthorizedAccess
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return claims, ErrUnauthorizedAccess
	}

	return claims, nil
}
