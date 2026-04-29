// Package observability — indexer-specific custom metrics (M4 #94).
//
// Auto-instrumentation already covers HTTP / DB / runtime metrics.
// This file adds the indexer-specific instruments that the OTel
// auto-instrumentation can't know about because they describe
// indexer-internal state:
//
//   - indexer_events_processed_total   — cumulative count of decoded
//                                        + published chain events.
//                                        Labelled by event_name and
//                                        chain_id. Cardinality is
//                                        bounded by the contract ABI
//                                        (≈ 30 event names) × the
//                                        configured chain (1).
//   - indexer_last_block_indexed       — gauge of the highest confirmed
//                                        block number the subscriber
//                                        has acknowledged. Operators
//                                        diff this against the chain
//                                        head to alert on stalls.
//   - indexer_reorgs_total             — cumulative reorg detections.
//                                        Sudden upticks indicate chain
//                                        instability or an indexer
//                                        misconfiguration.
//   - indexer_publish_errors_total     — cumulative NATS publish
//                                        failures. Labelled by
//                                        error_class to distinguish
//                                        wire timeouts from JetStream
//                                        rejection.
//
// Cardinality contract (enforced by Test_NoPIILabels):
//
//   - chain_id          — single int per process
//   - event_name        — bounded by ABI
//   - error_class       — bounded by hand: "timeout", "stream_full",
//                         "unmapped", "unknown"
//
// We do NOT label by block hash, transaction hash, or contract
// address. Those are unbounded; a Prometheus instance ingesting them
// would explode in TSDB head-block memory.
package observability

import (
	"context"
	"sync"

	"go.opentelemetry.io/otel/metric"
)

// indexerMeterName lands on the otel_scope_name label. Stable per
// service so Grafana dashboards can isolate indexer metrics from any
// gateway metrics that share an instrument name.
const indexerMeterName = "github.com/0xDevNinja/titular/services/indexer-go"

// Bounded label values for error_class on indexer_publish_errors_total.
// Anything that doesn't match one of these falls into "unknown" so the
// label cardinality stays fixed regardless of the underlying error
// surface.
const (
	ErrClassTimeout    = "timeout"
	ErrClassStreamFull = "stream_full"
	ErrClassUnmapped   = "unmapped"
	ErrClassUnknown    = "unknown"
)

var (
	indexerInstrumentsOnce sync.Once
	indexerInstrumentsErr  error

	indexerEventsCounter        metric.Int64Counter
	indexerLastBlockGauge       metric.Int64Gauge
	indexerReorgsCounter        metric.Int64Counter
	indexerPublishErrorsCounter metric.Int64Counter
)

// initIndexerInstruments builds the custom instruments lazily on first
// access. Errors land in indexerInstrumentsErr; the public Get…
// helpers return no-op fallbacks rather than nil so call sites stay
// nil-check-free.
func initIndexerInstruments() {
	indexerInstrumentsOnce.Do(func() {
		m := Meter(indexerMeterName)

		var err error
		indexerEventsCounter, err = m.Int64Counter(
			"indexer_events_processed_total",
			metric.WithDescription("Cumulative count of chain events decoded and published by the indexer."),
			metric.WithUnit("{event}"),
		)
		if err != nil {
			indexerInstrumentsErr = err
			return
		}

		indexerLastBlockGauge, err = m.Int64Gauge(
			"indexer_last_block_indexed",
			metric.WithDescription("Highest confirmed block number the subscriber has acknowledged."),
			metric.WithUnit("{block}"),
		)
		if err != nil {
			indexerInstrumentsErr = err
			return
		}

		indexerReorgsCounter, err = m.Int64Counter(
			"indexer_reorgs_total",
			metric.WithDescription("Cumulative reorg detections by the subscriber."),
			metric.WithUnit("{reorg}"),
		)
		if err != nil {
			indexerInstrumentsErr = err
			return
		}

		indexerPublishErrorsCounter, err = m.Int64Counter(
			"indexer_publish_errors_total",
			metric.WithDescription("Cumulative NATS publish failures, labelled by error_class."),
			metric.WithUnit("{error}"),
		)
		if err != nil {
			indexerInstrumentsErr = err
			return
		}
	})
}

// EventsProcessed returns the counter for decoded-and-published chain
// events. Call sites in internal/publisher and internal/decoder
// increment by one per successfully published event. Labels are
// `event_name` (bounded by ABI) and `chain_id` (bounded to 1 per
// process).
func EventsProcessed() metric.Int64Counter {
	initIndexerInstruments()
	if indexerEventsCounter == nil {
		return noopInt64Counter{}
	}
	return indexerEventsCounter
}

// LastBlockIndexed returns the gauge that tracks the indexer head.
// Call sites in internal/subscriber Record on every confirmed-block
// advance.
func LastBlockIndexed() metric.Int64Gauge {
	initIndexerInstruments()
	if indexerLastBlockGauge == nil {
		return noopInt64Gauge{}
	}
	return indexerLastBlockGauge
}

// Reorgs returns the cumulative reorg counter. Operators alert on
// rate(indexer_reorgs_total[5m]) > 0 to surface chain instability or
// a misconfigured indexer.
func Reorgs() metric.Int64Counter {
	initIndexerInstruments()
	if indexerReorgsCounter == nil {
		return noopInt64Counter{}
	}
	return indexerReorgsCounter
}

// PublishErrors returns the counter for NATS publish failures.
// Labelled by error_class (use the ErrClass… constants in this file).
func PublishErrors() metric.Int64Counter {
	initIndexerInstruments()
	if indexerPublishErrorsCounter == nil {
		return noopInt64Counter{}
	}
	return indexerPublishErrorsCounter
}

// noopInt64Counter is the safe fallback we hand back when instrument
// initialisation fails. Same justification as the gateway twin.
type noopInt64Counter struct{ metric.Int64Counter }

func (noopInt64Counter) Add(context.Context, int64, ...metric.AddOption) {}

// noopInt64Gauge is the gauge fallback.
type noopInt64Gauge struct{ metric.Int64Gauge }

func (noopInt64Gauge) Record(context.Context, int64, ...metric.RecordOption) {}
