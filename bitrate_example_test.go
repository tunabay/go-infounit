// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package infounit_test

import (
	"fmt"

	"github.com/tunabay/go-infounit"
)

//
func ExampleBitRate_Convert_siPrefix() {
	x := infounit.BitRate(123456789.0)

	fmt.Printf("%f\n", x.Convert(infounit.KilobitPerSecond))
	fmt.Printf("%f\n", x.Convert(infounit.MegabitPerSecond))
	fmt.Printf("%f\n", x.Convert(infounit.GigabitPerSecond))
	fmt.Printf("%f\n", x.Convert(infounit.TerabitPerSecond))
	// Output:
	// 123456.789000
	// 123.456789
	// 0.123457
	// 0.000123
}

//
func ExampleBitRate_Convert_binaryPrefix() {
	x := infounit.BitRate(123456789.0)

	fmt.Printf("%f\n", x.Convert(infounit.KibibitPerSecond))
	fmt.Printf("%f\n", x.Convert(infounit.MebibitPerSecond))
	fmt.Printf("%f\n", x.Convert(infounit.GibibitPerSecond))
	fmt.Printf("%f\n", x.Convert(infounit.TebibitPerSecond))
	// Output:
	// 120563.270508
	// 117.737569
	// 0.114978
	// 0.000112
}

//
func ExampleBitRate_ConvertRound_siPrefix() {
	x := infounit.BitRate(123456789.0)

	fmt.Printf("%f\n", x.ConvertRound(infounit.KilobitPerSecond, 2))
	fmt.Printf("%f\n", x.ConvertRound(infounit.KilobitPerSecond, 5))
	fmt.Printf("%f\n", x.ConvertRound(infounit.MegabitPerSecond, 0))
	fmt.Printf("%f\n", x.ConvertRound(infounit.MegabitPerSecond, 3))
	fmt.Printf("%f\n", x.ConvertRound(infounit.TerabitPerSecond, 6))
	// Output:
	// 123456.790000
	// 123456.789000
	// 123.000000
	// 123.457000
	// 0.000123
}

//
func ExampleBitRate_ConvertRound_binaryPrefix() {
	x := infounit.BitRate(123456789.0)

	fmt.Printf("%f\n", x.ConvertRound(infounit.KibibitPerSecond, 2))
	fmt.Printf("%f\n", x.ConvertRound(infounit.KibibitPerSecond, 5))
	fmt.Printf("%f\n", x.ConvertRound(infounit.MebibitPerSecond, 0))
	fmt.Printf("%f\n", x.ConvertRound(infounit.MebibitPerSecond, 3))
	fmt.Printf("%f\n", x.ConvertRound(infounit.TebibitPerSecond, 6))
	// Output:
	// 120563.270000
	// 120563.270510
	// 118.000000
	// 117.738000
	// 0.000112
}

//
func ExampleBitRate_Format_errorf() {
	limit := infounit.KilobitPerSecond * 100

	check := func(rate infounit.BitRate) error {
		if limit < rate {
			return fmt.Errorf("too fast: %v", rate)
		}
		return nil
	}

	for i := 0; i < 8; i++ {
		rate := infounit.BitRate(uint64(1 << (6 * i)))
		if err := check(rate); err != nil {
			fmt.Println("error:", err)
			continue
		}
		fmt.Println("rate:", rate)
	}
	// Output:
	// rate: 1.0 bit/s
	// rate: 64.0 bit/s
	// rate: 4.1 kbit/s
	// error: too fast: 262.1 kbit/s
	// error: too fast: 16.8 Mbit/s
	// error: too fast: 1.1 Gbit/s
	// error: too fast: 68.7 Gbit/s
	// error: too fast: 4.4 Tbit/s
}

//
func ExampleBitRate_Format_printf() {
	x := infounit.KilobitPerSecond * 123.456 // 123.456 kbit/s

	// Default format
	fmt.Printf("%v\n", x)  // default format, same as "% .1s"
	fmt.Printf("%#v\n", x) // Go syntax format

	// SI prefix
	fmt.Printf("%s\n", x)      // default precision
	fmt.Printf("% s\n", x)     // with space
	fmt.Printf("%014.2s\n", x) // zero padding, width 14, precision 2
	fmt.Printf("%#.2s\n", x)   // long unit
	fmt.Printf("%# .2s\n", x)  // long unit, with space
	fmt.Printf("[%16s]\n", x)  // width 16
	fmt.Printf("[%-16s]\n", x) // width 16, left-aligned
	fmt.Printf("% a\n", x)     // non-standard abbreviation "bps"

	// Binary prefix
	fmt.Printf("%S\n", x)     // default precision
	fmt.Printf("% S\n", x)    // with space
	fmt.Printf("%.1S\n", x)   // precision 1
	fmt.Printf("%# .2S\n", x) // precision 2, long unit, with space
	fmt.Printf("% A\n", x)    // non-standard abbreviation "bps"

	// Float
	fmt.Printf("%f\n", x) // floating point value in bit/s
	// Output:
	// 123.5 kbit/s
	// BitRate(123456)
	// 123.456kbit/s
	// 123.456 kbit/s
	// 00123.46kbit/s
	// 123.46kilobits per second
	// 123.46 kilobits per second
	// [   123.456kbit/s]
	// [123.456kbit/s   ]
	// 123.456 kbps
	// 120.5625Kibit/s
	// 120.5625 Kibit/s
	// 120.6Kibit/s
	// 120.56 kibibits per second
	// 120.5625 Kibps
	// 123456.000000
}
