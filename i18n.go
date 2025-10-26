// Package i18n is a simple lib that enables i18n translations using a language map.
// It loads JSON {langKey: langString} maps and mimics the functionality of the
// Javascript vue-i18n library enabling the same JSON language map to be reused in
// a Go backend application and a Vue frontend application. It can be used
// as a standalone lib in Go backends.
package i18n

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// I18n enables simple translation functions over a language map.
type I18n struct {
	code    string `json:"code"`
	name    string `json:"name"`
	langMap map[string]string
}

var reParam = regexp.MustCompile(`(?i)\{([a-z0-9-.]+)\}`)

// New returns an I18n instance from the given JSON language map bytes.
func New(jsonB []byte) (*I18n, error) {
	var l map[string]string
	if err := json.Unmarshal(jsonB, &l); err != nil {
		return nil, err
	}

	code, ok := l["_.code"]
	if !ok {
		return nil, errors.New("missing _.code field in language file")
	}

	name, ok := l["_.name"]
	if !ok {
		return nil, errors.New("missing _.name field in language file")
	}

	return &I18n{
		langMap: l,
		code:    code,
		name:    name,
	}, nil
}

// NewFromFile returns a I18n instance with the JSON language map read
// from the given file.
func NewFromFile(filepath string) (*I18n, error) {
	b, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return New(b)
}

// Load loads a JSON language map into the instance overwriting
// existing keys that conflict.
func (i *I18n) Load(b []byte) error {
	var l map[string]string
	if err := json.Unmarshal(b, &l); err != nil {
		return err
	}

	for k, v := range l {
		i.langMap[k] = v
	}

	return nil
}

// Name returns the canonical name of the language.
func (i *I18n) Name() string {
	return i.name
}

// Code returns the ISO code of the language.
func (i *I18n) Code() string {
	return i.code
}

// JSON returns the languagemap as raw JSON.
func (i *I18n) JSON() []byte {
	b, _ := json.Marshal(i.langMap)
	return b
}

// T returns the translation string for the given key.
func (i *I18n) T(key string) string {
	s, ok := i.langMap[key]
	if !ok {
		return key
	}

	return i.getSingular(s)
}

// Ts returns the translation for the given key similar to vue i18n's t()
// and substitutes the params in the given map in the translated value.
// In the language values, the substitutions are represented as: {key}
// The params and values are received as a pairs of succeeding values.
// That is, the number of these arguments should be an even number.
// Values can be of any type and will be converted to strings using fmt.Sprintf("%v", value).
// eg: Ts("globals.message.notFound",
//
//	"name", "campaigns",
//	"count", 123,
//	"error", err)
func (i *I18n) Ts(key string, params ...any) string {
	if len(params)%2 != 0 {
		return key + `: invalid arguments`
	}

	s, ok := i.langMap[key]
	if !ok {
		return key
	}

	s = i.getSingular(s)
	for n := 0; n < len(params); n += 2 {
		// Convert the key to string.
		paramKey, ok := params[n].(string)
		if !ok {
			paramKey = fmt.Sprintf("%v", params[n])
		}

		// If there are {params} in the param values, substitute them.
		val := i.subAllParams(params[n+1])
		s = strings.ReplaceAll(s, `{`+paramKey+`}`, val)
	}

	return s
}

// Tc returns the translation for the given key similar to vue i18n's tc().
// It expects the language string in the map to be of the form `Singular | Plural` and
// returns `Plural` if n > 1, or `Singular` otherwise.
func (i *I18n) Tc(key string, n int) string {
	s, ok := i.langMap[key]
	if !ok {
		return key
	}

	// Plural.
	if n > 1 {
		return i.getPlural(s)
	}

	return i.getSingular(s)
}

// S returns the singular form of a string that's represented as Singular|Plural.
func (i *I18n) S(key string) string {
	return i.Tc(key, 1)
}

// P returns the Plural form of a string that's represented as Singular|Plural.
func (i *I18n) P(key string) string {
	return i.Tc(key, 2)
}

// getSingular returns the singular term from the vuei18n pipe separated value.
// singular term | plural term
func (i *I18n) getSingular(s string) string {
	if !strings.Contains(s, "|") {
		return s
	}

	return strings.TrimSpace(strings.Split(s, "|")[0])
}

// getSingular returns the plural term from the vuei18n pipe separated value.
// singular term | plural term
func (i *I18n) getPlural(s string) string {
	if !strings.Contains(s, "|") {
		return s
	}

	chunks := strings.Split(s, "|")
	if len(chunks) == 2 {
		return strings.TrimSpace(chunks[1])
	}

	return strings.TrimSpace(chunks[0])
}

// subAllParams recursively resolves and replaces all {params} in a value.
// It converts any type to string and then processes nested translations.
func (i *I18n) subAllParams(v any) string {
	// Convert any value to string
	s, ok := v.(string)
	if !ok {
		s = fmt.Sprintf("%v", v)
	}

	if !strings.Contains(s, `{`) {
		return s
	}

	parts := reParam.FindAllStringSubmatch(s, -1)
	if len(parts) < 1 {
		return s
	}

	for _, p := range parts {
		s = strings.ReplaceAll(s, p[0], i.T(p[1]))
	}

	return i.subAllParams(s)
}
