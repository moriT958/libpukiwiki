package converter

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark/ast"
)

const maxIndentLevel = 3

type listItemInfo struct {
	level     int
	isOrdered bool
	text      string
}

func extractListItems(doc ast.Node, markdown []byte) (map[int]listItemInfo, error) {
	listLines := make(map[int]listItemInfo)

	lineNumber := func(offset int) int {
		return bytes.Count(markdown[:offset], []byte("\n"))
	}

	if err := ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		li, ok := node.(*ast.ListItem)
		if !ok {
			return ast.WalkContinue, nil
		}

		level := 0
		isOrdered := false
		for p := node.Parent(); p != nil; p = p.Parent() {
			if list, ok := p.(*ast.List); ok {
				level++
				if level == 1 {
					isOrdered = list.IsOrdered()
				}
			}
		}

		if level > maxIndentLevel {
			level = maxIndentLevel
		}

		itemText := extractListItemText(li, markdown)
		line := lineNumber(li.FirstChild().Lines().At(0).Start)
		listLines[line] = listItemInfo{
			level:     level,
			isOrdered: isOrdered,
			text:      itemText,
		}
		return ast.WalkContinue, nil
	}); err != nil {
		return listLines, fmt.Errorf("failed to walk markdown ast: %v", err)
	}

	return listLines, nil
}

func extractListItemText(li *ast.ListItem, markdown []byte) string {
	if p := li.FirstChild(); p != nil {
		return extractInlineText(p, markdown)
	}
	return ""
}
