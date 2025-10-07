package main

import "fmt"

type MapValue int

const (
	Male MapValue = iota
	Female
	Unknown
)

func main() {
	ck(Male, "XY")
	ck(Female, "XX")
	ck(Unknown, "XX|XY")
	ck(-127, "MapValue(-127)")
	ck(127, "MapValue(127)")
}

func ck(value MapValue, str string) {
	if fmt.Sprint(value) != str {
		panic("transform_map.go: " + str)
	}
}
