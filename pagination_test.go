package rest

import "testing"

func TestGetContentRange(t *testing.T) {
	type args struct {
		tpl       string
		offset    uint
		perPage   uint
		totalPage uint
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"test",
			args{
				RectJsAdminPageTpl,
				0,
				14,
				20,
			},
			"items 0-14/20",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetContentRange(tt.args.tpl, tt.args.offset, tt.args.perPage, tt.args.totalPage); got != tt.want {
				t.Errorf("GetContentRange() = %v, want %v", got, tt.want)
			}
		})
	}
}
