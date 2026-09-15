# PortuGo

[Português (Brasil)](README.pt-BR.md) | English

PortuGo runs Portugol programs in your terminal on Linux and macOS. You can
follow VisuAlg-based courses, solve exercises, or teach programming in
Portuguese using your preferred editor. It reads `.alg` files in the
VisuAlg 3.x dialect without needing Windows, a virtual machine, or Wine.

## Get started

### Download a release

Published versions are available from [GitHub Releases](https://github.com/ncode/PortuGo/releases).
Once a release is published, its **Assets** include the `portugo` executable
for Linux, macOS, and Windows. Downloading an executable does not require Go.

Choose the archive matching your computer:

| System | Archive name ends with |
| --- | --- |
| Linux, 64-bit Intel/AMD | `linux-amd64.tar.gz` |
| Linux, ARM64 | `linux-arm64.tar.gz` |
| macOS, Intel | `darwin-amd64.tar.gz` |
| macOS, Apple Silicon | `darwin-arm64.tar.gz` |
| Windows, 64-bit Intel/AMD | `windows-amd64.zip` |
| Windows, ARM64 | `windows-arm64.zip` |

Extract the archive and open a terminal in the extracted folder. Run
`./portugo --version` to check the installed version, then use
`./portugo run your-program.alg` to run an exercise. In Windows PowerShell,
use `.\portugo.exe` in place of `./portugo`. Each release also includes
`SHA256SUMS` for checking the downloaded archive's integrity.

### Build from source

You'll need [Git](https://git-scm.com/install/) and Go 1.27 or newer. The
[official Go installation guide](https://go.dev/doc/install) has instructions
for Linux and macOS. Go builds the command-line tool; you don't need to know
the language to use PortuGo.

Open a terminal and check that both tools are available:

```sh
git --version
go version
```

If either command isn't found, install the missing tool and open a new
terminal before continuing.

Then download the project, build it, and run your first example:

```sh
git clone https://github.com/ncode/PortuGo.git
cd PortuGo
go build -o portugo ./cmd/portugo
./portugo run examples/hello.alg
```

You should see:

```text
Ola, Portugol!
```

`portugo` is the executable you just built. The `./` tells your terminal to
find it in the current folder. The commands below assume you're still in the
`PortuGo` folder.

## Write something of your own

Open your text editor and save this as `boas-vindas.alg` in the
project folder, as a plain-text UTF-8 file:

```portugol
algoritmo "boas-vindas"
var
  nome: caractere
inicio
  escreva("Como voce se chama? ")
  leia(nome)
  escreval("Ola, ", nome, "! Vamos programar?")
fimalgoritmo
```

Run it:

```sh
./portugo run boas-vindas.alg
```

Type your name and press Enter. The program will greet you by name.
`leia` reads your answer, `escreva` prints the question, and `escreval` prints
the greeting with a newline. PortuGo also echoes input after reading it, so
seeing your answer repeated is expected.

Try changing the greeting and running the file again. Each `.alg` file is a
complete program, starting with `algoritmo` and ending with `fimalgoritmo`.

## Your everyday commands

| What you'd like to do | Command |
| --- | --- |
| Run a program | `./portugo run boas-vindas.alg` |
| Check syntax and types without running it | `./portugo check boas-vindas.alg` |
| Preview formatted code in the terminal | `./portugo fmt boas-vindas.alg` |
| Format and save the file | `./portugo fmt -w boas-vindas.alg` |
| Check whether a file is already formatted | `./portugo fmt --check boas-vindas.alg` |
| Start an interactive session | `./portugo repl` |
| See help for a command | `./portugo run -h` |
| Show the executable's version | `./portugo --version` |

Put options such as `-w` and `--check` **before** the filename.
`check` is silent when it succeeds. If something needs attention, diagnostics
point to a filename, line, and column so you can find it in your editor.
`fmt --check` exits with status 1 when formatting is needed.

In the interactive session, enter a **complete program**. The line containing
`fimalgoritmo` runs it; if it uses `leia`, enter the requested values next.
Type `:sair` on its own line to leave the session.

When experimenting with loops, you can limit the number of execution steps:

```sh
./portugo run --max-steps 10000 boas-vindas.alg
```

The program stops with a diagnostic if it reaches that limit. You can also
press Ctrl+C to stop a running command.

## Keep exploring

The [examples folder](examples/) is a good place to find your next exercise:

| Example | What to explore |
| --- | --- |
| [Hello](examples/hello.alg) | A first complete program |
| [Factorial](examples/fatorial.alg) | Input, variables, and a `para` loop; try entering `5` |
| [Console input](examples/console_input.alg) | Reading a name and an age, one input line at a time |
| [Logical conditions](examples/logical_conditions.alg) | Comparisons and logical expressions |
| [Vectors](examples/vector_bounds.alg) | Declaring and indexing vectors |

Run any example the same way:

```sh
./portugo run examples/fatorial.alg
```

Already have exercises saved from VisuAlg? Try them with
`./portugo run path/to/exercise.alg`. PortuGo reads UTF-8 and Windows-1252
source files, including the encoding commonly used by VisuAlg.

## Compatibility and project status

PortuGo targets the VisuAlg 3.x dialect, with behavior checked against
VisuAlg 3.0.7. It supports variables, conditions, loops, procedures, functions,
vectors, and numeric and text built-ins.

The [conformance report](docs/quality-acceptance-2026-09-14.md) records passing
checks for all accepted programs in the reference corpus. Some rejected-program
cases and behavior outside that corpus remain unverified, so these results
do not guarantee compatibility with every VisuAlg program.

This is a command-line interpreter. It doesn't recreate the VisuAlg desktop
interface or graphical debugger. Portugol Studio and other Portugol dialects
use different syntax and aren't the current target.

The [language reference](docs/language.md) describes supported syntax,
behavior, and known limitations. The [changelog](CHANGELOG.md) tracks changes.

## Questions, ideas, and contributions

Found a confusing instruction or an exercise that doesn't behave as expected?
[Open an issue](https://github.com/ncode/PortuGo/issues) with a small `.alg`
example, what you expected, and what happened. Questions and documentation
improvements are welcome too.

If you'd like to work on the implementation, start with the
[development guide](docs/development.md) and [contributor instructions](AGENTS.md).
You can run the test suite from the project folder with:

```sh
go test ./...
```
