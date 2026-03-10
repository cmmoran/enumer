package main

import (
	"go/ast"
	"testing"
)

func TestParseValueSpecComments(t *testing.T) {
	vspec := &ast.ValueSpec{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "//enumer:alias=available, ready"},
				{Text: "//enumer:aliasregexp=^avail(?:able)?$"},
			},
		},
		Comment: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "// Ready //enumer:alias=go,g //enumer:aliasregexp=^(?i:ready|go)$"},
			},
		},
	}

	aliases, regexAliases, lineCommentName, err := parseValueSpecComments(vspec)
	if err != nil {
		t.Fatalf("parseValueSpecComments returned error: %v", err)
	}

	wantAliases := []string{"available", "ready", "go", "g"}
	if len(aliases) != len(wantAliases) {
		t.Fatalf("got aliases %v, want %v", aliases, wantAliases)
	}
	for i := range wantAliases {
		if aliases[i] != wantAliases[i] {
			t.Fatalf("got aliases %v, want %v", aliases, wantAliases)
		}
	}
	wantRegexAliases := []string{"^avail(?:able)?$", "^(?i:ready|go)$"}
	if len(regexAliases) != len(wantRegexAliases) {
		t.Fatalf("got regex aliases %v, want %v", regexAliases, wantRegexAliases)
	}
	for i := range wantRegexAliases {
		if regexAliases[i] != wantRegexAliases[i] {
			t.Fatalf("got regex aliases %v, want %v", regexAliases, wantRegexAliases)
		}
	}

	if lineCommentName != "Ready" {
		t.Fatalf("got lineCommentName %q, want %q", lineCommentName, "Ready")
	}
}

func TestValidateParseKeysRejectsCrossedStreams(t *testing.T) {
	values := []Value{
		{originalName: "StatusOpen", name: "StatusOpen", aliases: []string{"closed"}},
		{originalName: "StatusClosed", name: "closed"},
	}

	err := validateParseKeys(values, "Status")
	if err == nil {
		t.Fatal("validateParseKeys returned nil, want conflict")
	}
}

func TestValidateParseKeysRejectsRegexpCrossedStreams(t *testing.T) {
	values := []Value{
		{originalName: "StatusOpen", name: "open", regexAliases: []string{"^closed$"}},
		{originalName: "StatusClosed", name: "closed"},
	}

	err := validateParseKeys(values, "Status")
	if err == nil {
		t.Fatal("validateParseKeys returned nil, want regexp conflict")
	}
}
