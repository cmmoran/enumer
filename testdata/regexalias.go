package main

import "fmt"

type RegexAlias int

const (
	RegexAliasOpen   RegexAlias = iota //enumer:alias=open,o //enumer:aliasregexp=^(?i:open(?:ed)?)$
	RegexAliasClosed                   //enumer:alias=closed,c //enumer:aliasregexp=^(?i:shut|close[d]?)$
)

func main() {
	tests := map[string]RegexAlias{
		"RegexAliasOpen":   RegexAliasOpen,
		"open":             RegexAliasOpen,
		"O":                RegexAliasOpen,
		"Opened":           RegexAliasOpen,
		"RegexAliasClosed": RegexAliasClosed,
		"closed":           RegexAliasClosed,
		"C":                RegexAliasClosed,
		"Shut":             RegexAliasClosed,
	}

	for input, want := range tests {
		got, err := RegexAliasString(input)
		if err != nil {
			panic(fmt.Sprintf("RegexAliasString(%q) returned error: %v", input, err))
		}
		if got != want {
			panic(fmt.Sprintf("RegexAliasString(%q) = %v, want %v", input, got, want))
		}
	}
}
