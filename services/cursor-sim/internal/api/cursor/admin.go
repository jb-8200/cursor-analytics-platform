package cursor

import (
	"net/http"

	"github.com/cursor-analytics-platform/services/cursor-sim/internal/api"
	"github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
)

// TeamMembersResponse matches the Cursor API response format for /teams/members.
type TeamMembersResponse struct {
	TeamMembers []TeamMember `json:"teamMembers"`
}

// TeamMember represents a team member in the Cursor API format.
type TeamMember struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// TeamMembers returns an HTTP handler for GET /teams/members.
// It retrieves all team members from the store and returns them in Cursor API format.
func TeamMembers(store storage.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get all developers from storage
		developers := store.ListDevelopers()

		// Convert to TeamMember format
		members := make([]TeamMember, 0, len(developers))
		for _, dev := range developers {
			members = append(members, TeamMember{
				Name:  dev.Name,
				Email: dev.Email,
				Role:  "member", // All developers are "member" role
			})
		}

		// Build response
		response := TeamMembersResponse{
			TeamMembers: members,
		}

		// Send JSON response
		api.RespondJSON(w, http.StatusOK, response)
	})
}
