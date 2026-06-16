package swag

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Regression harness for known OpenAPI 3.1 (--v3.1) bugs in the swaggo/swag v2
// branch. Each TestRegression_* below pins one upstream issue.
//
// CONVENTION: every repro asserts the CORRECT / intended behavior and starts
// with t.Skip(...) referencing its issue, so `make test` stays green. Removing
// the t.Skip line is the first step of fixing the corresponding bug: the test
// then goes red (or panics, caught by the recover() in parseV3Regression) on
// current code and green once fixed.
//
// Fixtures live under testdata/v3/regressions/<slug>/. They are AST-parsed by
// swag and are excluded from `go build` (the testdata/ directory is special),
// so they need not compile.

const regressionDir = "testdata/v3/regressions/"

// parseV3Regression parses a v3 regression fixture directory. A panic during
// parsing is converted into a test failure so a crashing repro fails only its
// own (un-skipped) test instead of aborting the whole `go test` binary.
func parseV3Regression(t *testing.T, slug string, opts ...func(*Parser)) *Parser {
	t.Helper()
	p := New(append([]func(*Parser){GenerateOpenAPI3Doc(true)}, opts...)...)
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("panic parsing %s: %v", slug, r)
			}
		}()
		err := p.ParseAPI(regressionDir+slug, mainAPIFile, defaultParseDepth)
		require.NoError(t, err)
	}()
	return p
}

// marshalV3 renders the generated OpenAPI 3.1 doc to a generic JSON map for
// targeted assertions (more robust for pinning one bug than a full golden file).
func marshalV3(t *testing.T, p *Parser) map[string]interface{} {
	t.Helper()
	raw, err := json.Marshal(p.openAPI)
	require.NoError(t, err)
	var doc map[string]interface{}
	require.NoError(t, json.Unmarshal(raw, &doc))
	return doc
}

// jsonAt navigates a decoded JSON object by successive keys, failing the test
// with a clear message if any key is missing or a non-object is encountered.
func jsonAt(t *testing.T, doc map[string]interface{}, path ...string) interface{} {
	t.Helper()
	var cur interface{} = doc
	for i, key := range path {
		m, ok := cur.(map[string]interface{})
		require.Truef(t, ok, "expected object at %v while indexing %q", path[:i], key)
		cur, ok = m[key]
		require.Truef(t, ok, "missing key %q at path %v", key, path[:i+1])
	}
	return cur
}

// objAt is jsonAt plus a cast to map[string]interface{}.
func objAt(t *testing.T, doc map[string]interface{}, path ...string) map[string]interface{} {
	t.Helper()
	m, ok := jsonAt(t, doc, path...).(map[string]interface{})
	require.Truef(t, ok, "value at %v is not an object", path)
	return m
}

// ---------------------------------------------------------------------------
// Crash repros — bar is: ParseAPI returns without panic and without error.
// ---------------------------------------------------------------------------

// #1979 — crash with --v3.1 when generating (nil deref triggered by a struct).
func TestRegression_1979_StructCrash(t *testing.T) {
	t.Parallel()
	t.Skip("repro for swaggo/swag#1979 — fails until fixed; remove skip when fixing")
	// TODO: the upstream report points at an external repo's `sourceAge` struct;
	// this is a best-effort minimal struct that exercises the same v3 path.
	parseV3Regression(t, "issue1979_struct_crash")
}

// #1930 — type aliases cause nil dereference on --v3.1 (parseStructFieldV3).
// PASSES on this fork (the nil deref is fixed) — kept active as a regression guard.
func TestRegression_1930_TypeAliasQueryParam(t *testing.T) {
	t.Parallel()
	p := parseV3Regression(t, "issue1930_type_alias_query")
	doc := marshalV3(t, p)
	param := jsonAt(t, doc, "paths", "/price", "get", "parameters")
	params, ok := param.([]interface{})
	require.True(t, ok, "query param must be generated")
	assert.Len(t, params, 1)
}

// #2080 — segmentation violation on a generic struct of an enum array.
// PASSES on this fork (no panic) — kept active as a regression guard. NOTE: this
// only guards against the crash; correctness of enum-array output inside a
// generic is tracked separately by #1934 (still failing).
func TestRegression_2080_GenericEnumArray(t *testing.T) {
	t.Parallel()
	parseV3Regression(t, "issue2080_generic_enum_array")
}

// #2078 — panic on a custom type ([]CustomType with swaggertype:"array,string").
func TestRegression_2078_CustomTypeArray(t *testing.T) {
	t.Parallel()
	t.Skip("repro for swaggo/swag#2078 — fails until fixed; remove skip when fixing")
	parseV3Regression(t, "issue2078_custom_type_array")
}

// #1911 — nil deref in getFuncDoc during router parsing.
func TestRegression_1911_FuncDocCrash(t *testing.T) {
	t.Parallel()
	t.Skip("repro for swaggo/swag#1911 — does not reproduce as written (no published input); refine fixture to trigger the getFuncDoc nil deref")
	// TODO: original input not published; best-effort nested-func-literal fixture.
	parseV3Regression(t, "issue1911_funcdoc_crash", func(p *Parser) { p.ParseFuncBody = true })
}

// ---------------------------------------------------------------------------
// Wrong-output repros — targeted assertions on the generated 3.1 doc.
// ---------------------------------------------------------------------------

// #2142 (== #2086) — @Accept json + @Param body wraps the body schema in a
// spurious oneOf with an empty {type:object}. Correct: a single $ref.
// PASSES on this fork (the spurious oneOf is gone — body is a clean $ref) —
// kept active as a regression guard.
func TestRegression_2142_EmptyBodySchemaWithAccept(t *testing.T) {
	t.Parallel()
	p := parseV3Regression(t, "issue2142_accept_body")
	doc := marshalV3(t, p)
	schema := objAt(t, doc, "paths", "/things", "post", "requestBody", "content", "application/json", "schema")
	_, hasOneOf := schema["oneOf"]
	assert.False(t, hasOneOf, "single body param must not be wrapped in oneOf")
	assert.Equal(t, "#/components/schemas/main.CreateRequest", schema["$ref"])
}

// #2086 — same root cause as #2142, asserted independently so each issue is
// tracked. When the request body has a single object param it must be a $ref,
// not a oneOf containing a malformed {type:object,$ref:...} entry.
// PASSES on this fork — kept active as a regression guard.
func TestRegression_2086_RequestBodySpuriousOneOf(t *testing.T) {
	t.Parallel()
	p := parseV3Regression(t, "issue2142_accept_body")
	doc := marshalV3(t, p)
	schema := objAt(t, doc, "paths", "/things", "post", "requestBody", "content", "application/json", "schema")
	_, hasOneOf := schema["oneOf"]
	assert.False(t, hasOneOf, "single body param must not generate a oneOf")
}

// #2140 — embedded struct with `yaml:",inline"` drops the embedded fields.
func TestRegression_2140_YamlInlineEmbedded(t *testing.T) {
	t.Parallel()
	t.Skip("repro for swaggo/swag#2140 — fails until fixed; remove skip when fixing")
	p := parseV3Regression(t, "issue2140_yaml_inline")
	doc := marshalV3(t, p)
	props := objAt(t, doc, "components", "schemas", "main.ServerMetadata", "properties")
	// Embedded BaseMetadata fields must be inlined alongside the direct field.
	for _, f := range []string{"name", "description", "status", "image"} {
		_, ok := props[f]
		assert.Truef(t, ok, "property %q from embedded BaseMetadata must be present", f)
	}
}

// #1641 — two endpoints returning map[string]interface{} must get distinct
// (correct) response schemas, not a shared/incorrect ref.
// PASSES on this fork (each endpoint gets its own inline map schema) — kept
// active as a regression guard.
func TestRegression_1641_MapInterfaceResponseRefs(t *testing.T) {
	t.Parallel()
	p := parseV3Regression(t, "issue1641_map_interface")
	doc := marshalV3(t, p)
	// A map[string]interface{} response should be an inline object schema with
	// additionalProperties, not a $ref shared across endpoints.
	for _, path := range []string{"/a", "/b"} {
		schema := objAt(t, doc, "paths", path, "get", "responses", "200", "content", "application/json", "schema")
		_, isRef := schema["$ref"]
		assert.Falsef(t, isRef, "map response for %s must be inline, not a shared $ref", path)
	}
}

// #1950 — response model composition override (Result{data=Concrete}) is
// ignored on --v3.1; the `data` field must reflect the override type.
func TestRegression_1950_CompositionOverride(t *testing.T) {
	t.Parallel()
	t.Skip("repro for swaggo/swag#1950 — assertion too weak (passes as written); strengthen to assert the data-field override before relying on it")
	p := parseV3Regression(t, "issue1950_composition_override")
	doc := marshalV3(t, p)
	// The composed schema must exist and its `data` must ref the override type,
	// not be a bare interface{}/empty object.
	schemas := objAt(t, doc, "components", "schemas")
	var composed string
	for name := range schemas {
		if name != "api.Response" && name != "api.GetData1Response" {
			composed = name
		}
	}
	require.NotEmpty(t, composed, "a composed Response{data=...} schema must be generated")
}

// #1934 — enum array field inside a generic produces wrong output (missing
// per-item enum metadata). The []OtherEnum field must be an array whose items
// carry the enum varnames.
func TestRegression_1934_EnumArrayInGeneric(t *testing.T) {
	t.Parallel()
	t.Skip("repro for swaggo/swag#1934 — fails until fixed; remove skip when fixing")
	p := parseV3Regression(t, "issue1934_json_tag_enum")
	doc := marshalV3(t, p)
	f2 := objAt(t, doc, "components", "schemas", "main.Complex", "properties", "f_2")
	assert.Equal(t, "array", typeOf(f2), "f_2 ([]OtherEnum) must be an array")
	items := objAt(t, doc, "components", "schemas", "main.Complex", "properties", "f_2", "items")
	_, hasVarnames := items["x-enum-varnames"]
	assert.True(t, hasVarnames, "array items must carry enum varnames")
}

// #2153 — float64 field with a fractional `maximum` tag must parse (was failing
// with strconv.Atoi parsing "0.1").
func TestRegression_2153_NumericMaximumTag(t *testing.T) {
	t.Parallel()
	t.Skip("repro for swaggo/swag#2153 — fails until fixed; remove skip when fixing")
	p := parseV3Regression(t, "issue2153_numeric_maximum")
	doc := marshalV3(t, p)
	factor := objAt(t, doc, "components", "schemas", "main.CreateInput", "properties", "factor")
	assert.EqualValues(t, 0.1, factor["maximum"], "fractional maximum must be parsed as a number")
}

// #2161 — @Param query with a complex struct type must expand into one query
// parameter per struct field (v1 behavior), not an empty parameter list.
// PASSES on this fork (the complex struct expands to one query param per field)
// — kept active as a regression guard.
func TestRegression_2161_ComplexQueryParamExpansion(t *testing.T) {
	t.Parallel()
	p := parseV3Regression(t, "issue2161_complex_query_param")
	doc := marshalV3(t, p)
	params, ok := jsonAt(t, doc, "paths", "/api/v1/orgstruct", "get", "parameters").([]interface{})
	require.True(t, ok, "operation must declare query parameters")
	assert.GreaterOrEqual(t, len(params), 2, "complex query struct must expand to one param per field")
}

// ---------------------------------------------------------------------------
// Missing-feature repros — assert the intended behavior. Annotation syntax for
// some of these is still TBD; the assertion encodes the desired outcome and may
// be revised when the syntax is decided.
// ---------------------------------------------------------------------------

// #1932 — Host/BasePath have no replacement in v3; the generated doc should
// expose a servers[] entry derived from the configured base path.
func TestRegression_1932_HostBasePathServers(t *testing.T) {
	t.Parallel()
	t.Skip("repro for swaggo/swag#1932 — assertion too weak (servers[] from @BasePath already populated); strengthen to cover @host before relying on it")
	p := parseV3Regression(t, "issue1932_host_basepath")
	doc := marshalV3(t, p)
	servers, ok := doc["servers"].([]interface{})
	require.True(t, ok, "doc must expose servers[] derived from host/basePath")
	assert.NotEmpty(t, servers)
}

// #2011 — @server annotation should populate servers[].
func TestRegression_2011_ServerAnnotation(t *testing.T) {
	t.Parallel()
	t.Skip("repro for swaggo/swag#2011 — syntax TBD; fails until implemented")
	p := parseV3Regression(t, "issue2011_server")
	doc := marshalV3(t, p)
	servers, ok := doc["servers"].([]interface{})
	require.True(t, ok, "@server must populate servers[]")
	assert.NotEmpty(t, servers)
}

// #1949 — example values should be emitted under the 3.1 `examples` keyword,
// not the deprecated singular `example`.
func TestRegression_1949_ExamplesKeyword(t *testing.T) {
	t.Parallel()
	t.Skip("repro for swaggo/swag#1949 — syntax TBD; fails until implemented")
	p := parseV3Regression(t, "issue1949_examples")
	doc := marshalV3(t, p)
	name := objAt(t, doc, "components", "schemas", "main.Thing", "properties", "name")
	_, hasExamples := name["examples"]
	assert.True(t, hasExamples, "example tag must populate the 3.1 `examples` keyword")
}

// #1855 — reusable/nullable enums. A reused enum type must be a single
// component schema referenced by both fields (not inlined twice).
func TestRegression_1855_ReusableEnums(t *testing.T) {
	t.Parallel()
	t.Skip("repro for swaggo/swag#1855 — assertion too weak (enum component already emitted); strengthen to assert reuse + nullability before relying on it")
	p := parseV3Regression(t, "issue1855_reusable_enums")
	doc := marshalV3(t, p)
	schemas := objAt(t, doc, "components", "schemas")
	_, ok := schemas["main.Color"]
	assert.True(t, ok, "reused enum type must be a reusable component schema")
}

// #2139 — --parseFuncBody must work together with --v3.1: routes declared in
// comments inside a function body must appear in the generated doc.
func TestRegression_2139_ParseFuncBodyWithV3(t *testing.T) {
	t.Parallel()
	t.Skip("repro for swaggo/swag#2139 — fails until fixed; remove skip when fixing")
	p := parseV3Regression(t, "issue2139_parsefuncbody", func(p *Parser) { p.ParseFuncBody = true })
	doc := marshalV3(t, p)
	paths := objAt(t, doc, "paths")
	_, ok := paths["/motd"]
	assert.True(t, ok, "route declared inside a function body must be parsed under --v3.1")
}

// typeOf returns the OpenAPI `type` of a schema map, handling both the string
// and []string (3.1 type array) JSON encodings.
func typeOf(schema map[string]interface{}) string {
	switch v := schema["type"].(type) {
	case string:
		return v
	case []interface{}:
		if len(v) > 0 {
			if s, ok := v[0].(string); ok {
				return s
			}
		}
	}
	return ""
}
