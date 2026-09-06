package httpapi

import (
	"net/http"

	"featureflags/internal/store"
)

// updateRequest is the JSON body accepted by PUT /flags/{key}. enabled is
// required; description and rollout_percent are optional and, when omitted,
// preserve the flag's current values.
type updateRequest struct {
	Enabled        *bool   `json:"enabled"`
	Description    *string `json:"description"`
	RolloutPercent *int    `json:"rollout_percent"`
}

// UpdateHandler handles PUT /flags/{key}. It decodes the body, validates it,
// and updates the flag in the store, returning 200 with the updated flag, 400
// on invalid input, or 404 when the key is unknown.
func UpdateHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")

		var req updateRequest
		if err := DecodeJSON(w, r, &req); err != nil {
			return
		}
		if req.Enabled == nil {
			WriteError(w, http.StatusBadRequest, "enabled is required")
			return
		}
		if req.RolloutPercent != nil && (*req.RolloutPercent < 0 || *req.RolloutPercent > 100) {
			WriteError(w, http.StatusBadRequest, "rollout_percent must be between 0 and 100")
			return
		}

		flag := store.Flag{
			Key:     key,
			Enabled: *req.Enabled,
		}
		if existing, ok := s.Get(key); ok {
			flag.Description = existing.Description
			flag.RolloutPercent = existing.RolloutPercent
		}
		if req.Description != nil {
			flag.Description = *req.Description
		}
		if req.RolloutPercent != nil {
			flag.RolloutPercent = *req.RolloutPercent
		}

		updated, ok := s.Update(key, flag)
		if !ok {
			WriteError(w, http.StatusNotFound, "flag not found")
			return
		}
		WriteJSON(w, http.StatusOK, updated)
	}
}
