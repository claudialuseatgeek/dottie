package ast

import (
	"reflect"

	"github.com/jippi/dottie/pkg/token"
)

// Comment node represents a comment statement.
type Comment struct {
	Value      string            `json:"value"`      // The actual comment value
	Annotation *token.Annotation `json:"annotation"` // If the comment was detected to be an annotation
	Group      *Group            `json:"-"`          // The (optional) group the comment belongs to
	Position   Position          `json:"position"`   // Information about position of the assignment in the file
}

func NewComment(value string) *Comment {
	return &Comment{
		Value: "# " + value,
	}
}

func (c *Comment) Is(other Statement) bool {
	return reflect.TypeOf(c) == reflect.TypeOf(other)
}

func (c *Comment) BelongsToGroup(name string) bool {
	if c.Group == nil && len(name) > 0 {
		return false
	}

	return c.Group == nil || c.Group.BelongsToGroup(name)
}

func (c *Comment) statementNode() {
}

func (c *Comment) String() string {
	return c.Value
}
