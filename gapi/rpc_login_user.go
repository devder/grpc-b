package gapi

import (
	"context"
	"errors"
	"fmt"

	db "github.com/devder/grpc-b/db/sqlc"
	"github.com/devder/grpc-b/pb"
	"github.com/devder/grpc-b/util"
	myValidator "github.com/devder/grpc-b/validator"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	if violations := validateLoginUserInput(req); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "incorrect password or user not found")
		}

		fmt.Println("err :>", err)
		return nil, status.Errorf(codes.Internal, "failed to find user")
	}

	err = util.CheckPassword(user.HashedPassword, req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "incorrect password or user not found")
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		req.GetUsername(),
		user.Role,
		server.config.AccessTokenDuration,
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token")
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		req.GetUsername(),
		user.Role,
		server.config.RefreshTokenDuration,
	)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token")
	}

	mtdt := server.extractMetadata(ctx)
	arg := db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIP,
		ExpiresAt:    refreshPayload.ExpiredAt,
	}

	session, err := server.store.CreateSession(ctx, arg)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session")
	}

	rsp := &pb.LoginUserResponse{
		User:                  convertUser(user),
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiredAt),
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiredAt),
	}

	return rsp, nil
}

func validateLoginUserInput(req *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := myValidator.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := myValidator.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}
