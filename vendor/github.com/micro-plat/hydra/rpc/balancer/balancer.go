package balancer

import "google.golang.org/grpc"

type CustomerBalancer interface {
	grpc.Balancer
	UpdateLimiter(map[string]int)
}
