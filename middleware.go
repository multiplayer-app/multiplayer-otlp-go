package pkg

import (
	"bytes"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"io"
	"net/http"
)

func WithRequestData(h http.Handler, options MiddlewareOptions) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span := trace.SpanFromContext(r.Context())
		traceId := span.SpanContext().TraceID().String()
		if !IsMultiplayerTrace(traceId) {
			h.ServeHTTP(w, r)
			return
		}

		span.SetAttributes(
			attribute.String(ATTR_MULTIPLAYER_HTTP_REQUEST_HEADERS, maskHeaders(r.Header, options.headersToMask)),
		)

		if r.Header.Get("Content-Type") == "application/json" && r.Body != nil {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				h.ServeHTTP(w, r)
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			if len(bodyBytes) != 0 {
				if IsDebugTrace(traceId) && options.maskDebSpanPayload {
					bodyBytes = MaskJSONIfValid(bodyBytes, MAX_MASK_DEPTH)
				} else if options.schemifyDocSpanPayload {
					bodyBytes = GenerateJSONSchema(bodyBytes)
				}
			}
			if bodyBytes != nil {
				mpResponseBody := truncateIfNeeded(string(bodyBytes), options.maxPayloadSize)
				span.SetAttributes(attribute.String(ATTR_MULTIPLAYER_HTTP_REQUEST_BODY, mpResponseBody))
			}
		}
		h.ServeHTTP(w, r)
	})
}

func WithResponseData(next http.Handler, options MiddlewareOptions) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		span := trace.SpanFromContext(r.Context())
		traceId := span.SpanContext().TraceID().String()
		if !IsMultiplayerTrace(traceId) {
			next.ServeHTTP(w, r)
			return
		}

		isDebugTrace := IsDebugTrace(traceId)
		if isDebugTrace {
			w.Header().Set("X-Trace-Id", traceId)
		}

		rww := NewResponseWriterWrapper(w)
		defer func() {
			header := attribute.String(ATTR_MULTIPLAYER_HTTP_RESPONSE_HEADERS, maskHeaders(w.Header(), options.headersToMask))
			span.SetAttributes(header)
			bodyBytes := rww.GetBody()
			if len(bodyBytes) != 0 {
				if isDebugTrace && options.maskDebSpanPayload {
					bodyBytes = MaskJSONIfValid(bodyBytes, MAX_MASK_DEPTH)
				} else if options.schemifyDocSpanPayload {
					bodyBytes = GenerateJSONSchema(bodyBytes)
				}
			}
			if bodyBytes != nil {
				mpResponseBody := truncateIfNeeded(string(bodyBytes), options.maxPayloadSize)
				span.SetAttributes(attribute.String(ATTR_MULTIPLAYER_HTTP_RESPONSE_BODY, mpResponseBody))
			}
		}()
		next.ServeHTTP(rww, r)
	})
}
