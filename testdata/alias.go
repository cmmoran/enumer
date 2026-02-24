package main

import "fmt"

type Alias int

const (
	AliasOpen   Alias = iota //enumer:alias=open,o
	AliasClosed              //enumer:alias=closed,c
)

func main() {
	tests := map[string]Alias{
		"AliasOpen": AliasOpen,
		"aliasopen": AliasOpen,
		"open":      AliasOpen,
		"O":         AliasOpen,
		"closed":    AliasClosed,
		"C":         AliasClosed,
	}

	for input, want := range tests {
		got, err := AliasString(input)
		if err != nil {
			panic(fmt.Sprintf("AliasString(%q) returned error: %v", input, err))
		}
		if got != want {
			panic(fmt.Sprintf("AliasString(%q) = %v, want %v", input, got, want))
		}
	}
}
