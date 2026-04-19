package converter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "H1",
			input:    []byte("# Heading 1"),
			expected: "* Heading 1",
		},
		{
			name:     "H2",
			input:    []byte("## Heading 2"),
			expected: "** Heading 2",
		},
		{
			name:     "H3",
			input:    []byte("### Heading 3"),
			expected: "*** Heading 3",
		},
		{
			name:     "複数種類の Heading がある場合も全て変換される",
			input:    []byte("# H1\n\n## H2\n\n### H3"),
			expected: "* H1\n\n** H2\n\n*** H3",
		},
		{
			name:     "H4 は変換されない",
			input:    []byte("# H1\n\n#### H4\n\n## H2"),
			expected: "* H1\n\n#### H4\n\n** H2",
		},
		{
			name:     "内容がそのまま保たれる",
			input:    []byte("# Title\n\nSome paragraph text.\n\n## Section\n\nPlain text."),
			expected: "* Title\n\nSome paragraph text.\n\n** Section\n\nPlain text.",
		},
		// Unordered List のテストケース
		{
			name:     "基本的な unordered list",
			input:    []byte("- item1\n- item2\n- item3"),
			expected: "-item1\n-item2\n-item3",
		},
		{
			name:     "ネストした unordered list",
			input:    []byte("- item1\n  - nested1\n  - nested2\n- item2"),
			expected: "-item1\n--nested1\n--nested2\n-item2",
		},
		{
			name:     "3レベルのネスト unordered list",
			input:    []byte("- level1\n  - level2\n    - level3"),
			expected: "-level1\n--level2\n---level3",
		},
		// Ordered List のテストケース
		{
			name:     "基本的な ordered list",
			input:    []byte("1. item1\n2. item2\n3. item3"),
			expected: "+item1\n+item2\n+item3",
		},
		{
			name:     "ネストした ordered list",
			input:    []byte("1. item1\n   1. nested1\n   2. nested2\n2. item2"),
			expected: "+item1\n++nested1\n++nested2\n+item2",
		},
		{
			name:     "3レベルのネスト ordered list",
			input:    []byte("1. level1\n   1. level2\n      1. level3"),
			expected: "+level1\n++level2\n+++level3",
		},
		// Codeblock のテストケース
		{
			name:     "基本的なコードブロック",
			input:    []byte("```\ncode\n```"),
			expected: "  code",
		},
		{
			name:     "言語指定付きコードブロック",
			input:    []byte("```go\nfunc main() {}\n```"),
			expected: "  func main() {}",
		},
		{
			name:     "複数行コードブロック",
			input:    []byte("```\nline1\nline2\nline3\n```"),
			expected: "  line1\n  line2\n  line3",
		},
		{
			name:     "インデント保持の確認",
			input:    []byte("```\n  indented\n    more indented\n```"),
			expected: "    indented\n      more indented",
		},
		{
			name:     "見出しとコードブロックの混在",
			input:    []byte("# Title\n\n```\ncode\n```\n\n## Section"),
			expected: "* Title\n\n  code\n\n** Section",
		},
		{
			name:     "リストとコードブロックの混在",
			input:    []byte("- item1\n\n```\ncode\n```\n\n- item2"),
			expected: "-item1\n\n  code\n\n-item2",
		},
		{
			name:     "空のコードブロック",
			input:    []byte("```\n```"),
			expected: "",
		},
		// Bold のテストケース
		{
			name:     "基本的なBold変換",
			input:    []byte("This is **bold** text"),
			expected: "This is ''bold'' text",
		},
		{
			name:     "複数のBoldが同一行にある場合",
			input:    []byte("**first** and **second** bold"),
			expected: "''first'' and ''second'' bold",
		},
		{
			name:     "見出し内のBold",
			input:    []byte("# Heading with **bold**"),
			expected: "* Heading with ''bold''",
		},
		{
			name:     "リスト内のBold",
			input:    []byte("- item with **bold**\n- normal item"),
			expected: "-item with ''bold''\n-normal item",
		},
		{
			name:     "コードブロック内のBoldは変換しない",
			input:    []byte("```\n**not bold**\n```"),
			expected: "  **not bold**",
		},
		{
			name:     "イタリックは変換しない",
			input:    []byte("This is *italic* text"),
			expected: "This is *italic* text",
		},
		{
			name:     "BoldとItalicが混在",
			input:    []byte("**bold** and *italic*"),
			expected: "''bold'' and *italic*",
		},
		// Link のテストケース
		{
			name:     "基本的なLink変換",
			input:    []byte("Click [here](https://example.com) for more"),
			expected: "Click [[here>https://example.com]] for more",
		},
		{
			name:     "複数のLinkが同一行にある場合",
			input:    []byte("[first](https://a.com) and [second](https://b.com)"),
			expected: "[[first>https://a.com]] and [[second>https://b.com]]",
		},
		{
			name:     "見出し内のLink",
			input:    []byte("# Heading with [link](https://example.com)"),
			expected: "* Heading with [[link>https://example.com]]",
		},
		{
			name:     "リスト内のLink",
			input:    []byte("- item with [link](https://example.com)\n- normal item"),
			expected: "-item with [[link>https://example.com]]\n-normal item",
		},
		{
			name:     "コードブロック内のLinkは変換しない",
			input:    []byte("```\n[not link](https://example.com)\n```"),
			expected: "  [not link](https://example.com)",
		},
		// Table のテストケース
		{
			name:     "基本的なテーブル変換",
			input:    []byte("| Col1 | Col2 |\n| --- | --- |\n| A | B |"),
			expected: "|~ Col1 |~ Col2 |\n| A | B |",
		},
		{
			name:     "複数行のテーブル",
			input:    []byte("| Column1 | Column2 | Column3 |\n| ------- | ------- | ------- |\n| Item1.1 | Item2.1 | Item3.1 |\n| Item1.2 | Item2.2 | Item3.2 |"),
			expected: "|~ Column1 |~ Column2 |~ Column3 |\n| Item1.1 | Item2.1 | Item3.1 |\n| Item1.2 | Item2.2 | Item3.2 |",
		},
		{
			name:     "見出しとテーブルの混在",
			input:    []byte("# Title\n\n| A | B |\n| - | - |\n| 1 | 2 |"),
			expected: "* Title\n\n|~ A |~ B |\n| 1 | 2 |",
		},
		{
			name:     "リストとテーブルの混在",
			input:    []byte("- item1\n\n| A | B |\n| - | - |\n| 1 | 2 |\n\n- item2"),
			expected: "-item1\n\n|~ A |~ B |\n| 1 | 2 |\n\n-item2",
		},
		{
			name:     "コードブロックとテーブルの混在",
			input:    []byte("```\ncode\n```\n\n| A | B |\n| - | - |\n| 1 | 2 |"),
			expected: "  code\n\n|~ A |~ B |\n| 1 | 2 |",
		},
		{
			name:     "複数テーブル",
			input:    []byte("| A | B |\n| - | - |\n| 1 | 2 |\n\n| X | Y |\n| - | - |\n| 3 | 4 |"),
			expected: "|~ A |~ B |\n| 1 | 2 |\n\n|~ X |~ Y |\n| 3 | 4 |",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Convert(tt.input)
			if err != nil {
				t.Fatalf("Convert returned error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestConvert_Integration(t *testing.T) {
	input, err := os.ReadFile(filepath.Join("testdata", "sample.md"))
	if err != nil {
		t.Fatalf("ReadFile(sample.md) returned error: %v", err)
	}

	expected, err := os.ReadFile(filepath.Join("testdata", "sample.golden.txt"))
	if err != nil {
		t.Fatalf("ReadFile(sample.golden.txt) returned error: %v", err)
	}

	result, err := Convert(input)
	if err != nil {
		t.Fatalf("Convert returned error: %v", err)
	}

	if result != string(expected) {
		t.Errorf("Expected %q, got %q", string(expected), result)
	}
}
