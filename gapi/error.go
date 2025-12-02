package gapi

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func fieldViolation(field string, err error) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: err.Error(),
	}
}

func invalidArgumentError(violations []*errdetails.BadRequest_FieldViolation) error {
	st := status.New(codes.InvalidArgument, "invalid parameters")
	br := &errdetails.BadRequest{FieldViolations: violations}
	st, err := st.WithDetails(br)
	if err != nil {
		return status.Errorf(codes.Internal, "fail to add details to the error: %s", err)
	}
	return st.Err()
}