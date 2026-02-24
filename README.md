# Enumer [![GoDoc](https://godoc.org/github.com/cmmoran/enumer?status.svg)](https://godoc.org/github.com/cmmoran/enumer) [![Go Report Card](https://goreportcard.com/badge/github.com/cmmoran/enumer)](https://goreportcard.com/report/github.com/cmmoran/enumer) [![GitHub Release](https://img.shields.io/github/release/dmarkham/enumer.svg)](https://github.com/cmmoran/enumer/releases)

Enumer is a tool to generate Go code that adds useful methods to Go enums (constants with a specific type).
It started as a fork of [Rob Pike’s Stringer tool](https://godoc.org/golang.org/x/tools/cmd/stringer)
maintained by [Álvaro López Espinosa](https://github.com/alvaroloes/enumer). 
This was again forked as (https://github.com/dmarkham/enumer) picking up where Álvaro left off.
And yet again this was forked here as (https://github.com/cmmoran/enumer) for my personal project needs. My intention is to submit PRs for my changes but current time constraints require this current fork in lieu of waiting through a potentially long PR process.

```
$ ./enumer --help
Enumer is a tool to generate Go code that adds useful methods to Go enums (constants with a specific type).
Usage of ./enumer:
        Enumer [flags] -type T [directory]
        Enumer [flags] -type T files... # Must be a single package
For more information, see:
        https://godoc.org/github.com/cmmoran/enumer
Flags:
  -addprefix string
        transform each item name by adding a prefix. Default: ""
  -comment value
        comments to include in generated code, can repeat. Default: ""
  -gqlgen
        if true, GraphQL marshaling methods for gqlgen will be generated. Default: false
  -json
        if true, json marshaling methods will be generated. Default: false
  -linecomment
        use line comment text as printed text when present
  -output string
        output file name; default srcdir/<type>_string.go
  -sql
        if true, the Scanner and Valuer interface will be implemented.
  -sql:int
        if true, the Scanner and Valuer interface will be implemented and Value() returns int
  -text
        if true, text marshaling methods will be generated. Default: false
  -transform string
        enum item name transformation method. Default: noop (default "noop")
  -trimprefix string
        transform each item name by removing a prefix. Default: ""
  -type string
        comma-separated list of type names; must be set
  -values
    	if true, alternative string values method will be generated. Default: false
  -yaml
        if true, yaml marshaling methods will be generated. Default: false
```


## Generated functions and methods

When Enumer is applied to a type, it will generate:

- The following basic methods/functions:

  - Method `String()`: returns the string representation of the enum value. This makes the enum conform
    the `Stringer` interface, so whenever you print an enum value, you'll get the string name instead of a number.
  - Function `<Type>String(s string)`: returns the enum value from its string representation. This is useful
    when you need to read enum values from command line arguments, from a configuration file, or
    from a REST API request... In short, from those places where using the real enum value (an integer) would
    be almost meaningless or hard to trace or use by a human. `s` string is Case Insensitive.
  - Function `<Type>Values()`: returns a slice with all the values of the enum
  - Function `<Type>Strings()`: returns a slice with all the Strings of the enum
  - Method `IsA<Type>()`: returns true only if the current value is among the values of the enum. Useful for validations.

- When the flag `json` is provided, two additional methods will be generated, `MarshalJSON()` and `UnmarshalJSON()`. These make
  the enum conform to the `json.Marshaler` and `json.Unmarshaler` interfaces. Very useful to use it in JSON APIs.
- When the flag `text` is provided, two additional methods will be generated, `MarshalText()` and `UnmarshalText()`. These make
  the enum conform to the `encoding.TextMarshaler` and `encoding.TextUnmarshaler` interfaces.
  **Note:** If you use your enum values as keys in a map and you encode the map as _JSON_, you need this flag set to true to properly
  convert the map keys to json (strings). If not, the numeric values will be used instead
- When the flag `yaml` is provided, two additional methods will be generated, `MarshalYAML()` and `UnmarshalYAML()`. These make
  the enum conform to the `gopkg.in/yaml.v2.Marshaler` and `gopkg.in/yaml.v2.Unmarshaler` interfaces.
- When the flag `sql` is provided, the methods for implementing the `Scanner` and `Valuer` interfaces.
  Useful when storing the enum in a database.
- When the flag `sql:int` is provided, the generated `Scanner` and `Valuer` interfaces are still implemented, but `Value()` returns the enum's integer value as `int64` instead of its string representation.


For example, if we have an enum type called `Pill`,

```go
type Pill int

const (
	Placebo Pill = iota
	Aspirin
	Ibuprofen
	Paracetamol
	Acetaminophen = Paracetamol
)
```

executing `enumer -type=Pill -json` will generate a new file with four basic methods and two extra for JSON:

```go
func (i Pill) String() string {
	//...
}

func PillString(s string) (Pill, error) {
	//...
}

func PillValues() []Pill {
	//...
}

func PillStrings() []string {
	//...
}

func (i Pill) IsAPill() bool {
	//...
}

func (i Pill) MarshalJSON() ([]byte, error) {
	//...
}

func (i *Pill) UnmarshalJSON(data []byte) error {
	//...
}
```

From now on, we can:

```go
// Convert any Pill value to string
var aspirinString string = Aspirin.String()
// (or use it in any place where a Stringer is accepted)
fmt.Println("I need ", Paracetamol) // Will print "I need Paracetamol"

// Convert a string with the enum name to the corresponding enum value
pill, err := PillString("Ibuprofen") // "ibuprofen" will also work.
if err != nil {
    fmt.Println("Unrecognized pill: ", err)
    return
}
// Now pill == Ibuprofen

// Get all the values of the string
allPills := PillValues()
fmt.Println(allPills) // Will print [Placebo Aspirin Ibuprofen Paracetamol]

// Check if a value belongs to the Pill enum values
var notAPill Pill = 42
if (notAPill.IsAPill()) {
	fmt.Println(notAPill, "is not a value of the Pill enum")
}

// Marshal/unmarshal to/from json strings, either directly or automatically when
// the enum is a field of a struct
pillJSON := Aspirin.MarshalJSON()
// Now pillJSON == `"Aspirin"`
```

The generated code is exactly the same as the Stringer tool plus the mentioned additions, so you can use
**Enumer** where you are already using **Stringer** without any code change.

## Parse aliases from enum comments

Enumer can extend the generated `<Type>String(s string)` parser with explicit aliases declared in enum comments.

This is useful when an enum has one canonical string form for output, but should accept additional input spellings from CLI flags, config files, JSON payloads, or legacy APIs.

Aliases are declared with a comment directive on the enum constant:

```go
type Status int

const (
	StatusOpen Status = iota //enumer:alias=opened,o
	StatusClosed             //enumer:alias=closed,c
)
```

With the definition above, all of the following resolve successfully:

```go
status, _ := StatusString("StatusOpen")
status, _ = StatusString("statusopen")
status, _ = StatusString("opened")
status, _ = StatusString("O")
```

### Alias rules

- Aliases are parse-only metadata. They affect `<Type>String(s string)` and any generated unmarshal methods that call it.
- `String()` still returns the canonical enum string only.
- `<Type>Strings()` still returns canonical strings only.
- Aliases are matched the same way canonical strings are matched: exact first, then case-insensitive via lowercase lookup.
- Aliases are explicit. They are not transformed by `-transform`, `-trimprefix`, `-trimsuffix`, `-addprefix`, or `-addsuffix`.

### Directive syntax

The directive format is:

```go
//enumer:alias=value0,value1,value2
```

Rules:

- Aliases are comma-separated.
- Surrounding whitespace is trimmed.
- Empty aliases are invalid.
- Alias directives may be placed in doc comments or line comments attached to the enum constant.

Example with a doc comment:

```go
const (
	//enumer:alias=running,wip
	JobStateInProgress JobState = iota
)
```

### Interaction with `-linecomment`

When `-linecomment` is enabled, alias directives are treated as generator metadata and are ignored when deriving the canonical display string.

Example:

```go
const (
	StatusOpen Status = iota // Open //enumer:alias=opened,o
	StatusClosed             //enumer:alias=closed,c
)
```

With `-linecomment`:

- `StatusOpen.String()` returns `"Open"`
- `StatusClosed.String()` continues to use the identifier-derived canonical string because its line comment contains only alias metadata
- `StatusString("opened")` and `StatusString("o")` still resolve to `StatusOpen`

### Collision rules

Enumer rejects any alias that collides with another parse token after normalization.

The following collisions are invalid:

- alias vs alias on different constants
- alias vs canonical name of another constant
- alias vs canonical lowercase name of another constant
- duplicate aliases on the same constant after normalization

Example:

```go
const (
	StatusOpen Status = iota //enumer:alias=closed
	StatusClosed
)
```

This fails because `"closed"` would resolve to two different enum values after normalization.

Enumer enforces this in two ways:

- The generator validates parse-token collisions and fails with a descriptive error.
- The generated `map[string]T` contains explicit alias entries, so duplicate keys still fail at compile time if validation ever misses a case.

## SQL and `sql:int`

The `sql` and `sql:int` flags both generate implementations of `database/sql.Scanner` and `database/sql/driver.Valuer`.

The difference is the representation used by `Value()`:

- `-sql` stores the enum as its canonical string form
- `-sql:int` stores the enum as its numeric value converted to `int64`

Conceptually, the generated `Value()` methods look like this:

```go
// with -sql
func (i MyEnum) Value() (driver.Value, error) {
	return i.String(), nil
}

// with -sql:int
func (i MyEnum) Value() (driver.Value, error) {
	return int64(i), nil
}
```

The generated `Scan` method is flexible in both modes:

- it accepts numeric database values such as `int64`, `int32`, `int`, `float64`, and `float32`
- it accepts strings and byte slices
- if given a string-like value, it first tries to parse it as an integer, then falls back to `<Type>String(...)`

That means both modes can read either integer-like or string-like SQL values in many cases. The main distinction is what Enumer writes back through `Value()`.

Use `-sql` when:

- your database column stores symbolic enum names
- you want values in the database to be human-readable
- you want the database representation to stay aligned with `String()`

Use `-sql:int` when:

- your schema stores enum values as integers
- you need compatibility with an existing numeric column
- you want the database representation to match the underlying Go enum value directly

Example:

```go
//go:generate go run github.com/cmmoran/enumer -type=Status -sql
```

or

```go
//go:generate go run github.com/cmmoran/enumer -type=Status -sql:int
```

In practice, choose one or the other based on the database representation you want to persist.

## Transforming the string representation of the enum value

By default, Enumer uses the same name of the enum value for generating the string representation (usually CamelCase in Go).

```go
type MyType int

 ...

name := MyTypeValue.String() // name => "MyTypeValue"
```

Sometimes you need to use some other string representation format than CamelCase (i.e. in JSON).

To transform it from CamelCase to another format, you can use the `transform` flag.

For example, the command `enumer -type=MyType -json -transform=snake` would generate the following string representation:

```go
name := MyTypeValue.String() // name => "my_type_value"
```

**Note**: The transformation only works from CamelCase to snake_case or kebab-case, not the other way around.

### Transformers

- snake
- snake-upper
- kebab
- kebab-upper
- lower (lowercase)
- upper (UPPERCASE)
- title (TitleCase)
- title-lower (titleCase)
- first (Use first character of string)
- first-lower (same as first only lower case)
- first-upper (same as first only upper case)
- whitespace
- map:K0=V0,K1=V1,K2=V2... (map specific enum identifiers to explicit output values)

### `map:` transform

The `map:` transform lets you define exact output strings for specific enum constants instead of applying a generic case conversion.

Syntax:

```text
-transform='map:Key0=Value0,Key1=Value1,Key2=Value2'
```

Each key must be the enum constant name after any `-trimprefix` and `-trimsuffix` processing, but before any `-addprefix` or `-addsuffix` processing.

Example:

```go
type MapValue int

const (
	Male MapValue = iota
	Female
	Unknown
)
```

Generation:

```bash
enumer -type=MapValue -transform='map:Male=XY,Female=XX,Unknown=XX|XY'
```

Result:

```go
Male.String()    // "XY"
Female.String()  // "XX"
Unknown.String() // "XX|XY"
```

Rules and behavior:

- Keys must match the enum names exactly at the point the map transform runs.
- Values are used exactly as written.
- If a key is not present in the map, that enum name is left unchanged.
- `map:` is a transform of the canonical output string. It affects `String()`, `<Type>Strings()`, marshaling output, and the canonical names accepted by `<Type>String(...)`.
- Alias directives are separate. Use `//enumer:alias=...` if you want extra accepted parse inputs without changing canonical output.

Ordering with other name transforms:

- `-trimprefix` and `-trimsuffix` run before `map:`
- `-addprefix` and `-addsuffix` run after `map:`

That means this:

```bash
enumer -type=Status -trimprefix=Status -transform='map:Open=opened,Closed=closed'
```

expects the map keys `Open` and `Closed`, not `StatusOpen` and `StatusClosed`.

Limitations:

- The `map:` syntax does not provide escaping for `,` or `=`, so mapped values should not contain those characters.
- `map:` is intended for explicit one-to-one renaming. If you want a general naming convention such as `snake` or `kebab`, use the corresponding built-in transform instead.

## How to use

For a module-aware repo with `enumer` in the `go.mod` file, generation can be called by adding the following to a `.go` source file:

```golang
//go:generate go run github.com/cmmoran/enumer -type=YOURTYPE
```

There are five optional generation flags: `json`, `text`, `yaml`, `sql`, and `sql:int`. You can use any combination that makes sense for your enum and storage format.

For enum string representation transformation the `transform` and `trimprefix` flags
were added (i.e. `enumer -type=MyType -json -transform=snake`).
Possible transform values are listed above in the [transformers](#transformers) section.
The default value for `transform` flag is `noop` which means no transformation will be performed.

Alias directives are independent from the transform flags. The transform flags define the canonical output string, while `//enumer:alias=...` adds extra accepted input strings for parsing only.

If a prefix is provided via the `trimprefix` flag, it will be trimmed from the start of each name (before
it is transformed). If a name doesn't have the prefix it will be passed unchanged.

If a prefix is provided via the `addprefix` flag, it will be added to the start of each name (after trimming and after transforming).

The boolean flag `values` will additionally create an alternative string values method `Values() []string` to fullfill the `EnumValues` interface of [ent](https://entgo.io/docs/schema-fields/#enum-fields).

## Inspiring projects

- [Álvaro López Espinosa](https://github.com/alvaroloes/enumer)
- [Stringer](https://godoc.org/golang.org/x/tools/cmd/stringer)
- [jsonenums](https://github.com/campoy/jsonenums)

## Credit

This repository is a complete copy of [github.com/dmarkham/enumer](https://github.com/dmarkham/enumer) with modifications to suit my personal project needs. 
