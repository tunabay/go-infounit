// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

/*
Package infounit provides information unit data types that can be formatted into
human-readable string representations. The following three data types are
implemented:

	ByteCount  non-negative number of bytes
	BitCount   non-negative number of bits
	BitRate    number of bits per unit of time

These types can be formatted into and scanned from human-readable string
representations with both SI and binary prefixes using the standard Printf and
Scanf family functions in the package fmt.

	import "github.com/tunabay/go-infounit"

	size := infounit.Megabyte * 2500
	fmt.Printf("% s\n", size)     // "2.5 GB"
	fmt.Printf("% .1S\n", size)   // "2.3 GiB"
	fmt.Printf("%# .2s\n", size)  // "2.50 gigabytes"
	fmt.Printf("%# .3S\n", size)  // "2.328 gibibytes"
	fmt.Printf("%d\n", size)      // "2500000000"

	bits := infounit.Kilobit * 32
	fmt.Printf("% s\n", bits)     // "32 kbit"
	fmt.Printf("% S\n", bits)     // "31.25 Kibit"
	fmt.Printf("%# .2s\n", bits)  // "32.00 kilobits"
	fmt.Printf("%# .3S\n", bits)  // "31.250 kibibits"
	fmt.Printf("%d\n", bits)      // "32000"

	rate := infounit.GigabitPerSecond * 123.45
	fmt.Printf("% s\n", rate)     // "123.45 Gbit/s"
	fmt.Printf("% .3S\n", rate)   // "114.972 Gibit/s"
	fmt.Printf("%# .1s\n", rate)  // "123.5 gigabits per second"
	fmt.Printf("%# .2S\n", rate)  // "114.97 gibibits per second"
	fmt.Printf("% .1a\n", rate)   // "123.5 Gbps"
	fmt.Printf("% .2A\n", rate)   // "114.97 Gibps"

Each data type also provides methods such as scanning and unit conversion.

ByteCount

ByteCount represents a number of bytes. It is a suitable type for storing file
size, transferred data amount, etc. It is internally uint64. Thus new values can
be declared in the usual ways and it is possible to convert between ByteCount
and uint64 values.

	var x infounit.ByteCount      // zero value is 0 bytes
	y := infounit.ByteCount(200)  // declared with value 200 bytes
	p := uint64(y)                // convert to uint64 with value 200

The common values for both SI and binary prefixes are defined as constants.
Values with units can be declared by multiplying by them in a similar way to
the Duration in the package time:

	infounit.ByteCount(30) * infounit.Kilobyte  // 30 kB
	infounit.ByteCount(50) * infounit.Kibibyte  // 50 KiB
	infounit.ByteCount(70) * infounit.Terabyte  // 70 TB
	infounit.ByteCount(90) * infounit.Pebibyte  // 90 PiB

Unit conversion is also possible by division:

	int(infounit.Mebibyte / infounit.Kibibyte)  // 1024

See the Constants section below for the complete list of defined constants.

ByteCount values can be flexibly formatted using the standard Printf family
functions in the package fmt, fmt.Printf, fmt.Fprintf, fmt.Sprintf,
fmt.Errorf, and functions deriverd from them.

For ByteCount type, two custom 'verbs' are implemented:

	%s	human readable format with SI prefix
	%S	human readable format with binary prefix

By using these, the same value can be formatted with both SI and binary
prefixes:

	x := infounit.Kilobyte * 100
	fmt.Printf("%.1s", x)        // 100.0kB
	fmt.Printf("%.1S", x)        //  97.7KiB
	y := infounit.Mebibyte * 100
	fmt.Printf("%.1s", y)        // 102.4MB
	fmt.Printf("%.1S", y)        // 100.0MiB

They also support flags ' '(space) and '#', width and precision, so the format
can be changed flexibly:

	z := infounit.ByteCount(987654321)
	fmt.Printf("%.1s",  z)       // 987.7MB
	fmt.Printf("% .1s", z)       // 987.7 MB
	fmt.Printf("%# .1s", z)      // 987.7 megabytes
	fmt.Printf("%.2S", z)        // 941.90MiB
	fmt.Printf("%# .3S", z)      // 941.901 mebibytes

A ' '(space) flag puts a spece between digits and the unit. A '#' flag uses an
alternative long unit name. See the Format method documentation bellow for
details on all supported verbs and flags.

ByteCount type also supports the standard Scanf family functions in the
package fmt, fmt.Scanf, fmt.Fscanf(), and fmt.Sscanf(). Human-readable string
representations with SI and binary prefixes can be "scanned" as ByteCount
values.

	var file string
	var size infounit.ByteCount
	_, _ = fmt.Sscanf("%s is %s", "test.png is 5 kB", &file, &size)
	fmt.Println(file, size)  // test.png 5.0 kB

Note that unlike Printf, the %s verb properly scans units with both SI and
binary prefixes. So it is usually fine to use only the %s verb to scan.
The %S verb treats the SI prefix representation as binary prefixes. See the
Scan method documentation bellow for details.

BitCount

BitCount represents a number of bits. It is internally uint64. Functions and
usage are almost the same as ByteCount, except that it represents a number of
bits instead of bytes.

BitRate

BitRate represents a number of bits that are transferred or processed per unit
of time. It is a suitable type for storing data transfer speed, processing
speed, etc. It is internally float64.
*/
package infounit
