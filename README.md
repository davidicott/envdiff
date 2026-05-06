# envdiff

> Compare `.env` files across environments and flag missing or mismatched keys.

---

## Installation

```bash
go install github.com/yourusername/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envdiff.git
cd envdiff
go build -o envdiff .
```

---

## Usage

```bash
envdiff [flags] <base-file> <compare-file> [compare-file...]
```

### Example

```bash
envdiff .env.example .env.production .env.staging
```

**Output:**

```
.env.production
  ✗ MISSING    DATABASE_URL
  ✗ MISSING    REDIS_HOST
  ~ MISMATCH   LOG_LEVEL (expected: "debug", got: "info")

.env.staging
  ✗ MISSING    REDIS_HOST
```

### Flags

| Flag | Description |
|------|-------------|
| `--keys-only` | Compare keys only, ignore values |
| `--strict` | Exit with non-zero code if any diff is found |
| `--json` | Output results as JSON |
| `--quiet` | Suppress output, use exit code only |
| `--ignore KEY` | Ignore a specific key (repeatable) |

---

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | No differences found |
| `1` | Differences found (when using `--strict`) |
| `2` | Error (e.g. file not found, parse error) |

---

## Why envdiff?

Keeping `.env` files in sync across environments is error-prone. `envdiff` makes it easy to catch missing or inconsistent configuration before it causes issues in production.

---

## License

MIT © [yourusername](https://github.com/yourusername)
