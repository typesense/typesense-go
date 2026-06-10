package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

// reads the OpenAPI spec, walks each wrapper .go file in
// wrapperDir, and rewrites doc comments on methods whose body issues a single
// direct apiClient.<Op>(...) call. The comment text comes verbatim from the
// spec's summary/description for the matched operation. Interface methods that
// match an annotated impl method get the same comment.
func injectWrapperDocs(specPath, wrapperDir string) {
	log.Println("Injecting wrapper doc comments from OpenAPI spec")

	ops, err := loadOpIndex(specPath)
	if err != nil {
		log.Printf("docgen: failed to load spec: %s", err)
		return
	}

	// Build map of Go func names (as oapi-codegen emits) → operation info.
	funcToOp := make(map[string]*opInfo, len(ops)*4)
	for _, op := range ops {
		for _, name := range goFuncNamesForOp(op.OperationID) {
			funcToOp[name] = op
		}
	}

	entries, err := os.ReadDir(wrapperDir)
	if err != nil {
		log.Printf("docgen: cannot read wrapper dir: %s", err)
		return
	}

	touched := 0
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		if strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		path := filepath.Join(wrapperDir, e.Name())
		changed, err := annotateFile(path, funcToOp)
		if err != nil {
			log.Printf("docgen: %s: %s", e.Name(), err)
			continue
		}
		if changed {
			touched++
		}
	}
	log.Printf("docgen: updated %d wrapper file(s)", touched)
}

type opInfo struct {
	OperationID string
	Summary     string
	Description string
	Method      string
	Path        string
	Tag         string
}

const docsBaseURL = "https://typesense.org/docs/latest/api/"

var tagToDocs = map[string]string{
	"collections":      "collections.html",
	"documents":        "documents.html",
	"keys":             "api-keys.html",
	"aliases":          "collection-alias.html",
	"synonyms":         "synonyms.html",
	"curation_sets":    "curation.html",
	"stopwords":        "stopwords.html",
	"presets":          "search.html#presets",
	"analytics":        "analytics-query-suggestions.html",
	"conversations":    "conversational-search-rag.html",
	"stemming":         "stemming.html",
	"nl_search_models": "natural-language-search.html",
	"debug":            "cluster-operations.html#debug",
	"health":           "cluster-operations.html#health",
	"operations":       "cluster-operations.html",
}

func loadOpIndex(specPath string) (map[string]*opInfo, error) {
	f, err := os.Open(specPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var spec yml
	if err := yaml.NewDecoder(f).Decode(&spec); err != nil {
		return nil, err
	}

	ops := map[string]*opInfo{}
	paths, _ := spec["paths"].(yml)
	for pathKey, pathVal := range paths {
		pathItem, ok := pathVal.(yml)
		if !ok {
			continue
		}
		for method, mVal := range pathItem {
			op, ok := mVal.(yml)
			if !ok {
				continue
			}
			opID, _ := op["operationId"].(string)
			if opID == "" {
				continue
			}
			info := &opInfo{
				OperationID: opID,
				Method:      strings.ToUpper(method),
				Path:        pathKey,
			}
			if s, ok := op["summary"].(string); ok {
				info.Summary = strings.TrimSpace(s)
			}
			if d, ok := op["description"].(string); ok {
				info.Description = strings.TrimSpace(d)
			}
			if tags, ok := op["tags"].([]interface{}); ok && len(tags) > 0 {
				if t, ok := tags[0].(string); ok {
					info.Tag = t
				}
			}
			ops[opID] = info
		}
	}
	return ops, nil
}

func goFuncNamesForOp(opID string) []string {
	base := upperFirst(opID)
	return []string{
		base,
		base + "WithBody",
		base + "WithResponse",
		base + "WithBodyWithResponse",
	}
}

func upperFirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

type funcSite struct {
	decl *ast.FuncDecl
	op   *opInfo
}

type ifaceSite struct {
	field *ast.Field
	op    *opInfo
}

type docPatch struct {
	startLine int // 1-indexed, inclusive
	endLine   int // 1-indexed, exclusive (the line of the decl itself)
	lines     []string
}

// parses path, finds methods that resolve to exactly one OpenAPI
// operation, and rewrites doc comments above those methods and any matching
// interface method declarations in the same file. Returns true if the file was
// modified.
func annotateFile(path string, funcToOp map[string]*opInfo) (bool, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, src, parser.ParseComments)
	if err != nil {
		return false, err
	}

	sites, methodOp := collectMethodSites(file, funcToOp)
	ifaceSites := collectInterfaceSites(file, methodOp)

	if len(sites) == 0 && len(ifaceSites) == 0 {
		return false, nil
	}

	patches := buildPatches(fset, sites, ifaceSites)
	out, err := applyPatches(src, patches)
	if err != nil {
		return false, err
	}
	if string(out) == string(src) {
		return false, nil
	}

	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	// #nosec G703 it's not an external file
	if err := os.WriteFile(path, out, info.Mode().Perm()); err != nil {
		return false, err
	}

	return true, nil
}

// collectMethodSites returns the set of method declarations that resolve to a
// unique OpenAPI operation, along with a map from method name to its op. The
// map drops names that resolve to different ops across receivers so the
// interface-side annotation pass can skip ambiguous matches.
func collectMethodSites(file *ast.File, funcToOp map[string]*opInfo) ([]funcSite, map[string]*opInfo) {
	sites := make([]funcSite, 0, len(file.Decls))
	methodOp := map[string]*opInfo{}
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv == nil || fn.Body == nil {
			continue
		}
		op := uniqueAPIOp(fn.Body, funcToOp)
		if op == nil {
			continue
		}
		sites = append(sites, funcSite{decl: fn, op: op})
		if existing, dup := methodOp[fn.Name.Name]; dup && existing != op {
			methodOp[fn.Name.Name] = nil
		} else if !dup {
			methodOp[fn.Name.Name] = op
		}
	}
	return sites, methodOp
}

// collectInterfaceSites returns interface method fields whose name matches an
// impl method in methodOp, so they can be annotated with the same doc.
func collectInterfaceSites(file *ast.File, methodOp map[string]*opInfo) []ifaceSite {
	sites := make([]ifaceSite, 0, len(methodOp))
	for _, decl := range file.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}
		for _, spec := range gd.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			it, ok := ts.Type.(*ast.InterfaceType)
			if !ok || it.Methods == nil {
				continue
			}
			for _, field := range it.Methods.List {
				if len(field.Names) != 1 {
					continue
				}
				op := methodOp[field.Names[0].Name]
				if op == nil {
					continue
				}
				sites = append(sites, ifaceSite{field: field, op: op})
			}
		}
	}
	return sites
}

// buildPatches converts func and interface sites into ordered docPatches,
// sorted by descending start line so applying them keeps earlier line numbers
// valid.
func buildPatches(fset *token.FileSet, sites []funcSite, ifaceSites []ifaceSite) []docPatch {
	patches := make([]docPatch, 0, len(sites)+len(ifaceSites))
	for _, s := range sites {
		startLine, endLine, indent := commentSpan(fset, s.decl.Pos(), s.decl.Doc)
		patches = append(patches, docPatch{
			startLine: startLine,
			endLine:   endLine,
			lines:     renderDoc(s.decl.Name.Name, s.op, indent),
		})
	}
	for _, s := range ifaceSites {
		startLine, endLine, indent := commentSpan(fset, s.field.Pos(), s.field.Doc)
		patches = append(patches, docPatch{
			startLine: startLine,
			endLine:   endLine,
			lines:     renderDoc(s.field.Names[0].Name, s.op, indent),
		})
	}
	sort.Slice(patches, func(i, j int) bool { return patches[i].startLine > patches[j].startLine })
	return patches
}

// applyPatches rewrites src by replacing the [startLine, endLine) range of each
// patch with its rendered lines. Patches must already be sorted in descending
// startLine order.
func applyPatches(src []byte, patches []docPatch) ([]byte, error) {
	lines := strings.Split(string(src), "\n")
	for _, p := range patches {
		s := p.startLine - 1
		e := p.endLine - 1
		if s < 0 || e > len(lines) || s > e {
			return nil, fmt.Errorf("patch out of bounds: start=%d end=%d total=%d", p.startLine, p.endLine, len(lines))
		}
		lines = append(lines[:s], append(append([]string{}, p.lines...), lines[e:]...)...)
	}
	return []byte(strings.Join(lines, "\n")), nil
}

// returns the OpenAPI operation that a function body calls iff it
// contains exactly one apiClient.<Name>(...) call where <Name> maps to a known
// operation. Returns nil otherwise (zero matches, multiple matches, or
// ambiguous matches).
func uniqueAPIOp(body *ast.BlockStmt, funcToOp map[string]*opInfo) *opInfo {
	var found *opInfo
	multiple := false
	ast.Inspect(body, func(n ast.Node) bool {
		if multiple {
			return false
		}
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		// Match <x>.apiClient.<Name>(...)
		inner, ok := sel.X.(*ast.SelectorExpr)
		if !ok || inner.Sel.Name != "apiClient" {
			return true
		}
		op, ok := funcToOp[sel.Sel.Name]
		if !ok {
			return true
		}
		if found != nil && found != op {
			multiple = true
			return false
		}
		found = op
		return true
	})
	if multiple {
		return nil
	}
	return found
}

// returns the [startLine, endLine) range of lines that should be
// replaced for a declaration, and the leading whitespace indent on the decl
// line. startLine equals the line of an existing doc comment (if any) or the
// decl line itself when no doc is present. endLine is always the decl line.
func commentSpan(fset *token.FileSet, declPos token.Pos, doc *ast.CommentGroup) (int, int, string) {
	declLine := fset.Position(declPos).Line
	startLine := declLine
	if doc != nil {
		startLine = fset.Position(doc.Pos()).Line
	}
	// Read the decl line text to figure out indentation.
	indent := ""
	col := fset.Position(declPos).Column - 1
	if col > 0 {
		indent = strings.Repeat(" ", col)
	}
	return startLine, declLine, indent
}

func renderDoc(_ string, op *opInfo, indent string) []string {
	var lines []string

	summary := firstSentence(op.Summary)
	if summary == "" {
		summary = firstSentence(op.Description)
	}
	if summary != "" {
		lines = append(lines, indent+"// "+ensureTrailingPeriod(upperFirst(summary)))
	}

	// include the full description only when it adds information beyond the summary.
	if op.Description != "" && !strings.EqualFold(strings.TrimSpace(op.Description), strings.TrimSpace(op.Summary)) {
		desc := strings.TrimSpace(op.Description)
		if firstSentence(desc) != firstSentence(op.Summary) {
			lines = append(lines, indent+"//")
			for _, l := range strings.Split(desc, "\n") {
				l = strings.TrimRight(l, " \t")
				if l == "" {
					lines = append(lines, indent+"//")
				} else {
					lines = append(lines, indent+"// "+l)
				}
			}
		}
	}

	lines = append(lines, indent+"//")
	lines = append(lines, indent+fmt.Sprintf("// HTTP: %s %s", op.Method, op.Path))

	if section := tagToDocs[op.Tag]; section != "" {
		lines = append(lines, indent+"//")
		lines = append(lines, indent+"// See: "+docsBaseURL+section)
	}

	return lines
}

func firstSentence(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	// take everything up to the first newline, that's usually the lead sentence.
	if idx := strings.IndexByte(s, '\n'); idx >= 0 {
		s = s[:idx]
	}
	return strings.TrimSpace(s)
}

func ensureTrailingPeriod(s string) string {
	if s == "" {
		return s
	}
	if strings.HasSuffix(s, ".") || strings.HasSuffix(s, "?") || strings.HasSuffix(s, "!") {
		return s
	}
	return s + "."
}
