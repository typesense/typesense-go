package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testSpec = `paths:
  /things/{id}:
    get:
      operationId: getAlias
      tags: [collections]
      summary: Retrieve a thing
      description: Fetch a thing by its id
    delete:
      operationId: deleteThing
      tags: [collections]
      summary: Delete a thing
      description: |
        delete a thing

        It cannot be undone
`

const testWrapper = `package typesense

type ThingInterface interface {
	// Retrieve is hand-written.
	Retrieve() error
	Delete() error
	Both() error
}

type thing struct{}

// Retrieve is hand-written.
func (t *thing) Retrieve() error { return t.apiClient.GetAlias() }

// Old text.
//
// Deprecated: use Retrieve.
//
//nolint:revive
func (t *thing) Delete() error { return t.apiClient.DeleteThing() }

// Both calls two operations, so it is left alone.
func (t *thing) Both() error { _ = t.apiClient.GetAlias(); return t.apiClient.DeleteThing() }
`

func TestInjectWrapperDocs(t *testing.T) {
	dir := t.TempDir()
	spec := filepath.Join(dir, "openapi.yml")
	wrapper := filepath.Join(dir, "thing.go")
	if err := os.WriteFile(spec, []byte(testSpec), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(wrapper, []byte(testWrapper), 0o600); err != nil {
		t.Fatal(err)
	}

	injectWrapperDocs(spec, dir)
	first, err := os.ReadFile(wrapper)
	if err != nil {
		t.Fatal(err)
	}
	got := string(first)

	// Each block ends at the declaration, so no trailing "//" is left behind.
	// Interface comments are rendered unindented; format.Source adds the tab.
	// The description's first line repeats the summary, so only the rest stays.
	// getAlias's See: link comes from opToDocs, deleteThing's from its tag.
	// Both calls two operations, so neither of its docs is touched.
	for _, want := range []string{
		"\t// Retrieve a thing.\n\t//\n\t// Fetch a thing by its id.\n\t//\n\t// HTTP: GET /things/{id}\n\t//\n\t// See: " + docsBaseURL + "collection-alias.html\n\tRetrieve() error",
		"\t// Delete a thing.\n\t//\n\t// It cannot be undone.\n\t//\n\t// HTTP: DELETE /things/{id}\n\t//\n\t// See: " + docsBaseURL + "collections.html\n\tDelete() error\n\tBoth() error",
		"\n// Retrieve a thing.\n//\n// Fetch a thing by its id.\n//\n// HTTP: GET /things/{id}\n//\n// See: " + docsBaseURL + "collection-alias.html\nfunc (t *thing) Retrieve",
		"\n// Delete a thing.\n//\n// It cannot be undone.\n//\n// HTTP: DELETE /things/{id}\n//\n// See: " + docsBaseURL + "collections.html\n//\n// Deprecated: use Retrieve.\n//\n//nolint:revive\nfunc (t *thing) Delete",
		"\n// Both calls two operations, so it is left alone.\nfunc (t *thing) Both",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("output is missing:\n%s\n\ngot:\n%s", want, got)
		}
	}
	for _, unwanted := range []string{"hand-written", "Old text", "// # ", "delete a thing"} {
		if strings.Contains(got, unwanted) {
			t.Errorf("output contains %q:\n%s", unwanted, got)
		}
	}

	injectWrapperDocs(spec, dir)
	second, err := os.ReadFile(wrapper)
	if err != nil {
		t.Fatal(err)
	}
	if string(second) != got {
		t.Errorf("second run changed the file:\n%s", second)
	}
}
