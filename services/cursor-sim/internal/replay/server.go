package replay

import (
	"net/http"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/services"
)

// ReplayServer serves research data from a pre-loaded corpus.
type ReplayServer struct {
	index *CorpusIndex
}

// NewReplayServer creates a new replay server from a corpus index.
func NewReplayServer(index *CorpusIndex) *ReplayServer {
	return &ReplayServer{
		index: index,
	}
}

// DatasetHandler returns an HTTP handler for GET /research/dataset.
func (s *ReplayServer) DatasetHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		from, to, err := parseDateRange(r)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		dataPoints := s.index.QueryByTimeRange(from, to)

		response := map[string]interface{}{
			"data": dataPoints,
			"params": map[string]string{
				"from": from.Format("2006-01-02"),
				"to":   to.Format("2006-01-02"),
				"mode": "replay",
			},
		}
		api.RespondJSON(w, http.StatusOK, response)
	})
}

// VelocityMetricsHandler returns an HTTP handler for GET /research/metrics/velocity.
func (s *ReplayServer) VelocityMetricsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		from, to, err := parseDateRange(r)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		dataPoints := s.index.QueryByTimeRange(from, to)
		svc := services.NewResearchMetricsService(dataPoints)
		period := from.Format("2006-01")
		metrics := svc.CalculateVelocityMetrics(period)

		response := map[string]interface{}{
			"data": metrics,
			"params": map[string]string{
				"from":   from.Format("2006-01-02"),
				"to":     to.Format("2006-01-02"),
				"period": period,
				"mode":   "replay",
			},
		}
		api.RespondJSON(w, http.StatusOK, response)
	})
}

// ReviewCostMetricsHandler returns an HTTP handler for GET /research/metrics/review-costs.
func (s *ReplayServer) ReviewCostMetricsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		from, to, err := parseDateRange(r)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		dataPoints := s.index.QueryByTimeRange(from, to)
		svc := services.NewResearchMetricsService(dataPoints)
		period := from.Format("2006-01")
		metrics := svc.CalculateReviewCostMetrics(period)

		response := map[string]interface{}{
			"data": metrics,
			"params": map[string]string{
				"from":   from.Format("2006-01-02"),
				"to":     to.Format("2006-01-02"),
				"period": period,
				"mode":   "replay",
			},
		}
		api.RespondJSON(w, http.StatusOK, response)
	})
}

// QualityMetricsHandler returns an HTTP handler for GET /research/metrics/quality.
func (s *ReplayServer) QualityMetricsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		from, to, err := parseDateRange(r)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		dataPoints := s.index.QueryByTimeRange(from, to)
		svc := services.NewResearchMetricsService(dataPoints)
		period := from.Format("2006-01")
		metrics := svc.CalculateQualityMetrics(period)

		response := map[string]interface{}{
			"data": metrics,
			"params": map[string]string{
				"from":   from.Format("2006-01-02"),
				"to":     to.Format("2006-01-02"),
				"period": period,
				"mode":   "replay",
			},
		}
		api.RespondJSON(w, http.StatusOK, response)
	})
}

// parseDateRange extracts and validates from/to date parameters.
func parseDateRange(r *http.Request) (from, to time.Time, err error) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	// Default to last 30 days if not specified
	if fromStr == "" {
		from = time.Now().AddDate(0, 0, -30)
	} else {
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}

	if toStr == "" {
		to = time.Now()
	} else {
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}

	// Add time to include full day
	to = to.Add(24*time.Hour - time.Second)

	return from, to, nil
}

// Corpus File Format Documentation
//
// The replay mode supports two corpus file formats:
//
// 1. JSON Array Format (.json)
//    A single JSON array containing all ResearchDataPoint objects:
//    [
//      {"commit_hash": "abc123", "ai_ratio": 0.5, ...},
//      {"commit_hash": "def456", "ai_ratio": 0.8, ...}
//    ]
//
// 2. NDJSON Format (.ndjson)
//    One JSON object per line (newline-delimited JSON):
//    {"commit_hash": "abc123", "ai_ratio": 0.5, ...}
//    {"commit_hash": "def456", "ai_ratio": 0.8, ...}
//
// To generate a corpus file:
//    1. Run simulator in runtime mode: ./cursor-sim -mode runtime -days 90
//    2. Export dataset: curl http://localhost:8080/research/dataset?format=json > corpus.json
//
// To use replay mode:
//    ./cursor-sim -mode replay -corpus path/to/corpus.json
//
// This allows researchers to:
//    - Share reproducible datasets
//    - Test analysis pipelines without simulation
//    - Compare results across different analysis runs
