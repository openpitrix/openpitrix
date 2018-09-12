// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// +build ignore

package iam

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/nats-io/nuid"
)

type TokenType string

const (
	TokenType_ID            TokenType = "id"
	TokenType_ResetPassword TokenType = "reset-password"
	TokenType_RefreshToken  TokenType = "refresh-token"
)

func (t TokenType) Valid() error {
	switch t {
	case TokenType_ID, TokenType_ResetPassword:
		return nil
	default:
		return fmt.Errorf("jwt-token-type: unknown %v", t)
	}
}

type TokenOptions func(opt *JwtToken)

type JwtToken struct {
	UserId    string    `json:"user-id"`
	ClientId  string    `json:"client-id"`
	TokenType TokenType `json:"token-type"`

	jwt.StandardClaims
}

func (c JwtToken) Valid() error {
	if err := c.StandardClaims.Valid(); err != nil {
		return err
	}

	if c.UserId == "" || c.ClientId == "" {
		return fmt.Errorf("jwt-token: UserId and ClientId is empty")
	}
	if err := c.TokenType.Valid(); err != nil {
		return err
	}

	return nil
}

func MakeJwtToken(secret string, opts ...TokenOptions) (jwtToken string, err error) {
	var claims JwtToken

	for _, fn := range opts {
		fn(&claims)
	}

	if claims.ExpiresAt == 0 {
		claims.ExpiresAt = time.Now().UTC().Add(time.Second).Unix()
	}
	if claims.NotBefore == 0 {
		claims.ExpiresAt = time.Now().UTC().Unix()
	}
	if claims.Id == "" {
		claims.Id = nuid.Next()
	}

	if err := claims.Valid(); err != nil {
		return "", nil
	}

	jwtToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func ValidateJwtToken(tokenStr, secret string) (*JwtToken, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JwtToken{},
		func(t *jwt.Token) (interface{}, error) {
			return secret, nil
		},
	)
	if err != nil {
		return nil, err
	}

	tok, _ := token.Claims.(*JwtToken)
	if err := tok.Valid(); err != nil {
		return tok, nil
	}

	return tok, nil
}

func WithUserId(id string) func(opt *JwtToken) {
	return func(opt *JwtToken) {
		opt.UserId = id
	}
}

func WithClientId(id string) func(opt *JwtToken) {
	return func(opt *JwtToken) {
		opt.ClientId = id
	}
}

func WithTokenType(typ TokenType) func(opt *JwtToken) {
	return func(opt *JwtToken) {
		opt.TokenType = typ
	}
}

func WithTokenExpiresAt(d time.Duration) func(opt *JwtToken) {
	return func(opt *JwtToken) {
		opt.ExpiresAt = time.Now().UTC().Add(d).Unix()
	}
}

func jwtMakeToken(id, secret string, d time.Duration, m map[string]interface{}) (tokenStr string, err error) {
	now := time.Now().UTC()
	exp := now.Add(d)

	m["sub"] = id
	m["iss"] = "sts" //issuer
	m["nbf"] = 0     // not before
	m["aud"] = "iam"
	m["azp"] = "iam"
	m["typ"] = "ID" // idTokenType
	m["acr"] = "1"
	m["iat"] = time.Now().UTC().Unix()
	m["exp"] = exp.Unix()
	m["jti"] = nuid.Next()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(m))
	tokenStr, err = token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func jwtValidateToken(tokenStr, secret string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	return claims, nil
}
