package cli

import (
	"net/http"
	"os"

	"github.com/kenriortega/goproxy/pkg/config"
)

func StartStaticServer(
	config config.Config,
) {
	frontEnd := os.DirFS(config.StaticServer.StaticFile)
	if config.ServerSSL.Enable {

		http.Handle("/", http.FileServer(http.FS(frontEnd)))

		portSSL := config.ServerSSL.SSLPort

		server := NewServerSSL(config.StaticServer.Host, portSSL)
		server.StartSSL(
			config.ServerSSL.CrtFile,
			config.ServerSSL.KeyFile,
		)

	} else {
		http.Handle("/", http.FileServer(http.FS(frontEnd)))

		server := NewServer(
			config.StaticServer.Host,
			config.StaticServer.Port,
		)
		server.Start()
	}
}
