package rpc

import (
	"github.com/bwmarrin/discordgo"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/plally/foxbot_gateway/proto"
	"google.golang.org/grpc"
	"log"
	"net"
)

func StartRPCServer(session *discordgo.Session) {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(authFunc)),
	)

	proto.RegisterFoxbotGatewayServer(s, &GatewayRPC{discordSession: session})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
