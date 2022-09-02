package tools

// RangeSearch used for finding bookend indices of a given lo and hi
type RangeSearch struct {
	Slice []int64
	Lower int64
	Upper int64
	Left  int64
	Right int64
}

func NewRangeSearch(s []int64, start int64, end int64) *RangeSearch {
	search := new(RangeSearch)
	search.Slice = s
	search.Lower = start
	search.Upper = end
	return search
}

func (s *RangeSearch) Execute() {
	hi := int64(len(s.Slice)) - 1
	left := binRangeSearch(s.Slice, 0, hi, s.Lower, Left)
	right := binRangeSearch(s.Slice, 0, hi, s.Upper, Right)
	if left > right {
		left = -1
		right = -1
	}
	if left == -1 && right != -1 {
		left = 0
	}
	s.Left = left
	s.Right = right
}

type Position int

const (
	Left  Position = 0
	Right Position = 1
)

func binRangeSearch(slice []int64, lo int64, hi int64, target int64, pos Position) int64 {
	for lo <= hi {
		mid := (lo + hi) / 2
		if target < slice[mid] {
			hi = mid - 1
			continue
		}
		if target > slice[mid] {
			lo = mid + 1
			continue
		}
		switch {
		case slice[mid] == target && pos == Left:
			if mid != 0 && slice[mid-1] == target {
				hi = mid - 1
			} else {
				return mid
			}
		case slice[mid] == target && pos == Right:
			last := int64(len(slice)) - 1
			if mid != last && slice[mid+1] == target {
				lo = mid + 1
			} else {
				return mid
			}
		}
	}

	if pos == Left {
		return lo
	} else {
		return hi
	}
}

func BinInsertInt64(slice []int64, num int64) []int64 {
	if len(slice) == 0 {
		return append(slice, num)
	}

	lo := 0
	hi := len(slice) - 1
	for lo < hi {
		mid := (lo + hi) / 2
		if slice[mid] < num {
			lo = mid + 1
		}
		if slice[mid] > num {
			hi = mid - 1
		}
		if num == slice[mid] {
			lo = mid
			break
		}
	}

	if len(slice)-1 == lo {
		return append(slice, num)
	}

	index := lo
	if num > slice[lo] {
		index = lo + 1
	}
	slice = append(slice[:index+1], slice[index:]...)
	slice[index] = num
	return slice
}
