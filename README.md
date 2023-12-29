# SCADFormat

SCADFormat is a source code formatter / beautifier for [OpenSCAD](https://openscad.org/).

SCADFormat is, shall we say, "opinionated" in the way that it formats OpenSCAD code. In other words, there are no configuration options that alter the way code is formatted. That's not because I feel strongly that OpenSCAD code should be formatted a certain way - it's just that I haven't had time to implement options.

## Installation

The easiest way to install is to download one of the pre-built binary releases:

https://github.com/hugheaves/scadformat/releases

Download the zip file for your operating system (windows.zip, linux,zip, macos.zip) and extract the contents.

## Usage

SCADFormat is a command line tool.

### Specifying a filename

SCADFormat can be run directly on a file by specifying the filename on the command line:

```bash
scadformat my-source.scad
```
```
INFO	formatting file my-source.scad
```
In this mode, SCADFormat will overwrite the existing code with the formatted version. Note that SCADFormat creates a backup of the original file (with a .scadbak extension) before overwriting it.

### Read from stdin / write to stdout

SCADFormat can also read from stdin and write to stdout as follows:

```bash
scadformat <my-source.scad >my-source-formatted.scad
```

### Format all .scad recursively

Format all .scad files in the directory "." recursively. Note that if the scadformat command is not in your search PATH, you'll need to specify the full path to `scadformat` after the `-exec-` option. (e.g. `-exec $HOME\scasformat\scadformat`) 

```bash
find . -type f -name "*.scad" -exec scadformat "{}" \;
```

If you are ok with the result, you can delete all backup files (.scadbak)
```bash
find $directory -type f -name "*.scadbak" -exec rm "{}" \;
```

## Building

### Install Prerequisites

SCADFormat is written in Go, and uses the ANTLR v4 parser generator. You'll need to install both tools to build the sourcecode.

See https://go.dev/doc/install to install Go (v1.21 or later) and https://github.com/antlr/antlr4/blob/master/doc/getting-started.md to install ANTLR.

```bash
python3 -m venv venv
. ./venv/bin/activate
pip install antlr4-tools
```

After installation, run the "antlr4" command to verify the command is available in your search path:

```bash
antlr4
```
should display
```
ANTLR Parser Generator  Version 4.13.1
-o ___              specify output directory where all output is generated
-lib ___            specify location of grammars, tokens files
-atn                generate rule augmented transition network diagrams
...
```

### Checkout Source
Checkout the source code and "cd" into the "scadformat" directory:
```bash
git clone https://github.com/hugheaves/scadformat
cd scadformat
```

#### Build with Make
If you have GNU Make (or a compatible make utility installed), you can build the program just by running the "make" command:
```bash
make
```

You should see output similar to the following:
```
go generate ./...
patching file internal/parser/openscad_base_visitor.go
go test ./...
?   	github.com/hugheaves/scadformat	[no test files]
?   	github.com/hugheaves/scadformat/cmd	[no test files]
?   	github.com/hugheaves/scadformat/internal/logutil	[no test files]
?   	github.com/hugheaves/scadformat/internal/parser	[no test files]
ok  	github.com/hugheaves/scadformat/internal/formatter	(cached)
go build cmd/scadformat.go
```

#### Build without Make
If you don't have make installed (i.e. on Windows), you can still build the program by running the necessary commands manually:

Generate ANTLR parser
```bash
go generate ./...
```

Build the executable
```bash
go build -o scadformat cmd/scadformat.go
```

run tests
```bash
go test -v ./...
```


run on a file
```bash
./scadformat ./internal/formatter/testdata/solo_adapter.scad
```
