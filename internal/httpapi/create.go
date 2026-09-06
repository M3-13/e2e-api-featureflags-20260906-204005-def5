package httpapi

import (
	"net/http"

	"featureflags/internal/store"
)

// createRequest is the decoded body of POST /flags. Enabled is a pointer so a
// missing field is distinguishable from an explicit false, and
// RolloutPercent is a pointer so an omitted value (default 100) is
// distinguishable from an explicit 0.
type createRequest struct {
	Key            string `json:"key"`
	Description    string `json:"description"`
	Enabled        *bool  `json:"enabled"`
	RolloutPercent *int   `json:"rollout_percent"`
}

// CreateHandler handles POST /flags: it validates the body, upserts the flag
// and answers 201 with the stored flag, or 400 on invalid input.
func CreateHandler(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createRequest
		if err := DecodeJSON(w, r, &req); err != nil {
			return
		}
		if req.Key == "" {
			WriteError(w, http.StatusBadRequest, "key is required")
			return
		}
		if req.Enabled == nil {
			WriteError(w, http.StatusBadRequest, "enabled is required")
			return
		}
		rollout := 100
		if req.RolloutPercent != nil {
			rollout = *req.RolloutPercent
		}
		if rollout < 0 || rollout > 100 {
			WriteError(w, http.StatusBadRequest, "rollout_percent must be between 0 and 100")
			return
		}
		flag := store.Flag{
			Key:            req.Key,
			Description:    req.Description,
			Enabled:        *req.Enabled,
			RolloutPercent: rollout,
		}
		WriteJSON(w, http.StatusCreated, s.Create(flag))
	}
}
