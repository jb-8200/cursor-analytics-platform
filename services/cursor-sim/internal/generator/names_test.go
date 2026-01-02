package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewNameGenerator(t *testing.T) {
	gen := NewNameGenerator(12345)
	require.NotNil(t, gen)
	assert.NotNil(t, gen.rng)
}

func TestNameGenerator_GenerateName(t *testing.T) {
	gen := NewNameGenerator(12345)

	firstName, lastName := gen.GenerateName()
	assert.NotEmpty(t, firstName, "First name should not be empty")
	assert.NotEmpty(t, lastName, "Last name should not be empty")
}

func TestNameGenerator_GenerateName_Deterministic(t *testing.T) {
	// Same seed should produce same sequence
	gen1 := NewNameGenerator(99999)
	gen2 := NewNameGenerator(99999)

	for i := 0; i < 10; i++ {
		first1, last1 := gen1.GenerateName()
		first2, last2 := gen2.GenerateName()

		assert.Equal(t, first1, first2, "First names should match at iteration %d", i)
		assert.Equal(t, last1, last2, "Last names should match at iteration %d", i)
	}
}

func TestNameGenerator_GenerateName_Variety(t *testing.T) {
	gen := NewNameGenerator(12345)

	names := make(map[string]bool)
	for i := 0; i < 50; i++ {
		first, last := gen.GenerateName()
		fullName := first + " " + last
		names[fullName] = true
	}

	// Should have high variety (at least 40 unique names out of 50)
	assert.GreaterOrEqual(t, len(names), 40,
		"Should generate variety of names, got %d unique names", len(names))
}

func TestNameGenerator_GenerateEmail(t *testing.T) {
	gen := NewNameGenerator(12345)

	email := gen.GenerateEmail("John", "Doe")
	assert.Equal(t, "john.doe@company.com", email)
}

func TestNameGenerator_GenerateEmail_SpecialCharacters(t *testing.T) {
	gen := NewNameGenerator(12345)

	tests := []struct {
		firstName string
		lastName  string
		want      string
	}{
		{"John", "Doe", "john.doe@company.com"},
		{"Mary-Jane", "O'Connor", "mary-jane.o'connor@company.com"},
		{"José", "García", "josé.garcía@company.com"},
	}

	for _, tt := range tests {
		t.Run(tt.firstName+" "+tt.lastName, func(t *testing.T) {
			email := gen.GenerateEmail(tt.firstName, tt.lastName)
			assert.Equal(t, tt.want, email)
		})
	}
}
