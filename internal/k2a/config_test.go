package k2a

import (
	"testing"
)

func TestExtractTopicAndOp(t *testing.T) {
	tests := []struct {
		name      string
		param     string
		wantTopic string
		wantOp    string
	}{
		{
			name:      "Single parameter",
			param:     "topic1",
			wantTopic: "topic1",
			wantOp:    "Publish",
		},
		{
			name:      "Two parameters with publish operation",
			param:     "topic1:p",
			wantTopic: "topic1",
			wantOp:    "Publish",
		},
		{
			name:      "Two parameters with subscribe operation",
			param:     "topic1:s",
			wantTopic: "topic1",
			wantOp:    "Subscribe",
		},
		{
			name:      "Two parameters with unknown operation",
			param:     "topic1:x",
			wantTopic: "topic1",
			wantOp:    "Subscribe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTopic, gotOp, _ := extractTopicAndOp(tt.param)
			if gotTopic != tt.wantTopic || gotOp != tt.wantOp {
				t.Errorf("extractTopicAndOp() = (%v, %v), want (%v, %v)", gotTopic, gotOp, tt.wantTopic, tt.wantOp)
			}
		})
	}
}

func TestExtractTopicAndOp_Error(t *testing.T) {
	tests := []struct {
		name      string
		param     string
		wantTopic string
		wantOp    string
		wantErr   bool
	}{
		{
			name:      "Unsupported operation",
			param:     "topic1:x",
			wantTopic: "",
			wantOp:    "",
			wantErr:   true,
		},
		{
			name:      "Empty operation",
			param:     "topic1:",
			wantTopic: "",
			wantOp:    "",
			wantErr:   true,
		},
		{
			name:      "Empty topic and operation",
			param:     ":",
			wantTopic: "",
			wantOp:    "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTopic, gotOp, err := extractTopicAndOp(tt.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractTopicAndOp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTopic != tt.wantTopic || gotOp != tt.wantOp {
				t.Errorf("extractTopicAndOp() = (%v, %v), want (%v, %v)", gotTopic, gotOp, tt.wantTopic, tt.wantOp)
			}
		})
	}
}
