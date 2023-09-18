package gapi

import (
	"BankApplication/internal/db"
	"BankApplication/internal/pb"
	"BankApplication/internal/val"
	"context"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (server *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	violations := validateVerifyEmailRequest(req)
	if violations != nil {
		log.Printf("Invalid request: %v", violations)
		return nil, invalidArgumentError(violations)
	}
	txResults, err := server.store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		EmailId:    req.GetEmailId(),
		SecretCode: req.GetSecretCode(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to verify email")
	}
	rsp := &pb.VerifyEmailResponse{
		IsVerified: txResults.User.IsEmailVerified,
	}
	return rsp, nil
}

func validateVerifyEmailRequest(req *pb.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateEmailId(req.GetEmailId()); err != nil {
		violations = append(violations, filedViolation("email_id", err))
	}
	if err := val.ValidateSecretCode(req.GetSecretCode()); err != nil {
		violations = append(violations, filedViolation("secret_code", err))
	}
	return violations
}
