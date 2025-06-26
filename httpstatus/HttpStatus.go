package httpstatus

import "net/http"

var (
	OK          = http.StatusOK
	CREATED     = http.StatusCreated
	BAD_REQUEST = http.StatusBadRequest
	NOT_FOUND   = http.StatusNotFound
	CONFLICT    = http.StatusConflict
	INTERNAL    = http.StatusInternalServerError
)
