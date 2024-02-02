package k2a

import "testing"

func Test_extractTopicAndOp(t *testing.T) {
	type args struct {
		param string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		{
			name:    "Test with only topic",
			args:    args{param: "topic1"},
			want:    "topic1",
			want1:   "Publish",
			wantErr: false,
		},
		{
			name:    "Test with topic and publish operation",
			args:    args{param: "topic1:p"},
			want:    "topic1",
			want1:   "Publish",
			wantErr: false,
		},
		{
			name:    "Test with topic and subscribe operation",
			args:    args{param: "topic1:s"},
			want:    "topic1",
			want1:   "Subscribe",
			wantErr: false,
		},
		{
			name:    "Test with unsupported operation",
			args:    args{param: "topic1:x"},
			want:    "",
			want1:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := extractTopicAndOp(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractTopicAndOp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("extractTopicAndOp() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("extractTopicAndOp() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
