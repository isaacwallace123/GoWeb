package exception

import (
	"github.com/isaacwallace123/GoUtils/timeutil"
	"github.com/isaacwallace123/GoWeb/app/types"
	"github.com/isaacwallace123/GoWeb/pkg/HttpStatus"
	"github.com/isaacwallace123/GoWeb/pkg/ResponseEntity"
)

func GenericHTTPError(status int, message string, extras ...any) *types.ResponseEntity {
	payload := map[string]any{
		"status":    status,
		"message":   message,
		"timestamp": timeutil.NowUTC(),
	}

	if len(extras) > 0 {
		payload["error"] = extras[0]
	}

	return ResponseEntity.Status(status).Body(payload)
}

func BadRequestException(message string) *types.ResponseEntity {
	return GenericHTTPError(HttpStatus.BAD_REQUEST, message)
}

func NotFoundException(message string) *types.ResponseEntity {
	return GenericHTTPError(HttpStatus.NOT_FOUND, message)
}

func InternalServerException(message string) *types.ResponseEntity {
	return GenericHTTPError(HttpStatus.INTERNAL_SERVER_ERR, message)
}

func MethodNotAllowed(message string) *types.ResponseEntity {
	return GenericHTTPError(HttpStatus.METHOD_NOT_ALLOWED, message)
}
