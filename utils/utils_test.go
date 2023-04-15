package utils

import (
	"errors"
	"reflect"
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	testFn := func(i int) string { return strconv.Itoa(i) }
	inSlice := []int{4, 3, 5, 2}
	want := []string{"4", "3", "5", "2"}

	t.Run("int to string", func(t *testing.T) {
		if got := Map(inSlice, testFn); !reflect.DeepEqual(got, want) {
			t.Errorf("Map() = %v, want %v", got, want)
		}
	})
}

func TestMapWithErr(t *testing.T) {
	testFnErr := func(i int) (string, error) { return "ignore me", errors.New("error") }
	testFnNoErr := func(i int) (string, error) { return strconv.Itoa(i) + "a", nil }
	inSlice := []int{4, 3, 5, 2}
	wantNoErr := []string{"4a", "3a", "5a", "2a"}
	wantErr := errors.New("error")

	t.Run("no error", func(t *testing.T) {
		got, err := MapWithErr(inSlice, testFnNoErr)
		if err != nil {
			t.Fatal("got an error, didn't want one")
		}

		if !reflect.DeepEqual(got, wantNoErr) {
			t.Errorf("got %v, wanted %v", got, wantNoErr)
		}
	})

	t.Run("return error", func(t *testing.T) {
		_, err := MapWithErr(inSlice, testFnErr)
		if err == nil {
			t.Fatal("didn't get an error, wanted one")
		}

		if err.Error() != wantErr.Error() {
			t.Errorf("got errStr %q, wanted errStr %q", err, wantErr)
		}
	})
}

func TestFilter(t *testing.T) {
	testPredicat := func(i int) bool { return i%2 == 0 }
	inSlice := []int{1, 2, 3, 4, 5, 6}
	wantPos := []int{2, 4, 6}
	wantNeg := []int{1, 3, 5}

	t.Run("correct filter", func(t *testing.T) {
		gotPos, gotNeg := Filter(inSlice, testPredicat)
		if !reflect.DeepEqual(gotPos, wantPos) {
			t.Errorf("got %v, wanted %v", gotPos, wantPos)
		}
		if !reflect.DeepEqual(gotNeg, wantNeg) {
			t.Errorf("got %v, wanted %v", gotNeg, wantNeg)
		}
	})
}

func TestNot(t *testing.T) {
	testFn := func(b bool) bool { return b }
	arg := true
	want := false

	t.Run("correct behavior", func(t *testing.T) {
		if got := Not(testFn)(arg); got != want {
			t.Errorf("got %v, wanted %v", got, want)
		}
	})
}
