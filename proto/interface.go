package proto

import (
	"github.com/x-insane/ngrokex/conn"
)

type Protocol interface {
	GetName() string
	WrapConn(conn.Conn, interface{}) conn.Conn
}
