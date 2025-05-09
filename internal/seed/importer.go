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
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"gopkg.in/yaml.v3"
)

// SeedConfig represents the top-level structure of the seed YAML file
type SeedConfig struct {
	Description         string             `yaml:"description"`
	DB                  string             `yaml:"db"`
	DefaultPasswordHash string             `yaml:"default_password_hash"`
	Collections         []CollectionConfig `yaml:"collections"`
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

// processEnvVarToken processes a token in the format __::env::ENV_VAR_NAME::default_value::__
// It returns the value from the environment variable if set, otherwise returns the default value
func processEnvVarToken(token string, defaultPasswordHash string) string {
	// Extract environment variable name and default value
	parts := strings.Split(token, "::")
	if len(parts) != 5 || parts[0] != "__" || parts[4] != "__" || parts[1] != "env" {
		return token // Return original token if format doesn't match
	}

	envVarName := parts[2]
	defaultValue := parts[3]

	// If default value is "default_password_hash", use the config's default password hash
	if defaultValue == "default_password_hash" {
		defaultValue = defaultPasswordHash
	}

	// Get value from environment variable
	if value := os.Getenv(envVarName); value != "" {
		// Add info line for TEST_USER_PASSWORD_HASH
		if envVarName == "TEST_USER_PASSWORD_HASH" {
			log.Printf("Using TEST_USER_PASSWORD_HASH from environment")
		}
		return value
	}

	return defaultValue
}

// processItemValues processes all values in an item map, handling special tokens
func processItemValues(itemMap map[string]interface{}, defaultPasswordHash string) {
	for key, value := range itemMap {
		if strValue, ok := value.(string); ok {
			if strings.HasPrefix(strValue, "__::env::") {
				itemMap[key] = processEnvVarToken(strValue, defaultPasswordHash)
			}
		}
	}
}

// seedCollection seeds a single collection with its items
func seedCollection(app core.App, collection CollectionConfig, defaultPasswordHash string) error {
	log.Printf("Seeding collection: %s", collection.Name)

	for i, item := range collection.Items {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return fmt.Errorf("invalid item format in collection %s", collection.Name)
		}

		// Process environment variable tokens in the item
		processItemValues(itemMap, defaultPasswordHash)

		// Get the ID for checking existence
		id, ok := itemMap["id"]
		if !ok {
			return fmt.Errorf("item %d in collection %s has no id field", i, collection.Name)
		}

		// Check if item exists
		var count int
		err := app.DB().NewQuery(collection.Select).Bind(dbx.Params{"id": id}).Row(&count)
		if err != nil {
			return fmt.Errorf("failed to check if item exists in %s: %w", collection.Name, err)
		}

		if count == 0 {
			// Extract values from item map in the correct order for the INSERT statement
			values, err := extractInsertValues(collection.Insert, itemMap)
			if err != nil {
				return fmt.Errorf("failed to extract values for item %s in %s: %w", id, collection.Name, err)
			}

			// Create params map for the insert query
			params := make(dbx.Params)
			for key, value := range itemMap {
				params[key] = value
			}

			// Insert item
			_, err = app.DB().NewQuery(collection.Insert).Bind(params).Execute()
			if err != nil {
				// Check if it's a column count mismatch error
				if strings.Contains(err.Error(), "values for") && strings.Contains(err.Error(), "columns") {
					// Parse the expected number of columns from INSERT statement
					expectedCols := strings.Count(collection.Insert, "{:")
					actualCols := len(values)
					tableName := extractTableName(collection.Insert)

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
			log.Printf("Skipping existing item %s in collection %s", id, collection.Name)
		}
	}

	return nil
}

// extractInsertValues extracts values from an item map in the order specified by the INSERT statement
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

// extractTableName extracts the table name from an INSERT statement
func extractTableName(insertStmt string) string {
	// Extract table name using regex - making it case-insensitive with (?i)
	// This pattern matches both "INSERT INTO table (" and "INSERT INTO table" formats
	tableRegex := regexp.MustCompile(`(?i)INSERT\s+INTO\s+(\w+)(?:\s*\(|\s+|$)`)
	tableMatches := tableRegex.FindStringSubmatch(insertStmt)
	if len(tableMatches) < 2 {
		return ""
	}
	return tableMatches[1]
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

// SeedDatabaseFromYAML seeds the database using the provided YAML configuration
func SeedDatabaseFromYAML(app core.App, configPath string) error {
	config, err := LoadSeedConfig(configPath)
	if err != nil {
		return err
	}

	// Resolve references in the config
	if err := ResolveReferences(config); err != nil {
		return fmt.Errorf("failed to resolve references: %w", err)
	}

	// Process all collections in a single transaction
	return app.RunInTransaction(func(txApp core.App) error {
		for _, collection := range config.Collections {
			if err := seedCollection(txApp, collection, config.DefaultPasswordHash); err != nil {
				return err
			}
		}
		return nil
	})
}

// RunSeedFromYAML is the entry point for the YAML-based seeding tool
func RunSeedFromYAML(app core.App, configPath string) error {
	// If configPath is not provided, use default
	if configPath == "" {
		workDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		configPath = filepath.Join(workDir, "test", "data", "seed_data.yaml")
	}

	// Seed the database
	if err := SeedDatabaseFromYAML(app, configPath); err != nil {
		return fmt.Errorf("failed to seed database: %w", err)
	}

	log.Println("Database seeded successfully from YAML configuration!")
	return nil
}
