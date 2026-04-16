package config
t"errors"
	"os"
)

// Config holds the runtime configuration for vntype Config struct {
	VaultAddr  string
	VaultToken string
	SecretPath string
	OutputFile string
}

// Load reads configuration from environment variables and applies overrides.
func Load(secretPath, outputFile string) (*Config, error) {
	addr := os.Getenv("VAULT_ADDR")
	if addr == "" {
		addr = "http://127.0.0.1:8200"
	}

	token := os.Getenv("VAULT_TOKEN")
	if token == "" {
		return nil, errors.New("VAULT_TOKEN environment variable is not set")
	}

	if secretPath == "" {
		return nil, errors.New("secret path must not be empty")
	}

	if outputFile == "" {
		outputFile = ".env"
	}

	return &Config{
		VaultAddr:  addr,
		VaultToken: token,
		SecretPath: secretPath,
		OutputFile: outputFile,
	}, nil
}
