package converter

import (
	"strings"
)

type tokenType int

const (
	tokenHeading       tokenType = iota // * ** ***
	tokenUnorderedList                  // - -- ---
	tokenOrderedList                    // + ++ +++
	tokenCodeLine                       // 行頭2スペース
	tokenTableRow                       // |cell|cell| (パイプ2個以上)
	tokenHR                             // ---- 水平線
	tokenBoldOpen                       // '' (2シングルクォート)
	tokenBoldClose                      // ''
	tokenItalicOpen                     // ''' (3シングルクォート)
	tokenItalicClose                    // '''
	tokenLinkOpen                       // [[
	tokenLinkText                       // [[text>url]] の text 部分
	tokenLinkSep                        // > (リンク内)
	tokenLinkURL                        // [[text>url]] の url 部分
	tokenLinkClose                      // ]]
	tokenText                           // プレーンテキスト
	tokenNewline                        // 空行
	tokenEOF
)

type token struct {
	typ   tokenType
	value string
	level int // Heading/List のレベル (1-3), その他は 0
}

type lexer struct {
	lines []string
}

func newLexer(input string) *lexer {
	return &lexer{
		lines: strings.Split(strings.TrimRight(input, "\n"), "\n"),
	}
}

func (l *lexer) tokenize() []token {
	var tokens []token

	for _, line := range l.lines {
		switch {
		case line == "":
			tokens = append(tokens, token{typ: tokenNewline})

		case strings.HasPrefix(line, "  "):
			tokens = append(tokens, token{typ: tokenCodeLine, value: line[2:]})

		case strings.HasPrefix(line, "|") && strings.Count(line, "|") >= 2:
			tokens = append(tokens, token{typ: tokenTableRow, value: line})

		case strings.HasPrefix(line, "----"):
			tokens = append(tokens, token{typ: tokenHR})

		case strings.HasPrefix(line, "***") && !strings.HasPrefix(line, "****"):
			tokens = append(tokens, token{typ: tokenHeading, value: strings.TrimSpace(line[3:]), level: 3})

		case strings.HasPrefix(line, "**") && !strings.HasPrefix(line, "***"):
			tokens = append(tokens, token{typ: tokenHeading, value: strings.TrimSpace(line[2:]), level: 2})

		case strings.HasPrefix(line, "*") && !strings.HasPrefix(line, "**"):
			tokens = append(tokens, token{typ: tokenHeading, value: strings.TrimSpace(line[1:]), level: 1})

		case strings.HasPrefix(line, "---") && !strings.HasPrefix(line, "----"):
			tokens = append(tokens, token{typ: tokenUnorderedList, value: strings.TrimSpace(line[3:]), level: 3})

		case strings.HasPrefix(line, "--") && !strings.HasPrefix(line, "---"):
			tokens = append(tokens, token{typ: tokenUnorderedList, value: strings.TrimSpace(line[2:]), level: 2})

		case strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "--"):
			tokens = append(tokens, token{typ: tokenUnorderedList, value: strings.TrimSpace(line[1:]), level: 1})

		case strings.HasPrefix(line, "+++") && !strings.HasPrefix(line, "++++"):
			tokens = append(tokens, token{typ: tokenOrderedList, value: strings.TrimSpace(line[3:]), level: 3})

		case strings.HasPrefix(line, "++") && !strings.HasPrefix(line, "+++"):
			tokens = append(tokens, token{typ: tokenOrderedList, value: strings.TrimSpace(line[2:]), level: 2})

		case strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "++"):
			tokens = append(tokens, token{typ: tokenOrderedList, value: strings.TrimSpace(line[1:]), level: 1})

		default:
			tokens = append(tokens, tokenizeInline(line)...)
		}
	}

	tokens = append(tokens, token{typ: tokenEOF})
	return tokens
}

func tokenizeInline(line string) []token {
	var tokens []token
	var buf strings.Builder
	inBold := false
	inItalic := false
	i := 0

	flush := func() {
		if buf.Len() > 0 {
			tokens = append(tokens, token{typ: tokenText, value: buf.String()})
			buf.Reset()
		}
	}

	for i < len(line) {
		switch {
		case strings.HasPrefix(line[i:], "'''"):
			flush()
			if !inItalic {
				tokens = append(tokens, token{typ: tokenItalicOpen, value: "'''"})
				inItalic = true
			} else {
				tokens = append(tokens, token{typ: tokenItalicClose, value: "'''"})
				inItalic = false
			}
			i += 3

		case strings.HasPrefix(line[i:], "''"):
			flush()
			if !inBold {
				tokens = append(tokens, token{typ: tokenBoldOpen, value: "''"})
				inBold = true
			} else {
				tokens = append(tokens, token{typ: tokenBoldClose, value: "''"})
				inBold = false
			}
			i += 2

		case strings.HasPrefix(line[i:], "[["):
			flush()
			tokens = append(tokens, token{typ: tokenLinkOpen, value: "[["})
			i += 2
			rest := line[i:]
			j := strings.Index(rest, "]]")
			if j == -1 {
				tokens = append(tokens, token{typ: tokenText, value: rest})
				i = len(line)
				break
			}
			inner := rest[:j]
			label, link, exist := strings.Cut(inner, ">")
			if !exist {
				tokens = append(tokens, token{typ: tokenLinkText, value: inner})
			} else {
				tokens = append(tokens, token{typ: tokenLinkText, value: label})
				tokens = append(tokens, token{typ: tokenLinkSep, value: ">"})
				tokens = append(tokens, token{typ: tokenLinkURL, value: link})
			}
			tokens = append(tokens, token{typ: tokenLinkClose, value: "]]"})
			i += j + 2

		default:
			buf.WriteByte(line[i])
			i++
		}
	}

	flush()
	return tokens
}
