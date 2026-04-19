package converter

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark/ast"
)

type boldInfo struct {
	originalText  string // "**text**"
	convertedText string // "''text''"
}

func extractBolds(doc ast.Node, markdown []byte) ([]boldInfo, error) {
	var bolds []boldInfo

	err := ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		em, ok := node.(*ast.Emphasis)
		if !ok {
			return ast.WalkContinue, nil
		}

		// Level 2 is bold (**text**), Level 1 is italic (*text*)
		if em.Level != 2 {
			return ast.WalkContinue, nil
		}

		text := extractEmphasisText(em, markdown)
		bolds = append(bolds, boldInfo{
			originalText:  "**" + text + "**",
			convertedText: "''" + text + "''",
		})

		return ast.WalkContinue, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk markdown ast: %v", err)
	}

	return bolds, nil
}

func extractEmphasisText(em *ast.Emphasis, markdown []byte) string {
	var buf bytes.Buffer
	for child := em.FirstChild(); child != nil; child = child.NextSibling() {
		if t, ok := child.(*ast.Text); ok {
			buf.Write(t.Segment.Value(markdown))
		}
	}
	return buf.String()
}
