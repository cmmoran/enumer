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
			},
		},
		Comment: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "// Ready //enumer:alias=go,g"},
			},
		},
	}

	aliases, lineCommentName, err := parseValueSpecComments(vspec)
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
