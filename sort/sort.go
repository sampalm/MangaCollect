package sort

import (
	"log"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type By func(name1, name2 string) bool

func (by By) Sort(names []string) {
	ps := &nameSorter{
		names: names,
		by:    by,
	}
	sort.Sort(ps)
}

type nameSorter struct {
	names []string
	by    func(name1, name2 string) bool
}

func GetNumber(name string) (num float64) {
	isNum := func() bool {
		for _, element := range name {
			if unicode.IsNumber(element) {
				return true
			}
		}
		return false
	}
	if !isNum() {
		return
	}

	if !strings.Contains(name, " ") {
		return
	}

	newstr := strings.FieldsFunc(name, split)[1]
	num, err := strconv.ParseFloat(newstr, 64)
	if err != nil {
		log.Fatalln(err)
	}
	return num
}

// closure to sort by
var NameNumber = func(name1, name2 string) bool {
	return GetNumber(name1) < GetNumber(name2)
}

func (s *nameSorter) Swap(i, j int) {
	s.names[i], s.names[j] = s.names[j], s.names[i]
}
func (s *nameSorter) Len() int {
	return len(s.names)
}
func (s *nameSorter) Less(i, j int) bool {
	return s.by(s.names[i], s.names[j])
}

func split(r rune) bool {
	return r == ' ' || r == '.'
}
