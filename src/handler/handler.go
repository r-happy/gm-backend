package handler

import (
	"back/model"

	"github.com/google/uuid"
)

// generateUniqueID generates a unique ID for a space.
func generateUniqueID() (string, error) {
	for {
		tmpUUID := uuid.NewString()
		s := model.FindSpace(&model.Space{ID: tmpUUID})
		if s.ID == "" {
			return tmpUUID, nil
		}
	}
}

// IsUserMemberOfSpace checks if a user is a member of a specific space.
func IsUserMemberOfSpace(email string, spaceID string) bool {
    members := model.FindMembers(&model.Member{Email: email, Space: spaceID})
    return len(members) > 0
}

