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
			gotTopic, gotOp := extractTopicAndOp(tt.param)
			if gotTopic != tt.wantTopic || gotOp != tt.wantOp {
				t.Errorf("extractTopicAndOp() = (%v, %v), want (%v, %v)", gotTopic, gotOp, tt.wantTopic, tt.wantOp)
			}
		})
	}
}
