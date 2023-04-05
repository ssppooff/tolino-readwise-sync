package utils

import (
	"reflect"
	"testing"
)

func TestMap(t *testing.T) {
	type args struct {
		entries []E
		fn      func(E) F
	}
	tests := []struct {
		name string
		args args
		want []F
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Map(tt.args.entries, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Fail("not implemented")
				// t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}
