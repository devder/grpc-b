package gapi

import (
	"context"
	"time"

	db "github.com/devder/grpc-b/db/sqlc"
	"github.com/devder/grpc-b/pb"
	"github.com/devder/grpc-b/util"
	myValidator "github.com/devder/grpc-b/validator"
	"github.com/devder/grpc-b/worker"
	"github.com/hibiken/asynq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if violations := validateCreateUserInput(req); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := util.HashPassword(req.GetPassword()) // GetPassword() is auto generated

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password : %s", err)
	}

	arg := db.CreateUserTxParams{
		CreateUserParams: db.CreateUserParams{
			Username:       req.GetUsername(),
			HashedPassword: hashedPassword,
			FullName:       req.GetFullName(),
			Email:          req.GetEmail(),
		},
		AfterCreate: func(user db.User) error {
			// send verify email to user
			taskPayload := &worker.PayloadSendVerifyEmail{
				Username: user.Username,
			}
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(5 * time.Second),  // add delay
				asynq.Queue(worker.QueueCritical), // tell the task to use the critical queue
			}

			return server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
		},
	}

	createUserTxRes, err := server.store.CreateUserTx(ctx, arg)
	if err != nil {
		if db.ErrCode(err) == db.UniqueViolation {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}

		return nil, status.Errorf(codes.Internal, "failed to create user : %s", err)
	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(createUserTxRes.User),
	}

	return rsp, nil
}

func validateCreateUserInput(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := myValidator.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}

	if err := myValidator.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	if err := myValidator.ValidateFullname(req.GetFullName()); err != nil {
		// full_name as defined in the proto file
		violations = append(violations, fieldViolation("full_name", err))
	}

	if err := myValidator.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}
