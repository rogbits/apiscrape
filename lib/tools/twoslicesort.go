package tools

type TwoSliceSort struct {
	Primary   interface{}
	Secondary interface{}
}

func NewTwoSliceSort(primary, secondary interface{}) *TwoSliceSort {
	tss := new(TwoSliceSort)
	tss.Primary = primary
	tss.Secondary = secondary
	return tss
}

func (tss TwoSliceSort) Len() int {
	switch tss.Primary.(type) {
	case []int64:
		s := tss.Primary.([]int64)
		return len(s)
	case []string:
		s := tss.Primary.([]string)
		return len(s)
	default:
		panic("not implemented")
	}
	return 0
}

func (tss TwoSliceSort) Swap(i, j int) {
	// primary
	switch tss.Primary.(type) {
	case []int64:
		s := tss.Primary.([]int64)
		s[i], s[j] = s[j], s[i]
	case []string:
		s := tss.Primary.([]string)
		s[i], s[j] = s[j], s[i]
	default:
		panic("not implemented")
	}

	// secondary
	switch tss.Secondary.(type) {
	case []int64:
		s := tss.Secondary.([]int64)
		s[i], s[j] = s[j], s[i]
	case []string:
		s := tss.Secondary.([]string)
		s[i], s[j] = s[j], s[i]
	default:
		panic("not implemented")
	}
}

func (tss TwoSliceSort) Less(i, j int) bool {
	switch tss.Primary.(type) {
	case []int64:
		s := tss.Primary.([]int64)
		return s[i] < s[j]
	case []string:
		s := tss.Primary.([]string)
		return s[i] < s[j]
	default:
		panic("not implemented")
	}

	return false
}
