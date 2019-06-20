package server

import (
	"crypto/tls"
	"github.com/x-insane/ngrokex/server/db"
	"math/rand"
	"os"
	"runtime/debug"
	"time"

	"github.com/x-insane/ngrokex/conn"
	"github.com/x-insane/ngrokex/log"
	"github.com/x-insane/ngrokex/msg"
	"github.com/x-insane/ngrokex/server/admin_http_handler"
	"github.com/x-insane/ngrokex/server/config"
	"github.com/x-insane/ngrokex/util"
)

const (
	registryCacheSize uint64        = 1024 * 1024 // 1 MB
	connReadTimeout   time.Duration = 10 * time.Second
)

// GLOBALS
var (
	tunnelRegistry  *TunnelRegistry
	controlRegistry *ControlRegistry

	// XXX: kill these global variables - they're only used in tunnel.go for constructing forwarding URLs
	conf      *config.Configuration
	listeners map[string]*conn.Listener
)

func NewProxy(pxyConn conn.Conn, regPxy *msg.RegProxy) {
	// fail gracefully if the proxy connection fails to register
	defer func() {
		if r := recover(); r != nil {
			_ = pxyConn.Warn("Failed with error: %v", r)
			_ = pxyConn.Close()
		}
	}()

	// set logging prefix
	pxyConn.SetType("pxy")

	// look up the control connection for this proxy
	pxyConn.Info("Registering new proxy for %s", regPxy.ClientId)
	ctl := controlRegistry.Get(regPxy.ClientId)

	if ctl == nil {
		panic("No client found for identifier: " + regPxy.ClientId)
	}

	ctl.RegisterProxy(pxyConn)
}

// Listen for incoming control and proxy connections
// We listen for incoming control and proxy connections on the same port
// for ease of deployment. The hope is that by running on port 443, using
// TLS and running all connections over the same port, we can bust through
// restrictive firewalls.
func tunnelListener(addr string, tlsConfig *tls.Config) {
	// listen for incoming connections
	listener, err := conn.Listen(addr, "tun", tlsConfig)
	if err != nil {
		panic(err)
	}

	log.Info("Listening for control and proxy connections on %s", listener.Addr.String())
	for c := range listener.Conns {
		go func(tunnelConn conn.Conn) {
			// don't crash on panics
			defer func() {
				if r := recover(); r != nil {
					tunnelConn.Info("tunnelListener failed with error %v: %s", r, debug.Stack())
				}
			}()

			_ = tunnelConn.SetReadDeadline(time.Now().Add(connReadTimeout))
			var rawMsg msg.Message
			if rawMsg, err = msg.ReadMsg(tunnelConn); err != nil {
				_ = tunnelConn.Warn("Failed to read message: %v", err)
				_ = tunnelConn.Close()
				return
			}

			// don't timeout after the initial read
			// tunnel heartbeat will kill dead connections
			_ = tunnelConn.SetReadDeadline(time.Time{})

			switch m := rawMsg.(type) {
			case *msg.Auth:
				NewControl(tunnelConn, m)

			case *msg.RegProxy:
				NewProxy(tunnelConn, m)

			default:
				_ = tunnelConn.Close()
			}
		}(c)
	}
}

func Main() {
	// parse conf
	conf = config.ParseConfig()

	// init database
	db.InitTables()

	// init logging
	log.LogTo(conf.LogTo, conf.LogLevel)

	// seed random number generator
	seed, err := util.RandomSeed()
	if err != nil {
		panic(err)
	}
	rand.Seed(seed)

	// init tunnel/control registry
	registryCacheFile := os.Getenv("REGISTRY_CACHE_FILE")
	tunnelRegistry = NewTunnelRegistry(registryCacheSize, registryCacheFile)
	controlRegistry = NewControlRegistry()

	// start listeners
	listeners = make(map[string]*conn.Listener)

	// load server tls configuration
	serverTlsConfig, err := LoadServerTLSConfig(conf.Tunnel.CACert.CrtPath,
		conf.Tunnel.TlsCert.CrtPath, conf.Tunnel.TlsCert.KeyPath)
	if err != nil {
		panic(err)
	}

	// listen for http
	if conf.Http.HttpAddr != "" {
		listeners["http"] = startHttpListener(conf.Http.HttpAddr, nil)
	}

	// listen for https
	if conf.Http.HttpsAddr != "" {
		// load public tls configuration
		publicTlsConfig, err := LoadPublicTLSConfig(conf.Http.TlsCert.CrtPath, conf.Http.TlsCert.KeyPath)
		if err != nil {
			panic(err)
		}
		listeners["https"] = startHttpListener(conf.Http.HttpsAddr, publicTlsConfig)
	}

	// admin http listener
	go admin_http_handler.AdminHttpListener(conf.Admin)

	// ngrok clients
	tunnelListener(conf.Tunnel.TunnelAddr, serverTlsConfig)
}
