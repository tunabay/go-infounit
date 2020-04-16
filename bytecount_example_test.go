// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package infounit_test

import (
	"fmt"

	"github.com/tunabay/go-infounit"
)

//
func ExampleByteCount_Convert_siPrefix() {
	x := infounit.ByteCount(123456789)

	fmt.Printf("%f\n", x.Convert(infounit.Kilobyte))
	fmt.Printf("%f\n", x.Convert(infounit.Megabyte))
	fmt.Printf("%f\n", x.Convert(infounit.Gigabyte))
	fmt.Printf("%f\n", x.Convert(infounit.Terabyte))
	// Output:
	// 123456.789000
	// 123.456789
	// 0.123457
	// 0.000123
}

//
func ExampleByteCount_Convert_binaryPrefix() {
	x := infounit.ByteCount(123456789)

	fmt.Printf("%f\n", x.Convert(infounit.Kibibyte))
	fmt.Printf("%f\n", x.Convert(infounit.Mebibyte))
	fmt.Printf("%f\n", x.Convert(infounit.Gibibyte))
	fmt.Printf("%f\n", x.Convert(infounit.Tebibyte))
	// Output:
	// 120563.270508
	// 117.737569
	// 0.114978
	// 0.000112
}

//
func ExampleByteCount_ConvertRound_siPrefix() {
	x := infounit.ByteCount(123456789)

	fmt.Printf("%f\n", x.ConvertRound(infounit.Kilobyte, 2))
	fmt.Printf("%f\n", x.ConvertRound(infounit.Kilobyte, 5))
	fmt.Printf("%f\n", x.ConvertRound(infounit.Megabyte, 0))
	fmt.Printf("%f\n", x.ConvertRound(infounit.Megabyte, 3))
	fmt.Printf("%f\n", x.ConvertRound(infounit.Terabyte, 6))
	// Output:
	// 123456.790000
	// 123456.789000
	// 123.000000
	// 123.457000
	// 0.000123
}

//
func ExampleByteCount_ConvertRound_binaryPrefix() {
	x := infounit.ByteCount(123456789)

	fmt.Printf("%f\n", x.ConvertRound(infounit.Kibibyte, 2))
	fmt.Printf("%f\n", x.ConvertRound(infounit.Kibibyte, 5))
	fmt.Printf("%f\n", x.ConvertRound(infounit.Mebibyte, 0))
	fmt.Printf("%f\n", x.ConvertRound(infounit.Mebibyte, 3))
	fmt.Printf("%f\n", x.ConvertRound(infounit.Tebibyte, 6))
	// Output:
	// 120563.270000
	// 120563.270510
	// 118.000000
	// 117.738000
	// 0.000112
}

//
func ExampleByteCount_Format_errorf() {
	limit := infounit.Kilobyte * 100

	check := func(size infounit.ByteCount) error {
		if limit < size {
			return fmt.Errorf("too large: %v", size)
		}
		return nil
	}

	for i := 0; i < 8; i++ {
		size := infounit.ByteCount(1 << (6 * i))
		if err := check(size); err != nil {
			fmt.Println("error:", err)
			continue
		}
		fmt.Println("processed:", size)
	}
	// Output:
	// processed: 1 B
	// processed: 64 B
	// processed: 4.1 kB
	// error: too large: 262.1 kB
	// error: too large: 16.8 MB
	// error: too large: 1.1 GB
	// error: too large: 68.7 GB
	// error: too large: 4.4 TB
}

//
func ExampleByteCount_Format_printf() {
	x := infounit.Kilobyte * 100 // 100000 bytes

	// Default format
	fmt.Printf("%v\n", x)  // default format, same as "% .1s"
	fmt.Printf("%#v\n", x) // Go syntax format

	// SI prefix
	fmt.Printf("%s\n", x)      // default precision
	fmt.Printf("% s\n", x)     // with space
	fmt.Printf("%010.2s\n", x) // zero padding, width 10, precision 2
	fmt.Printf("%#.2s\n", x)   // long unit
	fmt.Printf("%# .2s\n", x)  // long unit, with space
	fmt.Printf("[%10s]\n", x)  // width 10
	fmt.Printf("[%-10s]\n", x) // width 10, left-aligned

	// Binary prefix
	fmt.Printf("%S\n", x)     // default precision
	fmt.Printf("% S\n", x)    // with space
	fmt.Printf("%.1S\n", x)   // precision 1
	fmt.Printf("%# .2S\n", x) // precision 2, long unit, with space

	// Integer
	fmt.Printf("%d\n", x) // decimal in bytes
	fmt.Printf("%x\n", x) // hexadecimal in bytes
	// Output:
	// 100.0 kB
	// ByteCount(100000)
	// 100kB
	// 100 kB
	// 00100.00kB
	// 100.00kilobytes
	// 100.00 kilobytes
	// [     100kB]
	// [100kB     ]
	// 97.65625KiB
	// 97.65625 KiB
	// 97.7KiB
	// 97.66 kibibytes
	// 100000
	// 186a0
}

//
func ExampleByteCount_Scan_sscanf() {
	src := []string{
		"kawaii.png is 1.23MB",
		"work.xlsx is 234.56 kibibytes",
		"huge.zip is 999 TB",
	}
	for _, s := range src {
		var file string
		var size infounit.ByteCount
		_, _ = fmt.Sscanf(s, "%s is %s", &file, &size)
		fmt.Println(file, size)
	}
	// Output:
	// kawaii.png 1.2 MB
	// work.xlsx 240.2 kB
	// huge.zip 999.0 TB
}
