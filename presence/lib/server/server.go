package server

import "context"

type Server interface {
	// Listen blocks until an error occurs or the context closes
	Listen(ctx context.Context) error
}
