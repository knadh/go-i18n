# go-i18n

go-i18n is a tiny i18n library that enables `{langKey: langValue}` style JSON language maps to be loaded and used in Go programs. It is modeled after [vue-i18n](https://kazupon.github.io/vue-i18n/), optionally enabling interoperability of the same language files in Go backends a Vue frontends. It is used by [listmonk](https://github.com/knadh/listmonk) and [dictpress](https://github.com/knadh/dictpress).


## Usage

### Sample JSON language file

A JSON language file looks is a simple map of `key: value` pairs. Singular/plural terms are represented as `Singular|Plural`. `_.code` and `_.name` are mandatory special keys. Check [listmonk translations](https://github.com/knadh/listmonk/tree/master/i18n) for complex examples.

```json
{
	"_.code": "en",
	"_.name": "English",

	"pageTitle": "Welcome to the page",
	"page": "Single page|Many pages",
	"pageVars": "The page is named {name} and has {count} items"
}
```

```go

	i := i18n.NewFromFile("en.json")

	i.T("pageTitle") // Welcome to the page
	i.T("page") // Single Page
	i.S("page") // Single Page
	i.P("page") // Many pages
	i.Tc("page", 1) // Single Page (second param is a number. 1 is singular)
	i.Tc("page", 2) // Many pages (>= 1 is plural)
	i.Ts("pageVars", "name", "Foo", "count", "123") // The page is named Foo and has 123 items
```

Licensed under the MIT license.
