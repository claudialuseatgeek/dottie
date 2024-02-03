package ast

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	"github.com/jippi/dottie/pkg/token"
	"github.com/jippi/dottie/pkg/tui"
)

type Assignment struct {
	Name         string      `json:"key"`       // Name of the key (left hand side of the "=" sign)
	Literal      string      `json:"literal"`   // Value of the key (right hand side of the "=" sign)
	Interpolated string      `json:"value"`     // Value of the key (after interpolation)
	Complete     bool        `json:"complete"`  // The key/value had no value/content after the "=" sign
	Active       bool        `json:"commented"` // The assignment was commented out (#KEY=VALUE)
	Quote        token.Quote `json:"quote"`     // The style of quotes used for the assignment
	Group        *Group      `json:"-"`         // The (optional) group this assignment belongs to
	Comments     []*Comment  `json:"comments"`  // Comments attached to the assignment (e.g. doc block before it)
	Position     Position    `json:"position"`  // Information about position of the assignment in the file
}

func (a *Assignment) statementNode() {}

func (a *Assignment) Is(other Statement) bool {
	if other == nil || a == nil {
		return false
	}

	return reflect.TypeOf(a) == reflect.TypeOf(other)
}

func (a *Assignment) BelongsToGroup(config RenderSettings) bool {
	if a.Group == nil && len(config.FilterGroup) > 0 {
		return false
	}

	return a.Group == nil || a.Group.BelongsToGroup(config)
}

func (a *Assignment) HasComments() bool {
	return len(a.Comments) > 0
}

func (a *Assignment) RenderDocumentation(withoutPrefix bool) string {
	var buff bytes.Buffer

	for _, c := range a.Comments {
		val := c.Value

		if withoutPrefix {
			val = strings.TrimPrefix(val, "# ")
		}

		buff.WriteString(val)
		buff.WriteString("\n")
	}

	return buff.String()
}

func (a *Assignment) Render(config RenderSettings) string {
	if !config.Match(a) || !a.BelongsToGroup(config) {
		return ""
	}

	var buf bytes.Buffer

	if config.WithComments() {
		for _, c := range a.Comments {
			buf.WriteString(c.Render(config, true))
		}
	}

	if !a.Active {
		if config.WithColors() {
			out := tui.Theme.Danger.Printer(tui.RendererWithTTY(&buf))
			out.Print("#")
		} else {
			buf.WriteString("#")
		}
	}

	buf.WriteString(a.Assignment(config))
	buf.WriteString("\n")

	return buf.String()
}

func (a *Assignment) SetQuote(in string) {
	switch in {
	case "\"", "double":
		a.Quote = token.DoubleQuotes
	case "'", "single":
		a.Quote = token.SingleQuotes
	case "none":
		a.Quote = token.NoQuotes
	}
}

func (a *Assignment) ValidationRules() string {
	for _, comment := range a.Comments {
		if comment.Annotation == nil {
			continue
		}

		if comment.Annotation.Key == "dottie/validate" {
			return comment.Annotation.Value
		}
	}

	return ""
}

func (a *Assignment) Assignment(config RenderSettings) string {
	val := a.Literal

	if config.Interpolate {
		val = a.Interpolated
	}

	if config.WithColors() {
		var buf bytes.Buffer

		tui.Theme.Primary.Printer(tui.RendererWithTTY(&buf)).Print(a.Name)
		tui.Theme.Dark.Printer(tui.RendererWithTTY(&buf)).Print("=")
		tui.Theme.Success.Printer(tui.RendererWithTTY(&buf)).Print(a.Quote)
		tui.Theme.Warning.Printer(tui.RendererWithTTY(&buf)).Print(val)
		tui.Theme.Success.Printer(tui.RendererWithTTY(&buf)).Print(a.Quote)

		return buf.String()
	}

	// panic(a.Quote)
	return fmt.Sprintf("%s=%s%s%s", a.Name, a.Quote, val, a.Quote)
}
