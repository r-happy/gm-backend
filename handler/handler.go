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
