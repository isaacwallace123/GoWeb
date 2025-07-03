package exception

import (
	"github.com/isaacwallace123/GoUtils/timeutil"
	"github.com/isaacwallace123/GoWeb/httpstatus"
	"github.com/isaacwallace123/GoWeb/response"
)

func GenericHTTPError(status int, message string) *response.ResponseEntity {
	return response.Status(status).Body(map[string]any{"status": status, "message": message, "timestamp": timeutil.NowUTC()})
}

func BadRequestException(message string) *response.ResponseEntity {
	return GenericHTTPError(httpstatus.BAD_REQUEST, message)
}

func NotFoundException(message string) *response.ResponseEntity {
	return GenericHTTPError(httpstatus.NOT_FOUND, message)
}

func InternalServerException(message string) *response.ResponseEntity {
	return GenericHTTPError(httpstatus.INTERNAL_SERVER_ERR, message)
}

func MethodNotAllowed(message string) *response.ResponseEntity {
	return GenericHTTPError(httpstatus.METHOD_NOT_ALLOWED, message)
}
