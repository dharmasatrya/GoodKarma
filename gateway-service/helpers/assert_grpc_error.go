package helpers

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrResponse struct {
	Status int    `json:"status" extensions:"x-order=0"`
	Type   string `json:"type" extensions:"x-order=1"`
	Detail string `json:"detail" extensions:"x-order=2"`
}

var (
	ErrBadRequest = ErrResponse{
		Status: http.StatusBadRequest,
		Type:   "Bad Request",
	}

	ErrUnauthorized = ErrResponse{
		Status: http.StatusUnauthorized,
		Type:   "Unauthorized",
	}

	ErrForbidden = ErrResponse{
		Status: http.StatusForbidden,
		Type:   "Forbidden",
	}

	ErrNotFound = ErrResponse{
		Status: http.StatusNotFound,
		Type:   "Not Found",
	}

	ErrConflict = ErrResponse{
		Status: http.StatusConflict,
		Type:   "Conflict",
	}

	ErrUnprocessable = ErrResponse{
		Status: http.StatusUnprocessableEntity,
		Type:   "Unprocessable Content",
	}

	ErrInternalServer = ErrResponse{
		Status: http.StatusInternalServerError,
		Type:   "Internal server error",
	}
)

func (er ErrResponse) New(detail string) *ErrResponse {
	return &ErrResponse{
		Status: er.Status,
		Type:   er.Type,
		Detail: detail,
	}
}

func (er *ErrResponse) EchoFormat() (int, ErrResponse) {
	return er.Status, *er
}

func (er ErrResponse) EchoFormatDetails(detail string) (int, ErrResponse) {
	er.Detail = detail
	return er.Status, er
}

func AssertGrpcStatus(err error) int {
	if e, ok := status.FromError(err); ok {
		switch e.Code() {
		case codes.Unavailable:
			return http.StatusServiceUnavailable // 503
		case codes.InvalidArgument:
			return http.StatusBadRequest // 400
		case codes.FailedPrecondition:
			return http.StatusBadRequest // 400
		case codes.NotFound:
			return http.StatusNotFound // 404
		case codes.PermissionDenied:
			return http.StatusForbidden // 403
		case codes.Unauthenticated:
			return http.StatusUnauthorized // 401
		case codes.Internal:
			return http.StatusInternalServerError // 500
		default:
			return http.StatusInternalServerError // 500 for unhandled cases
		}
	}
	// Return 500 if the error is not a gRPC error
	return http.StatusInternalServerError
}

func AssertJSONStatus(status int) map[string]interface{} {
	switch status {
	case http.StatusServiceUnavailable: // 503
		return map[string]interface{}{
			"Status":  http.StatusServiceUnavailable,
			"Message": "Unavailable",
		}
	case http.StatusBadRequest: // 400
		return map[string]interface{}{
			"Status":  http.StatusBadRequest,
			"Message": "Invalid argument or failed precondition",
		}
	case http.StatusNotFound: // 404
		return map[string]interface{}{
			"Status":  http.StatusNotFound,
			"Message": "ID Event not found!",
		}
	case http.StatusForbidden: // 403
		return map[string]interface{}{
			"Status":  http.StatusForbidden,
			"Message": "Status forbidden!",
		}
	case http.StatusUnauthorized: // 401
		return map[string]interface{}{
			"Status":  http.StatusUnauthorized,
			"Message": "You are not authorized!",
		}
	case http.StatusInternalServerError: // 500
		return map[string]interface{}{
			"Status":  http.StatusInternalServerError,
			"Message": "Internal server error!",
		}
	default:
		// Fallback for any unexpected status code
		return map[string]interface{}{
			"Status":  http.StatusInternalServerError, // Default to 500
			"Message": "Internal server error!",
		}
	}
}
