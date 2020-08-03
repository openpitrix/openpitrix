package helm

import (
	"reflect"
	"testing"
)

func Test_trimStringInMap(t *testing.T) {
	type args struct {
		values map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			"", args{map[string]interface{}{
				"":   "",
				"x2": "y",
				"x1": "y",
				"x4": map[string]interface{}{},
				"x3": map[string]interface{}{"x": ""},
				"x6": map[string]interface{}{"x": 42},
				"x5": map[string]interface{}{"x": map[string]interface{}{"x": ""}},
				"x7": map[string]interface{}{"x": map[string]interface{}{"x": map[string]interface{}{"x": map[string]interface{}{"x": ""}}}},
			}}, map[string]interface{}{
				"x2": "y",
				"x1": "y",
				"x4": map[string]interface{}{},
				"x3": map[string]interface{}{},
				"x6": map[string]interface{}{"x": 42},
				"x5": map[string]interface{}{"x": map[string]interface{}{}},
				"x7": map[string]interface{}{"x": map[string]interface{}{"x": map[string]interface{}{"x": map[string]interface{}{}}}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trimStringInMap(tt.args.values); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("trimStringInMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
