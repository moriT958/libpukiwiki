package converter

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark/ast"
)

type linkInfo struct {
	originalText  string // "[text](url)"
	convertedText string // "[[text>url]]"
}

func extractLinks(doc ast.Node, markdown []byte) ([]linkInfo, error) {
	var links []linkInfo

	err := ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		link, ok := node.(*ast.Link)
		if !ok {
			return ast.WalkContinue, nil
		}

		text := extractLinkText(link, markdown)
		url := string(link.Destination)
		links = append(links, linkInfo{
			originalText:  "[" + text + "](" + url + ")",
			convertedText: "[[" + text + ">" + url + "]]",
		})

		return ast.WalkContinue, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk markdown ast: %v", err)
	}

	return links, nil
}

func extractLinkText(link *ast.Link, markdown []byte) string {
	var buf bytes.Buffer
	for child := link.FirstChild(); child != nil; child = child.NextSibling() {
		if t, ok := child.(*ast.Text); ok {
			buf.Write(t.Segment.Value(markdown))
		}
	}
	return buf.String()
}
