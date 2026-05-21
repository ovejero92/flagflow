package api

import (
	"hash/fnv"
)

// evaluateRollout returns true if the user should see the enabled flag
// based on a consistent hash of flagName + userID.
func evaluateRollout(flagName, userID string, rolloutPercentage int) bool {
	if rolloutPercentage >= 100 {
		return true
	}
	if rolloutPercentage <= 0 {
		return false
	}
	if userID == "" {
		userID = "anonymous"
	}
	h := fnv.New32a()
	_, _ = h.Write([]byte(flagName + ":" + userID))
	bucket := int(h.Sum32() % 100)
	return bucket < rolloutPercentage
}
