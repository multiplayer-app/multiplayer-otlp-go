package multiplayer

const (
	MULTIPLAYER_TRACE_DOC_PREFIX                 = "d0cd0c"
	MULTIPLAYER_TRACE_DEBUG_PREFIX               = "debdeb"
	MULTIPLAYER_OTEL_DEFAULT_TRACES_EXPORTER_URL = "https://api.multiplayer.app/v1/traces"
	MULTIPLAYER_OTEL_DEFAULT_LOGS_EXPORTER_URL   = "https://api.multiplayer.app/v1/logs"
	MULTIPLAYER_ATTRIBUTE_PREFIX                 = "multiplayer."
	ATTR_MULTIPLAYER_DEBUG_SESSION               = "multiplayer.debug_session._id"
	ATTR_MULTIPLAYER_HTTP_REQUEST_BODY           = "multiplayer.http.request.body"
	ATTR_MULTIPLAYER_HTTP_RESPONSE_BODY          = "multiplayer.http.response.body"
	ATTR_MULTIPLAYER_HTTP_REQUEST_HEADERS        = "multiplayer.http.request.headers"
	ATTR_MULTIPLAYER_HTTP_RESPONSE_HEADERS       = "multiplayer.http.response.headers"
	ATTR_MULTIPLAYER_HTTP_RESPONSE_BODY_ENCODING = "multiplayer.http.response.body.encoding"
	MASK_PLACEHOLDER                             = "***MASKED***"
	MAX_MASK_DEPTH                               = 8
	MULTIPLAYER_MAX_HTTP_REQUEST_RESPONSE_SIZE   = 50000
)
