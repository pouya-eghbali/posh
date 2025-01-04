# PoSH

PoSH is a modern language designed to address the complexities and limitations
of traditional bash scripting for large automation tasks. By providing better
error handling, flexibility, and portability, PoSH offers developers a powerful
tool to write maintainable and robust scripts that can be compiled and shipped
as a single binary.

> [!NOTE]
> PoSH is in active development. Features and syntax are subject to change. Stay
> tuned for updates!

## Why PoSH?

While bash is great for small, interactive tasks and quick commands, it often
becomes unwieldy and error-prone as scripts grow in size and complexity. PoSH
aims to fix this by:

- **Reducing Errors:** Strong typing and immutability help catch issues at
  compile-time instead of runtime.
- **Improving Portability:** Scripts can be compiled into standalone binaries,
  removing dependencies on interpreters or environment setups.
- **Enhancing Readability:** PoSH syntax is clean, modern, and intuitive, making
  scripts easier to understand and maintain.
- **Maintaining Compatibility:** PoSH integrates seamlessly with existing Unix
  tools and pipelines.

## Key Features

- **Modern Syntax:**

```posh
fn main(name string, age int) {
  message = io.Format(
    "Hello, %s! You are %d years old!",
    name,
    age
  )

  result =
    message
      | echo()
      | tr("[:lower:]", "[:upper:]")
      | lolcat(-f)

  io.Line(result)
}
```

- **Compilation:** Write your scripts once and compile them into portable
  binaries that can run anywhere.
- **Pipeline Support:** Unix-style pipelines (`|`) allow for chaining operations
  seamlessly.
- **Safety by Design:** Features like immutability and strong typing reduce
  common scripting errors.

## Installation

To get started with PoSH:

1. **Download the PoSH Compiler:**

```bash
go install github.com/pouya-eghbali/posh
```

2. **Verify Installation:**

```bash
posh --version
```

## Getting Started

1. Create a PoSH script:

```posh
fn main() {
  io.Line("Hello, World!")
}
```

2. Compile the script:

```bash
posh -i script.posh -o script
```

3. Run the compiled binary:

```bash
./script
```

## Examples

> [!NOTE]
> PoSH is currently under active development. The following examples are
> conceptual and may not work as expected in the current version.
> See [./examples](./examples) for working examples.

### Flag Parsing

```posh
fn main(name string) {
  io.Line(name)
}
```

To pass the `name` parameter to the compiled binary, use the following command:

```bash
./compiled --name "Pouya"
```

### File Manipulation

```posh
fn main(filename string) {
    const content = io.ReadFile(filename)
    const processed = content | tr("foo", "bar")
    io.WriteFile("output.txt", processed)
}
```

### API Calls

```posh
fn main() {
    const response = http.Get("https://api.example.com/data")
    io.Line(response.Body | json.PrettyPrint())
}
```

## TODO

- [x] Lexer, parser, and code generator
- [x] Proper scope management
- [x] Imports
- [x] Control statements (if/elif/else)
- [x] Add loops
- [ ] Add arrays, hashmaps
- [ ] Pipe-friendly functions
- [ ] Make a syntax diagram
- [ ] Proper type tracking
- [ ] Dead code elimination
- [ ] Standard library (`io`, `exec`, `http`, `os`...)
- [ ] Commands (Improved CLI args handling)
- [ ] Source mapping for runtime errors
- [ ] Proper error handling (internals)
- [ ] Semantic checks
- [ ] Fix all in-code TODOs
- [ ] Unit tests
- [ ] Interactive shell (REPL)
