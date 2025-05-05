// Package seed provides utilities for seeding the database with test data.
//
// This file implements a YAML-based seeding system that allows you to define your seed data
// in a YAML configuration file. The system supports relationships between different collections
// through references and special tokens like timestamps.
//
// The YAML-based approach makes it easy to maintain and update your seed data as your
// application evolves, without having to modify Go code. Just update the YAML file with
// new collections, items, or relationships.
//
// Key features:
// - Define seed data in YAML format
// - Support for references between collections
// - Special tokens like timestamps
// - Transactional processing to ensure data consistency
// - Skip existing items to avoid duplicates
package seed

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v3"
)

// SeedConfig represents the top-level structure of the seed YAML file
type SeedConfig struct {
	Description string             `yaml:"description"`
	DB          string             `yaml:"db"`
	Collections []CollectionConfig `yaml:"collections"`
}

// CollectionConfig represents a collection's configuration for seeding
type CollectionConfig struct {
	Name   string        `yaml:"name"`
	Select string        `yaml:"select"`
	Insert string        `yaml:"insert"`
	Items  []interface{} `yaml:"items"`
}

// LoadSeedConfig loads the seed configuration from a YAML file
func LoadSeedConfig(configPath string) (*SeedConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Replace special token __::currentTimestamp::__ with the current time
	timestamp := time.Now().UTC().Format("2006-01-02 15:04:05.000Z")
	re := regexp.MustCompile(`__::currentTimestamp::__`)
	data = re.ReplaceAll(data, []byte(timestamp))

	var config SeedConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Validate the configuration
	if err := validateConfig(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// validateConfig validates the seed configuration
func validateConfig(config *SeedConfig) error {
	if config.DB == "" {
		return fmt.Errorf("database path is required")
	}

	if len(config.Collections) == 0 {
		return fmt.Errorf("at least one collection is required")
	}

	// Create a map of collection names and their items for reference validation
	collectionMap := make(map[string]bool)
	itemsRegistry := make(map[string]map[string]bool)

	// Build dependency graph for circular reference detection
	dependencyGraph := make(map[string][]string)

	for _, collection := range config.Collections {
		collectionMap[collection.Name] = true
		itemsRegistry[collection.Name] = make(map[string]bool)
		dependencyGraph[collection.Name] = []string{}

		for _, item := range collection.Items {
			itemData, ok := item.(map[string]interface{})
			if !ok {
				continue
			}

			if id, ok := itemData["id"].(string); ok {
				itemsRegistry[collection.Name][id] = true
			}

			// Build dependency graph by finding references
			for _, value := range itemData {
				strValue, ok := value.(string)
				if ok && strings.HasPrefix(strValue, "__ref::") {
					parts := strings.Split(strValue, "::")
					if len(parts) == 3 {
						refCollection := parts[1]
						// Add dependency
						found := false
						for _, dep := range dependencyGraph[collection.Name] {
							if dep == refCollection {
								found = true
								break
							}
						}
						if !found {
							dependencyGraph[collection.Name] = append(dependencyGraph[collection.Name], refCollection)
						}
					}
				}
			}
		}
	}

	// Check for circular references
	if err := detectCircularReferences(dependencyGraph); err != nil {
		return err
	}

	// Validate each collection
	for i, collection := range config.Collections {
		if collection.Name == "" {
			return fmt.Errorf("collection %d has no name", i)
		}

		if collection.Select == "" {
			return fmt.Errorf("collection %s has no select statement", collection.Name)
		}

		if collection.Insert == "" {
			return fmt.Errorf("collection %s has no insert statement", collection.Name)
		}

		if len(collection.Items) == 0 {
			return fmt.Errorf("collection %s has no items", collection.Name)
		}

		// Validate each item has an ID and check references
		for j, item := range collection.Items {
			itemMap, ok := item.(map[string]interface{})
			if !ok {
				return fmt.Errorf("item %d in collection %s is not a map", j, collection.Name)
			}

			if _, ok := itemMap["id"]; !ok {
				return fmt.Errorf("item %d in collection %s has no id field", j, collection.Name)
			}

			// Validate references
			for key, value := range itemMap {
				strValue, ok := value.(string)
				if ok && strings.HasPrefix(strValue, "__ref::") {
					// Parse reference: __ref::collection::id
					parts := strings.Split(strValue, "::")
					if len(parts) != 3 {
						return fmt.Errorf("invalid reference format in %s.%s: %s", collection.Name, key, strValue)
					}

					refCollection := parts[1]
					refID := parts[2]

					// Check if referenced collection exists
					if !collectionMap[refCollection] {
						return fmt.Errorf("reference to non-existent collection in %s.%s: %s", collection.Name, key, refCollection)
					}

					// For collections that appear before this one in the YAML, we can validate item existence
					if refColls, ok := itemsRegistry[refCollection]; ok {
						if !refColls[refID] {
							return fmt.Errorf("reference to non-existent item in %s.%s: %s::%s", collection.Name, key, refCollection, refID)
						}
					}
					// For collections that appear later, we can't validate yet (will be caught during ResolveReferences)
				}
			}
		}
	}

	return nil
}

// detectCircularReferences checks for circular dependencies in the collection references
func detectCircularReferences(graph map[string][]string) error {
	// Keep track of visited nodes and current path
	visited := make(map[string]bool)
	path := make(map[string]bool)

	// DFS to detect cycles
	var dfs func(node string) error
	dfs = func(node string) error {
		if !visited[node] {
			visited[node] = true
			path[node] = true

			for _, neighbor := range graph[node] {
				if !visited[neighbor] {
					if err := dfs(neighbor); err != nil {
						return err
					}
				} else if path[neighbor] {
					return fmt.Errorf("circular reference detected: %s -> %s", node, neighbor)
				}
			}
		}
		path[node] = false
		return nil
	}

	// Start DFS from each node
	for node := range graph {
		if !visited[node] {
			if err := dfs(node); err != nil {
				return err
			}
		}
	}

	return nil
}

// SeedDatabaseFromYAML seeds the database using the provided YAML configuration
func SeedDatabaseFromYAML(configPath string) error {
	config, err := LoadSeedConfig(configPath)
	if err != nil {
		return err
	}

	// Resolve references in the config
	if err := ResolveReferences(config); err != nil {
		return fmt.Errorf("failed to resolve references: %w", err)
	}

	// Determine database path
	dbPath := config.DB
	if !filepath.IsAbs(dbPath) {
		workDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		dbPath = filepath.Join(workDir, dbPath)
	}

	// Check if the database exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return fmt.Errorf("database file not found at: %s", dbPath)
	}

	// Connect to the database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	// Validate schema before processing
	if err := validateDatabaseSchema(db, config); err != nil {
		return fmt.Errorf("schema validation failed: %w", err)
	}

	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Process each collection in order
	for _, collection := range config.Collections {
		if err := seedCollection(tx, collection); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// validateDatabaseSchema checks if columns in INSERT statements exist in the database
func validateDatabaseSchema(db *sql.DB, config *SeedConfig) error {
	for _, collection := range config.Collections {
		// Extract table name and columns from INSERT statement
		tableName, columns, err := extractTableAndColumns(collection.Insert)
		if err != nil {
			return fmt.Errorf("failed to parse insert statement for %s: %w", collection.Name, err)
		}

		// Get actual table columns from database
		tableColumns, err := getTableColumns(db, tableName)
		if err != nil {
			return fmt.Errorf("failed to get columns for table %s: %w", tableName, err)
		}

		// Validate each column exists in the table
		for _, col := range columns {
			if !columnExists(tableColumns, col) {
				return fmt.Errorf("table %s has no column named %s", tableName, col)
			}
		}

		log.Printf("Validated schema for table %s", tableName)
	}
	return nil
}

// extractTableAndColumns extracts the table name and column names from an INSERT statement
func extractTableAndColumns(insertStmt string) (string, []string, error) {
	// Extract table name using regex - making it case-insensitive with (?i)
	tableRegex := regexp.MustCompile(`(?i)INSERT\s+INTO\s+(\w+)\s*\(`)
	tableMatches := tableRegex.FindStringSubmatch(insertStmt)
	if len(tableMatches) < 2 {
		return "", nil, fmt.Errorf("failed to extract table name from insert statement")
	}
	tableName := tableMatches[1]

	// Extract column names using regex - making it case-insensitive with (?i)
	colRegex := regexp.MustCompile(`(?i)INSERT\s+INTO\s+\w+\s*\((.*?)\)`)
	colMatches := colRegex.FindStringSubmatch(insertStmt)
	if len(colMatches) < 2 {
		return "", nil, fmt.Errorf("failed to extract column names from insert statement")
	}

	// Process column names
	colsStr := colMatches[1]
	colsParts := strings.Split(colsStr, ",")
	columns := make([]string, len(colsParts))
	for i, col := range colsParts {
		columns[i] = strings.TrimSpace(col)
	}

	return tableName, columns, nil
}

// getTableColumns retrieves the column names for a given table from the database
func getTableColumns(db *sql.DB, tableName string) ([]string, error) {
	// SQLite-specific query to get column information
	rows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var cid int
		var name string
		var type_ string
		var notnull int
		var dflt_value interface{}
		var pk int

		if err := rows.Scan(&cid, &name, &type_, &notnull, &dflt_value, &pk); err != nil {
			return nil, err
		}
		columns = append(columns, name)
	}

	if len(columns) == 0 {
		return nil, fmt.Errorf("table %s does not exist or has no columns", tableName)
	}

	return columns, nil
}

// columnExists checks if a column exists in the list of table columns
func columnExists(tableColumns []string, column string) bool {
	for _, col := range tableColumns {
		if col == column {
			return true
		}
	}
	return false
}

// seedCollection seeds a single collection based on its configuration
func seedCollection(tx *sql.Tx, collection CollectionConfig) error {
	log.Printf("Seeding collection: %s", collection.Name)

	for i, item := range collection.Items {
		// Convert item to map
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return fmt.Errorf("item %d in collection %s is not a map", i, collection.Name)
		}

		// Get the ID for checking existence
		id, ok := itemMap["id"]
		if !ok {
			return fmt.Errorf("item %d in collection %s has no id field", i, collection.Name)
		}

		// Check if item exists
		var count int
		err := tx.QueryRow(collection.Select, id).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to check if item exists in %s: %w", collection.Name, err)
		}

		if count == 0 {
			// Extract values from item map in the correct order for the INSERT statement
			values, err := extractInsertValues(collection.Insert, itemMap)
			if err != nil {
				return fmt.Errorf("failed to extract values for item %s in %s: %w", id, collection.Name, err)
			}

			// Insert item
			_, err = tx.Exec(collection.Insert, values...)
			if err != nil {
				// Check if it's a column count mismatch error
				if strings.Contains(err.Error(), "values for") && strings.Contains(err.Error(), "columns") {
					// Parse the expected number of columns from INSERT statement
					expectedCols, _ := countPlaceholdersInStatement(collection.Insert)
					actualCols := len(values)
					tableName, _, _ := extractTableAndColumns(collection.Insert)

					return fmt.Errorf("column count mismatch for table %s: statement expects %d columns but %d values provided\n"+
						"SQL: %s\n"+
						"Values count: %d\n"+
						"Please check your insert statement and ensure it matches the item fields",
						tableName, expectedCols, actualCols, collection.Insert, actualCols)
				}
				return fmt.Errorf("failed to insert item %s in %s: %w", id, collection.Name, err)
			}

			log.Printf("Inserted item %s in collection %s", id, collection.Name)
		} else {
			log.Printf("Skipped existing item %s in collection %s", id, collection.Name)
		}
	}

	return nil
}

// countPlaceholdersInStatement counts the number of placeholders (?) in an SQL statement
func countPlaceholdersInStatement(stmt string) (int, error) {
	// Count the number of question marks (placeholders)
	count := strings.Count(stmt, "?")

	// Try to verify by parsing the statement - count commas in values section + 1
	valuesRegex := regexp.MustCompile(`(?i)VALUES\s*\((.*?)\)`)
	matches := valuesRegex.FindStringSubmatch(stmt)

	if len(matches) >= 2 {
		valuesPart := matches[1]
		commaCount := strings.Count(valuesPart, ",")
		expectedPlaceholders := commaCount + 1

		if expectedPlaceholders != count {
			return count, fmt.Errorf("mismatch between placeholders count (%d) and values section format (%d)", count, expectedPlaceholders)
		}
	}

	return count, nil
}

// extractInsertValues extracts values from the item map in the order required by the insert statement
func extractInsertValues(insertStmt string, itemMap map[string]interface{}) ([]interface{}, error) {
	// Extract column names from insert statement - making it case-insensitive with (?i)
	re := regexp.MustCompile(`(?i)INSERT\s+INTO\s+\w+\s+\((.*?)\)`)
	matches := re.FindStringSubmatch(insertStmt)
	if len(matches) < 2 {
		return nil, fmt.Errorf("failed to extract column names from insert statement")
	}

	// Get column names
	columnsStr := matches[1]
	columns := strings.Split(columnsStr, ",")

	// Clean up column names - trim whitespace and convert to YAML field names if needed
	fieldNames := make([]string, len(columns))
	for i, col := range columns {
		col = strings.TrimSpace(col)
		// Convert snake_case to camelCase for YAML field mapping if needed
		fieldNames[i] = convertDBColumnToFieldName(col)
	}

	// Extract values in column order
	values := make([]interface{}, len(fieldNames))
	for i, field := range fieldNames {
		val, ok := itemMap[field]
		if !ok {
			return nil, fmt.Errorf("missing value for field %s", field)
		}
		values[i] = val
	}

	return values, nil
}

// convertDBColumnToFieldName converts a database column name to the appropriate field name
// This handles differences between the database column names and the YAML field names
func convertDBColumnToFieldName(colName string) string {
	// Special case mappings - add more as needed
	switch colName {
	case "emailVisibility":
		return "emailVisibility"
	case "tokenKey":
		return "tokenKey"
	default:
		// By default, keep as is - this assumes YAML fields match DB column names
		return colName
	}
}

// ResolveReferences resolves reference placeholders in the seed data
func ResolveReferences(config *SeedConfig) error {
	// Create a map to store all items by their collection name and ID
	refMap := make(map[string]map[string]interface{})

	// First pass: collect all items by ID
	for _, collection := range config.Collections {
		collectionMap := make(map[string]interface{})
		for _, item := range collection.Items {
			itemMap, ok := item.(map[string]interface{})
			if !ok {
				continue
			}

			id, ok := itemMap["id"].(string)
			if !ok {
				continue
			}

			collectionMap[id] = itemMap
		}
		refMap[collection.Name] = collectionMap
	}

	// Second pass: resolve references
	for _, collection := range config.Collections {
		for _, item := range collection.Items {
			itemMap, ok := item.(map[string]interface{})
			if !ok {
				continue
			}

			for key, value := range itemMap {
				strValue, ok := value.(string)
				if ok && strings.HasPrefix(strValue, "__ref::") {
					// Parse reference: __ref::collection::id
					parts := strings.Split(strValue, "::")
					if len(parts) != 3 {
						return fmt.Errorf("invalid reference format: %s", strValue)
					}

					refCollection := parts[1]
					refID := parts[2]

					// Find the referenced item
					if collMap, ok := refMap[refCollection]; ok {
						if refItem, ok := collMap[refID]; ok {
							// Replace with actual ID
							refItemMap := refItem.(map[string]interface{})
							if id, ok := refItemMap["id"]; ok {
								itemMap[key] = id
							} else {
								return fmt.Errorf("referenced item has no id: %s", strValue)
							}
						} else {
							return fmt.Errorf("referenced id not found: %s", strValue)
						}
					} else {
						return fmt.Errorf("referenced collection not found: %s", strValue)
					}
				}
			}
		}
	}

	return nil
}

// RunSeedFromYAML is the entry point for the YAML-based seeding tool
func RunSeedFromYAML(configPath string) {
	// If configPath is not provided, use default
	if configPath == "" {
		workDir, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get current directory: %v", err)
		}
		configPath = filepath.Join(workDir, "test", "data", "seed_data.yaml")
	}

	// Seed the database
	if err := SeedDatabaseFromYAML(configPath); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	log.Println("Database seeded successfully from YAML configuration!")
}
