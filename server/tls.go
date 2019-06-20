package server

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/x-insane/ngrokex/server/assets"
)

const (
	PublicCrtAssetFile = "assets/server/tls/public.crt"
	PublicKeyAssetFile = "assets/server/tls/public.key"
	CACrtAssetFile     = "assets/server/tls/ca.crt"
	ServerCrtAssetFile = "assets/server/tls/server.crt"
	ServerKeyAssetFile = "assets/server/tls/server.key"
)

func LoadPublicTLSConfig(crtPath string, keyPath string) (tlsConfig *tls.Config, err error) {
	return loadTLSConfig(crtPath, keyPath, PublicCrtAssetFile, PublicKeyAssetFile)
}

func LoadServerTLSConfig(caPath string, crtPath string, keyPath string) (tlsConfig *tls.Config, err error) {
	tlsConfig, err = loadTLSConfig(crtPath, keyPath, ServerCrtAssetFile, ServerKeyAssetFile)
	if err != nil {
		return
	}

	// load ca cert
	var caBytes []byte
	caBytes, err = fileOrAsset(caPath, CACrtAssetFile)
	if err != nil {
		return
	}
	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caBytes)

	tlsConfig.ClientCAs = caPool
	tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	return
}

func fileOrAsset(path string, defaultPath string) ([]byte, error) {
	loadFn := ioutil.ReadFile
	if path == "" {
		loadFn = assets.Asset
		path = defaultPath
	}
	return loadFn(path)
}

func loadTLSConfig(crtPath string, keyPath string, defaultAssetCrt string, defaultAssetKey string) (tlsConfig *tls.Config, err error) {
	var (
		crt  []byte
		key  []byte
		cert tls.Certificate
	)

	if crt, err = fileOrAsset(crtPath, defaultAssetCrt); err != nil {
		return
	}

	if key, err = fileOrAsset(keyPath, defaultAssetKey); err != nil {
		return
	}

	if cert, err = tls.X509KeyPair(crt, key); err != nil {
		return
	}

	tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	return
}
