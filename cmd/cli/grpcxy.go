package cli

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"

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
		var opts []grpc.ServerOption

		lis, err := net.Listen("tcp", configFromYaml.GrpcProxy.Listener)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		logger.LogInfo(fmt.Sprintf("Proxy running at %q\n", configFromYaml.GrpcProxy.Listener))
		simpleBackendGen := func(hostname string) proxy.Backend {
			return &proxy.SingleBackend{
				GetConn: func(ctx context.Context) (context.Context, *grpc.ClientConn, error) {
					md, _ := metadata.FromIncomingContext(ctx)

					outCtx := metadata.NewOutgoingContext(ctx, md.Copy())
					if configFromYaml.GrpcSSL.Enable {
						creds, sslErr := credentials.NewClientTLSFromFile(
							configFromYaml.GrpcClientCert, "")
						if sslErr != nil {
							log.Fatalf("Failed to parse credentials: %v", sslErr)
						}
						conn, err := grpc.DialContext(
							ctx,
							hostname,
							grpc.WithTransportCredentials(creds),
							grpc.WithCodec(proxy.Codec()),
						) //nolint: staticcheck
						return outCtx, conn, err
					}
					conn, err := grpc.DialContext(
						ctx,
						hostname,
						grpc.WithInsecure(),
						grpc.WithCodec(proxy.Codec()),
					) //nolint: staticcheck

					return outCtx, conn, err
				},
			}
		}

		director = func(ctx context.Context, fullMethodName string) (proxy.Mode, []proxy.Backend, error) {
			for _, bkd := range configFromYaml.GrpcEndpoints {
				// Make sure we never forward internal services.
				if !strings.HasPrefix(fullMethodName, bkd.Name) {
					return proxy.One2One, nil, status.Errorf(codes.Unimplemented, "Unknown method")
				}
				md, ok := metadata.FromIncomingContext(ctx)
				if ok {
					if _, exists := md[":authority"]; exists {
						return proxy.One2One, []proxy.Backend{
							simpleBackendGen(bkd.HostURI),
						}, nil
					}
				}
			}
			return proxy.One2One, nil, status.Errorf(codes.Unimplemented, "Unknown method")
		}
		opts = append(opts,
			grpc.CustomCodec(proxy.Codec()),
			grpc.UnknownServiceHandler(proxy.TransparentHandler(director)),
		)
		// SSL
		if configFromYaml.GrpcSSL.Enable {

			creds, sslErr := credentials.NewServerTLSFromFile(
				configFromYaml.GrpcSSL.CrtFile,
				configFromYaml.GrpcSSL.KeyFile,
			)
			if sslErr != nil {
				log.Fatalf("Failed to parse credentials: %v", sslErr)
				return
			}
			opts = append(opts, grpc.Creds(creds))
		}

		server := grpc.NewServer(opts...)

		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	},
}

func init() {

	rootCmd.AddCommand(grpcCmd)
}
