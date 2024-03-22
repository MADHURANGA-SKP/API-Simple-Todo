package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatwayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader = "user-agent"
	xForwarededForHeader = "x-forwarded-for"
)

type Metadata struct{
	UserAgent string
	ClientIP string
}

func (server *Server) extractMetadata(ctx context.Context) *Metadata{
	mtd := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgents := md.Get(grpcGatwayUserAgentHeader); len(userAgents)>0{
			mtd.UserAgent = userAgents[0]
		}

		if userAgents := md.Get(userAgentHeader); len(userAgents)>0{
			mtd.UserAgent = userAgents[0]
		}

		if clientIPs := md.Get(xForwarededForHeader); len(clientIPs) >0 {
			mtd.ClientIP = clientIPs[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		mtd.ClientIP = p.Addr.String()
	}

	return mtd
}