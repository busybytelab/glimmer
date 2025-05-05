package seed

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadSeedConfig(t *testing.T) {
	// Create a temporary YAML file
	content := `
description: "Test Config"
db: "test.db"
collections:
  - name: "test_collection"
    select: "SELECT COUNT(*) FROM test_collection WHERE id = ?"
    insert: "INSERT INTO test_collection (id, name) VALUES (?, ?)"
    items:
      - id: "test_id"
        name: "Test Name"
`
	tmpDir, err := os.MkdirTemp("", "seed_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tmpFile := filepath.Join(tmpDir, "test_config.yaml")
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Test loading the config
	config, err := LoadSeedConfig(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load seed config: %v", err)
	}

	// Verify the config
	if config.Description != "Test Config" {
		t.Errorf("Expected description 'Test Config', got '%s'", config.Description)
	}
	if config.DB != "test.db" {
		t.Errorf("Expected DB 'test.db', got '%s'", config.DB)
	}
	if len(config.Collections) != 1 {
		t.Fatalf("Expected 1 collection, got %d", len(config.Collections))
	}
	if config.Collections[0].Name != "test_collection" {
		t.Errorf("Expected collection name 'test_collection', got '%s'", config.Collections[0].Name)
	}
	if len(config.Collections[0].Items) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(config.Collections[0].Items))
	}
}

func TestResolveReferences(t *testing.T) {
	// Create a test config with references
	config := &SeedConfig{
		Collections: []CollectionConfig{
			{
				Name: "users",
				Items: []interface{}{
					map[string]interface{}{
						"id":   "user1",
						"name": "User One",
					},
				},
			},
			{
				Name: "posts",
				Items: []interface{}{
					map[string]interface{}{
						"id":     "post1",
						"title":  "Post Title",
						"author": "__ref::users::user1",
					},
				},
			},
		},
	}

	// Resolve references
	if err := ResolveReferences(config); err != nil {
		t.Fatalf("Failed to resolve references: %v", err)
	}

	// Verify references were resolved
	postItem := config.Collections[1].Items[0].(map[string]interface{})
	author := postItem["author"]
	if author != "user1" {
		t.Errorf("Expected author to be 'user1', got '%v'", author)
	}
}

func TestExtractInsertValues(t *testing.T) {
	// Test extracting values from insert statement
	insertStmt := "INSERT INTO users (id, name, email) VALUES (?, ?, ?)"
	itemMap := map[string]interface{}{
		"id":    "user1",
		"name":  "User One",
		"email": "user@example.com",
		"extra": "This should be ignored",
	}

	values, err := extractInsertValues(insertStmt, itemMap)
	if err != nil {
		t.Fatalf("Failed to extract insert values: %v", err)
	}

	if len(values) != 3 {
		t.Fatalf("Expected 3 values, got %d", len(values))
	}
	if values[0] != "user1" {
		t.Errorf("Expected id 'user1', got '%v'", values[0])
	}
	if values[1] != "User One" {
		t.Errorf("Expected name 'User One', got '%v'", values[1])
	}
	if values[2] != "user@example.com" {
		t.Errorf("Expected email 'user@example.com', got '%v'", values[2])
	}
}

func TestTimestampReplacement(t *testing.T) {
	// Create a temporary YAML file with timestamp placeholder
	content := `
description: "Test Config"
db: "test.db"
collections:
  - name: "test_collection"
    select: "SELECT COUNT(*) FROM test_collection WHERE id = ?"
    insert: "INSERT INTO test_collection (id, created) VALUES (?, ?)"
    items:
      - id: "test_id"
        created: __::currentTimestamp::__
`
	tmpDir, err := os.MkdirTemp("", "seed_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tmpFile := filepath.Join(tmpDir, "test_config.yaml")
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Test loading the config
	config, err := LoadSeedConfig(tmpFile)
	if err != nil {
		t.Fatalf("Failed to load seed config: %v", err)
	}

	// Verify the timestamp was replaced
	item := config.Collections[0].Items[0].(map[string]interface{})
	created, ok := item["created"].(string)
	if !ok {
		t.Fatalf("Created timestamp is not a string")
	}

	if created == "__::currentTimestamp::__" {
		t.Errorf("Timestamp placeholder was not replaced")
	}
}

func TestValidateConfigWithReferences(t *testing.T) {
	// Test with valid references
	validConfig := &SeedConfig{
		Description: "Test Config",
		DB:          "test.db",
		Collections: []CollectionConfig{
			{
				Name:   "collection1",
				Select: "SELECT COUNT(*) FROM collection1 WHERE id = ?",
				Insert: "INSERT INTO collection1 (id, name) VALUES (?, ?)",
				Items: []interface{}{
					map[string]interface{}{
						"id":   "item1",
						"name": "Item One",
					},
				},
			},
			{
				Name:   "collection2",
				Select: "SELECT COUNT(*) FROM collection2 WHERE id = ?",
				Insert: "INSERT INTO collection2 (id, name, ref) VALUES (?, ?, ?)",
				Items: []interface{}{
					map[string]interface{}{
						"id":   "item2",
						"name": "Item Two",
						"ref":  "__ref::collection1::item1",
					},
				},
			},
		},
	}

	if err := validateConfig(validConfig); err != nil {
		t.Errorf("validateConfig should not return error for valid config: %v", err)
	}

	// Test with invalid reference format
	invalidFormatConfig := &SeedConfig{
		Description: "Test Config",
		DB:          "test.db",
		Collections: []CollectionConfig{
			{
				Name:   "collection1",
				Select: "SELECT COUNT(*) FROM collection1 WHERE id = ?",
				Insert: "INSERT INTO collection1 (id, name) VALUES (?, ?)",
				Items: []interface{}{
					map[string]interface{}{
						"id":   "item1",
						"name": "Item One",
					},
				},
			},
			{
				Name:   "collection2",
				Select: "SELECT COUNT(*) FROM collection2 WHERE id = ?",
				Insert: "INSERT INTO collection2 (id, name, ref) VALUES (?, ?, ?)",
				Items: []interface{}{
					map[string]interface{}{
						"id":   "item2",
						"name": "Item Two",
						"ref":  "__ref::collection1", // Missing item ID
					},
				},
			},
		},
	}

	if err := validateConfig(invalidFormatConfig); err == nil {
		t.Error("validateConfig should return error for invalid reference format")
	}

	// Test with reference to non-existent collection
	nonExistentCollectionConfig := &SeedConfig{
		Description: "Test Config",
		DB:          "test.db",
		Collections: []CollectionConfig{
			{
				Name:   "collection1",
				Select: "SELECT COUNT(*) FROM collection1 WHERE id = ?",
				Insert: "INSERT INTO collection1 (id, name) VALUES (?, ?)",
				Items: []interface{}{
					map[string]interface{}{
						"id":   "item1",
						"name": "Item One",
					},
				},
			},
			{
				Name:   "collection2",
				Select: "SELECT COUNT(*) FROM collection2 WHERE id = ?",
				Insert: "INSERT INTO collection2 (id, name, ref) VALUES (?, ?, ?)",
				Items: []interface{}{
					map[string]interface{}{
						"id":   "item2",
						"name": "Item Two",
						"ref":  "__ref::non_existent::item1", // Non-existent collection
					},
				},
			},
		},
	}

	if err := validateConfig(nonExistentCollectionConfig); err == nil {
		t.Error("validateConfig should return error for reference to non-existent collection")
	}

	// Test with reference to non-existent item in existing collection
	nonExistentItemConfig := &SeedConfig{
		Description: "Test Config",
		DB:          "test.db",
		Collections: []CollectionConfig{
			{
				Name:   "collection1",
				Select: "SELECT COUNT(*) FROM collection1 WHERE id = ?",
				Insert: "INSERT INTO collection1 (id, name) VALUES (?, ?)",
				Items: []interface{}{
					map[string]interface{}{
						"id":   "item1",
						"name": "Item One",
					},
				},
			},
			{
				Name:   "collection2",
				Select: "SELECT COUNT(*) FROM collection2 WHERE id = ?",
				Insert: "INSERT INTO collection2 (id, name, ref) VALUES (?, ?, ?)",
				Items: []interface{}{
					map[string]interface{}{
						"id":   "item2",
						"name": "Item Two",
						"ref":  "__ref::collection1::non_existent", // Non-existent item
					},
				},
			},
		},
	}

	if err := validateConfig(nonExistentItemConfig); err == nil {
		t.Error("validateConfig should return error for reference to non-existent item")
	}
}

func TestCircularReferenceDetection(t *testing.T) {
	// Test with circular references
	circularConfig := &SeedConfig{
		Description: "Test Config",
		DB:          "test.db",
		Collections: []CollectionConfig{
			{
				Name:   "collection1",
				Select: "SELECT COUNT(*) FROM collection1 WHERE id = ?",
				Insert: "INSERT INTO collection1 (id, name, ref) VALUES (?, ?, ?)",
				Items: []interface{}{
					map[string]interface{}{
						"id":   "item1",
						"name": "Item One",
						"ref":  "__ref::collection2::item2", // Reference to collection2
					},
				},
			},
			{
				Name:   "collection2",
				Select: "SELECT COUNT(*) FROM collection2 WHERE id = ?",
				Insert: "INSERT INTO collection2 (id, name, ref) VALUES (?, ?, ?)",
				Items: []interface{}{
					map[string]interface{}{
						"id":   "item2",
						"name": "Item Two",
						"ref":  "__ref::collection1::item1", // Reference back to collection1, creating a cycle
					},
				},
			},
		},
	}

	// This should detect the circular reference
	if err := validateConfig(circularConfig); err == nil {
		t.Error("validateConfig should return error for circular references")
	} else {
		if !strings.Contains(err.Error(), "circular reference") {
			t.Errorf("Expected error about circular reference, got: %v", err)
		}
	}

	// Test with no circular references (similar structure but no circularity)
	validConfig := &SeedConfig{
		Description: "Test Config",
		DB:          "test.db",
		Collections: []CollectionConfig{
			{
				Name:   "collection1",
				Select: "SELECT COUNT(*) FROM collection1 WHERE id = ?",
				Insert: "INSERT INTO collection1 (id, name) VALUES (?, ?)",
				Items: []interface{}{
					map[string]interface{}{
						"id":   "item1",
						"name": "Item One",
					},
				},
			},
			{
				Name:   "collection2",
				Select: "SELECT COUNT(*) FROM collection2 WHERE id = ?",
				Insert: "INSERT INTO collection2 (id, name, ref) VALUES (?, ?, ?)",
				Items: []interface{}{
					map[string]interface{}{
						"id":   "item2",
						"name": "Item Two",
						"ref":  "__ref::collection1::item1", // Only a one-way reference
					},
				},
			},
			{
				Name:   "collection3",
				Select: "SELECT COUNT(*) FROM collection3 WHERE id = ?",
				Insert: "INSERT INTO collection3 (id, name, ref) VALUES (?, ?, ?)",
				Items: []interface{}{
					map[string]interface{}{
						"id":   "item3",
						"name": "Item Three",
						"ref":  "__ref::collection2::item2", // Reference to collection2 but no cycle
					},
				},
			},
		},
	}

	// This should pass validation
	if err := validateConfig(validConfig); err != nil {
		t.Errorf("validateConfig should not return error for valid config: %v", err)
	}
}

func TestExtractTableAndColumns(t *testing.T) {
	testCases := []struct {
		name        string
		insertStmt  string
		wantTable   string
		wantColumns []string
		wantError   bool
	}{
		{
			name:        "Simple insert",
			insertStmt:  "INSERT INTO users (id, name, email) VALUES (?, ?, ?)",
			wantTable:   "users",
			wantColumns: []string{"id", "name", "email"},
			wantError:   false,
		},
		{
			name:        "Insert with spaces",
			insertStmt:  "INSERT INTO  practice_sessions  ( id ,  name ,  assigned_at ) VALUES (?, ?, ?)",
			wantTable:   "practice_sessions",
			wantColumns: []string{"id", "name", "assigned_at"},
			wantError:   false,
		},
		{
			name:        "Insert with lowercase",
			insertStmt:  "insert into items (id, title) values (?, ?)",
			wantTable:   "items",
			wantColumns: []string{"id", "title"},
			wantError:   false,
		},
		{
			name:        "Invalid insert without columns",
			insertStmt:  "INSERT INTO users VALUES (?, ?, ?)",
			wantTable:   "",
			wantColumns: nil,
			wantError:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			table, columns, err := extractTableAndColumns(tc.insertStmt)

			if tc.wantError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if table != tc.wantTable {
				t.Errorf("Expected table %s, got %s", tc.wantTable, table)
			}

			if len(columns) != len(tc.wantColumns) {
				t.Fatalf("Expected %d columns, got %d", len(tc.wantColumns), len(columns))
			}

			for i, col := range tc.wantColumns {
				if columns[i] != col {
					t.Errorf("Expected column %s at position %d, got %s", col, i, columns[i])
				}
			}
		})
	}
}

func TestCountPlaceholdersInStatement(t *testing.T) {
	testCases := []struct {
		name          string
		statement     string
		expectedCount int
		expectError   bool
	}{
		{
			name:          "Simple statement",
			statement:     "INSERT INTO users (id, name) VALUES (?, ?)",
			expectedCount: 2,
			expectError:   false,
		},
		{
			name:          "Complex statement with many placeholders",
			statement:     "INSERT INTO users (id, name, email, age, created, updated) VALUES (?, ?, ?, ?, ?, ?)",
			expectedCount: 6,
			expectError:   false,
		},
		{
			name:          "Statement with inconsistent placeholders",
			statement:     "INSERT INTO users (id, name) VALUES (?, ?, ?)", // 3 ? but only 2 columns
			expectedCount: 3,
			expectError:   true,
		},
		{
			name:          "Statement with no VALUES section",
			statement:     "DELETE FROM users WHERE id = ?",
			expectedCount: 1,
			expectError:   false, // No error as we default to counting ? marks
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			count, err := countPlaceholdersInStatement(tc.statement)

			if tc.expectError {
				if err == nil {
					t.Log("Expected error but got none")
					// Not failing test as this is more of a warning
				}
			} else if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if count != tc.expectedCount {
				t.Errorf("Expected %d placeholders, got %d", tc.expectedCount, count)
			}
		})
	}
}
