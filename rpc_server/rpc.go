package rpc

import (
	"context"
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/plally/foxbot_gateway/proto"
	"github.com/plally/foxbot_gateway/role_persister"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

type GatewayRPC struct {
	proto.FoxbotGatewayServer
	discordSession *discordgo.Session
}

func (rpc *GatewayRPC) UnpersistRole(ctx context.Context, data *proto.UnpersistRoleData) (*proto.Empty, error) {
	err := role_persister.DeletePersistedRole(data.GuildId, data.UserId, data.RoleId)
	if err != nil {
		return &proto.Empty{}, status.Error(codes.InvalidArgument, err.Error())
	}

	if !data.RemoveRole {
		return &proto.Empty{}, nil
	}

	err = rpc.discordSession.GuildMemberRoleRemove(data.GuildId, data.UserId, data.RoleId)
	if err != nil {
		return &proto.Empty{}, status.Error(codes.Unknown, err.Error())
	}

	return &proto.Empty{}, nil
}

func (rpc *GatewayRPC) PersistRole(ctx context.Context, data *proto.PersistRoleData) (*proto.Empty, error){
	if time.Now().After(data.Expiration.AsTime())  {
		return &proto.Empty{}, status.Error(codes.InvalidArgument, "expiration before current time")
	}

	err := rpc.discordSession.GuildMemberRoleAdd(data.GuildId, data.UserId, data.RoleId)

	var restError discordgo.RESTError
	if errors.As(err, &restError) {
		if restError.Response.StatusCode == http.StatusForbidden {
			return &proto.Empty{}, status.Error(codes.FailedPrecondition, "bot cannot add role")
		}
	}

	if err != nil {
		return &proto.Empty{}, status.Error(codes.Unknown, err.Error())
	}

	err = role_persister.CreateRolePersist(data.GuildId, data.UserId, data.RoleId, data.Expiration.AsTime())
	if err != nil {
		return &proto.Empty{}, status.Error(codes.Unknown, err.Error())
	}

	return &proto.Empty{}, nil
}