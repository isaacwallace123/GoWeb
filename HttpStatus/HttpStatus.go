package HttpStatus

import "net/http"

const (
	OK                  = http.StatusOK
	CREATED             = http.StatusCreated
	NO_CONTENT          = http.StatusNoContent
	BAD_REQUEST         = http.StatusBadRequest
	NOT_FOUND           = http.StatusNotFound
	INTERNAL_SERVER_ERR = http.StatusInternalServerError
	METHOD_NOT_ALLOWED  = http.StatusMethodNotAllowed
)
