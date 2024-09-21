package api

import (
	database "github.com/DEVunderdog/user-management-gRPC/database/sqlc"
	"github.com/DEVunderdog/user-management-gRPC/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user database.User)  *pb.User {
	return &pb.User{
		Email: user.Email,
		CreatedAt: timestamppb.New(user.CreatedAt.Time),
		UpdatedAt: timestamppb.New(user.CreatedAt.Time),
	}
}