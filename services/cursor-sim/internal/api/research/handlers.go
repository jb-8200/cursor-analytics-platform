package research

import (
	"bytes"
	"net/http"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/export"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/generator"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/services"
)

// DatasetHandler returns an HTTP handler for GET /research/dataset.
// Supports format query parameter: json (default), csv.
func DatasetHandler(gen *generator.ResearchGenerator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse date range
		from, to, err := parseDateRange(r)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Generate dataset
		dataPoints, err := gen.GenerateDataset(from, to)
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, "failed to generate dataset: "+err.Error())
			return
		}

		// Get format from query
		format := r.URL.Query().Get("format")
		if format == "" {
			format = "json"
		}

		switch format {
		case "csv":
			w.Header().Set("Content-Type", "text/csv")
			w.Header().Set("Content-Disposition", "attachment; filename=research_dataset.csv")
			exporter := export.NewCSVExporter(w)
			if err := exporter.ExportDataPoints(dataPoints); err != nil {
				api.RespondError(w, http.StatusInternalServerError, "failed to export CSV: "+err.Error())
				return
			}
		case "json":
			fallthrough
		default:
			response := map[string]interface{}{
				"data": dataPoints,
				"params": map[string]string{
					"from":   from.Format("2006-01-02"),
					"to":     to.Format("2006-01-02"),
					"format": format,
				},
			}
			api.RespondJSON(w, http.StatusOK, response)
		}
	})
}

// VelocityMetricsHandler returns an HTTP handler for GET /research/metrics/velocity.
func VelocityMetricsHandler(gen *generator.ResearchGenerator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse date range
		from, to, err := parseDateRange(r)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Generate dataset
		dataPoints, err := gen.GenerateDataset(from, to)
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, "failed to generate dataset: "+err.Error())
			return
		}

		// Calculate metrics
		svc := services.NewResearchMetricsService(dataPoints)
		period := from.Format("2006-01")
		metrics := svc.CalculateVelocityMetrics(period)

		response := map[string]interface{}{
			"data": metrics,
			"params": map[string]string{
				"from":   from.Format("2006-01-02"),
				"to":     to.Format("2006-01-02"),
				"period": period,
			},
		}
		api.RespondJSON(w, http.StatusOK, response)
	})
}

// ReviewCostMetricsHandler returns an HTTP handler for GET /research/metrics/review-costs.
func ReviewCostMetricsHandler(gen *generator.ResearchGenerator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse date range
		from, to, err := parseDateRange(r)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Generate dataset
		dataPoints, err := gen.GenerateDataset(from, to)
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, "failed to generate dataset: "+err.Error())
			return
		}

		// Calculate metrics
		svc := services.NewResearchMetricsService(dataPoints)
		period := from.Format("2006-01")
		metrics := svc.CalculateReviewCostMetrics(period)

		response := map[string]interface{}{
			"data": metrics,
			"params": map[string]string{
				"from":   from.Format("2006-01-02"),
				"to":     to.Format("2006-01-02"),
				"period": period,
			},
		}
		api.RespondJSON(w, http.StatusOK, response)
	})
}

// QualityMetricsHandler returns an HTTP handler for GET /research/metrics/quality.
func QualityMetricsHandler(gen *generator.ResearchGenerator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse date range
		from, to, err := parseDateRange(r)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Generate dataset
		dataPoints, err := gen.GenerateDataset(from, to)
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, "failed to generate dataset: "+err.Error())
			return
		}

		// Calculate metrics
		svc := services.NewResearchMetricsService(dataPoints)
		period := from.Format("2006-01")
		metrics := svc.CalculateQualityMetrics(period)

		response := map[string]interface{}{
			"data": metrics,
			"params": map[string]string{
				"from":   from.Format("2006-01-02"),
				"to":     to.Format("2006-01-02"),
				"period": period,
			},
		}
		api.RespondJSON(w, http.StatusOK, response)
	})
}

// DatasetNDJSONHandler returns an HTTP handler for GET /research/dataset.ndjson.
// Returns data in NDJSON (JSON Lines) format for large datasets.
func DatasetNDJSONHandler(gen *generator.ResearchGenerator) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse date range
		from, to, err := parseDateRange(r)
		if err != nil {
			api.RespondError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Generate dataset
		dataPoints, err := gen.GenerateDataset(from, to)
		if err != nil {
			api.RespondError(w, http.StatusInternalServerError, "failed to generate dataset: "+err.Error())
			return
		}

		// Stream as NDJSON
		w.Header().Set("Content-Type", "application/x-ndjson")
		w.Header().Set("Content-Disposition", "attachment; filename=research_dataset.ndjson")

		var buf bytes.Buffer
		exporter := export.NewJSONExporter(&buf)
		if err := exporter.ExportDataPointsStream(dataPoints); err != nil {
			api.RespondError(w, http.StatusInternalServerError, "failed to export NDJSON: "+err.Error())
			return
		}

		w.Write(buf.Bytes())
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
