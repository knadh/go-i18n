package i18n

import (
	"fmt"
	"strings"
	"testing"
)

func TestGeneral(t *testing.T) {
	j := `
{
	"_.code": "en",
	"_.name": "English",

	"pageTitle": "Welcome to the page",
	"foo": "Foo",
	"page": "Single page|Many pages",
	"pageVars": "The page is named {name} and has {count} items",
	"nested": "This is {nested}",
	"priceMsg": "The price is ${price} with {discount}% discount",
	"statusMsg": "The system is {status} and auto-save is {autosave}",
	"mixedMsg": "User {user} has {count} items worth ${total} (active: {active})",
	"app.name": "MyApp",
	"nestedParams": "Welcome to {appName}"
}
`

	i, err := New([]byte(j))
	if err != nil {
		t.Fatal(err)
	}

	assert(t, i.Code(), "en")
	assert(t, i.Name(), "English")
	assert(t, i.T("pageTitle"), "Welcome to the page")
	assert(t, i.T("page"), "Single page")
	assert(t, i.S("page"), "Single page")
	assert(t, i.P("page"), "Many pages")
	assert(t, i.Tc("page", 0), "Single page")
	assert(t, i.Tc("page", 1), "Single page")
	assert(t, i.Tc("page", 2), "Many pages")
	assert(t, i.S("foo"), "Foo")
	assert(t, i.P("foo"), "Foo")
	assert(t, i.T("pageVars"), "The page is named {name} and has {count} items")
	assert(t, i.Ts("pageVars"), "The page is named {name} and has {count} items")
	assert(t, i.Ts("pageVars", "name", "Foo"), "The page is named Foo and has {count} items")
	assert(t, i.Ts("pageVars", "name", "Foo", "count", "1234"), "The page is named Foo and has 1234 items")

	// Test with integer values
	assert(t, i.Ts("pageVars", "name", "Product", "count", 42), "The page is named Product and has 42 items")

	// Test with float values
	assert(t, i.Ts("priceMsg", "price", 19.99, "discount", 15.5), "The price is $19.99 with 15.5% discount")

	// Test with boolean values
	assert(t, i.Ts("statusMsg", "status", "online", "autosave", true), "The system is online and auto-save is true")
	assert(t, i.Ts("statusMsg", "status", "offline", "autosave", false), "The system is offline and auto-save is false")

	// Test with mixed types in a single call
	assert(t, i.Ts("mixedMsg", "user", "John", "count", 5, "total", 99.99, "active", true),
		"User John has 5 items worth $99.99 (active: true)")

	// Test nested translations with non-string values
	assert(t, i.Ts("nestedParams", "appName", "{app.name}"), "Welcome to MyApp")

	// Test with nil value (should convert to "<nil>")
	var nilValue *string = nil
	assert(t, i.Ts("pageVars", "name", nilValue, "count", 0), "The page is named <nil> and has 0 items")
}

func TestTypes(t *testing.T) {
	j := `
{
	"_.code": "en",
	"_.name": "English",
	"template": "Value: {val}, Key: {key}",
	"numbers": "Int: {int}, Float: {float}, Negative: {neg}",
	"complex": "{a} {b} {c} {d} {e}"
}
`

	i, err := New([]byte(j))
	if err != nil {
		t.Fatal(err)
	}

	// Test with struct (uses its string representation).
	type TestStruct struct {
		Name string
		Age  int
	}
	ts := TestStruct{Name: "Alice", Age: 30}
	res := i.Ts("template", "val", ts, "key", "test")
	// The exact output depends on Go's fmt.Sprintf("%v", struct).
	if !strings.Contains(res, "Alice") || !strings.Contains(res, "30") {
		t.Errorf("Expected struct to be formatted, got: %s", res)
	}

	// Test with slice.
	slice := []int{1, 2, 3}
	assert(t, i.Ts("template", "val", slice, "key", "array"), "Value: [1 2 3], Key: array")

	// Test with map.
	m := map[string]int{"a": 1, "b": 2}
	res = i.Ts("template", "val", m, "key", "map")
	// Maps may have non-deterministic order, so just check it contains expected parts.
	if !strings.Contains(res, "map[") || !strings.Contains(res, "Key: map") {
		t.Errorf("Expected map to be formatted, got: %s", res)
	}

	// Test with various numeric types.
	var (
		i8  int8    = -123
		i16 int16   = 12345
		i32 int32   = 123456789
		u8  uint8   = 123
		f32 float32 = 1.2345
		f64 float64 = 1.1235
	)

	assert(t, i.Ts("numbers", "int", i8, "float", f32, "neg", i16), fmt.Sprintf("Int: %v, Float: %v, Negative: %v", i8, f32, i16))
	assert(t, i.Ts("complex", "a", i32, "b", u8, "c", f64, "d", true, "e", "text"), fmt.Sprintf("%v %v %v %v %v", i32, u8, f64, true, "text"))

	// Test runes.
	var r rune = '世'
	res = i.Ts("template", "val", r, "key", "rune")
	if !strings.Contains(res, "19990") {
		t.Errorf("Expected rune to be formatted as number, got: %s", res)
	}

	// Test with byte.
	var b byte = 65
	assert(t, i.Ts("template", "val", b, "key", "byte"), "Value: 65, Key: byte")

	// Test error.
	testErr := fmt.Errorf("test error")
	assert(t, i.Ts("template", "val", testErr, "key", "error"), "Value: test error, Key: error")

	// Test invalid arguments / odd number of params.
	assert(t, i.Ts("template", "val"), "template: invalid arguments")
}

func assert(t *testing.T, a, v interface{}) {
	t.Helper()
	if fmt.Sprintf("%v", a) != fmt.Sprintf("%v", v) {
		t.Fatalf("expected '%v', got '%v'", a, v)
	}
}
