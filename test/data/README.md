# Database Seeding System
The database will evolve over time, and its schema must be backward compatible with the previous versions.

Seed data is used for these purposes:
1. To populate the database with data for testing and development purposes
2. To provide a way to migrate data from one version of the database to another
3. To verify backward compatibility of the database schema between releases

## YAML-Based Seeding

The YAML-based seeding system allows you to define your seed data in a YAML file, making it easier to maintain and update.

### Usage

```bash
# Run with default configuration file (test/data/seed_data.yaml)
go run cmd/seed/main.go

# Run with a custom configuration file
go run cmd/seed/main.go -config path/to/your/seed_config.yaml
```

### YAML Configuration Format

The YAML configuration file has the following structure:

```yaml
description: "Description of the seed data"
db: "path/to/database.db"  # Relative or absolute path to the database
collections:
  - name: "collection_name"
    select: "SQL query to check if an item exists"
    insert: "SQL query to insert an item"
    items:
      - id: "item_id_1"
        field1: "value1"
        field2: "value2"
      - id: "item_id_2"
        field1: "value3"
        field2: "value4"
```

### Special Tokens

The YAML-based seeding system supports several special tokens:

1. `__::currentTimestamp::__` - Replaced with the current UTC timestamp (format: `2006-01-02 15:04:05.000Z`)
2. `__ref::collection::id` - References another item's ID. For example, `__ref::users::user1` references the ID of the item with ID "user1" in the "users" collection.

### Example

```yaml
description: "Test Seed Data"
db: "pb_data/data.db"
collections:
  - name: "users"
    select: "SELECT COUNT(*) FROM users WHERE id = ?"
    insert: "INSERT INTO users (id, email, name, created) VALUES (?, ?, ?, ?)"
    items:
      - id: "user1"
        email: "user1@example.com"
        name: "User One"
        created: __::currentTimestamp::__

  - name: "posts"
    select: "SELECT COUNT(*) FROM posts WHERE id = ?"
    insert: "INSERT INTO posts (id, title, author, created) VALUES (?, ?, ?, ?)"
    items:
      - id: "post1"
        title: "Post Title"
        author: "__ref::users::user1"
        created: __::currentTimestamp::__
```

### Order of Insertion

The collections are processed in the order they appear in the YAML file. This is important for handling references between collections. Make sure to list collections that are referenced by other collections first.

### Reference Validation

The system validates references during the configuration loading process:

1. It checks for correctly formatted references (`__ref::collection::id`)
2. It verifies that referenced collections exist
3. For collections that appear earlier in the YAML file, it verifies that referenced items exist
4. For forward references (to collections defined later in the file), full validation happens during the seeding process

Additionally, the system detects circular references between collections (e.g., collection A referencing collection B, which references back to collection A). Circular references are not allowed as they would create logical contradictions in the dependency graph and could lead to issues during seeding.

### Database Schema Validation

Before executing any insert operations, the system validates that the column names specified in your YAML configuration match the actual database schema:

1. It extracts table and column names from your INSERT statements
2. Queries the database for the actual column names in each table
3. Verifies that each column in your INSERT statement exists in the database

This validation catches common errors like:
- Typos in column names
- Using outdated column names that have been renamed
- Including columns that don't exist in the table
- Using wrong case for column names in case-sensitive databases

If any column is missing, the system will provide a clear error message identifying the table and column causing the issue, allowing you to fix the problem before any data is inserted.

### Error Handling

The seeding process is transactional. If any error occurs during the seeding process, the entire transaction is rolled back, ensuring the database remains in a consistent state.

The system provides detailed error messages for common issues:

1. **Schema Validation Errors**: When columns in your INSERT statements don't exist in the database tables
   ```
   Failed to seed database: schema validation failed: table practice_sessions has no column named instructor
   ```

2. **Column Count Mismatches**: When the number of values doesn't match the number of placeholders in your SQL
   ```
   Failed to seed database: column count mismatch for table practice_sessions: statement expects 11 columns but 12 values provided
   SQL: INSERT INTO practice_sessions (...) VALUES (?, ?, ...)
   Values count: 12
   ```

3. **Reference Errors**: When references between collections are invalid
   ```
   Failed to seed database: failed to resolve references: referenced id not found: __ref::users::non_existent
   ```

These detailed error messages help you quickly identify and fix issues in your seed configuration.

## Legacy Hardcoded Seeding

For backward compatibility, the package also provides a hardcoded seeding approach. You can use the `RunSeed()` function from `seed.go` to seed the database with the hardcoded data.

```go
import (
    seedtool "github.com/busybytelab.com/glimmer/test/data"
)

func main() {
    seedtool.RunSeed()
} 