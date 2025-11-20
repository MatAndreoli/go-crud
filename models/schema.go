package models

import (
	"encoding/json"
	"fmt"
	"os"
)

type Schema struct {
	TableName string  `json:"table_name"`
	Fields    []Field `json:"fields"`
}

type Field struct {
	Name       string     `json:"name"`
	Type       string     `json:"type"`
	PrimaryKey bool       `json:"primary_key"`
	Required   bool       `json:"required"`
	Validation Validation `json:"validation"`
	Mask       string     `json:"mask"`
}

type Validation struct {
	Type       string      `json:"type"`
	RegexRules []RegexRule `json:"regex_rules"`
}

type RegexRule struct {
	Pattern string `json:"pattern"`
	Message string `json:"message"`
}

func LoadSchema(path string) (*Schema, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo schema: %w", err)
	}

	var schema Schema
	if err := json.Unmarshal(file, &schema); err != nil {
		return nil, fmt.Errorf("erro ao parsear JSON schema: %w", err)
	}

	return &schema, nil
}
