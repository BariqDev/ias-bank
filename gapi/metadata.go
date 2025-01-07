package gapi

import (
	"context"
	"log"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcClient = "grpc-client"
	httpUserAgent = "grpcgateway-user-agent"
	httpForwardedFor = "x-forwarded-for"
)

type MetaData struct {
	ClientIp  string
	UserAgent string
}

func (server *Server) extractMetaData(ctx context.Context) *MetaData {
	mtd := &MetaData{}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Print(md)

		// get  meta for grpc client
		if grpcClients := md.Get(grpcClient); len(grpcClients) > 0 {
			mtd.UserAgent = grpcClients[0]
			log.Println("UserAgent: ", mtd.UserAgent)
		}
		log.Println("check 1")

		if peer, ok := peer.FromContext(ctx); ok {
			mtd.ClientIp = peer.Addr.String()
		}
		log.Println("check 2")

		// get meta for http client

		if userAgents := md.Get(httpUserAgent); len(userAgents) > 0 {
			mtd.UserAgent = userAgents[0]
		}
		log.Println("check 3")

		if ClientIps := md.Get(httpForwardedFor); len(ClientIps) > 0 {
			mtd.ClientIp = ClientIps[0]
			log.Println("UserAgent: ", mtd.UserAgent)
		}
		log.Println("check 4")

	}
	return mtd
}
