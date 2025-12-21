# grab
`grab` is a grep-like command-line tool written in Go.

---

## Installation

### Build from source
```sh
git clone "https://github.com/albertcheong/grab.git"
cd grab
make build
```

The binary will be available at
```sh
bin/grab
```

---

## Usage

### From files
```sh
grab [options] <pattern> [file1 file2 ...]
```

### From standard input (pipe)
```sh
cat dummy.txt | grab <pattern>
```

---

## Exit Codes
| Exit Code | Meaning |
|:-|:-|
|0|**Success**; Success (at least one match found) |
|1|**No match**; No matches found |
|2|**Error**; such as invalid options or a non existing input file|

---

### License
MIT License. See [LICENSE](LICENSE) for details.