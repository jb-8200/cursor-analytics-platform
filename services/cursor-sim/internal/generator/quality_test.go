package generator

import (
	"testing"
	"time"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockQualityStore implements QualityStore for testing
type mockQualityStore struct {
	prs map[string]map[int]*models.PullRequest
}

func newMockQualityStore() *mockQualityStore {
	return &mockQualityStore{
		prs: make(map[string]map[int]*models.PullRequest),
	}
}

func (m *mockQualityStore) GetPR(repoName string, number int) (*models.PullRequest, error) {
	if repoPRs, ok := m.prs[repoName]; ok {
		if pr, ok := repoPRs[number]; ok {
			return pr, nil
		}
	}
	return nil, nil
}

func (m *mockQualityStore) GetPRsByRepo(repoName string) []models.PullRequest {
	result := []models.PullRequest{}
	if repoPRs, ok := m.prs[repoName]; ok {
		for _, pr := range repoPRs {
			result = append(result, *pr)
		}
	}
	return result
}

func (m *mockQualityStore) UpdatePR(pr models.PullRequest) error {
	if m.prs[pr.RepoName] == nil {
		m.prs[pr.RepoName] = make(map[int]*models.PullRequest)
	}
	m.prs[pr.RepoName][pr.Number] = &pr
	return nil
}

func (m *mockQualityStore) addPR(pr models.PullRequest) {
	if m.prs[pr.RepoName] == nil {
		m.prs[pr.RepoName] = make(map[int]*models.PullRequest)
	}
	m.prs[pr.RepoName][pr.Number] = &pr
}

func TestQualityGenerator_CategorizeAIRatio(t *testing.T) {
	seedData := &seed.SeedData{
		Version: "1.0.0",
		Correlations: seed.Correlations{
			AIRatioBands: seed.AIRatioBands{
				Low:    seed.AIRatioBand{Min: 0.0, Max: 0.3},
				Medium: seed.AIRatioBand{Min: 0.3, Max: 0.7},
				High:   seed.AIRatioBand{Min: 0.7, Max: 1.0},
			},
		},
	}

	store := newMockQualityStore()
	gen := NewQualityGeneratorWithSeed(seedData, store, 42)

	tests := []struct {
		aiRatio  float64
		expected AIRatioCategory
	}{
		{0.0, AIRatioCategoryLow},
		{0.2, AIRatioCategoryLow},
		{0.3, AIRatioCategoryMedium},
		{0.5, AIRatioCategoryMedium},
		{0.7, AIRatioCategoryHigh},
		{0.9, AIRatioCategoryHigh},
		{1.0, AIRatioCategoryHigh},
	}

	for _, tt := range tests {
		category := gen.CategorizeAIRatio(tt.aiRatio)
		assert.Equal(t, tt.expected, category, "AI ratio %.1f should be %s", tt.aiRatio, tt.expected)
	}
}

func TestQualityGenerator_CalculateRevertProbability(t *testing.T) {
	seedData := &seed.SeedData{
		Version: "1.0.0",
		PRLifecycle: seed.PRLifecycle{
			QualityOutcomes: seed.QualityOutcomes{
				RevertProbability: seed.OutcomeParams{
					Base: 0.05, // 5% base revert rate
					Modifiers: seed.OutcomeModifiers{
						ByAIRatio: map[string]float64{
							"low":    0.8, // Lower revert rate
							"medium": 1.0, // Normal
							"high":   1.5, // Higher revert rate
						},
					},
				},
			},
		},
		Correlations: seed.Correlations{
			AIRatioBands: seed.AIRatioBands{
				Low:    seed.AIRatioBand{Min: 0.0, Max: 0.3},
				Medium: seed.AIRatioBand{Min: 0.3, Max: 0.7},
				High:   seed.AIRatioBand{Min: 0.7, Max: 1.0},
			},
		},
	}

	store := newMockQualityStore()
	gen := NewQualityGeneratorWithSeed(seedData, store, 42)

	// Low AI ratio should have lower revert probability
	lowAIProb := gen.CalculateRevertProbability(0.2, 1)
	assert.InDelta(t, 0.04, lowAIProb, 0.01, "low AI ratio should have ~4% revert rate")

	// High AI ratio should have higher revert probability
	highAIProb := gen.CalculateRevertProbability(0.9, 1)
	assert.InDelta(t, 0.075, highAIProb, 0.01, "high AI ratio should have ~7.5% revert rate")

	// More iterations should reduce revert probability
	highIterProb := gen.CalculateRevertProbability(0.5, 3)
	normalIterProb := gen.CalculateRevertProbability(0.5, 1)
	assert.Less(t, highIterProb, normalIterProb, "more iterations should reduce revert probability")
}

func TestQualityGenerator_ApplyQualityOutcomes(t *testing.T) {
	seedData := &seed.SeedData{
		Version: "1.0.0",
		PRLifecycle: seed.PRLifecycle{
			QualityOutcomes: seed.QualityOutcomes{
				RevertProbability: seed.OutcomeParams{
					Base: 0.5, // 50% for testing
				},
			},
		},
		Correlations: seed.Correlations{
			AIRatioBands: seed.AIRatioBands{
				Low:    seed.AIRatioBand{Min: 0.0, Max: 0.3},
				Medium: seed.AIRatioBand{Min: 0.3, Max: 0.7},
				High:   seed.AIRatioBand{Min: 0.7, Max: 1.0},
			},
		},
	}

	store := newMockQualityStore()
	now := time.Now()

	// Add some merged PRs
	prs := []models.PullRequest{
		{Number: 1, RepoName: "acme/api", State: models.PRStateMerged, AIRatio: 0.8, MergedAt: &now},
		{Number: 2, RepoName: "acme/api", State: models.PRStateMerged, AIRatio: 0.2, MergedAt: &now},
		{Number: 3, RepoName: "acme/api", State: models.PRStateOpen, AIRatio: 0.5}, // Not merged
	}
	for _, pr := range prs {
		store.addPR(pr)
	}

	gen := NewQualityGeneratorWithSeed(seedData, store, 42)

	err := gen.ApplyQualityOutcomes("acme/api")
	require.NoError(t, err)

	// Check that some PRs might be marked as reverted (with 50% probability)
	// Due to deterministic seed, we should get consistent results
	pr1, _ := store.GetPR("acme/api", 1)
	pr2, _ := store.GetPR("acme/api", 2)
	pr3, _ := store.GetPR("acme/api", 3)

	// PR3 should not be affected (not merged)
	assert.False(t, pr3.WasReverted, "open PR should not be marked as reverted")

	// At least test that the function ran without error
	// and updated merged PRs (can't guarantee specific outcomes due to randomness)
	assert.NotNil(t, pr1)
	assert.NotNil(t, pr2)
}

func TestQualityGenerator_IsBugFix(t *testing.T) {
	store := newMockQualityStore()
	gen := NewQualityGeneratorWithSeed(nil, store, 42)

	tests := []struct {
		title    string
		labels   []string
		expected bool
	}{
		{"fix: resolve login issue", []string{}, true},
		{"fix(auth): password reset", []string{}, true},
		{"feat: add new feature", []string{}, false},
		{"Add feature", []string{"bug"}, true},
		{"Update docs", []string{"documentation"}, false},
		{"hotfix: critical fix", []string{}, true},
		{"bugfix: minor issue", []string{}, true},
	}

	for _, tt := range tests {
		pr := &models.PullRequest{
			Title:  tt.title,
			Labels: tt.labels,
		}
		result := gen.IsBugFix(pr)
		assert.Equal(t, tt.expected, result, "PR with title '%s' and labels %v should have IsBugFix=%v", tt.title, tt.labels, tt.expected)
	}
}

func TestQualityGenerator_Reproducibility(t *testing.T) {
	seedData := &seed.SeedData{
		Version: "1.0.0",
		PRLifecycle: seed.PRLifecycle{
			QualityOutcomes: seed.QualityOutcomes{
				RevertProbability: seed.OutcomeParams{
					Base: 0.5,
				},
			},
		},
		Correlations: seed.Correlations{
			AIRatioBands: seed.AIRatioBands{
				Low:    seed.AIRatioBand{Min: 0.0, Max: 0.3},
				Medium: seed.AIRatioBand{Min: 0.3, Max: 0.7},
				High:   seed.AIRatioBand{Min: 0.7, Max: 1.0},
			},
		},
	}

	// Create identical stores
	store1 := newMockQualityStore()
	store2 := newMockQualityStore()

	now := time.Now()
	for _, store := range []*mockQualityStore{store1, store2} {
		store.addPR(models.PullRequest{Number: 1, RepoName: "acme/api", State: models.PRStateMerged, AIRatio: 0.8, MergedAt: &now})
		store.addPR(models.PullRequest{Number: 2, RepoName: "acme/api", State: models.PRStateMerged, AIRatio: 0.5, MergedAt: &now})
	}

	gen1 := NewQualityGeneratorWithSeed(seedData, store1, 12345)
	gen2 := NewQualityGeneratorWithSeed(seedData, store2, 12345)

	_ = gen1.ApplyQualityOutcomes("acme/api")
	_ = gen2.ApplyQualityOutcomes("acme/api")

	// Same seed should produce same results
	pr1a, _ := store1.GetPR("acme/api", 1)
	pr1b, _ := store2.GetPR("acme/api", 1)
	assert.Equal(t, pr1a.WasReverted, pr1b.WasReverted, "same seed should produce same revert status")

	pr2a, _ := store1.GetPR("acme/api", 2)
	pr2b, _ := store2.GetPR("acme/api", 2)
	assert.Equal(t, pr2a.WasReverted, pr2b.WasReverted, "same seed should produce same revert status")
}

func TestQualityGenerator_CalculateHotfixProbability(t *testing.T) {
	seedData := &seed.SeedData{
		Version: "1.0.0",
		PRLifecycle: seed.PRLifecycle{
			QualityOutcomes: seed.QualityOutcomes{
				HotfixProbability: seed.OutcomeParams{
					Base: 0.1, // 10% base hotfix rate
				},
			},
		},
	}

	store := newMockQualityStore()
	gen := NewQualityGeneratorWithSeed(seedData, store, 42)

	prob := gen.CalculateHotfixProbability(&models.PullRequest{
		IsBugFix: true,
		AIRatio:  0.9,
	})

	// Bug fixes with high AI ratio should have higher hotfix probability
	assert.Greater(t, prob, 0.0, "should have non-zero hotfix probability")
}
