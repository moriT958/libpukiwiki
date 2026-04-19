package converter

import (
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []token
	}{
		// --- Heading ---
		{
			name:  "H1 with space",
			input: "* Heading",
			want:  []token{{tokenHeading, "Heading", 1}, {tokenEOF, "", 0}},
		},
		{
			name:  "H1 without space",
			input: "*Heading",
			want:  []token{{tokenHeading, "Heading", 1}, {tokenEOF, "", 0}},
		},
		{
			name:  "H2",
			input: "** Heading",
			want:  []token{{tokenHeading, "Heading", 2}, {tokenEOF, "", 0}},
		},
		{
			name:  "H3",
			input: "*** Heading",
			want:  []token{{tokenHeading, "Heading", 3}, {tokenEOF, "", 0}},
		},
		{
			name:  "H4 is plain text",
			input: "**** Heading",
			want:  []token{{tokenText, "**** Heading", 0}, {tokenEOF, "", 0}},
		},

		// --- Unordered List ---
		{
			name:  "unordered list level 1",
			input: "- item",
			want:  []token{{tokenUnorderedList, "item", 1}, {tokenEOF, "", 0}},
		},
		{
			name:  "unordered list level 1 no space",
			input: "-item",
			want:  []token{{tokenUnorderedList, "item", 1}, {tokenEOF, "", 0}},
		},
		{
			name:  "unordered list level 2",
			input: "-- item",
			want:  []token{{tokenUnorderedList, "item", 2}, {tokenEOF, "", 0}},
		},
		{
			name:  "unordered list level 3",
			input: "--- item",
			want:  []token{{tokenUnorderedList, "item", 3}, {tokenEOF, "", 0}},
		},

		// --- Ordered List ---
		{
			name:  "ordered list level 1",
			input: "+ item",
			want:  []token{{tokenOrderedList, "item", 1}, {tokenEOF, "", 0}},
		},
		{
			name:  "ordered list level 2",
			input: "++ item",
			want:  []token{{tokenOrderedList, "item", 2}, {tokenEOF, "", 0}},
		},
		{
			name:  "ordered list level 3",
			input: "+++ item",
			want:  []token{{tokenOrderedList, "item", 3}, {tokenEOF, "", 0}},
		},

		// --- Code Line ---
		{
			name:  "code line (2 spaces)",
			input: "  code here",
			want:  []token{{tokenCodeLine, "code here", 0}, {tokenEOF, "", 0}},
		},

		// --- Table Row ---
		{
			name:  "table row",
			input: "|~ col1 |~ col2 |",
			want:  []token{{tokenTableRow, "|~ col1 |~ col2 |", 0}, {tokenEOF, "", 0}},
		},
		{
			name:  "table data row",
			input: "| data | data |",
			want:  []token{{tokenTableRow, "| data | data |", 0}, {tokenEOF, "", 0}},
		},
		{
			name:  "single pipe is plain text",
			input: "|",
			want:  []token{{tokenText, "|", 0}, {tokenEOF, "", 0}},
		},

		// --- Horizontal Rule ---
		{
			name:  "horizontal rule",
			input: "----",
			want:  []token{{tokenHR, "", 0}, {tokenEOF, "", 0}},
		},
		{
			name:  "horizontal rule with trailing text (text is discarded)",
			input: "----text after",
			want:  []token{{tokenHR, "", 0}, {tokenEOF, "", 0}},
		},

		// --- Newline ---
		{
			name:  "empty line",
			input: "",
			want:  []token{{tokenNewline, "", 0}, {tokenEOF, "", 0}},
		},

		// --- Plain Text ---
		{
			name:  "plain text",
			input: "hello world",
			want:  []token{{tokenText, "hello world", 0}, {tokenEOF, "", 0}},
		},

		// --- Bold ---
		{
			name:  "bold",
			input: "''bold''",
			want: []token{
				{tokenBoldOpen, "''", 0},
				{tokenText, "bold", 0},
				{tokenBoldClose, "''", 0},
				{tokenEOF, "", 0},
			},
		},
		{
			name:  "bold with surrounding text",
			input: "hello ''world'' end",
			want: []token{
				{tokenText, "hello ", 0},
				{tokenBoldOpen, "''", 0},
				{tokenText, "world", 0},
				{tokenBoldClose, "''", 0},
				{tokenText, " end", 0},
				{tokenEOF, "", 0},
			},
		},
		{
			name:  "two bold spans",
			input: "''a'' and ''b''",
			want: []token{
				{tokenBoldOpen, "''", 0},
				{tokenText, "a", 0},
				{tokenBoldClose, "''", 0},
				{tokenText, " and ", 0},
				{tokenBoldOpen, "''", 0},
				{tokenText, "b", 0},
				{tokenBoldClose, "''", 0},
				{tokenEOF, "", 0},
			},
		},

		// --- Italic ---
		{
			name:  "italic",
			input: "'''italic'''",
			want: []token{
				{tokenItalicOpen, "'''", 0},
				{tokenText, "italic", 0},
				{tokenItalicClose, "'''", 0},
				{tokenEOF, "", 0},
			},
		},

		// --- Link ---
		{
			name:  "link with separator",
			input: "[[label>https://example.com]]",
			want: []token{
				{tokenLinkOpen, "[[", 0},
				{tokenLinkText, "label", 0},
				{tokenLinkSep, ">", 0},
				{tokenLinkURL, "https://example.com", 0},
				{tokenLinkClose, "]]", 0},
				{tokenEOF, "", 0},
			},
		},
		{
			name:  "link without separator",
			input: "[[nodest]]",
			want: []token{
				{tokenLinkOpen, "[[", 0},
				{tokenLinkText, "nodest", 0},
				{tokenLinkClose, "]]", 0},
				{tokenEOF, "", 0},
			},
		},
		{
			name:  "link with surrounding text",
			input: "see [[x>y]] now",
			want: []token{
				{tokenText, "see ", 0},
				{tokenLinkOpen, "[[", 0},
				{tokenLinkText, "x", 0},
				{tokenLinkSep, ">", 0},
				{tokenLinkURL, "y", 0},
				{tokenLinkClose, "]]", 0},
				{tokenText, " now", 0},
				{tokenEOF, "", 0},
			},
		},
		{
			name:  "malformed link (no closing brackets)",
			input: "[[broken",
			want: []token{
				{tokenLinkOpen, "[[", 0},
				{tokenText, "broken", 0},
				{tokenEOF, "", 0},
			},
		},

		// --- Multi-line ---
		{
			name:  "heading then blank then text",
			input: "* H\n\ntext",
			want: []token{
				{tokenHeading, "H", 1},
				{tokenNewline, "", 0},
				{tokenText, "text", 0},
				{tokenEOF, "", 0},
			},
		},
		{
			name:  "mixed block types",
			input: "* H\n- item\n+ ordered\n  code",
			want: []token{
				{tokenHeading, "H", 1},
				{tokenUnorderedList, "item", 1},
				{tokenOrderedList, "ordered", 1},
				{tokenCodeLine, "code", 0},
				{tokenEOF, "", 0},
			},
		},
		{
			name:  "table row among text",
			input: "text\n| A | B |",
			want: []token{
				{tokenText, "text", 0},
				{tokenTableRow, "| A | B |", 0},
				{tokenEOF, "", 0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newLexer(tt.input)
			got := l.tokenize()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot:  %v\nwant: %v", got, tt.want)
			}
		})
	}
}
