package tortoise

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/spacemeshos/go-spacemesh/metrics"
)

const namespace = "tortoise"

var (
	ballotsNumber = metrics.NewGauge(
		"ballots",
		namespace,
		"Number of ballots in the state",
		[]string{},
	).WithLabelValues()
	delayedBallots = metrics.NewGauge(
		"delayed_ballots",
		namespace,
		"Number of ballots that are delayed due to the wrong beacon",
		[]string{},
	).WithLabelValues()
	blocksNumber = metrics.NewGauge(
		"blocks",
		namespace,
		"Number of blocks in the state",
		[]string{},
	).WithLabelValues()
	layersNumber = metrics.NewGauge(
		"layers",
		namespace,
		"Number of layers in the state",
		[]string{},
	).WithLabelValues()
	epochsNumber = metrics.NewGauge(
		"epochs",
		namespace,
		"Number of epochs in the state",
		[]string{},
	).WithLabelValues()
	atxsNumber = metrics.NewGauge(
		"atxs",
		namespace,
		"number of atxs in the tortoise state",
		[]string{},
	).WithLabelValues()
	malfeasantNumber = metrics.NewGauge(
		"malfeasant",
		namespace,
		"number of malfeasant identities",
		[]string{},
	).WithLabelValues()
)

var (
	stateLayers = metrics.NewGauge(
		"state_layers",
		namespace,
		"Layers in the state (evicted, verified, latest, processed)",
		[]string{"kind"},
	)
	evictedLayer   = stateLayers.WithLabelValues("evicted")
	verifiedLayer  = stateLayers.WithLabelValues("verified")
	processedLayer = stateLayers.WithLabelValues("processed")
	lastLayer      = stateLayers.WithLabelValues("last")
)

var modeGauge = metrics.NewGauge(
	"mode",
	namespace,
	"0 for verifying, 1 for full",
	[]string{},
).WithLabelValues()

var errorsCounter = metrics.NewCounter(
	"errors",
	namespace,
	"Counter for all errors",
	[]string{},
).WithLabelValues()

var (
	onBallotHist = metrics.NewHistogramWithBuckets(
		"tortoise_on_ballot_ns",
		namespace,
		"Time to process ballot in ns.",
		[]string{"step"},
		prometheus.ExponentialBuckets(100_000, 2, 10),
	)
	waitBallotDuration   = onBallotHist.WithLabelValues("wait")
	decodeBallotDuration = onBallotHist.WithLabelValues("decode")
	fcountBallotDuration = onBallotHist.WithLabelValues("full_count")
	vcountBallotDuration = onBallotHist.WithLabelValues("verifying_count")
)

var (
	onBlockHist = metrics.NewHistogramWithBuckets(
		"tortoise_on_block_ns",
		namespace,
		"Time to process block in ns.",
		[]string{"step"},
		prometheus.ExponentialBuckets(100_000, 2, 10),
	)
	waitBlockDuration = onBlockHist.WithLabelValues("wait")
	addBlockDuration  = onBlockHist.WithLabelValues("add")
	lateBlockDuration = onBlockHist.WithLabelValues("late")
)

var (
	onHareOutputHist = metrics.NewHistogramWithBuckets(
		"tortoise_on_hare_output_ns",
		namespace,
		"Time to process hare output in ns.",
		[]string{"step"},
		prometheus.ExponentialBuckets(100_000, 2, 10),
	)
	waitHareOutputDuration = onHareOutputHist.WithLabelValues("wait")
	addHareOutput          = onHareOutputHist.WithLabelValues("add")
)

var (
	onAtxHist = metrics.NewHistogramWithBuckets(
		"tortoise_on_atx_ns",
		namespace,
		"Time to add atx in ns.",
		[]string{"step"},
		prometheus.ExponentialBuckets(100_000, 2, 10),
	)
	waitAtxDuration = onAtxHist.WithLabelValues("wait")
	addAtxDuration  = onAtxHist.WithLabelValues("add")
)

var (
	tallyVotesHist = metrics.NewHistogramWithBuckets(
		"tortoise_tally_votes_ns",
		namespace,
		"Time for tally votes to complete in ns.",
		[]string{"step"},
		prometheus.ExponentialBuckets(100_000, 2, 10),
	)
	waitTallyVotes    = tallyVotesHist.WithLabelValues("wait")
	executeTallyVotes = tallyVotesHist.WithLabelValues("execute")
)

var (
	encodeVotesHist = metrics.NewHistogramWithBuckets(
		"tortoise_encode_votes_ns",
		namespace,
		"Time for encode votes to complete in ns.",
		[]string{"step"},
		prometheus.ExponentialBuckets(100_000, 2, 10),
	)
	waitEncodeVotes    = encodeVotesHist.WithLabelValues("wait")
	executeEncodeVotes = encodeVotesHist.WithLabelValues("execute")
)

func enableCollector(t *Tortoise) {
	prometheus.MustRegister(&collector{
		tortoise: t,
		opinions: prometheus.NewDesc(
			prometheus.BuildFQName(metrics.Namespace, namespace, "opinion"),
			"opinion on the layer",
			[]string{"layer", "opinion"}, nil),
		applied: prometheus.NewDesc(
			prometheus.BuildFQName(metrics.Namespace, namespace, "applied"),
			"applied block",
			[]string{"layer", "block"}, nil),
		blocksByLayer: prometheus.NewDesc(
			prometheus.BuildFQName(metrics.Namespace, namespace, "blocks_by_layer"),
			"count blocks by layer",
			[]string{"layer"}, nil),
	})
}

type collector struct {
	tortoise      *Tortoise
	opinions      *prometheus.Desc
	applied       *prometheus.Desc
	blocksByLayer *prometheus.Desc
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.opinions
	ch <- c.applied
	ch <- c.blocksByLayer
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	for _, rst := range c.tortoise.latestsResults(c.tortoise.cfg.CollectDetails) {
		ch <- prometheus.MustNewConstMetric(c.opinions, prometheus.GaugeValue, 1,
			rst.Layer.String(), rst.Opinion.ShortString())
		ch <- prometheus.MustNewConstMetric(c.applied, prometheus.GaugeValue, 1,
			rst.Layer.String(), rst.FirstValid().String())
		ch <- prometheus.MustNewConstMetric(c.blocksByLayer, prometheus.GaugeValue, float64(len(rst.Blocks)),
			rst.Layer.String())
	}
}
