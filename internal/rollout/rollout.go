package rollout

// fnvOffset64 and fnvPrime64 are the FNV-1a 64-bit constants.
const (
	fnvOffset64 = 14695981039346656037
	fnvPrime64  = 1099511628211
)

// Evaluate deterministically decides whether a feature is enabled for a user.
// It computes the FNV-1a 64-bit hash of key+":"+userID and maps it onto the
// range [0, 100); the result is true when that value is below percent. The
// function is seed-free and therefore stable: the same key and user always
// produce the same result.
func Evaluate(key, userID string, percent int) bool {
	h := uint64(fnvOffset64)
	for _, b := range []byte(key + ":" + userID) {
		h ^= uint64(b)
		h *= fnvPrime64
	}
	return int(h%100) < percent
}
