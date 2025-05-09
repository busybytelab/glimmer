# Database Seeding System

A flexible YAML-based seeding system for populating your database with test data. Perfect for development, testing, and ensuring backward compatibility.

## Quick Start

```bash
# Build the application first
make build

# Use default sample data (examples for common platforms)
./build/glimmer-darwin-amd64 seed    # macOS Intel
./build/glimmer-darwin-arm64 seed    # macOS Apple Silicon
./build/glimmer-linux-amd64 seed     # Linux x86_64
./build/glimmer-linux-arm64 seed     # Linux ARM64

# Use custom YAML config
./build/glimmer-darwin-arm64 seed --config path/to/your/seed_config.yaml

# Generate password hash for seed files
./build/glimmer-darwin-arm64 seed password-hash your_password
```

## YAML Configuration

Define your seed data in YAML format:

```yaml
description: "Your seed data description"
db: "path/to/database.db"
default_password_hash: "your_default_password_hash"  # Optional
collections:
  - name: "users"
    select: "SELECT COUNT(*) FROM users WHERE id = {:id}"
    insert: "INSERT INTO users (id, email, name) VALUES ({:id}, {:email}, {:name})"
    items:
      - id: "user1"
        email: "user1@example.com"
        name: "User One"
```

## Special Features

- **References**: Link items between collections using `__ref::collection::id`
- **Timestamps**: Use `__::currentTimestamp::__` for current UTC time
- **Environment Variables**: Access env vars with `__::env::VAR_NAME::default_value::__`
  - For password hashes, use `__::env::TEST_USER_PASSWORD_HASH::default_password_hash::__` to use either the environment variable or fallback to default_password_hash
  - Note: When setting bcrypt hashes in environment variables, escape each `$` with another `$` (e.g., `$$2a$$10$$...` instead of `$2a$10$...`)
- **Transaction Safety**: All operations are wrapped in transactions
- **Duplicate Prevention**: Skips existing items automatically
- **Schema Validation**: Validates column names against database schema
- **Reference Validation**: Checks for valid references and circular dependencies

## Best Practices

1. Order collections by dependencies (referenced collections first)
2. Use meaningful IDs for easy reference
3. Keep your YAML file organized and well-documented
4. Use the password-hash command to generate secure password hashes

## Example

See `sample-v0.1.yaml` for a complete example with all features. 