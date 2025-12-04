package gapi

import (
	db "SimpleBank/db/sqlc"
	"SimpleBank/pb"
	"SimpleBank/util"
	"SimpleBank/val"
	"context"

	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	vilolations := validateCreateUserRequest(req)
	if len(vilolations) > 0 {
		return nil, invalidArgumentError(vilolations)
	}
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "fail to hash password: %s", err)
	}
	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":

				return nil, status.Errorf(codes.AlreadyExists, "username %s already exists %s", req.GetUsername(), pqErr.Detail)
			}
		}
		return nil, status.Errorf(codes.Internal, "fail to create user: %s", err)
	}
	rsp := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	return rsp, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := val.ValidatePassword(req.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	if err := val.ValidateFullName(req.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}
	if err := val.ValidateEmail(req.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}
	return violations
}