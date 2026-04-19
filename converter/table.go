package converter

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
)

type tableRowInfo struct {
	isHeader    bool     // ヘッダー行か
	isSeparator bool     // セパレータ行か（削除対象）
	cells       []string // セル内容
}

func extractTables(doc ast.Node, markdown []byte) (map[int]tableRowInfo, error) {
	tableLines := make(map[int]tableRowInfo)

	lineNumber := func(offset int) int {
		return bytes.Count(markdown[:offset], []byte("\n"))
	}

	err := ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		table, ok := node.(*east.Table)
		if !ok {
			return ast.WalkContinue, nil
		}

		// Process table children (TableHeader and TableRow)
		for child := table.FirstChild(); child != nil; child = child.NextSibling() {
			switch row := child.(type) {
			case *east.TableHeader:
				// Process header row
				cells := extractTableCells(row, markdown)
				line := getRowLineNumber(row, lineNumber)
				if line >= 0 {
					tableLines[line] = tableRowInfo{
						isHeader:    true,
						isSeparator: false,
						cells:       cells,
					}
					// Mark the separator line (next line after header)
					tableLines[line+1] = tableRowInfo{
						isHeader:    false,
						isSeparator: true,
						cells:       nil,
					}
				}
			case *east.TableRow:
				// Process data row
				cells := extractTableCells(row, markdown)
				line := getRowLineNumber(row, lineNumber)
				if line >= 0 {
					tableLines[line] = tableRowInfo{
						isHeader:    false,
						isSeparator: false,
						cells:       cells,
					}
				}
			}
		}

		return ast.WalkContinue, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk markdown ast: %v", err)
	}

	return tableLines, nil
}

// getRowLineNumber extracts line number from the first cell's text segment
func getRowLineNumber(row ast.Node, lineNumber func(int) int) int {
	for cell := row.FirstChild(); cell != nil; cell = cell.NextSibling() {
		if _, ok := cell.(*east.TableCell); ok {
			for child := cell.FirstChild(); child != nil; child = child.NextSibling() {
				if txt, ok := child.(*ast.Text); ok {
					return lineNumber(txt.Segment.Start)
				}
			}
		}
	}
	return -1
}

func extractTableCells(row ast.Node, markdown []byte) []string {
	var cells []string
	for cell := row.FirstChild(); cell != nil; cell = cell.NextSibling() {
		if _, ok := cell.(*east.TableCell); ok {
			cellText := extractCellText(cell, markdown)
			cells = append(cells, cellText)
		}
	}
	return cells
}

func extractCellText(cell ast.Node, markdown []byte) string {
	var buf bytes.Buffer
	for child := cell.FirstChild(); child != nil; child = child.NextSibling() {
		switch c := child.(type) {
		case *ast.Text:
			buf.Write(c.Segment.Value(markdown))
		default:
			buf.WriteString(extractInlineText(child, markdown))
		}
	}
	return strings.TrimSpace(buf.String())
}

func convertTableRow(tr tableRowInfo) string {
	var buf bytes.Buffer
	for i, cell := range tr.cells {
		if i == 0 {
			buf.WriteString("|")
		}
		if tr.isHeader {
			buf.WriteString("~ ")
		} else {
			buf.WriteString(" ")
		}
		buf.WriteString(cell)
		buf.WriteString(" |")
	}
	return buf.String()
}
