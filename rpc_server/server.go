package rpc

import (
	"github.com/bwmarrin/discordgo"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/plally/foxbot_gateway/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
)

func StartRPCServer(session *discordgo.Session) {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	creds, err := credentials.NewServerTLSFromFile("rpc-cert.pem", "rpc-key.pem")
	if err != nil {
		log.Fatal("error loading tls: ", err)
	}

	s := grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(authFunc)),
	)

	proto.RegisterFoxbotGatewayServer(s, &GatewayRPC{discordSession: session})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
