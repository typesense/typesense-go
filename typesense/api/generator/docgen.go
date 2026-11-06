package main

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"
)

// injectWrapperDocs reads the OpenAPI spec, walks each wrapper .go file in
// wrapperDir, and rewrites doc comments on methods whose body issues a single
// direct apiClient.<Op>(...) call. The comment text comes from the spec's
// summary and description for the matched operation. Interface methods that
// match an annotated impl method get the same comment. An existing doc on a
// matched method is replaced; only its Deprecated: paragraph and directive
// lines are kept. A doc with a //docgen:keep line is left as written, for when
// the spec text is wrong or doesn't describe what the wrapper does. Like the
// rest of the generator, it exits on any error.
func injectWrapperDocs(specPath, wrapperDir string) {
	log.Println("Injecting wrapper doc comments from OpenAPI spec")

	ops, err := loadOpIndex(specPath)
	if err != nil {
		log.Fatalf("docgen: failed to load spec: %s", err)
	}

	funcToOp := make(map[string]*opInfo, len(ops)*4)
	for _, op := range ops {
		for _, name := range goFuncNamesForOp(op.OperationID) {
			funcToOp[name] = op
		}
	}

	entries, err := os.ReadDir(wrapperDir)
	if err != nil {
		log.Fatalf("docgen: cannot read wrapper dir: %s", err)
	}

	touched := 0
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		path := filepath.Join(wrapperDir, e.Name())
		changed, err := annotateFile(path, funcToOp)
		if err != nil {
			log.Fatalf("docgen: %s: %s", e.Name(), err)
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

// tagToDocs maps an operation's first spec tag to its docs page.
var tagToDocs = map[string]string{
	"collections":      "collections.html",
	"documents":        "documents.html",
	"keys":             "api-keys.html",
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

// opToDocs overrides tagToDocs for operations documented on another page.
var opToDocs = map[string]string{
	"getAliases":       "collection-alias.html",
	"getAlias":         "collection-alias.html",
	"upsertAlias":      "collection-alias.html",
	"deleteAlias":      "collection-alias.html",
	"multiSearch":      "federated-multi-search.html",
	"retrieveAPIStats": "cluster-operations.html#api-stats",
	"retrieveMetrics":  "cluster-operations.html#cluster-metrics",
}

// loadOpIndex reads the spec at specPath and indexes its operations by operationId.
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

	paths, ok := spec["paths"].(yml)
	if !ok {
		return nil, fmt.Errorf("%s has no paths", specPath)
	}
	ops := map[string]*opInfo{}
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

// goFuncNamesForOp returns the client method names oapi-codegen emits for opID.
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

// site is a method or interface method declaration to document with op.
type site struct {
	pos token.Pos
	doc *ast.CommentGroup
	op  *opInfo
}

type docPatch struct {
	startLine int // 1-indexed, inclusive
	endLine   int // 1-indexed, exclusive (the line of the decl itself)
	lines     []string
}

// annotateFile parses path, finds methods that resolve to exactly one OpenAPI
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
	sites = append(sites, collectInterfaceSites(file, methodOp)...)
	sites = slices.DeleteFunc(sites, func(s site) bool { return keepsDoc(s.doc) })
	if len(sites) == 0 {
		return false, nil
	}

	out, err := format.Source(applyPatches(src, buildPatches(fset, sites)))
	if err != nil {
		return false, fmt.Errorf("gofmt: %w", err)
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
// unique OpenAPI operation, along with a map from method name to its op. Names
// that resolve to different ops across receivers map to nil, so the
// interface-side pass skips them.
func collectMethodSites(file *ast.File, funcToOp map[string]*opInfo) ([]site, map[string]*opInfo) {
	sites := make([]site, 0, len(file.Decls))
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
		sites = append(sites, site{pos: fn.Pos(), doc: fn.Doc, op: op})
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
func collectInterfaceSites(file *ast.File, methodOp map[string]*opInfo) []site {
	sites := make([]site, 0, len(methodOp))
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
				sites = append(sites, site{pos: field.Pos(), doc: field.Doc, op: op})
			}
		}
	}
	return sites
}

// buildPatches converts sites into docPatches, sorted by descending start line
// so applying them keeps earlier line numbers valid.
func buildPatches(fset *token.FileSet, sites []site) []docPatch {
	patches := make([]docPatch, 0, len(sites))
	for _, s := range sites {
		declLine := fset.Position(s.pos).Line
		startLine := declLine
		if s.doc != nil {
			startLine = fset.Position(s.doc.Pos()).Line
		}
		patches = append(patches, docPatch{
			startLine: startLine,
			endLine:   declLine,
			lines:     renderDoc(s.op, s.doc),
		})
	}
	sort.Slice(patches, func(i, j int) bool { return patches[i].startLine > patches[j].startLine })
	return patches
}

// applyPatches rewrites src by replacing the [startLine, endLine) range of each
// patch with its rendered lines. Patches must already be sorted in descending
// startLine order.
func applyPatches(src []byte, patches []docPatch) []byte {
	lines := strings.Split(string(src), "\n")
	for _, p := range patches {
		lines = slices.Concat(lines[:p.startLine-1], p.lines, lines[p.endLine-1:])
	}
	return []byte(strings.Join(lines, "\n"))
}

// uniqueAPIOp returns the operation that every apiClient.<Name>(...) call in
// body maps to, or nil when body calls no known operation or more than one.
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

// directiveRe matches the lines go/ast treats as directives, such as
// //go:generate or //nolint:revive. gofmt rewrites any other //text line to
// "// text", so carrying those over as directives would change the next run.
var directiveRe = regexp.MustCompile(`^//(line |extern |export |[a-z0-9]+:[a-z0-9])`)

// keepDirective marks a hand-written doc that docgen must leave alone.
const keepDirective = "//docgen:keep"

func keepsDoc(doc *ast.CommentGroup) bool {
	return doc != nil && slices.ContainsFunc(doc.List, func(c *ast.Comment) bool { return c.Text == keepDirective })
}

// keptLines returns the parts of an existing doc that a rewrite must carry
// over: its Deprecated: paragraph and directive lines.
func keptLines(doc *ast.CommentGroup) (deprecated, directives []string) {
	if doc == nil {
		return nil, nil
	}
	inDeprecated := false
	for _, c := range doc.List {
		t := c.Text
		switch {
		case directiveRe.MatchString(t):
			directives = append(directives, t)
		case strings.HasPrefix(t, "// Deprecated:"), inDeprecated && t != "//":
			inDeprecated = true
			deprecated = append(deprecated, t)
		default:
			inDeprecated = false
		}
	}
	return deprecated, directives
}

// descriptionLines splits desc into lines and ends each paragraph with
// punctuation, or Go renders a one-line paragraph as a heading. Markdown lists
// and code fences pass through as plain text; no spec description is
// multi-line today.
func descriptionLines(desc string) []string {
	lines := strings.Split(desc, "\n")
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], " \t")
		if i+1 == len(lines) || strings.TrimSpace(lines[i+1]) == "" {
			lines[i] = ensureTrailingPeriod(lines[i])
		}
	}
	return lines
}

// renderDoc returns the doc comment lines for op, keeping the Deprecated:
// paragraph and directives of the existing doc. Lines are unindented;
// format.Source indents them to match the declaration.
func renderDoc(op *opInfo, doc *ast.CommentGroup) []string {
	var lines []string

	summary := firstLine(op.Summary)
	if summary == "" {
		summary = firstLine(op.Description)
	}
	if summary != "" {
		lines = append(lines, "// "+ensureTrailingPeriod(upperFirst(summary)))
	}

	desc := trimSummary(op.Description, summary)
	if desc != "" {
		lines = append(lines, "//")
		for _, l := range descriptionLines(desc) {
			if l == "" {
				lines = append(lines, "//")
			} else {
				lines = append(lines, "// "+l)
			}
		}
	}

	lines = append(lines, "//")
	lines = append(lines, fmt.Sprintf("// HTTP: %s %s", op.Method, op.Path))

	section := tagToDocs[op.Tag]
	if s, ok := opToDocs[op.OperationID]; ok {
		section = s
	}
	if section != "" {
		lines = append(lines, "//")
		lines = append(lines, "// See: "+docsBaseURL+section)
	}

	deprecated, directives := keptLines(doc)
	if len(deprecated) > 0 {
		lines = append(lines, "//")
		lines = append(lines, deprecated...)
	}
	// Directives go last, after a blank line, which is where gofmt puts them.
	if len(directives) > 0 {
		lines = append(lines, "//")
		lines = append(lines, directives...)
	}

	return lines
}

// trimSummary drops summary from the start of desc, ignoring case, when desc
// opens with it as its whole first line or as a whole sentence.
func trimSummary(desc, summary string) string {
	s := ensureTrailingPeriod(summary)
	if first := firstLine(desc); strings.EqualFold(ensureTrailingPeriod(first), s) {
		return strings.TrimSpace(strings.TrimPrefix(desc, first))
	}
	if len(desc) > len(s) && strings.EqualFold(desc[:len(s)], s) && unicode.IsSpace(rune(desc[len(s)])) {
		return strings.TrimSpace(desc[len(s):])
	}
	return desc
}

func firstLine(s string) string {
	s = strings.TrimSpace(s)
	if idx := strings.IndexByte(s, '\n'); idx >= 0 {
		s = s[:idx]
	}
	return strings.TrimSpace(s)
}

// ensureTrailingPeriod appends a period to s unless it is empty or already
// ends in '.', '?', '!' or ':', ignoring closing parens.
func ensureTrailingPeriod(s string) string {
	if s == "" {
		return s
	}
	last := strings.TrimRight(s, ")")
	if strings.HasSuffix(last, ".") || strings.HasSuffix(last, "?") || strings.HasSuffix(last, "!") || strings.HasSuffix(last, ":") {
		return s
	}
	return s + "."
}
