// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package infounit_test

import (
	"fmt"

	"github.com/tunabay/go-infounit"
)

//
func Example_printf() {
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

	// Output:
	// 2.5 GB
	// 2.3 GiB
	// 2.50 gigabytes
	// 2.328 gibibytes
	// 2500000000
	// 32 kbit
	// 31.25 Kibit
	// 32.00 kilobits
	// 31.250 kibibits
	// 32000
	// 123.45 Gbit/s
	// 114.972 Gibit/s
	// 123.5 gigabits per second
	// 114.97 gibibits per second
	// 123.5 Gbps
	// 114.97 Gibps
}
