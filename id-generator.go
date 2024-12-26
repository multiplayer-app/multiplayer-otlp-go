package pkg

import (
	"context"
	crand "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	otelTrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"math/rand"
	"sync"
)

type ratioDependentIdGenerator struct {
	debugSessionShortID string
	traceIDUpperBound   uint64
	sync.Mutex
	randSource *rand.Rand
}

var _ otelTrace.IDGenerator = &ratioDependentIdGenerator{}

func NewRatioDependentIdGenerator(autoDocTracesRatio float64) otelTrace.IDGenerator {
	var rngSeed int64
	_ = binary.Read(crand.Reader, binary.LittleEndian, &rngSeed)

	return &ratioDependentIdGenerator{
		traceIDUpperBound: uint64(autoDocTracesRatio * (1 << 63)),
		randSource:        rand.New(rand.NewSource(rngSeed)),
	}
}

func (gen *ratioDependentIdGenerator) SetDebugSessionShortID(debugSessionShortID string) {
	gen.debugSessionShortID = debugSessionShortID
}

func (gen *ratioDependentIdGenerator) NewIDs(ctx context.Context) (trace.TraceID, trace.SpanID) {
	gen.Lock()
	defer gen.Unlock()
	traceId := trace.TraceID{}
	spanId := trace.SpanID{}
	for {
		_, _ = gen.randSource.Read(traceId[:])
		if traceId.IsValid() {
			break
		}
	}
	for {
		_, _ = gen.randSource.Read(spanId[:])
		if spanId.IsValid() {
			break
		}
	}

	var prefix []byte
	if gen.debugSessionShortID != "" {
		prefix, _ = hex.DecodeString(MULTIPLAYER_TRACE_DEBUG_PREFIX)
	} else {
		x := binary.BigEndian.Uint64(traceId[8:16]) >> 1
		if x < gen.traceIDUpperBound {
			prefix, _ = hex.DecodeString(MULTIPLAYER_TRACE_DOC_PREFIX)
		}
	}
	for i := 0; i < len(prefix); i += 1 {
		traceId[i] = prefix[i]
	}
	return traceId, spanId
}

func (gen *ratioDependentIdGenerator) NewSpanID(ctx context.Context, traceID trace.TraceID) trace.SpanID {
	gen.Lock()
	defer gen.Unlock()
	sid := trace.SpanID{}
	for {
		gen.randSource.Read(sid[:])
		if sid.IsValid() {
			break
		}
	}
	return sid
}
