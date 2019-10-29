// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package pilotutil

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"openpitrix.io/openpitrix/pkg/logger"
	pbpilot "openpitrix.io/openpitrix/pkg/pb/metadata/pilot"
	pbtypes "openpitrix.io/openpitrix/pkg/pb/metadata/types"
	"openpitrix.io/openpitrix/pkg/util/tlsutil"
)

func MustLoadPilotConfig(path string) *pbtypes.PilotConfig {
	p, err := LoadPilotConfig(path)
	if err != nil {
		logger.Critical(nil, "%+v", err)
		os.Exit(1)
	}
	return p
}

func LoadPilotConfig(path string) (*pbtypes.PilotConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	p := new(pbtypes.PilotConfig)
	if err := json.Unmarshal(data, p); err != nil {
		return nil, err
	}

	return p, nil
}

func DialPilotService(ctx context.Context, host string, port int) (
	client pbpilot.PilotServiceClient,
	conn *grpc.ClientConn,
	err error,
) {
	conn, err = grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())
	if err != nil {
		return
	}

	client = pbpilot.NewPilotServiceClient(conn)
	return
}

func DialPilotServiceForFrontgate_TLS(
	ctx context.Context, host string, port int,
	tlsConfig *tls.Config,
) (
	client pbpilot.PilotServiceForFrontgateClient,
	conn *grpc.ClientConn,
	err error,
) {
	creds := credentials.NewTLS(tlsConfig)
	conn, err = grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithTransportCredentials(creds))
	if err != nil {
		return
	}

	client = pbpilot.NewPilotServiceForFrontgateClient(conn)
	return
}

func LoadPilotTLSConfig(
	serverCertFile, serverKeyFile string,
	clientCertFile, clientKeyFile string,
	caCertFile string,
	tlsServerName string,
) (p *pbtypes.PilotTLSConfig, err error) {

	caData, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		return nil, err
	}

	serverCrtData, err := ioutil.ReadFile(serverCertFile)
	if err != nil {
		return nil, err
	}
	serverKeyData, err := ioutil.ReadFile(serverKeyFile)
	if err != nil {
		return nil, err
	}

	clientCrtData, err := ioutil.ReadFile(clientCertFile)
	if err != nil {
		return nil, err
	}
	clientKeyData, err := ioutil.ReadFile(clientKeyFile)
	if err != nil {
		return nil, err
	}

	p = &pbtypes.PilotTLSConfig{
		CaCrtData:       string(caData),
		ServerCrtData:   string(serverCrtData),
		ServerKeyData:   string(serverKeyData),
		ClientCrtData:   string(clientCrtData),
		ClientKeyData:   string(clientKeyData),
		PilotServerName: tlsServerName,
	}

	return p, nil
}

func LoadPilotClientTLSConfig(
	certFile, keyFile, caCertFile, tlsServerName string,
) (p *pbtypes.PilotClientTLSConfig, err error) {

	caData, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		return nil, err
	}

	clientCrtData, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, err
	}
	clientKeyData, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	p = &pbtypes.PilotClientTLSConfig{
		CaCrtData:       string(caData),
		ClientCrtData:   string(clientCrtData),
		ClientKeyData:   string(clientKeyData),
		PilotServerName: tlsServerName,
	}

	return p, nil
}

func NewServerTLSConfigFromPbConfig(pbcfg *pbtypes.PilotTLSConfig) (*tls.Config, error) {
	return tlsutil.NewServerTLSConfigFromString(
		pbcfg.ServerCrtData,
		pbcfg.ServerKeyData,
		pbcfg.CaCrtData,
	)
}

func NewClientTLSConfigFromPbConfig(pbcfg *pbtypes.PilotClientTLSConfig) (*tls.Config, error) {
	return tlsutil.NewClientTLSConfigFromString(
		pbcfg.ClientCrtData,
		pbcfg.ClientKeyData,
		pbcfg.CaCrtData,
		pbcfg.PilotServerName,
	)
}
