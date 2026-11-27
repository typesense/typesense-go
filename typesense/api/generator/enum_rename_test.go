package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"testing"
)

// scripts/v4-enum-rename.sh must rename each pre-v4 bare constant to one that
// types_gen.go still defines, named the way always-prefix-enum-values names it.
func TestEnumRenameTableMatchesGeneratedTypes(t *testing.T) {
	typeOf := enumConstTypes(t, "../types_gen.go")
	pairs := renamePairs(t, "../../../scripts/v4-enum-rename.sh")
	if len(pairs) == 0 {
		t.Fatal("no rename pairs found in the script")
	}
	for _, p := range pairs {
		typ, ok := typeOf[p[1]]
		if !ok {
			t.Errorf("%s -> %s: types_gen.go has no constant %s", p[0], p[1], p[1])
			continue
		}
		if p[1] != typ+p[0] {
			t.Errorf("%s -> %s: want %s%s", p[0], p[1], typ, p[0])
		}
	}
}

func enumConstTypes(t *testing.T, path string) map[string]string {
	t.Helper()
	file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	types := map[string]string{}
	for _, decl := range file.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.CONST {
			continue
		}
		for _, spec := range gen.Specs {
			value := spec.(*ast.ValueSpec)
			if ident, ok := value.Type.(*ast.Ident); ok {
				for _, name := range value.Names {
					types[name.Name] = ident.Name
				}
			}
		}
	}
	return types
}

func renamePairs(t *testing.T, path string) [][2]string {
	t.Helper()
	src, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	_, table, found := strings.Cut(string(src), "<<'EOF'\n")
	if !found {
		t.Fatal("no heredoc table in the script")
	}
	table, _, _ = strings.Cut(table, "\nEOF")
	lines := strings.Split(table, "\n")
	pairs := make([][2]string, 0, len(lines))
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) != 2 {
			t.Fatalf("malformed table line %q", line)
		}
		pairs = append(pairs, [2]string{fields[0], fields[1]})
	}
	return pairs
}
