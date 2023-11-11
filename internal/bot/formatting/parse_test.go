package formatting_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xopoww/standup/internal/bot/formatting"
)

func TestRenderTemplate(t *testing.T) {
	cc := []struct {
		name string
		tmpl string
		want string
	}{
		{"simple", `Just a text`, "Just a text"},
		{"macro", `Text with {{ mono "monospace item" }}`, "Text with `monospace item`"},
		{"escape", `Text with a dot.`, "Text with a dot\\."},
		{"if", `Outer text. {{ if true }}Inner text.{{end}}`, "Outer text\\. Inner text\\."},
	}

	for _, c := range cc {
		t.Run(c.name, func(t *testing.T) {
			s, err := formatting.RenderTemplate(c.tmpl, nil)
			require.NoError(t, err)
			require.Equal(t, c.want, s)
		})
	}
}
