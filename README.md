# vaultpull

> A CLI tool to sync secrets from HashiCorp Vault into local `.env` files with namespace filtering and audit logging.

---

## Installation

```bash
go install github.com/yourusername/vaultpull@latest
```

Or download a pre-built binary from the [Releases](https://github.com/yourusername/vaultpull/releases) page.

---

## Usage

Set your Vault address and token, then run `vaultpull` with a path and output target:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.your-token-here"

# Pull secrets from a specific path into a .env file
vaultpull pull --path secret/myapp/prod --out .env

# Filter by namespace
vaultpull pull --path secret/myapp --namespace prod --out .env.prod

# Enable audit logging
vaultpull pull --path secret/myapp/prod --out .env --audit audit.log
```

### Flags

| Flag | Description |
|-------------|--------------------------------------|
| `--path` | Vault secret path to read from |
| `--out` | Output `.env` file path |
| `--namespace` | Filter secrets by namespace prefix |
| `--audit` | Path to write audit log entries |
| `--dry-run` | Preview output without writing to disk |

---

## Configuration

`vaultpull` respects standard Vault environment variables:

- `VAULT_ADDR` — Vault server address
- `VAULT_TOKEN` — Authentication token
- `VAULT_NAMESPACE` — (optional) Vault enterprise namespace

---

## License

[MIT](LICENSE) © 2024 yourusername