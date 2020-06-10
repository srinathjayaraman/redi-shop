package order

import (
	"fmt"
	"strconv"
	"strings"
)

func itemStringToJSONString(items string) string {
	if items == "[]" {
		return "[]"
	}

	m := itemStringToMap(items)
	res := ""
	for k := range m {
		res = fmt.Sprintf("%s\"%s\",", res, k)
	}

	return fmt.Sprintf("[%s]", res[:len(res)-1])
}

func itemStringToMap(items string) map[string]int {
	m := map[string]int{}

	if items == "[]" {
		return m
	}

	itemSplit := strings.Split(items[1:len(items)-1], ",")
	for i := range itemSplit {
		item := strings.Split(itemSplit[i], "->")
		val, err := strconv.Atoi(item[1])
		if err != nil {
			panic(fmt.Sprintf("invalid string representation of item, %s", itemSplit[i]))
		}
		m[item[0]] = val
	}

	return m
}

func mapToItemString(items map[string]int) string {
	s := ""

	for k, v := range items {
		s = fmt.Sprintf("%s%s->%d,", s, k, v)
	}

	if s == "" {
		return "[]"
	}

	return fmt.Sprintf("[%s]", s[:len(s)-1])
}
