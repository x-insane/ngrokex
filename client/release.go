// +build release

package client

var (
	rootCrtPaths = []string{"assets/client/tls/ca.crt"}
)

func useInsecureSkipVerify() bool {
	return false
}
