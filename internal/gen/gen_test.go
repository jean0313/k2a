package gen

import (
	"reflect"
	"testing"
)

func TestExtractMessageNameAndSchemaName(t *testing.T) {
	tests := []struct {
		name           string
		message        *Message
		mapping        map[string]string
		expectedMsg    string
		expectedSchema string
	}{
		{
			"Message with reference",
			&Message{
				Reference: &Reference{Ref: "#/components/messages/MessageName"},
			},
			map[string]string{"MessageName": "SchemaName"},
			"MessageName",
			"SchemaName",
		},
		{
			"Message with reference but no mapping",
			&Message{
				Reference: &Reference{Ref: "#/components/messages/MessageName"},
			},
			map[string]string{},
			"MessageName",
			"",
		},
		{
			"Message with MessageEntity",
			&Message{
				MessageEntity: &MessageEntity{
					Name:    "MessageName",
					Payload: map[string]interface{}{"$ref": "#/components/schemas/SchemaName"},
				},
			},
			map[string]string{},
			"MessageName",
			"SchemaName",
		},
		{
			"Message without reference or MessageEntity",
			&Message{},
			map[string]string{},
			"",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msgName, schemaName := extractMessageNameAndSchemaName(tt.message, tt.mapping)
			if msgName != tt.expectedMsg || schemaName != tt.expectedSchema {
				t.Errorf("extractMessageNameAndSchemaName(%v, %v) = (%q, %q), want (%q, %q)", tt.message, tt.mapping, msgName, schemaName, tt.expectedMsg, tt.expectedSchema)
			}
		})
	}
}

func TestExtractSchemaName(t *testing.T) {
	tests := []struct {
		name     string
		input    *MessageEntity
		expected string
	}{
		{
			"Payload with $ref",
			&MessageEntity{
				Payload: map[string]interface{}{"$ref": "#/components/schemas/SchemaName"},
			},
			"SchemaName",
		},
		{
			"Payload with name",
			&MessageEntity{
				Payload: map[string]interface{}{"name": "SchemaName"},
			},
			"SchemaName",
		},
		{
			"Payload without $ref or name",
			&MessageEntity{
				Name:    "SchemaName",
				Payload: map[string]interface{}{},
			},
			"SchemaName",
		},
		{
			"Nil payload",
			&MessageEntity{
				Name: "SchemaName",
			},
			"SchemaName",
		},
		{
			"Nil MessageEntity",
			nil,
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := extractSchemaName(tt.input)
			if output != tt.expected {
				t.Errorf("extractSchemaName(%v) = %q, want %q", tt.input, output, tt.expected)
			}
		})
	}
}

func TestConvertMapI2MapS(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			"Convert map[interface{}]interface{} to map[string]interface{}",
			map[interface{}]interface{}{"key": "value", 1: 2},
			map[string]interface{}{"key": "value", "1": 2},
		},
		{
			"Convert nested map[interface{}]interface{} to map[string]interface{}",
			map[interface{}]interface{}{"key": map[interface{}]interface{}{1: "value"}},
			map[string]interface{}{"key": map[string]interface{}{"1": "value"}},
		},
		{
			"Convert []interface{} containing map[interface{}]interface{}",
			[]interface{}{map[interface{}]interface{}{1: "value"}},
			[]interface{}{map[string]interface{}{"1": "value"}},
		},
		{
			"Convert map[string]interface{}",
			map[string]interface{}{"key": "value"},
			map[string]interface{}{"key": "value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := convertMapI2MapS(tt.input)
			if !reflect.DeepEqual(output, tt.expected) {
				t.Errorf("convertMapI2MapS(%v) = %v, want %v", tt.input, output, tt.expected)
			}
		})
	}
}

func TestCapitalize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"First letter lowercase", "hello", "Hello"},
		{"First letter uppercase", "World", "World"},
		{"Empty string", "", ""},
		{"Non-alphabetic first character", "123abc", "123abc"},
		{"String with one character", "a", "A"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := capitalize(tt.input)
			if output != tt.expected {
				t.Errorf("capitalize(%q) = %q, want %q", tt.input, output, tt.expected)
			}
		})
	}
}

func TestUncapitalize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"First letter uppercase", "Hello", "hello"},
		{"First letter lowercase", "world", "world"},
		{"Empty string", "", ""},
		{"Non-alphabetic first character", "123ABC", "123ABC"},
		{"String with one character", "A", "a"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := uncapitalize(tt.input)
			if output != tt.expected {
				t.Errorf("uncapitalize(%q) = %q, want %q", tt.input, output, tt.expected)
			}
		})
	}
}
