package exception

import (
	"github.com/isaacwallace123/GoUtils/timeutil"
	"github.com/isaacwallace123/GoWeb/HttpStatus"
	"github.com/isaacwallace123/GoWeb/ResponseEntity"
)

func GenericHTTPError(status int, message string) *ResponseEntity.ResponseEntity {
	return ResponseEntity.Status(status).Body(map[string]any{"status": status, "message": message, "timestamp": timeutil.NowUTC()})
}

func BadRequestException(message string) *ResponseEntity.ResponseEntity {
	return GenericHTTPError(HttpStatus.BAD_REQUEST, message)
}

func NotFoundException(message string) *ResponseEntity.ResponseEntity {
	return GenericHTTPError(HttpStatus.NOT_FOUND, message)
}

func InternalServerException(message string) *ResponseEntity.ResponseEntity {
	return GenericHTTPError(HttpStatus.INTERNAL_SERVER_ERR, message)
}

func MethodNotAllowed(message string) *ResponseEntity.ResponseEntity {
	return GenericHTTPError(HttpStatus.METHOD_NOT_ALLOWED, message)
}
