package healthchecks

import (
	"net"
	"net/http"
	"time"
)

func CheckServerAlive(host string, r *http.Request) bool {
	const retries = 4

	conn, err := net.DialTimeout("tcp", host, 2*time.Second)

	if err != nil {
		return false
	}

	conn.Close()
	return true
}
