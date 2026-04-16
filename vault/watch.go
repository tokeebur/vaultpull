package vault

import (
	"context"
	"log"
	"time"
)

// WatchOptions configures the secret watcher.
type WatchOptions struct {
	Client     *Client
	SecretPath string
	OutputFile string
	Interval   time.Duration
	OnChange   func(diff DiffResult)
}

// DiffResult holds the result of a secrets comparison.
type DiffResult struct {
	Added   map[string]string
	Removed map[string]string
	Changed map[string]string
}

// WatchSecrets polls Vault at the given interval and triggers OnChange when secrets differ.
func WatchSecrets(ctx context.Context, opts WatchOptions) error {
	if opts.Interval <= 0 {
		opts.Interval = 30 * time.Second
	}

	ticker := time.NewTicker(opts.Interval)
	defer ticker.Stop()

	log.Printf("[watch] starting watch on %s every %s", opts.SecretPath, opts.Interval)

	for {
		select {
		case <-ctx.Done():
			log.Println("[watch] stopped")
			return ctx.Err()
		case <-ticker.C:
			remote, err := opts.Client.FetchSecrets(opts.SecretPath)
			if err != nil {
				log.Printf("[watch] fetch error: %v", err)
				continue
			}

			local, err := ParseEnvFile(opts.OutputFile)
			if err != nil {
				local = map[string]string{}
			}

			added, removed, changed := ComputeDiff(local, remote)
			if len(added)+len(removed)+len(changed) == 0 {
				continue
			}

			log.Printf("[watch] change detected: +%d -%d ~%d", len(added), len(removed), len(changed))

			if err := opts.Client.WriteEnvFile(opts.OutputFile, remote); err != nil {
				log.Printf("[watch] write error: %v", err)
				continue
			}

			if opts.OnChange != nil {
				opts.OnChange(DiffResult{Added: added, Removed: removed, Changed: changed})
			}
		}
	}
}
