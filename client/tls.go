package client

import (
	_ "crypto/sha512"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"github.com/x-insane/ngrokex/client/assets"
)

const (
	ClientCrtAssetFile = "assets/client/tls/client.crt"
	ClientKeyAssetFile = "assets/client/tls/client.key"
)

func LoadTLSConfig(caPaths []string, crtPath string, keyPath string) (*tls.Config, error) {
	// load ca cert
	pool := x509.NewCertPool()
	for _, certPath := range caPaths {
		rootCrt, err := assets.Asset(certPath)
		if err != nil {
			return nil, err
		}

		pemBlock, _ := pem.Decode(rootCrt)
		if pemBlock == nil {
			return nil, fmt.Errorf("Bad PEM data")
		}

		certs, err := x509.ParseCertificates(pemBlock.Bytes)
		if err != nil {
			return nil, err
		}

		pool.AddCert(certs[0])
	}

	// load client cert
	var (
		crt  []byte
		key  []byte
		err  error
		cert tls.Certificate
	)
	if crt, err = fileOrAsset(crtPath, ClientCrtAssetFile); err != nil {
		return nil, err
	}
	if key, err = fileOrAsset(keyPath, ClientKeyAssetFile); err != nil {
		return nil, err
	}
	if cert, err = tls.X509KeyPair(crt, key); err != nil {
		return nil, err
	}

	return &tls.Config{
		RootCAs: pool,
		Certificates: []tls.Certificate{cert},
	}, nil
}

func fileOrAsset(path string, defaultPath string) ([]byte, error) {
	loadFn := ioutil.ReadFile
	if path == "" {
		loadFn = assets.Asset
		path = defaultPath
	}
	return loadFn(path)
}