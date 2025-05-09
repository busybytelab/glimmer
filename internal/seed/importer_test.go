package seed

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase/tests"
)

func TestLoadSeedConfig(t *testing.T) {
	// Create a temporary YAML file
	content := `
description: "Test Config"
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

func TestSeedDatabaseFromYAML(t *testing.T) {
	// Create a test app instance
	testApp, err := tests.NewTestApp()
	if err != nil {
		t.Fatalf("Failed to create test app: %v", err)
	}
	defer testApp.Cleanup()

	// Create test table
	_, err = testApp.DB().NewQuery(`
		CREATE TABLE test_collection (
			id TEXT PRIMARY KEY,
			name TEXT,
			created_at TEXT,
			updated_at TEXT
		)
	`).Execute()
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	// Create a temporary YAML file with test configuration
	content := `
collections:
  - name: test_collection
    select: SELECT COUNT(*) FROM test_collection WHERE id = {:id}
    insert: INSERT INTO test_collection (id, name, created_at, updated_at) VALUES ({:id}, {:name}, {:created_at}, {:updated_at})
    items:
      - id: "test1"
        name: "Test Item 1"
        created_at: "__::currentTimestamp::__"
        updated_at: "__::currentTimestamp::__"
      - id: "test2"
        name: "Test Item 2"
        created_at: "__::currentTimestamp::__"
        updated_at: "__::currentTimestamp::__"
`
	tmpFile, err := os.CreateTemp("", "test-seed-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Run the seed
	if err := SeedDatabaseFromYAML(testApp, tmpFile.Name()); err != nil {
		t.Fatalf("Failed to seed database: %v", err)
	}

	// Verify the seeded data
	var count int
	err = testApp.DB().NewQuery("SELECT COUNT(*) FROM test_collection").Row(&count)
	if err != nil {
		t.Fatalf("Failed to verify seeded data: %v", err)
	}

	if count != 2 {
		t.Errorf("Expected 2 items in test_collection, got %d", count)
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

func TestExtractTableName(t *testing.T) {
	tests := []struct {
		name     string
		insert   string
		expected string
	}{
		{
			name:     "Simple insert",
			insert:   "INSERT INTO users (id, name) VALUES (?, ?)",
			expected: "users",
		},
		{
			name:     "Insert with spaces",
			insert:   "INSERT INTO  users  (id, name) VALUES (?, ?)",
			expected: "users",
		},
		{
			name:     "Insert with lowercase",
			insert:   "insert into users (id, name) values (?, ?)",
			expected: "users",
		},
		{
			name:     "Invalid insert without columns",
			insert:   "INSERT INTO users",
			expected: "users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractTableName(tt.insert)
			if got != tt.expected {
				t.Errorf("Expected table %s, got %s", tt.expected, got)
			}
		})
	}
}

func TestProcessEnvVarToken(t *testing.T) {
	// Save original environment and restore after test
	originalEnv := os.Getenv("TEST_ENV_VAR")
	defer os.Setenv("TEST_ENV_VAR", originalEnv)

	tests := []struct {
		name                string
		token               string
		defaultPasswordHash string
		envValue            string
		expected            string
	}{
		{
			name:                "valid env var with value",
			token:               "__::env::TEST_ENV_VAR::default_value::__",
			defaultPasswordHash: "default_hash",
			envValue:            "env_value",
			expected:            "env_value",
		},
		{
			name:                "valid env var without value",
			token:               "__::env::TEST_ENV_VAR::default_value::__",
			defaultPasswordHash: "default_hash",
			envValue:            "",
			expected:            "default_value",
		},
		{
			name:                "default password hash reference",
			token:               "__::env::TEST_ENV_VAR::default_password_hash::__",
			defaultPasswordHash: "default_hash",
			envValue:            "",
			expected:            "default_hash",
		},
		{
			name:                "invalid token format",
			token:               "invalid_token",
			defaultPasswordHash: "default_hash",
			envValue:            "",
			expected:            "invalid_token",
		},
		{
			name:                "invalid env token format",
			token:               "__::env::TEST_ENV_VAR::__",
			defaultPasswordHash: "default_hash",
			envValue:            "",
			expected:            "__::env::TEST_ENV_VAR::__",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("TEST_ENV_VAR", tt.envValue)
			result := processEnvVarToken(tt.token, tt.defaultPasswordHash)
			if result != tt.expected {
				t.Errorf("processEnvVarToken() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestProcessItemValues(t *testing.T) {
	// Save original environment and restore after test
	originalEnv := os.Getenv("TEST_ENV_VAR")
	defer os.Setenv("TEST_ENV_VAR", originalEnv)

	tests := []struct {
		name                string
		itemMap             map[string]interface{}
		defaultPasswordHash string
		envValue            string
		expected            map[string]interface{}
	}{
		{
			name: "process env var token",
			itemMap: map[string]interface{}{
				"field1": "__::env::TEST_ENV_VAR::default_value::__",
				"field2": "normal_value",
			},
			defaultPasswordHash: "default_hash",
			envValue:            "env_value",
			expected: map[string]interface{}{
				"field1": "env_value",
				"field2": "normal_value",
			},
		},
		{
			name: "process default password hash",
			itemMap: map[string]interface{}{
				"password": "__::env::TEST_ENV_VAR::default_password_hash::__",
			},
			defaultPasswordHash: "default_hash",
			envValue:            "",
			expected: map[string]interface{}{
				"password": "default_hash",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("TEST_ENV_VAR", tt.envValue)
			processItemValues(tt.itemMap, tt.defaultPasswordHash)
			if !reflect.DeepEqual(tt.itemMap, tt.expected) {
				t.Errorf("processItemValues() = %v, want %v", tt.itemMap, tt.expected)
			}
		})
	}
}
