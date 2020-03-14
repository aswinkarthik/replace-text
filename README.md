# replace-text

<a href="https://github.com/aswinkarthik/replace-text/actions"><img alt="GitHub Actions status" src="https://github.com/aswinkarthik/replace-text/workflows/Go/badge.svg"></a>

Find & Replace multiple text in files.

## Why?

Currently, the best tool for finding and replacing text efficiently is `sed` (and tools similar to that). This tool started as a learning project on using [Tries](https://en.wikipedia.org/wiki/Trie) datastructure as finite state machines to simulate the behavior of replacing multiple texts simultaneously like `sed`. It does not support regexes (yet?). The different find & replace patterns can be input as a JSON file.

## Usage

```bash
NAME:
   replace-text - Find & Replace multiple texts in files

USAGE:
   replace-text [global options] [PATH ...]

GLOBAL OPTIONS:
   --patterns-file value, -p value  Load find & replace patterns from a JSON file [$PATTERNS_FILE, $REPLACE_TEXT_PATTERNS_FILE]
   --help, -h                       show help (default: false)
```

## Examples

```bash
## Build before doing this

./replace-text -p examples/patterns.json examples/input1.txt examples/input1.txt
```

```bash
# Patterns file

$ cat examples/patterns.json
{
   "key1": "value1",
   "key2": "value2"
}
```

## Development

Clone the repository

**To run tests**

```bash
./scripts/test
```

**To Lint**

```bash
./scripts/ling
```

**To build**

```bash
./scripts/build
```
