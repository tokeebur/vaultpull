# vaultpull

> CLI tool to sync secrets from HashiCorp Vault into local `.env` files safely

---

## Installation

```bash
go install github.com/yourusername/vaultpull@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/vaultpull/releases).

---

## Usage

Set your Vault address and token, then run `vaultpull` pointing at a secret path:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.xxxxxxxx"

vaultpull --path secret/data/myapp --output .env
```

This will fetch all key/value pairs from the specified Vault path and write them to `.env`:

```
DB_HOST=db.example.com
DB_PASSWORD=supersecret
API_KEY=abc123
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--path` | Vault secret path | *(required)* |
| `--output` | Output file path | `.env` |
| `--merge` | Merge with existing file instead of overwriting | `false` |
| `--dry-run` | Print secrets to stdout without writing | `false` |

### Example with merge

```bash
vaultpull --path secret/data/myapp --output .env --merge
```

---

## Requirements

- Go 1.21+
- A running HashiCorp Vault instance
- A valid `VAULT_TOKEN` or other supported auth method

---

## License

[MIT](LICENSE) © 2024 yourusername