package tracing

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"

	"sourcegraph.com/sourcegraph/appdash"
	appdashpot "sourcegraph.com/sourcegraph/appdash/opentracing"
)

// for zipkin implement
const (
	// Host + port of our service.
	hostPort = "127.0.0.1:0"
	// Debug mode.
	debug = false
	// same span can be set to true for RPC style spans (Zipkin V1) vs Node style (OpenTracing)
	sameSpan = true
	// make Tracer generate 128 bit traceID's for root spans.
	traceID128Bit = true
)

func initTracerImplAppdash(config *TracerConfig) error {
	switch config.CollectorType {
	case COLLECTOR_TCP:
		coll := appdash.NewChunkedCollector(appdash.NewRemoteCollector(config.CollTcp.Addr))
		coll.MinInterval = time.Duration(uint64(config.CollInfo.CommitInterval)) * time.Second

		tracer := appdashpot.NewTracer(coll)
		opentracing.InitGlobalTracer(tracer)
	default:
		return ErrUnsupportColl
	}
	return nil
}

func initTracerImplZipkin(config *TracerConfig) error {
	var coll zipkin.Collector
	var err error

	switch config.CollectorType {
	case COLLECTOR_KAFKA:
		coll, err = zipkin.NewKafkaCollector(strings.Split(config.CollKafka.Addr, ","))
		if err != nil {
			return err
		}
	case COLLECTOR_HTTP:
		coll, err = zipkin.NewHTTPCollector(config.CollHttp.Addr)
		if err != nil {
			return err
		}
	default:
		return ErrUnsupportColl
	}

	recorder := zipkin.NewRecorder(coll, debug, hostPort, config.AppName)
	tracer, err := zipkin.NewTracer(
		recorder,
		zipkin.ClientServerSameSpan(sameSpan),
		zipkin.TraceID128Bit(traceID128Bit),
	)
	if err != nil {
		return err
	}
	opentracing.InitGlobalTracer(tracer)
	return nil
}

func initGolbalTracer(config *TracerConfig) error {
	switch config.TraceImplement {
	case TRACE_IMPL_APPDASH:
		return initTracerImplAppdash(config)
	case TRACE_IMPL_ZIPKIN:
		return initTracerImplZipkin(config)
	default:
		return ErrUnknownTraceImpl
	}
}

func TracerInit(configFile string) error {
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, GTraceConfig)
	if err != nil {
		return err
	}
	err = GTraceConfig.Check()
	if err != nil {
		return err
	}
	return initGolbalTracer(GTraceConfig)
}

func GenOpName(opName string) string {
	return GTraceConfig.AppName + "|" + opName
}

func InjectSpanToHttpHeader(sp opentracing.Span, head http.Header) error {
	return sp.Tracer().Inject(sp.Context(), opentracing.TextMap, head)
}

func StartSpanFromHttpHeader(opName string, head http.Header) opentracing.Span {
	wireContext, err := opentracing.GlobalTracer().Extract(
		opentracing.TextMap,
		opentracing.HTTPHeadersCarrier(head))
	if err != nil {
		return opentracing.StartSpan(opName)
	} else {
		return opentracing.StartSpan(opName, opentracing.ChildOf(wireContext))
	}
}

func StartChildSpan(opName string, sp opentracing.Span) opentracing.Span {
	return opentracing.StartSpan(opName, opentracing.ChildOf(sp.Context()))
}
