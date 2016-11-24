package tracing

import (
	"errors"
)

const (
	COLLECTOR_KAFKA = "kafka"
	COLLECTOR_HTTP  = "http"
	COLLECTOR_TCP   = "tcp"

	TRACE_IMPL_APPDASH = "appdash"
	TRACE_IMPL_ZIPKIN  = "zipkin"

	DEFAULT_COLL_INTERVAL   = 30
	DEFAULT_COLL_MIN_TRACES = 100
)

var (
	ErrLostAppName      = errors.New("Lost App Name")
	ErrUnknownTraceImpl = errors.New("Unknown trace implement")
	ErrUnknownCollector = errors.New("Unknown Collector")
	ErrLostAddr         = errors.New("Lost Addr in collector")
	ErrLostTopic        = errors.New("Lost Topic in collector")
	ErrUnsupportColl    = errors.New("Unsupport collecotr type on this tracer's implement")
)

type CollectorHttp struct {
	Addr string `json:"addr"`
	User string `json:"user,omitempty"`
	Pwd  string `json:"pwd,omitempty"`
}

type CollectorKafka struct {
	Addr  string `json:"addr"`
	User  string `json:"user,omitempty"`
	Pwd   string `json:"pwd,omitempty"`
	Topic string `json:"topic,omitempty"`
}

type CollectorTcp struct {
	Addr string `json:"addr"`
}

type CollectorInfo struct {
	CommitInterval  int `json:"commit_interval,omitempty"`
	CommitMinTraces int `json:"commit_min_traces,omitempty"`
}

type TracerConfig struct {
	AppName        string `json:"app_name"`
	TraceImplement string `json:"trace_implement"`

	CollectorType string         `json:"collector_type"`
	CollTcp       CollectorTcp   `json:"collector_tcp,,omitempty"`
	CollHttp      CollectorHttp  `json:"collector_http,omitempty"`
	CollKafka     CollectorKafka `json:"collector_kafka,omitempty"`
	CollInfo      CollectorInfo  `json:"collector_info,omitempty"`
}

func (t *TracerConfig) checkAppName() error {
	if t.AppName == "" {
		return ErrLostAppName
	}
	return nil
}

func (t *TracerConfig) checkImpl() error {
	switch t.TraceImplement {
	case TRACE_IMPL_APPDASH:
	case TRACE_IMPL_ZIPKIN:
	default:
		return ErrUnknownTraceImpl
	}
	return nil
}

func (t *TracerConfig) checkCollector() error {
	switch t.CollectorType {
	case COLLECTOR_TCP:
		if t.CollTcp.Addr == "" {
			return ErrLostAddr
		}
	case COLLECTOR_KAFKA:
		if t.CollKafka.Addr == "" {
			return ErrLostAddr
		}
		if t.CollKafka.Topic == "" {
			return ErrLostTopic
		}
	case COLLECTOR_HTTP:
		if t.CollHttp.Addr == "" {
			return ErrLostAddr
		}
	default:
		return ErrUnknownCollector
	}
	return nil
}

func (t *TracerConfig) checkCollInfo() {
	if t.CollInfo.CommitInterval <= 0 {
		t.CollInfo.CommitInterval = DEFAULT_COLL_INTERVAL
	}
	if t.CollInfo.CommitMinTraces <= 0 {
		t.CollInfo.CommitMinTraces = DEFAULT_COLL_MIN_TRACES
	}
}

func (t *TracerConfig) Check() error {
	if err := t.checkAppName(); err != nil {
		return err
	}
	if err := t.checkImpl(); err != nil {
		return err
	}
	if err := t.checkCollector(); err != nil {
		return err
	}
	t.checkCollInfo()
	return nil
}

var GTraceConfig *TracerConfig = &TracerConfig{}
