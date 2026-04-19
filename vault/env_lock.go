package vault

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

// LockEntry records a lock on a secret key.
type LockEntry struct {
	Key       string    `json:"key"`
	LockedBy  string    `json:"locked_by"`
	LockedAt  time.Time `json:"locked_at"`
	Reason    string    `json:"reason,omitempty"`
}

// LockFile maps key -> LockEntry.
type LockFile map[string]LockEntry

// LoadLockFile reads a lock file from disk. Returns empty map if not found.
func LoadLockFile(path string) (LockFile, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return LockFile{}, nil
	}
	if err != nil {
		return nil, err
	}
	var lf LockFile
	if err := json.Unmarshal(data, &lf); err != nil {
		return nil, err
	}
	return lf, nil
}

// SaveLockFile writes the lock file to disk.
func SaveLockFile(path string, lf LockFile) error {
	data, err := json.MarshalIndent(lf, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// LockKey adds a lock entry for the given key.
func LockKey(lf LockFile, key, lockedBy, reason string) LockFile {
	lf[key] = LockEntry{
		Key:      key,
		LockedBy: lockedBy,
		LockedAt: time.Now().UTC(),
		Reason:   reason,
	}
	return lf
}

// UnlockKey removes a lock entry for the given key.
func UnlockKey(lf LockFile, key string) (LockFile, bool) {
	_, exists := lf[key]
	delete(lf, key)
	return lf, exists
}

// IsLocked reports whether the given key is locked.
func IsLocked(lf LockFile, key string) bool {
	_, ok := lf[key]
	return ok
}

// FilterLocked returns only secrets whose keys are NOT locked.
func FilterLocked(secrets map[string]string, lf LockFile) map[string]string {
	out := make(map[string]string)
	for k, v := range secrets {
		if !IsLocked(lf, k) {
			out[k] = v
		}
	}
	return out
}
