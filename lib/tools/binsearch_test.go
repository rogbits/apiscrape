package tools

import (
	"reflect"
	"testing"
)

func TestBinRangeSearch(t *testing.T) {
	s := []int64{10, 20, 30}
	rs := NewRangeSearch(s, 10, 30)
	rs.Execute()
	left := rs.Left
	right := rs.Right
	if left != 0 || right != 2 {
		t.FailNow()
	}

	s = []int64{1, 5, 9}
	rs = NewRangeSearch(s, 2, 4)
	rs.Execute()
	left = rs.Left
	right = rs.Right
	if left != -1 || right != -1 {
		t.FailNow()
	}

	s = []int64{10, 20, 30}
	rs = NewRangeSearch(s, 10, 20)
	rs.Execute()
	left = rs.Left
	right = rs.Right
	if left != 0 || right != 1 {
		t.FailNow()
	}

	s = []int64{10, 20, 30}
	rs = NewRangeSearch(s, 20, 30)
	rs.Execute()
	left = rs.Left
	right = rs.Right
	if left != 1 || right != 2 {
		t.FailNow()
	}

	s = []int64{9, 11, 20, 25, 30}
	rs = NewRangeSearch(s, 10, 24)
	rs.Execute()
	left = rs.Left
	right = rs.Right
	if left != 1 || right != 2 {
		t.FailNow()
	}

	s = []int64{1, 2, 3, 4, 5, 6}
	rs = NewRangeSearch(s, 1, 3)
	rs.Execute()
	left = rs.Left
	right = rs.Right
	if left != 0 && right != 2 {
		t.FailNow()
	}

	s = []int64{1, 2, 2, 2, 2, 2, 3}
	rs = NewRangeSearch(s, 2, 2)
	rs.Execute()
	left = rs.Left
	right = rs.Right
	if left != 1 && right != 5 {
		t.FailNow()
	}

	s = []int64{2, 2, 2, 2, 2, 2, 3}
	rs = NewRangeSearch(s, 2, 2)
	rs.Execute()
	left = rs.Left
	right = rs.Right
	if left != 0 && right != 5 {
		t.FailNow()
	}

	s = []int64{2, 2, 2, 2, 2, 2, 3, 3, 3, 3}
	rs = NewRangeSearch(s, 3, 3)
	rs.Execute()
	left = rs.Left
	right = rs.Right
	if left != 6 && right != 9 {
		t.FailNow()
	}

	s = []int64{4, 4, 4, 5, 5, 5}
	rs = NewRangeSearch(s, 3, 4)
	rs.Execute()
	left = rs.Left
	right = rs.Right
	if left != 0 || right != 2 {
		t.FailNow()
	}

	s = []int64{4, 4, 4, 5, 5, 5}
	rs = NewRangeSearch(s, 4, 6)
	rs.Execute()
	left = rs.Left
	right = rs.Right
	if left != 0 || right != 5 {
		t.FailNow()
	}

	s = []int64{4, 4, 4, 5, 5, 5}
	rs = NewRangeSearch(s, 2, 3)
	rs.Execute()
	left = rs.Left
	right = rs.Right
	if left != -1 || right != -1 {
		t.FailNow()
	}

	s = []int64{4, 4, 4, 5, 5, 5}
	rs = NewRangeSearch(s, 8, 8)
	rs.Execute()
	left = rs.Left
	right = rs.Right
	if left != -1 || right != -1 {
		t.FailNow()
	}

	s = []int64{4, 4, 4, 8, 8, 8}
	rs = NewRangeSearch(s, 5, 5)
	rs.Execute()
	left = rs.Left
	right = rs.Right
	if left != -1 && right != -1 {
		t.FailNow()
	}
}

func TestBinInsertInt64(t *testing.T) {
	in := []int64{}
	actual := BinInsertInt64(in, 1)
	expected := []int64{1}
	if !reflect.DeepEqual(actual, expected) {
		t.FailNow()
	}

	in = []int64{1, 2, 3}
	actual = BinInsertInt64(in, 0)
	expected = []int64{0, 1, 2, 3}
	if !reflect.DeepEqual(actual, expected) {
		t.FailNow()
	}

	in = []int64{1, 2, 3}
	actual = BinInsertInt64(in, 4)
	expected = []int64{1, 2, 3, 4}
	if !reflect.DeepEqual(actual, expected) {
		t.FailNow()
	}

	in = []int64{1, 2, 3}
	actual = BinInsertInt64(in, 2)
	expected = []int64{1, 2, 2, 3}
	if !reflect.DeepEqual(actual, expected) {
		t.FailNow()
	}

	in = []int64{1, 2, 2, 3}
	actual = BinInsertInt64(in, 2)
	expected = []int64{1, 2, 2, 2, 3}
	if !reflect.DeepEqual(actual, expected) {
		t.FailNow()
	}

	in = []int64{1, 1}
	actual = BinInsertInt64(in, 2)
	expected = []int64{1, 1, 2}
	if !reflect.DeepEqual(actual, expected) {
		t.FailNow()
	}

	in = []int64{1, 1}
	actual = BinInsertInt64(in, 0)
	expected = []int64{0, 1, 1}
	if !reflect.DeepEqual(actual, expected) {
		t.FailNow()
	}
}
