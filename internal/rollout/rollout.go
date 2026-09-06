package rollout

// Evaluate deterministically decides whether a feature is enabled for a user.
// Scaffold stub: always returns false until the evaluate ticket implements the
// FNV-1a rollout hash.
func Evaluate(key, userID string, percent int) bool {
	return false
}
