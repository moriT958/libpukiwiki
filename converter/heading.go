package converter

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark/ast"
)

const maxHeadingLevel = 3

type headingInfo struct {
	level int
	text  string
}

func extractHeadings(doc ast.Node, markdown []byte) (map[int]headingInfo, error) {
	headingLines := make(map[int]headingInfo)

	lineNumber := func(offset int) int {
		return bytes.Count(markdown[:offset], []byte("\n"))
	}

	err := ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		h, ok := node.(*ast.Heading)
		if !ok {
			return ast.WalkContinue, nil
		}

		headingText := extractHeadingText(h, markdown)
		if t, ok := h.FirstChild().(*ast.Text); ok {
			line := lineNumber(t.Segment.Start)
			headingLines[line] = headingInfo{
				level: h.Level,
				text:  headingText,
			}
		}
		return ast.WalkContinue, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk markdown ast: %v", err)
	}

	return headingLines, nil
}

func extractHeadingText(h *ast.Heading, markdown []byte) string {
	return extractInlineText(h, markdown)
}

func extractInlineText(node ast.Node, markdown []byte) string {
	var buf bytes.Buffer
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		switch c := child.(type) {
		case *ast.Text:
			buf.Write(c.Segment.Value(markdown))
		case *ast.Emphasis:
			marker := "*"
			if c.Level == 2 {
				marker = "**"
			}
			buf.WriteString(marker)
			buf.WriteString(extractInlineText(c, markdown))
			buf.WriteString(marker)
		case *ast.Link:
			text := extractInlineText(c, markdown)
			url := string(c.Destination)
			buf.WriteString("[" + text + "](" + url + ")")
		default:
			buf.WriteString(extractInlineText(child, markdown))
		}
	}
	return buf.String()
}
