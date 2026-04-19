# libpukiwiki

f-lab の PukiWiki 操作用 Go ライブラリ

## Install

```bash
go get github.com/moriT958/libpukiwiki
```

## Usage

[Examples](./examples)

## Converter

Work in progress...

## md2pw

`md2pw` は Markdown を PukiWiki 記法へ変換する CLI ツールです。

### インストール

`go install` を使う場合:

```bash
go install github.com/moriT958/libpukiwiki/cmd/md2pw@latest
```

リリース済みバイナリを使う場合:

```bash
curl -sSfL https://raw.githubusercontent.com/moriT958/libpukiwiki/main/install.sh | sh
```

### 使い方

基本形:

```bash
md2pw [options] [<file.md>|-]
```

オプション:

- `-o <path>`: 変換結果を標準出力ではなくファイルへ書き出します

#### Example

標準入力から読む:

```bash
echo '# heading' | md2pw
cat input.md | md2pw
md2pw -
```

ファイルを読む:

```bash
md2pw input.md
md2pw -o output.txt input.md
```

リダイレクトで渡す:

```bash
md2pw < input.md
```

### 対応記法

実装上、現在変換対象になっているのは以下です。

- 見出し
- リスト
- コードブロック
- 太字
- リンク
- テーブル

#### 見出し

`#` から `###` までを PukiWiki の見出しへ変換します。`####` 以降はそのまま残ります。

```markdown
# H1

## H2

### H3

#### H4
```

```text
* H1
** H2
*** H3
#### H4
```

#### リスト

順不同リストは `-`、番号付きリストは `+` に変換します。ネストは 3 段までサポートされ、4 段目以降は 3 段として扱われます。

```markdown
- item1
  - nested1
  - nested2

1. ordered1
   1. nested1
   2. nested2
```

```text
-item1
--nested1
--nested2

+ordered1
++nested1
++nested2
```

#### コードブロック

フェンス付きコードブロックを、各行の先頭に半角スペース 2 つを付けた PukiWiki 形式へ変換します。開始・終了のフェンス行は出力しません。言語指定は無視されます。

````markdown
```go
func main() {}
```
````

````

```text
  func main() {}
````

#### 太字

`**text**` を `''text''` へ変換します。`*italic*` は変換しません。

```markdown
This is **bold** text.
```

```text
This is ''bold'' text.
```

#### リンク

インラインリンクを `[[text>url]]` へ変換します。

```markdown
[example](https://example.com)
```

```text
[[example>https://example.com]]
```

#### テーブル

Markdown テーブルを PukiWiki テーブルへ変換します。ヘッダー行には `~` を付け、区切り行は出力しません。

```markdown
| Column1 | Column2 |
| ------- | ------- |
| A       | B       |
```

```text
|~ Column1 |~ Column2 |
| A | B |
```
