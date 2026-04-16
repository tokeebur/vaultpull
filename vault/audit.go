package vault

import (
	"encoding/json"
	"os"
	"time"
)

// AuditEntry records a single sync operation.
type AuditEntry struct {
	Timestamp  time.Time `json:"timestamp"`
	SecretPath string    `json:"secret_path"`
	OutputFile string    `json:"output_file"`
	Added      int       `json:"added"`
	Modified   int       `json:"modified"`
	Removed    int       `json:"removed"`
	DryRun     bool      `json:"dry_run"`
}

// AppendAuditLog appends an AuditEntry as a JSON line to the given log file.
func AppendAuditLog(logPath string, entry AuditEntry) error {
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	entry.Timestamp = time.Now().UTC()
	line, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	_, err = f.Write(append(line, '\n'))
	return err
}

// ReadAuditLog reads all AuditEntry records from the given log file.
func ReadAuditLog(logPath string) ([]AuditEntry, error) {
	data, err := os.ReadFile(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var entries []AuditEntry
	dec := json.NewDecoder(
		// wrap bytes in a reader via strings trick
		newBytesReader(data),
	)
	for dec.More() {
		var e AuditEntry
		if err := dec.Decode(&e); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}
