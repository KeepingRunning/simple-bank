package gapi

import (
	"SimpleBank/pb"
	"SimpleBank/val"
	"context"


	"google.golang.org/genproto/googleapis/rpc/errdetails"
)

func (server *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	vilolations := validateEmailRequest(req)
	if len(vilolations) > 0 {
		return nil, invalidArgumentError(vilolations)
	}

	rsp := &pb.VerifyEmailResponse{}
	return rsp, nil
}

func validateEmailRequest(req *pb.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateEmailId(req.GetEmailId()); err != nil {
		violations = append(violations, fieldViolation("email_id", err))
	}
	
	if err := val.ValidateSecretCode(req.GetSecretCode()); err != nil {
		violations = append(violations, fieldViolation("secret_code", err))
	}
	return violations
}