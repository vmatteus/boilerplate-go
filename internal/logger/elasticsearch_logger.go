package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/your-org/boilerplate-go/internal/config"
)

type ElasticsearchLogger struct {
	config    config.LoggerConfig
	client    *elasticsearch.Client
	fields    map[string]interface{}
	indexName string
}

func NewElasticsearchLogger(cfg config.LoggerConfig) *ElasticsearchLogger {
	// Elasticsearch client configuration
	esCfg := elasticsearch.Config{
		Addresses: []string{cfg.Url},
	}

	// Add authentication if provided
	if cfg.Username != "" && cfg.Password != "" {
		esCfg.Username = cfg.Username
		esCfg.Password = cfg.Password
	}

	if cfg.ApiKey != "" {
		esCfg.APIKey = cfg.ApiKey
	}

	client, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		fmt.Printf("Failed to create Elasticsearch client: %v\n", err)
		return nil
	}

	indexName := cfg.Index
	if indexName == "" {
		indexName = "boilerplate-go-logs"
	}

	return &ElasticsearchLogger{
		config:    cfg,
		client:    client,
		fields:    make(map[string]interface{}),
		indexName: indexName,
	}
}

func (l *ElasticsearchLogger) AddField(key string, value interface{}) {
	if l.fields == nil {
		l.fields = make(map[string]interface{})
	}
	l.fields[key] = value
}

func (l *ElasticsearchLogger) Log(ctx context.Context, level, message string, fields map[string]interface{}) {
	if l.client == nil {
		fmt.Printf("Elasticsearch client not initialized\n")
		return
	}

	// Merge stored fields with provided fields
	mergedFields := make(map[string]interface{})
	for k, v := range l.fields {
		mergedFields[k] = v
	}
	for k, v := range fields {
		mergedFields[k] = v
	}

	// Create ECS-compliant log document
	doc := map[string]interface{}{
		"@timestamp": time.Now().UTC().Format(time.RFC3339),
		"log": map[string]interface{}{
			"level": level,
		},
		"message": message,
		"service": map[string]interface{}{
			"name": "boilerplate-go",
		},
		"ecs": map[string]interface{}{
			"version": "8.0.0",
		},
	}

	// Add custom fields
	for key, value := range mergedFields {
		doc[key] = value
	}

	// Extract request ID from context if available
	if requestID := extractRequestID(ctx); requestID != "" {
		doc["trace"] = map[string]interface{}{
			"id": requestID,
		}
	}

	// Convert document to JSON
	docBytes, err := json.Marshal(doc)
	if err != nil {
		fmt.Printf("Failed to marshal log document: %v\n", err)
		return
	}

	// Create index name with date for daily rotation
	indexName := fmt.Sprintf("%s-%s", l.indexName, time.Now().UTC().Format("2006.01.02"))

	// Send document to Elasticsearch
	res, err := l.client.Index(
		indexName,
		bytes.NewReader(docBytes),
		l.client.Index.WithContext(ctx),
		l.client.Index.WithRefresh("true"),
	)
	if err != nil {
		fmt.Printf("Failed to index log document: %v\n", err)
		return
	}
	defer res.Body.Close()

	if res.IsError() {
		fmt.Printf("Elasticsearch indexing failed: %s\n", res.String())
	}
}

// extractRequestID extracts request ID from context
func extractRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	// Try to extract from common context keys
	if requestID := ctx.Value("x-request-id"); requestID != nil {
		if id, ok := requestID.(string); ok {
			return id
		}
	}

	if requestID := ctx.Value("request-id"); requestID != nil {
		if id, ok := requestID.(string); ok {
			return id
		}
	}

	return ""
}
