package rest

import (
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"testing"
)

const (
	case01 = `filter={"category":"646cceafccc54408c8ccb308","category_id":"646cceafccc54408c8ccb308","status":false}&order=ASC&page=1&perPage=10&sort=id`

	case02 = "postgres://user:pass@host.com:5432/path?k=v#f"
)

func TestParams_Load(t *testing.T) {
	type args struct {
		payload string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"test",
			args{
				case01,
			},
		},

		{
			"test2",
			args{
				case02,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Params{}
			p.Load(tt.args.payload)
		})
	}
}

func TestParams_ToMongoOptions(t *testing.T) {
	type fields struct {
		Page    int64
		PerPage int64
		Sort    string
		Order   string
		Filter  map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   *options.FindOptions
	}{
		{
			"test",
			fields{
				1,
				15,
				"id",
				"ASC",
				map[string]interface{}{
					"name":   "frank",
					"status": true,
				},
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Params{
				Page:    tt.fields.Page,
				PerPage: tt.fields.PerPage,
				Sort:    tt.fields.Sort,
				Order:   tt.fields.Order,
				Filter:  tt.fields.Filter,
			}
			if got := p.ToMongoOptions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToMongoOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toSkip(t *testing.T) {
	type args struct {
		a int64
		b int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			"test",
			args{
				16,
				15,
			},
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toSkip(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("toSkip() = %v, want %v", got, tt.want)
			}
		})
	}
}
