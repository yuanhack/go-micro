package grpc

import (
	"crypto/tls"

	gc "github.com/macheal/go-micro/v2/client/grpc"
	gs "github.com/macheal/go-micro/v2/server/grpc"
	"github.com/macheal/go-micro/v2/service"
)

// WithTLS sets the TLS config for the service
func WithTLS(t *tls.Config) service.Option {
	return func(o *service.Options) {
		o.Client.Init(
			gc.AuthTLS(t),
		)
		o.Server.Init(
			gs.AuthTLS(t),
		)
	}
}
