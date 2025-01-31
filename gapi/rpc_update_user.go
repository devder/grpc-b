package gapi

import (
	"context"
	"errors"
	"time"

	db "github.com/devder/grpc-b/db/sqlc"
	"github.com/devder/grpc-b/pb"
	"github.com/devder/grpc-b/util"
	myValidator "github.com/devder/grpc-b/validator"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	authPayload, err := server.authorizeUser(ctx, []string{util.BankerRole, util.DepositorRole})
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	if violations := validateUpdateUserInput(req); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	arg := db.UpdateUserParams{
		Username: authPayload.Username,
		FullName: pgtype.Text{
			Valid:  req.FullName != nil,
			String: req.GetFullName(),
		},
		Email: pgtype.Text{
			Valid:  req.Email != nil,
			String: req.GetEmail(),
		},
	}

	if req.Password != nil {
		hashedPassword, err := util.HashPassword(req.GetPassword()) // GetPassword() is auto generated
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password : %s", err)
		}

		arg.HashedPassword = pgtype.Text{
			String: hashedPassword,
			Valid:  true,
		}

		arg.PasswordChangedAt = pgtype.Timestamptz{
			Valid: true,
			Time:  time.Now(),
		}

	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update user : %s", err)
	}

	rsp := &pb.UpdateUserResponse{
		User: convertUser(user),
	}

	return rsp, nil
}

func validateUpdateUserInput(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if req.Password != nil {
		if err := myValidator.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	}

	if req.FullName != nil {
		if err := myValidator.ValidateFullname(req.GetFullName()); err != nil {
			// full_name as defined in the proto file
			violations = append(violations, fieldViolation("full_name", err))
		}
	}

	if req.Email != nil {
		if err := myValidator.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	}

	return violations
}
