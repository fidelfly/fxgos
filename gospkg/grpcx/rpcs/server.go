package rpcs

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Server struct {
	*grpc.Server
}

type TokenValidator interface {
	Validate(authorization string) bool
}

type FunctionValidator func(string) bool

func (fv FunctionValidator) Validate(authorization string) bool {
	return fv(authorization)
}

func NewServer(validator TokenValidator) *Server {
	return &Server{grpc.NewServer(grpc.UnaryInterceptor(newInterceptor(validator)))}
}

func (s *Server) ListenAndServe(address string) {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("failed to listen:%v", err)
	}

	if err := s.Serve(listen); err != nil {
		fmt.Printf("RPC Server failed to serve :%v", err)
	}
}

func newInterceptor(validator TokenValidator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if validator != nil {
			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
			}
			if tokens, ok := md["authorization"]; ok {
				if len(tokens) > 0 {
					if !validator.Validate(tokens[0]) {
						return nil, status.Errorf(codes.Unauthenticated, "invalid token")
					}
				} else {
					return nil, status.Errorf(codes.Unauthenticated, "invalid token")
				}
			} else {
				return nil, status.Errorf(codes.Unauthenticated, "invalid token")
			}

		}
		// Continue execution of handler after ensuring a valid token.
		return handler(ctx, req)
	}
}
