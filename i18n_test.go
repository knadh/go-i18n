package i18n

import (
	"fmt"
	"testing"
)

func TestTestXxx(t *testing.T) {
	j := `
{
	"_.code": "en",
	"_.name": "English",

	"pageTitle": "Welcome to the page",
	"foo": "Foo",
	"page": "Single page|Many pages",
	"pageVars": "The page is named {name} and has {count} items"
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
}

func assert(t *testing.T, a, v interface{}) {
	t.Helper()
	if fmt.Sprintf("%v", a) != fmt.Sprintf("%v", v) {
		t.Fatalf("expected '%v', got '%v'", a, v)
	}
}
