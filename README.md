# md2pw

Convert markdown to pukiwiki notation. (markdown -> pukiwiki)

**Supported Notations**

- [x] Headers
- [x] List
- [x] Codeblock
- [x] Bold
- [x] Link
- [x] Table

## Install

**cURL**

```bash
curl -sSfL https://raw.githubusercontent.com/moriT958/libpukiwiki/main/install.sh | sh
```

## Example

Example usage

### from pipe

```bash
echo '# heading' | md2pw
cat input.md | md2pw
```

### Redirect

```bash
md2pw < input.md
```

### Explicit Stdin

```bash
md2pw -
```

### Read file

```bash
md2pw input.md
md2pw -o output.txt input.md
```

## pukiWiki notaion coverage

Supported notaitions.

### Headers

pukiwiki は 3 段階ある。

**PukiWiki**

```text
* H1
** H2
*** H3
```

**Markdown**

```markdown
# H1

## H2

### H3
```

### List

インデントは 3 Level まで対応

#### Ordered

同じ

**PukiWiki / Markdown**

```text
- list1
- list2
- list3
```

#### Unordered

**PukiWiki**

```text
+ ordered1
+ ordered2
+ ordered3
```

**Mardown**

```markdown
1. ordered1
2. ordered2
3. ordered3
```

### Code Block

**PukiWiki**

```text
  this is sample code.
  need 2 spaces
```

**Markdown**

```markdown
\`\`\`(filetype)
this is sample code.
\`\`\`
```

### Bold

**PukiWiki**

```text
''text''
```

**Markdown**

```markdown
**text**
```

### Link

**PukiWiki**

```text
[[this is link>https://example.com]]
```

**Markdown**

```markdown
[this is link](https://example.com)
```

### Table

**Pukiwiki**

```text
|~ Column1 |~ Column2 |~ Column3 |
| Item1.1  | Item2.1  | Item3.1  |
| Item1.2  | Item2.2  | Item3.2  |
| Item1.2  | Item2.2  | Item3.2  |
```

**Markdown**

```markdown
| Column1 | Column2 | Column3 |
| ------- | ------- | ------- |
| Item1.1 | Item2.1 | Item3.1 |
| Item1.2 | Item2.2 | Item3.2 |
| Item1.2 | Item2.2 | Item3.2 |
```

## Development

- deps
  - Go
  - golangci-lint
  - task

- dev commands
  - Run command: `task run`
    - with args: `task run -- <args>`
  - Build: `task build`
  - Test: `task test`
  - Lint: `task lint`
  - Release check: `task release`
  - Deploy: `task deploy -- --tag <v0.0.0>`
