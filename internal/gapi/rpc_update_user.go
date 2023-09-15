package gapi

import (
	"context"
	"errors"
	"log"

	"BankApplication/internal/db"
	"BankApplication/internal/pb"
	"BankApplication/internal/util"
	"BankApplication/internal/val"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	violations := validateUpdateUserRequest(req)
	if violations != nil {
		log.Printf("Invalid request: %v", violations)
		return nil, invalidArgumentError(violations)
	}

	arg := db.UpdateUserParams{
		Username: req.GetUsername(),
		FullName: pgtype.Text{String: req.GetFullName(), Valid: req.FullName != nil},
		Email:    pgtype.Text{String: req.GetEmail(), Valid: req.Email != nil},
	}

	if req.Password != nil {
		hashedPassword, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
		}
		arg.HashedPassword = pgtype.Text{String: hashedPassword, Valid: true}
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "user %s not found", req.Username)
		}
		return nil, status.Errorf(codes.Internal, "failed to update user: %s", err)
	}
	rsp := &pb.UpdateUserResponse{
		User: convertUser(user),
	}
	return rsp, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.GetUsername()); err != nil {
		violations = append(violations, filedViolation("username", err))
	}
	if req.Password != nil {
		if err := val.ValidatePassword(req.GetPassword()); err != nil {
			violations = append(violations, filedViolation("password", err))
		}
	}
	if req.FullName != nil {
		if err := val.ValidateUserFullName(req.GetFullName()); err != nil {
			violations = append(violations, filedViolation("full_name", err))
		}
	}
	if req.Email != nil {
		if err := val.ValidateEmail(req.GetEmail()); err != nil {
			violations = append(violations, filedViolation("email", err))
		}
	}
	return violations
}
