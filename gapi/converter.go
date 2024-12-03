package gapi

import (
	db "github.com/devder/grpc-b/db/sqlc"
	"github.com/devder/grpc-b/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func convertUser(user db.User) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		CreatedAt:         timestamppb.New(user.CreatedAt),
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
	}
}
