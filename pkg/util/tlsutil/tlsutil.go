// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package tlsutil

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
)

func NewTLSConfigFromFile(certFile, keyFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	cfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	return cfg, nil
}

func NewServerTLSConfigFromFile(server_crt, server_key, ca_crt string) (*tls.Config, error) {
	certificate, err := tls.LoadX509KeyPair(server_crt, server_key)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(ca_crt)
	if err != nil {
		return nil, err
	}

	// Append the client certificates from the CA
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		return nil, fmt.Errorf("failed to append client certs")
	}

	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert, // NOTE: this is optional!
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
	}
	return tlsConfig, nil
}

func NewServerTLSConfigFromString(certData, keyData, caCertData string) (*tls.Config, error) {
	certificate, err := tls.X509KeyPair([]byte(certData), []byte(keyData))
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM([]byte(caCertData)); !ok {
		err = fmt.Errorf("failed to append ca certs")
		return nil, err
	}

	// Append the client certificates from the CA
	if ok := certPool.AppendCertsFromPEM([]byte(caCertData)); !ok {
		return nil, fmt.Errorf("failed to append client certs")
	}

	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert, // NOTE: this is optional!
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
	}
	return tlsConfig, nil
}

func NewClientTLSConfigFromFile(certFile, keyFile, caCertFile, tlsServerName string) (*tls.Config, error) {
	certificate, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		return nil, err
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		err = fmt.Errorf("failed to append ca certs")
		return nil, err
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,         // NOTE: this is required!
		ServerName:         tlsServerName, // NOTE: this is required!
		Certificates:       []tls.Certificate{certificate},
		RootCAs:            certPool,
	}
	return tlsConfig, nil
}

func NewClientTLSConfigFromString(certData, keyData, caCertData, tlsServerName string) (*tls.Config, error) {
	certificate, err := tls.X509KeyPair([]byte(certData), []byte(keyData))
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM([]byte(caCertData)); !ok {
		err = fmt.Errorf("failed to append ca certs")
		return nil, err
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         tlsServerName, // NOTE: this is required!
		Certificates:       []tls.Certificate{certificate},
		RootCAs:            certPool,
	}

	return tlsConfig, nil
}
