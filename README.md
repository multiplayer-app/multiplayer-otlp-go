# multiplayer-otlp-go
## Introduction
The `multiplayer-otlp-go` module integrates OpenTelemetry with the Multiplayer platform to enable seamless trace collection and analysis. This exporter helps developers monitor, debug, and document application performance with detailed trace data. It supports flexible trace ID generation, sampling strategies, and provides middleware to capture API request/response data.

## Installation

To install the `multiplayer-otlp-go` module, use the following command:

```bash
go get github.com/multiplayer-app/multiplayer-otlp-go
```

### Multiplayer Exporter
To get started, initialize the Multiplayer Exporter in your OpenTelemetry setup. Below is a detailed example:
1. Set up the environment variable `MULTIPLAYER_OTLP_KEY`. 
2. Add multiplayer.Exporter to your trace provider
```go
multiplayerOtlpKey := os.Getenv("MULTIPLAYER_OTLP_KEY")

traceExporter := multiplayer.NewExporter(multiplayerOtlpKey)
traceProvider := trace.NewTracerProvider(
    trace.WithIDGenerator(multiplayer.NewRatioDependentIdGenerator(1)),
    trace.WithSampler(multiplayer.NewSampler(trace.AlwaysSample())),
    trace.WithBatcher(traceExporter, trace.WithBatchTimeout(time.Second)),
)
```

### ID Generator

The Multiplayer listens for specific trace IDs. You can configure the ratio of traces that will be used for documentation. A ratio of `1` corresponds to the behavior of `trace.AlwaysSample()`, and a ratio of `0` corresponds to `trace.NeverSample()`.

```go
trace.WithIDGenerator(multiplayer.NewRatioDependentIdGenerator(0.5))
```

### Sampler

The Multiplayer sampler samples traces with specific IDs and other traces sampled by the provided sampler.

```go
multiplayer.NewSampler(trace.AlwaysSample())
```

### Middlewares

This module provides two middlewares to integrate your REST API request/response data into Multiplayer Debugger/Documentation.

```go
options := multiplayer.NewMiddlewareOptions(
    multiplayer.WithHeadersToMask([]string{"AUTH_HEADER_NAME"})
)
multiplayer.WithRequestData(handler, options)
multiplayer.WithResponseData(handler, options)
```

These middlewares add the request and response data to the trace span attributes. By default, the following headers are masked: `set-cookie`, `cookie`, `authorization`, and `proxy-authorization`. To customize the list of masked headers, use the `multiplayer.WithHeadersToMask` option.

For debug traces, the body values are masked by default. This behavior can be changed using the `multiplayer.WithMaskDebSpanPayload(false)` option. For documentation traces, the body is converted into a corresponding JSON schema, which can be disabled using `multiplayer.WithSchemifyDocSpanPayload(false)`.

---
## Example
For a complete OpenTelemetry configuration example, refer to [dice_example](example/dice).

## License
This module is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

## Support
If you encounter any issues or have questions, please open an issue in the GitHub repository or contact the maintainers.
