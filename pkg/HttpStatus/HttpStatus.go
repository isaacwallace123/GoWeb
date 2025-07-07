package HttpStatus

import "net/http"

// 1xx: Informational
const (
	CONTINUE            = http.StatusContinue
	SWITCHING_PROTOCOLS = http.StatusSwitchingProtocols
	PROCESSING          = http.StatusProcessing
	EARLY_HINTS         = http.StatusEarlyHints
)

// 2xx: Success
const (
	OK                = http.StatusOK
	CREATED           = http.StatusCreated
	ACCEPTED          = http.StatusAccepted
	NON_AUTHORITATIVE = http.StatusNonAuthoritativeInfo
	NO_CONTENT        = http.StatusNoContent
	RESET_CONTENT     = http.StatusResetContent
	PARTIAL_CONTENT   = http.StatusPartialContent
	MULTI_STATUS      = http.StatusMultiStatus
	ALREADY_REPORTED  = http.StatusAlreadyReported
	IM_USED           = http.StatusIMUsed
)

// 3xx: Redirection
const (
	MULTIPLE_CHOICES   = http.StatusMultipleChoices
	MOVED_PERMANENTLY  = http.StatusMovedPermanently
	FOUND              = http.StatusFound
	SEE_OTHER          = http.StatusSeeOther
	NOT_MODIFIED       = http.StatusNotModified
	USE_PROXY          = http.StatusUseProxy
	TEMPORARY_REDIRECT = http.StatusTemporaryRedirect
	PERMANENT_REDIRECT = http.StatusPermanentRedirect
)

// 4xx: Client Errors
const (
	BAD_REQUEST                     = http.StatusBadRequest
	UNAUTHORIZED                    = http.StatusUnauthorized
	PAYMENT_REQUIRED                = http.StatusPaymentRequired
	FORBIDDEN                       = http.StatusForbidden
	NOT_FOUND                       = http.StatusNotFound
	METHOD_NOT_ALLOWED              = http.StatusMethodNotAllowed
	NOT_ACCEPTABLE                  = http.StatusNotAcceptable
	PROXY_AUTH_REQUIRED             = http.StatusProxyAuthRequired
	REQUEST_TIMEOUT                 = http.StatusRequestTimeout
	CONFLICT                        = http.StatusConflict
	GONE                            = http.StatusGone
	LENGTH_REQUIRED                 = http.StatusLengthRequired
	PRECONDITION_FAILED             = http.StatusPreconditionFailed
	REQUEST_ENTITY_TOO_LARGE        = http.StatusRequestEntityTooLarge
	REQUEST_URI_TOO_LONG            = http.StatusRequestURITooLong
	UNSUPPORTED_MEDIA_TYPE          = http.StatusUnsupportedMediaType
	REQUESTED_RANGE_NOT_SATISFIABLE = http.StatusRequestedRangeNotSatisfiable
	EXPECTATION_FAILED              = http.StatusExpectationFailed
	IM_A_TEAPOT                     = http.StatusTeapot
	MISDIRECTED_REQUEST             = http.StatusMisdirectedRequest
	UNPROCESSABLE_ENTITY            = http.StatusUnprocessableEntity
	LOCKED                          = http.StatusLocked
	FAILED_DEPENDENCY               = http.StatusFailedDependency
	TOO_EARLY                       = http.StatusTooEarly
	UPGRADE_REQUIRED                = http.StatusUpgradeRequired
	PRECONDITION_REQUIRED           = http.StatusPreconditionRequired
	TOO_MANY_REQUESTS               = http.StatusTooManyRequests
	REQUEST_HEADER_FIELDS_TOO_LARGE = http.StatusRequestHeaderFieldsTooLarge
	UNAVAILABLE_FOR_LEGAL_REASONS   = http.StatusUnavailableForLegalReasons
)

// 5xx: Server Errors
const (
	INTERNAL_SERVER_ERR             = http.StatusInternalServerError
	NOT_IMPLEMENTED                 = http.StatusNotImplemented
	BAD_GATEWAY                     = http.StatusBadGateway
	SERVICE_UNAVAILABLE             = http.StatusServiceUnavailable
	GATEWAY_TIMEOUT                 = http.StatusGatewayTimeout
	HTTP_VERSION_NOT_SUPPORTED      = http.StatusHTTPVersionNotSupported
	VARIANT_ALSO_NEGOTIATES         = http.StatusVariantAlsoNegotiates
	INSUFFICIENT_STORAGE            = http.StatusInsufficientStorage
	LOOP_DETECTED                   = http.StatusLoopDetected
	NOT_EXTENDED                    = http.StatusNotExtended
	NETWORK_AUTHENTICATION_REQUIRED = http.StatusNetworkAuthenticationRequired
)
