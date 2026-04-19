package converter

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark/ast"
)

type codeblockLineInfo struct {
	isFence bool   // fence行（削除対象）かどうか
	content string // 変換後のコンテンツ（2スペース + 元のコード）
}

func extractCodeblocks(doc ast.Node, markdown []byte) (map[int]codeblockLineInfo, error) {
	codeblockLines := make(map[int]codeblockLineInfo)

	lineNumber := func(offset int) int {
		return bytes.Count(markdown[:offset], []byte("\n"))
	}

	// Find all fence line positions in markdown
	mdLines := bytes.Split(markdown, []byte("\n"))
	fenceLineNums := make([]int, 0)
	for i, line := range mdLines {
		trimmed := bytes.TrimSpace(line)
		if bytes.HasPrefix(trimmed, []byte("```")) {
			fenceLineNums = append(fenceLineNums, i)
		}
	}

	fenceIndex := 0

	err := ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		fcb, ok := node.(*ast.FencedCodeBlock)
		if !ok {
			return ast.WalkContinue, nil
		}

		var startLine, endLine int

		if fcb.Lines().Len() > 0 {
			// Content exists - fence is one line before the first content line
			startLine = lineNumber(fcb.Lines().At(0).Start) - 1

			// Process content lines
			lines := fcb.Lines()
			for i := 0; i < lines.Len(); i++ {
				seg := lines.At(i)
				lineNum := lineNumber(seg.Start)
				content := string(seg.Value(markdown))
				// Remove trailing newline if present
				content = trimTrailingNewline(content)
				codeblockLines[lineNum] = codeblockLineInfo{
					isFence: false,
					content: "  " + content, // 2スペースプレフィックス
				}
				endLine = lineNum
			}
			endLine++
		} else {
			// Empty code block - use fence positions from scan
			if fenceIndex+1 < len(fenceLineNums) {
				startLine = fenceLineNums[fenceIndex]
				endLine = fenceLineNums[fenceIndex+1]
				fenceIndex += 2
			} else {
				return ast.WalkContinue, nil
			}
		}

		// Mark start fence line
		codeblockLines[startLine] = codeblockLineInfo{
			isFence: true,
			content: "",
		}

		// Mark end fence line
		codeblockLines[endLine] = codeblockLineInfo{
			isFence: true,
			content: "",
		}

		return ast.WalkContinue, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk markdown ast: %v", err)
	}

	return codeblockLines, nil
}

func trimTrailingNewline(s string) string {
	if len(s) > 0 && s[len(s)-1] == '\n' {
		return s[:len(s)-1]
	}
	return s
}
