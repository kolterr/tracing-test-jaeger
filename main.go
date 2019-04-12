package main

import (
    "github.com/gin-contrib/opengintracing"
    "github.com/gin-gonic/gin"
    "github.com/opentracing/opentracing-go"
    "github.com/uber/jaeger-client-go"
    "github.com/uber/jaeger-client-go/zipkin"
)

func main() {
    e := gin.Default()
    propagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
    //sender := transport.NewHTTPTransport("http://localhost:14268/api/v2/spans", transport.HTTPBatchSize(1))
    sender, _ := jaeger.NewUDPTransport("127.0.0.1:6831", 0)
    tracer, closer := jaeger.NewTracer(
        "api_gateway",
        jaeger.NewConstSampler(true),
        jaeger.NewRemoteReporter(sender),
        jaeger.TracerOptions.Injector(opentracing.HTTPHeaders, propagator),
        jaeger.TracerOptions.Extractor(opentracing.HTTPHeaders, propagator),
        jaeger.TracerOptions.ZipkinSharedRPCSpan(true),
    )
    defer closer.Close()
    opentracing.SetGlobalTracer(tracer)
    e.GET("/echo", opengintracing.NewSpan(tracer, "forward to service 1"), handler)
    e.Run(":8080")
}

func handler(c *gin.Context) {
    span, _ := opengintracing.GetSpan(c)
    defer span.Finish()
    res, _ := service(span)
    c.String(200, res)
}
// service
func service(span opentracing.Span) (string, error) {
    sp := opentracing.StartSpan("work 1", opentracing.ChildOf(span.Context()))
    sp.SetTag("func", "do home work")
    defer sp.Finish()
    return service1(sp)
}

func service1(span opentracing.Span) (string, error) {
    sp := opentracing.StartSpan("work 2", opentracing.ChildOf(span.Context()))
    sp.SetTag("func1", "do home work2")
    defer sp.Finish()
    return "Hello wolrd!", nil
}
