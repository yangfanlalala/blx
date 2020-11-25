package aliyun

import (
	"bytes"
	"sort"
)

type headerSorter struct {
	Keys []string
	Vals []string
}

func newHeaderSorter(m map[string]string) *headerSorter {
	l := len(m)
	hs := &headerSorter{
		Keys: make([]string, 0, l),
		Vals: make([]string, 0, l),
	}

	for k, v := range m {
		hs.Keys = append(hs.Keys, k)
		hs.Vals = append(hs.Vals, v)
	}
	return hs
}

func (hs *headerSorter) Sort() {
	sort.Sort(hs)
}

func (hs *headerSorter) Len() int {
	return len(hs.Keys)
}

func (hs *headerSorter) Less(i, j int) bool {
	return bytes.Compare([]byte(hs.Keys[i]), []byte(hs.Keys[j])) < 0
}

func (hs *headerSorter) Swap(i, j int) {
	hs.Keys[i], hs.Keys[j] = hs.Keys[j], hs.Keys[j]
	hs.Vals[i], hs.Vals[j] = hs.Vals[j], hs.Vals[i]
}