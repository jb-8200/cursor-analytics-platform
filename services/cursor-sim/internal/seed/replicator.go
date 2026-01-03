package seed

import (
	"fmt"
	"math/rand"
)

// ReplicateDevelopers creates N developers by replicating/sampling from seed data.
// Strategy:
//   - If N <= len(seed.Developers): Sample N developers randomly
//   - If N > len(seed.Developers): Cycle through seed developers, modifying IDs
func ReplicateDevelopers(seed *SeedData, targetCount int, rng *rand.Rand) ([]Developer, error) {
	if targetCount < 1 {
		return nil, fmt.Errorf("target count must be >= 1, got %d", targetCount)
	}

	if len(seed.Developers) == 0 {
		return nil, fmt.Errorf("seed data has no developers")
	}

	result := make([]Developer, 0, targetCount)

	// Case 1: Downsample (N <= seed count)
	if targetCount <= len(seed.Developers) {
		// Shuffle and take first N
		indices := rng.Perm(len(seed.Developers))
		for i := 0; i < targetCount; i++ {
			result = append(result, seed.Developers[indices[i]])
		}
		return result, nil
	}

	// Case 2: Replicate (N > seed count)
	for i := 0; i < targetCount; i++ {
		// Cycle through seed developers
		srcDev := seed.Developers[i%len(seed.Developers)]

		// Clone and modify ID/email to make unique
		clonedDev := srcDev
		cloneNum := i / len(seed.Developers)

		if cloneNum == 0 {
			// First iteration: use original
			result = append(result, clonedDev)
		} else {
			// Subsequent iterations: clone with modified IDs
			clonedDev.UserID = fmt.Sprintf("%s_clone%d", srcDev.UserID, cloneNum)
			clonedDev.Email = fmt.Sprintf("clone%d_%s", cloneNum, srcDev.Email)
			clonedDev.Name = fmt.Sprintf("%s (Clone %d)", srcDev.Name, cloneNum)
			result = append(result, clonedDev)
		}
	}

	return result, nil
}
