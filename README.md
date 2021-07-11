# go-infounit

[![Go Reference](https://pkg.go.dev/badge/github.com/tunabay/go-infounit.svg)](https://pkg.go.dev/github.com/tunabay/go-infounit)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat)](LICENSE)

go-infounit is a Go package providing three data types for units of information:

- `ByteCount` - number of bytes with both SI and binary prefixes
- `BitCount`  - number of bits with both SI and binary prefixes
- `BitRate`   - number of bits transferred or processed per unit of time

These values can be converted to human-readable string representations
using the standard `fmt.Printf` family functions:
```
import (
	"fmt"
	"github.com/tunabay/go-infounit"
)

func main() {
	size := infounit.Megabyte * 2500
	fmt.Printf("% s\n", size)    // "2.5 GB"
	fmt.Printf("% .1S\n", size)  // "2.3 GiB"
	fmt.Printf("%# .2s\n", size) // "2.50 gigabytes"
	fmt.Printf("%# .3S\n", size) // "2.328 gibibytes"
	fmt.Printf("%d\n", size)     // "2500000000"

	bits := infounit.Kilobit * 32
	fmt.Printf("% s\n", bits)    // "32 kbit"
	fmt.Printf("% S\n", bits)    // "31.25 Kibit"
	fmt.Printf("%# .2s\n", bits) // "32.00 kilobits"
	fmt.Printf("%# .3S\n", bits) // "31.250 kibibits"
	fmt.Printf("%d\n", bits)     // "32000"

	rate := infounit.GigabitPerSecond * 123.45
	fmt.Printf("% s\n", rate)    // "123.45 Gbit/s"
	fmt.Printf("% .3S\n", rate)  // "114.972 Gibit/s"
	fmt.Printf("%# .1s\n", rate) // "123.5 gigabits per second"
	fmt.Printf("%# .2S\n", rate) // "114.97 gibibits per second"
	fmt.Printf("% .1a\n", rate)  // "123.5 Gbps"
	fmt.Printf("% .2A\n", rate)  // "114.97 Gibps"
}
```
[Run in Go Playground](https://play.golang.org/p/aUvaP6JZXeV)

- `%s` is for SI prefixes and `%S` is for binary prefixes.
- `%a` and `%A` for `BitRate` use the non-standard abbreviation "bps".
- `' '`(space) flag puts a space between digits and the unit suffix.
- `'#'` flag uses long unit suffixes.

They also implement convenience methods for:

- Rounding to specified precision.
- Scanning from human-readable string representation using fmt.Scanf family functions.

## For more examples:

- Read the [documentation](http://godoc.org/github.com/tunabay/go-infounit).

## See also

- https://docs.microsoft.com/en-us/style-guide/a-z-word-list-term-collections/term-collections/bits-bytes-terms

## License

go-infounit is available under the MIT license. See the [LICENSE](LICENSE) file for more information.
