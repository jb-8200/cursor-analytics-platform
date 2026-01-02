package replay

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
)

// CorpusLoader loads research data from corpus files.
type CorpusLoader struct{}

// NewCorpusLoader creates a new corpus loader.
func NewCorpusLoader() *CorpusLoader {
	return &CorpusLoader{}
}

// LoadJSON loads research data points from a JSON array file.
func (l *CorpusLoader) LoadJSON(path string) ([]models.ResearchDataPoint, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read corpus file: %w", err)
	}

	var dataPoints []models.ResearchDataPoint
	if err := json.Unmarshal(data, &dataPoints); err != nil {
		return nil, fmt.Errorf("failed to parse JSON corpus: %w", err)
	}

	return dataPoints, nil
}

// LoadNDJSON loads research data points from a NDJSON (JSON Lines) file.
func (l *CorpusLoader) LoadNDJSON(path string) ([]models.ResearchDataPoint, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open corpus file: %w", err)
	}
	defer file.Close()

	var dataPoints []models.ResearchDataPoint
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if line == "" {
			continue
		}

		var dp models.ResearchDataPoint
		if err := json.Unmarshal([]byte(line), &dp); err != nil {
			return nil, fmt.Errorf("failed to parse line %d: %w", lineNum, err)
		}
		dataPoints = append(dataPoints, dp)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read corpus file: %w", err)
	}

	return dataPoints, nil
}

// CorpusIndex provides efficient querying of corpus data.
type CorpusIndex struct {
	dataPoints []models.ResearchDataPoint
	timeIndex  []timeEntry
	bandIndex  map[models.AIRatioBand][]int // indices into dataPoints
}

type timeEntry struct {
	timestamp time.Time
	index     int
}

// NewCorpusIndex creates a new index from data points.
func NewCorpusIndex(dataPoints []models.ResearchDataPoint) *CorpusIndex {
	if dataPoints == nil {
		dataPoints = []models.ResearchDataPoint{}
	}

	idx := &CorpusIndex{
		dataPoints: dataPoints,
		timeIndex:  make([]timeEntry, len(dataPoints)),
		bandIndex:  make(map[models.AIRatioBand][]int),
	}

	// Build time index
	for i, dp := range dataPoints {
		idx.timeIndex[i] = timeEntry{timestamp: dp.Timestamp, index: i}
	}

	// Sort time index
	sort.Slice(idx.timeIndex, func(i, j int) bool {
		return idx.timeIndex[i].timestamp.Before(idx.timeIndex[j].timestamp)
	})

	// Build AI ratio band index
	for i, dp := range dataPoints {
		band := dp.GetAIRatioBand()
		idx.bandIndex[band] = append(idx.bandIndex[band], i)
	}

	return idx
}

// QueryByTimeRange returns data points within the time range.
func (idx *CorpusIndex) QueryByTimeRange(from, to time.Time) []models.ResearchDataPoint {
	if len(idx.timeIndex) == 0 {
		return nil
	}

	var result []models.ResearchDataPoint

	for _, entry := range idx.timeIndex {
		if (entry.timestamp.Equal(from) || entry.timestamp.After(from)) &&
			(entry.timestamp.Equal(to) || entry.timestamp.Before(to)) {
			result = append(result, idx.dataPoints[entry.index])
		}
	}

	return result
}

// QueryByAIRatioBand returns data points with the specified AI ratio band.
func (idx *CorpusIndex) QueryByAIRatioBand(band models.AIRatioBand) []models.ResearchDataPoint {
	indices, ok := idx.bandIndex[band]
	if !ok {
		return nil
	}

	result := make([]models.ResearchDataPoint, len(indices))
	for i, dpIdx := range indices {
		result[i] = idx.dataPoints[dpIdx]
	}

	return result
}

// GetAll returns all data points.
func (idx *CorpusIndex) GetAll() []models.ResearchDataPoint {
	return idx.dataPoints
}

// Count returns the number of data points in the index.
func (idx *CorpusIndex) Count() int {
	return len(idx.dataPoints)
}
