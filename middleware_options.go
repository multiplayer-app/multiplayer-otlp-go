package multiplayer

import "strings"

type MiddlewareOptions struct {
	headersToMask          []string
	maxPayloadSize         int
	schemifyDocSpanPayload bool
	maskDebSpanPayload     bool
}

type Option func(*MiddlewareOptions)

func NewMiddlewareOptions(options ...Option) MiddlewareOptions {
	middleware := &MiddlewareOptions{
		headersToMask:          []string{},
		maxPayloadSize:         MULTIPLAYER_MAX_HTTP_REQUEST_RESPONSE_SIZE, // Default timeout
		schemifyDocSpanPayload: true,                                       // Default user agent
		maskDebSpanPayload:     true,
	}

	for _, opt := range options {
		opt(middleware)
	}

	return *middleware
}

func WithHeadersToMask(headersToMask []string) Option {
	return func(c *MiddlewareOptions) {
		normalizedHeaders := make([]string, len(headersToMask))
		for key, value := range headersToMask {
			normalizedHeaders[key] = strings.ToLower(value)
		}
		c.headersToMask = normalizedHeaders
	}
}

func WithMaxPayloadSize(maxPayloadSize int) Option {
	return func(c *MiddlewareOptions) {
		if maxPayloadSize < MULTIPLAYER_MAX_HTTP_REQUEST_RESPONSE_SIZE {
			c.maxPayloadSize = maxPayloadSize
		}
	}
}

func WithSchemifyDocSpanPayload(schemifyDocSpanPayload bool) Option {
	return func(c *MiddlewareOptions) {
		c.schemifyDocSpanPayload = schemifyDocSpanPayload
	}
}

func WithMaskDebSpanPayload(maskDebSpanPayload bool) Option {
	return func(c *MiddlewareOptions) {
		c.maskDebSpanPayload = maskDebSpanPayload
	}
}
