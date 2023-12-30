package formatting

import (
	"fmt"
	"strings"
	"text/template"
	"text/template/parse"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const ParseMode = tgbotapi.ModeMarkdownV2

// Escape makes arbitrary text valid when parsed with ParseMode.
func Escape(s string) string {
	return tgbotapi.EscapeText(ParseMode, s)
}

// Escapef is a formatted version of Escape, but it does not escape
// a. This means that these values must be valid when parsed via ParseMode.
func Escapef(format string, a ...any) string {
	return fmt.Sprintf(Escape(format), a...)
}

func Bold(s string) string {
	return "*" + s + "*"
}

func Mono(s string) string {
	return "`" + s + "`"
}

// RenderTemplate renders text/template.Template tmpl. All plaint text nodes
// of the template are escaped via Escape. The following functions are also
// availible in template context:
//
// - "bold" -> Bold
//
// - "mono" -> Mono
//
// - "esc"  -> Escape
func RenderTemplate(tmpl string, data any) (string, error) {
	t, err := template.New("").Funcs(template.FuncMap{
		"bold": Bold,
		"mono": Mono,
		"esc":  Escape,
	}).Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}
	escapeTemplate(t)

	bldr := &strings.Builder{}
	err = t.Execute(bldr, data)
	if err != nil {
		return "", fmt.Errorf("render template: %w", err)
	}
	return bldr.String(), nil
}

// MustRenderTemplate acts like RenderTemplate, but panics on errors.
func MustRenderTemplate(tmpl string, data any) string {
	s, err := RenderTemplate(tmpl, data)
	if err != nil {
		panic(fmt.Sprintf("standup: render template: %s", err))
	}
	return s
}

func escapeTemplate(t *template.Template) {
	if t.Root == nil {
		return
	}
	nodes := []parse.Node{t.Root}
	for i := 0; i < len(nodes); i++ {
		switch n := nodes[i].(type) {
		case *parse.TextNode:
			n.Text = []byte(Escape(string(n.Text)))
		case *parse.ListNode:
			if n != nil {
				nodes = append(nodes, n.Nodes...)
			}
		case *parse.IfNode:
			nodes = append(nodes, n.List, n.ElseList)
		case *parse.RangeNode:
			nodes = append(nodes, n.List, n.ElseList)
		case *parse.WithNode:
			nodes = append(nodes, n.List, n.ElseList)
		}
	}
}
