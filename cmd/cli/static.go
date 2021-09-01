package cli

import (
	"net/http"
	"os"

	"github.com/kenriortega/ngonx/pkg/httpsrv"
	"github.com/spf13/cobra"
)

var staticCmd = &cobra.Command{
	Use:   "static",
	Short: "Run ngonx as a static web server",
	Run: func(cmd *cobra.Command, args []string) {
		frontEnd := os.DirFS(configFromYaml.StaticServer.StaticFile)
		if configFromYaml.ServerSSL.Enable {

			http.Handle("/", http.FileServer(http.FS(frontEnd)))

			portSSL := configFromYaml.ServerSSL.SSLPort

			server := httpsrv.NewServerSSL(configFromYaml.StaticServer.Host, portSSL, nil)
			server.StartSSL(
				configFromYaml.ServerSSL.CrtFile,
				configFromYaml.ServerSSL.KeyFile,
			)

		} else {
			http.Handle("/", http.FileServer(http.FS(frontEnd)))

			server := httpsrv.NewServer(
				configFromYaml.StaticServer.Host,
				configFromYaml.StaticServer.Port,
				nil,
			)
			server.Start()
		}

	},
}

func init() {

	rootCmd.AddCommand(staticCmd)
}
