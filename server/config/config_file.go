package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Configuration struct {
	Admin    AdminConf  `yaml:"admin"`
	Tunnel   TunnelConf `yaml:"tunnel"`
	Http     HttpConf   `yaml:"http"`
	LogTo    string     `yaml:"log_to"`
	LogLevel string     `yaml:"log_level"`
}

// config for admin settings
type AdminConf struct {
	HttpPort  int64 `yaml:"http_port"`
	HttpsPort int64 `yaml:"https_port"`
}

// config for the tunnel between ngrok server & client
type TunnelConf struct {
	TunnelAddr string   `yaml:"tunnel_addr"`
	CACert     CertPath `yaml:"ca"`
	TlsCert    CertPath `yaml:"tls"`
}

// config for public HTTP connection
type HttpConf struct {
	Domain    string   `yaml:"domain"`
	HttpAddr  string   `yaml:"http_addr"`
	HttpsAddr string   `yaml:"https_addr"`
	TlsCert   CertPath `yaml:"https_tls"`
}

// rsa key pair file path
type CertPath struct {
	CrtPath string `yaml:"crt"`
	KeyPath string `yaml:"key"`
}

func ParseConfig() *Configuration {
	var config Configuration

	// set default values
	config.LogTo = "stdout" // log to stdout by default
	config.LogLevel = "DEBUG"

	if len(os.Args) < 2 {
		fmt.Printf("usage: %s <config_file>\n", os.Args[0])
		os.Exit(0)
	} else {
		bytes, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}

		err = yaml.Unmarshal(bytes, &config)
		if err != nil {
			panic(err)
		}
	}

	return &config
}
