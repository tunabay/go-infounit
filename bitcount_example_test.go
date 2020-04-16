// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package infounit_test

import (
	"fmt"

	"github.com/tunabay/go-infounit"
)

//
func ExampleBitCount_Convert_siPrefix() {
	x := infounit.BitCount(123456789)

	fmt.Printf("%f\n", x.Convert(infounit.Kilobit))
	fmt.Printf("%f\n", x.Convert(infounit.Megabit))
	fmt.Printf("%f\n", x.Convert(infounit.Gigabit))
	fmt.Printf("%f\n", x.Convert(infounit.Terabit))
	// Output:
	// 123456.789000
	// 123.456789
	// 0.123457
	// 0.000123
}

//
func ExampleBitCount_Convert_binaryPrefix() {
	x := infounit.BitCount(123456789)

	fmt.Printf("%f\n", x.Convert(infounit.Kibibit))
	fmt.Printf("%f\n", x.Convert(infounit.Mebibit))
	fmt.Printf("%f\n", x.Convert(infounit.Gibibit))
	fmt.Printf("%f\n", x.Convert(infounit.Tebibit))
	// Output:
	// 120563.270508
	// 117.737569
	// 0.114978
	// 0.000112
}

//
func ExampleBitCount_ConvertRound_siPrefix() {
	x := infounit.BitCount(123456789)

	fmt.Printf("%f\n", x.ConvertRound(infounit.Kilobit, 2))
	fmt.Printf("%f\n", x.ConvertRound(infounit.Kilobit, 5))
	fmt.Printf("%f\n", x.ConvertRound(infounit.Megabit, 0))
	fmt.Printf("%f\n", x.ConvertRound(infounit.Megabit, 3))
	fmt.Printf("%f\n", x.ConvertRound(infounit.Terabit, 6))
	// Output:
	// 123456.790000
	// 123456.789000
	// 123.000000
	// 123.457000
	// 0.000123
}

//
func ExampleBitCount_ConvertRound_binaryPrefix() {
	x := infounit.BitCount(123456789)

	fmt.Printf("%f\n", x.ConvertRound(infounit.Kibibit, 2))
	fmt.Printf("%f\n", x.ConvertRound(infounit.Kibibit, 5))
	fmt.Printf("%f\n", x.ConvertRound(infounit.Mebibit, 0))
	fmt.Printf("%f\n", x.ConvertRound(infounit.Mebibit, 3))
	fmt.Printf("%f\n", x.ConvertRound(infounit.Tebibit, 6))
	// Output:
	// 120563.270000
	// 120563.270510
	// 118.000000
	// 117.738000
	// 0.000112
}

//
func ExampleBitCount_Format_errorf() {
	limit := infounit.Kilobit * 100

	check := func(size infounit.BitCount) error {
		if limit < size {
			return fmt.Errorf("too large: %v", size)
		}
		return nil
	}

	for i := 0; i < 8; i++ {
		size := infounit.BitCount(1 << (6 * i))
		if err := check(size); err != nil {
			fmt.Println("error:", err)
			continue
		}
		fmt.Println("processed:", size)
	}
	// Output:
	// processed: 1 bit
	// processed: 64 bit
	// processed: 4.1 kbit
	// error: too large: 262.1 kbit
	// error: too large: 16.8 Mbit
	// error: too large: 1.1 Gbit
	// error: too large: 68.7 Gbit
	// error: too large: 4.4 Tbit
}

//
func ExampleBitCount_Format_printf() {
	x := infounit.Kilobit * 100 // 100000 bytes

	// Default format
	fmt.Printf("%v\n", x)  // default format, same as "% .1s"
	fmt.Printf("%#v\n", x) // Go syntax format

	// SI prefix
	fmt.Printf("%s\n", x)      // default precision
	fmt.Printf("% s\n", x)     // with space
	fmt.Printf("%012.2s\n", x) // zero padding, width 12, precision 2
	fmt.Printf("%#.2s\n", x)   // long unit
	fmt.Printf("%# .2s\n", x)  // long unit, with space
	fmt.Printf("[%12s]\n", x)  // width 12
	fmt.Printf("[%-12s]\n", x) // width 12, left-aligned

	// Binary prefix
	fmt.Printf("%S\n", x)     // default precision
	fmt.Printf("% S\n", x)    // with space
	fmt.Printf("%.1S\n", x)   // precision 1
	fmt.Printf("%# .2S\n", x) // precision 2, long unit, with space

	// Integer
	fmt.Printf("%d\n", x) // decimal in bytes
	fmt.Printf("%x\n", x) // hexadecimal in bytes
	// Output:
	// 100.0 kbit
	// BitCount(100000)
	// 100kbit
	// 100 kbit
	// 00100.00kbit
	// 100.00kilobits
	// 100.00 kilobits
	// [     100kbit]
	// [100kbit     ]
	// 97.65625Kibit
	// 97.65625 Kibit
	// 97.7Kibit
	// 97.66 kibibits
	// 100000
	// 186a0
}

//
func ExampleBitCount_Scan_sscanf() {
	src := []string{
		"kawaii.png is 1.23Mbit",
		"work.xlsx is 234.56 kibibits",
		"huge.zip is 999 Tbit",
	}
	for _, s := range src {
		var file string
		var size infounit.BitCount
		_, _ = fmt.Sscanf(s, "%s is %s", &file, &size)
		fmt.Println(file, size)
	}
	// Output:
	// kawaii.png 1.2 Mbit
	// work.xlsx 240.2 kbit
	// huge.zip 999.0 Tbit
}
