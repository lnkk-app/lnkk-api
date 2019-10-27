package metrics

import (
	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"

	"github.com/lnkk-ai/lnkk/pkg/errorreporting"
	"github.com/lnkk-ai/lnkk/pkg/store"
	"github.com/majordomusio/commons/pkg/util"
)

const (
	datastoreMetrics string = "METRICS"
)

type (

	// Metric is a generic data structure to store metrics
	Metric struct {
		Topic   string // describes the metric
		Label   string // additional context, e.g. an id, name/value pairs, comma separated labels etc
		Type    string // the type, e.g. count,
		Created int64
	}

	// CounterDS is a metric to collect integer values
	CounterDS struct {
		Metric
		Value int64
	}
)

// Count collects a numeric counter value
func Count(ctx context.Context, topic, label string, value int) {

	m := CounterDS{}
	m.Topic = topic
	m.Label = label
	m.Type = "count"
	m.Created = util.Timestamp()
	m.Value = int64(value)

	key := datastore.IncompleteKey(datastoreMetrics, nil)
	_, err := store.Client().Put(ctx, key, &m)

	if err != nil {
		errorreporting.Report(err)
	}
}
