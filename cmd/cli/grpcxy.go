package cli

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/kenriortega/ngonx/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/talos-systems/grpc-proxy/proxy"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var director proxy.StreamDirector
var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Run ngonx as a grpc proxy",
	Run: func(cmd *cobra.Command, args []string) {
		var s *grpc.Server
		port, err := cmd.Flags().GetInt(flagPort)
		if err != nil {
			logger.LogError(err.Error())
		}
		mode, err := cmd.Flags().GetString(flagModeGRPC)
		if err != nil {
			logger.LogError(err.Error())
		}
		lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		switch mode {
		case "transparent":

			if configFromYaml.GrpcProxy.GrpcSSL.Enable {
				creds, sslErr := credentials.NewServerTLSFromFile(
					configFromYaml.GrpcProxy.GrpcSSL.CrtFile,
					configFromYaml.GrpcProxy.GrpcSSL.KeyFile,
				)
				if sslErr != nil {
					log.Fatalf("Failed to parse credentials: %v", sslErr)
					return
				}
				simpleBackendGen := func(hostname string) proxy.Backend {
					return &proxy.SingleBackend{
						GetConn: func(ctx context.Context) (context.Context, *grpc.ClientConn, error) {
							md, _ := metadata.FromIncomingContext(ctx)

							// Copy the inbound metadata explicitly.
							outCtx := metadata.NewOutgoingContext(ctx, md.Copy())
							// Make sure we use DialContext so the dialing can be canceled/time out together with the context.
							conn, err := grpc.DialContext(ctx, hostname, grpc.WithCodec(proxy.Codec())) //nolint: staticcheck

							return outCtx, conn, err
						},
					}
				}

				director = func(ctx context.Context, fullMethodName string) (proxy.Mode, []proxy.Backend, error) {

					_, ok := metadata.FromIncomingContext(ctx)

					if ok {
						// Decide on which backend to dial
						hostUri := configFromYaml.GrpcProxy.GrpcEndpoints[0].HostURI
						return proxy.One2One, []proxy.Backend{
							simpleBackendGen(hostUri),
						}, nil
					}

					return proxy.One2One, nil, status.Errorf(codes.Unimplemented, "Unknown method")
				}
				opts := grpc.Creds(creds)
				s = grpc.NewServer(
					opts,
					grpc.CustomCodec(proxy.Codec()),
					grpc.UnknownServiceHandler(
						proxy.TransparentHandler(director),
					),
				)
				if err := s.Serve(lis); err != nil {
					logger.LogError(err.Error())
				}
			}
		}
	},
}

func init() {
	proxyCmd.Flags().Int(flagPort, 50_000, "Port to run grpc proxy")
	proxyCmd.Flags().String(flagModeGRPC, "transparent", "Action for generate hash for protected routes")

	rootCmd.AddCommand(grpcCmd)
}
