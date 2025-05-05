package llm

import (
	"fmt"

	"github.com/busybytelab.com/glimmer/internal/domain"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// ensures that the LLMResponseRecord struct satisfies the core.RecordProxy interface
var _ core.RecordProxy = (*LLMResponseRecord)(nil)

// LLMResponseRecord represents a cached LLM response in the database
type LLMResponseRecord struct {
	core.BaseRecordProxy
	Key              string  `db:"key" json:"key"`
	Prompt           string  `db:"prompt" json:"prompt"`
	SystemPrompt     string  `db:"system_prompt" json:"system_prompt"`
	Response         string  `db:"response" json:"response"`
	ModelName        string  `db:"model_name" json:"model_name"`
	Backend          string  `db:"backend" json:"backend"`
	TotalTokens      int     `db:"total_tokens" json:"total_tokens"`
	PromptTokens     int     `db:"prompt_tokens" json:"prompt_tokens"`
	CompletionTokens int     `db:"completion_tokens" json:"completion_tokens"`
	Cost             float64 `db:"cost" json:"cost"`
	TTL              int     `db:"ttl" json:"ttl"`
}

// LLMResponseRecordDao provides database access methods for LLMResponseRecord
type LLMResponseRecordDao struct {
	app core.App
}

// NewLLMResponseRecordDao creates a new DAO for LLMResponseRecord
func NewLLMResponseRecordDao(app core.App) *LLMResponseRecordDao {
	return &LLMResponseRecordDao{
		app: app,
	}
}

// FindLLMResponseRecordByKey finds a cached LLM response by its key
func (dao *LLMResponseRecordDao) FindLLMResponseRecordByKey(key string) (*LLMResponseRecord, error) {
	cache, err := dao.app.FindFirstRecordByData(domain.CollectionLLMResponses, "key", key)

	if err != nil {
		return nil, err
	}

	result := &LLMResponseRecord{
		BaseRecordProxy: core.BaseRecordProxy{
			Record: cache,
		},
	}

	// Map the fields from the record to the struct
	result.Key = cache.GetString("key")
	result.Prompt = cache.GetString("prompt")
	result.SystemPrompt = cache.GetString("system_prompt")
	result.Response = cache.GetString("response")
	result.ModelName = cache.GetString("model_name")
	result.Backend = cache.GetString("backend")
	result.TotalTokens = int(cache.GetInt("total_tokens"))
	result.PromptTokens = int(cache.GetInt("prompt_tokens"))
	result.CompletionTokens = int(cache.GetInt("completion_tokens"))
	result.Cost = cache.GetFloat("cost")
	result.TTL = int(cache.GetInt("ttl"))

	return result, nil
}

// SaveLLMResponseRecord saves or updates a cached LLM response
func (dao *LLMResponseRecordDao) SaveLLMResponseRecord(record *LLMResponseRecord) error {
	// Check if the collection exists
	collection, err := dao.app.FindCollectionByNameOrId(domain.CollectionLLMResponses)
	if err != nil {
		return fmt.Errorf("failed to find %s collection: %w", domain.CollectionLLMResponses, err)
	}

	// If record is new, create it
	if record.BaseRecordProxy.Record == nil {
		record.BaseRecordProxy = core.BaseRecordProxy{
			Record: core.NewRecord(collection),
		}
	}

	// Set all required fields using the Set method
	record.Set("key", record.Key)
	record.Set("prompt", record.Prompt)
	record.Set("system_prompt", record.SystemPrompt)
	record.Set("response", record.Response)
	record.Set("model_name", record.ModelName)
	record.Set("backend", record.Backend)
	if record.TotalTokens > 0 {
		record.Set("total_tokens", record.TotalTokens)
	}
	if record.Cost > 0 {
		record.Set("cost", record.Cost)
	}
	if record.PromptTokens > 0 {
		record.Set("prompt_tokens", record.PromptTokens)
	}
	if record.CompletionTokens > 0 {
		record.Set("completion_tokens", record.CompletionTokens)
	}
	if record.TTL > 0 {
		record.Set("ttl", record.TTL)
	}

	return dao.app.Save(record)
}

// DeleteLLMResponseRecord deletes a cached LLM response
func (dao *LLMResponseRecordDao) DeleteLLMResponseRecord(record *LLMResponseRecord) error {
	return dao.app.Delete(record)
}

// DeleteExpiredLLMResponseRecords deletes all expired cache entries
func (dao *LLMResponseRecordDao) DeleteExpiredLLMResponseRecords() error {
	var expired []*LLMResponseRecord

	err := dao.app.RecordQuery(domain.CollectionLLMResponses).
		AndWhere(dbx.NewExp("ttl > 0 AND created + ttl < NOW()")).
		All(&expired)

	if err != nil {
		return err
	}

	for _, cache := range expired {
		if err := dao.app.Delete(cache); err != nil {
			return err
		}
	}

	return nil
}
