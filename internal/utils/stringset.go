package utils

import (
	"fmt"
)

type StringSet struct {
	items    []string
	itemsMap map[string]bool
}

func (s *StringSet) Add(item string) {
	if _, exists := s.itemsMap[item]; !exists {
		s.items = append(s.items, item)
		s.itemsMap[item] = true
	}
}

func (s *StringSet) Contains(item string) bool {
	_, exists := s.itemsMap[item]
	return exists
}

func (s StringSet) String() string {
	return fmt.Sprintln(s.items)
}

func (s StringSet) Len() int {
	return len(s.items)
}

func (s StringSet) ToArray() []string {
	return s.items
}

func NewStringSet() *StringSet {
	return &StringSet{
		items:    make([]string, 0),
		itemsMap: make(map[string]bool),
	}
}
