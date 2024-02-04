package render

import (
	"github.com/jippi/dottie/pkg/ast"
)

func NewFormatter() *Renderer {
	settings := Settings{
		IncludeDisabled:       true,
		UseInterpolatedValues: false,
		ShowBlankLines:        true,
		ShowColors:            false,
		ShowComments:          true,
		ShowGroupBanners:      true,
	}

	return NewRenderer(settings, FormatterHandler)
}

// FormatterHandler is responsible for formatting an .env file according
// to our opinionated style.
func FormatterHandler(hi *HandlerInput) HandlerSignal {
	switch statement := hi.CurrentStatement.(type) {
	case *ast.Newline:
		// Ignore all existing newlines when doing formatting as
		// we will be injecting these ourself in other places.
		return hi.Stop()

	case *ast.Group:
		output := hi.Presenter.Group(statement)
		if len(output) == 0 {
			return hi.Stop()
		}

		buf := NewLineBuffer()

		// If the previous line is a Newline, don't add another one.
		// This could happen if a group is the *first* thing in the document
		if hi.PreviousStatement != nil && !hi.PreviousStatement.Is(&ast.Newline{}) {
			buf.AddNewline()
		}

		return hi.Return(buf.Add(output).AddNewline().Get())

	case *ast.Assignment:
		output := hi.Presenter.Assignment(statement)
		if len(output) == 0 {
			return hi.Stop()
		}

		buf := NewLineBuffer()

		// If the previous Statement was also an Assignment, detect if they should
		// be allowed to cuddle (without newline between them) or not.
		//
		// Statements are only allow cuddle if both have no comments
		if statement.Is(hi.PreviousStatement) && (statement.HasComments() || assignmentHasComments(hi.PreviousStatement)) {
			buf.AddNewline()
		}

		return hi.Return(buf.Add(output).Get())
	}

	return hi.Continue()
}

// assignmentHasComments checks if the Statement is an Assignment
// and if it has any comments attached to it
func assignmentHasComments(statement ast.Statement) bool {
	assignment, ok := statement.(*ast.Assignment)
	if !ok {
		return false
	}

	return assignment.HasComments()
}
