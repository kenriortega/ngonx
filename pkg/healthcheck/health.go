package healthcheck

import (
	"net"
	"net/url"
	"time"

	"github.com/kenriortega/ngonx/pkg/errors"
	"github.com/kenriortega/ngonx/pkg/logger"
)

// IsBackendAlive checks whether a backend is Alive by establishing a TCP connection
func IsBackendAlive(u *url.URL) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		logger.LogError(errors.ErrIsBackendAlive.Error())
		return false
	}
	_ = conn.Close()
	return true
}
