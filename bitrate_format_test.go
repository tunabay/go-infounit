// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package infounit_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/tunabay/go-infounit"
)

//
func TestBitRate_Format_1(t *testing.T) {
	t.Parallel()

	// %g, %G

	v := infounit.BitRate(987654321.2345)
	tc := []struct {
		f string
		s string
	}{
		{"%b", "8285044940342297p-23"},
		{"%e", "9.876543e+08"},
		{"%E", "9.876543E+08"},
		{"[%30.14e]", "[          9.87654321234500e+08]"},
		{"[%030.14E]", "[00000000009.87654321234500E+08]"},
		{"[%-30.14E]", "[9.87654321234500E+08          ]"},
		{"%f", "987654321.234500"},
		{"%F", "987654321.234500"},
		{"%.0f", "987654321"},
		{"%.3F", "987654321.235"},

		{"%v", "987.7 Mbit/s"},
		{"%#v", "BitRate(987654321.2345)"},
		{"%s", "987.6543212345Mbit/s"},
		{"%a", "987.6543212345Mbps"},
		{"% s", "987.6543212345 Mbit/s"},
		{"%.1s", "987.7Mbit/s"},
		{"%.1a", "987.7Mbps"},
		{"% .1s", "987.7 Mbit/s"},
		{"%#s", "987.6543212345megabits per second"},
		{"%# s", "987.6543212345 megabits per second"},
		{"%#.1s", "987.7megabits per second"},
		{"%# .1s", "987.7 megabits per second"},
		{"%.2s", "987.65Mbit/s"},
		{"% .3s", "987.654 Mbit/s"},
		{"%# .4s", "987.6543 megabits per second"},
		{"%14.0s", "     988Mbit/s"},
		{"%14.s", "     988Mbit/s"},
		{"%-14.s", "988Mbit/s     "},
		{"%-14.0s", "988Mbit/s     "},
		{"%S", "941.9005596489907Mibit/s"},
		{"%A", "941.9005596489907Mibps"},
		{"% S", "941.9005596489907 Mibit/s"},
		{"%.1S", "941.9Mibit/s"},
		{"% .1S", "941.9 Mibit/s"},
		{"%#S", "941.9005596489907mebibits per second"},
		{"%#A", "941.9005596489907mebibits per second"},
		{"%# S", "941.9005596489907 mebibits per second"},
		{"%#.1S", "941.9mebibits per second"},
		{"%# .1S", "941.9 mebibits per second"},
		{"%.3S", "941.901Mibit/s"},
		{"% .4S", "941.9006 Mibit/s"},
		{"%#.5S", "941.90056mebibits per second"},
		{"%# .6S", "941.900560 mebibits per second"},
	}

	for _, c := range tc {
		s := fmt.Sprintf(c.f, v)
		// t.Logf(`%s: "%s"`, c.f, s)
		if s != c.s {
			t.Errorf(`fmt "%s": want: "%s", got: "%s"`, c.f, c.s, s)
		}
	}
}

//
func TestBitRate_Format_2(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b float64
		s string // "% .1s", "% .1S", "%# .1s", "%# .1S"
	}{
		{0, "0.0 bit/s, 0.0 bit/s, 0.0 bits per second, 0.0 bits per second"},
		{1, "1.0 bit/s, 1.0 bit/s, 1.0 bit per second, 1.0 bit per second"},
		{777, "777.0 bit/s, 777.0 bit/s, 777.0 bits per second, 777.0 bits per second"},

		{1000, "1.0 kbit/s, 1000.0 bit/s, 1.0 kilobit per second, 1000.0 bits per second"},
		{1024, "1.0 kbit/s, 1.0 Kibit/s, 1.0 kilobits per second, 1.0 kibibit per second"},
		{777777, "777.8 kbit/s, 759.5 Kibit/s, 777.8 kilobits per second, 759.5 kibibits per second"},

		{1000 * 1000, "1.0 Mbit/s, 976.6 Kibit/s, 1.0 megabit per second, 976.6 kibibits per second"},
		{1024 * 1024, "1.0 Mbit/s, 1.0 Mibit/s, 1.0 megabits per second, 1.0 mebibit per second"},
		{777777000, "777.8 Mbit/s, 741.7 Mibit/s, 777.8 megabits per second, 741.7 mebibits per second"},

		{1000 * 1000 * 1000, "1.0 Gbit/s, 953.7 Mibit/s, 1.0 gigabit per second, 953.7 mebibits per second"},
		{1024 * 1024 * 1024, "1.1 Gbit/s, 1.0 Gibit/s, 1.1 gigabits per second, 1.0 gibibit per second"},
		{777777000000, "777.8 Gbit/s, 724.4 Gibit/s, 777.8 gigabits per second, 724.4 gibibits per second"},

		{1000 * 1000 * 1000 * 1000, "1.0 Tbit/s, 931.3 Gibit/s, 1.0 terabit per second, 931.3 gibibits per second"},
		{1024 * 1024 * 1024 * 1024, "1.1 Tbit/s, 1.0 Tibit/s, 1.1 terabits per second, 1.0 tebibit per second"},
		{777777000000000, "777.8 Tbit/s, 707.4 Tibit/s, 777.8 terabits per second, 707.4 tebibits per second"},

		{1000 * 1000 * 1000 * 1000 * 1000, "1.0 Pbit/s, 909.5 Tibit/s, 1.0 petabit per second, 909.5 tebibits per second"},
		{1024 * 1024 * 1024 * 1024 * 1024, "1.1 Pbit/s, 1.0 Pibit/s, 1.1 petabits per second, 1.0 pebibit per second"},
		{777777000000000000, "777.8 Pbit/s, 690.8 Pibit/s, 777.8 petabits per second, 690.8 pebibits per second"},

		{1000 * 1000 * 1000 * 1000 * 1000 * 1000, "1.0 Ebit/s, 888.2 Pibit/s, 1.0 exabit per second, 888.2 pebibits per second"},
		{1024 * 1024 * 1024 * 1024 * 1024 * 1024, "1.2 Ebit/s, 1.0 Eibit/s, 1.2 exabits per second, 1.0 exbibit per second"},
		{7777000000000000000, "7.8 Ebit/s, 6.7 Eibit/s, 7.8 exabits per second, 6.7 exbibits per second"},

		{18446744073709551615, "18.4 Ebit/s, 16.0 Eibit/s, 18.4 exabits per second, 16.0 exbibits per second"},
	}

	for _, c := range tc {
		b := infounit.BitRate(c.b)
		es := strings.Split(c.s, ", ")
		for i, f := range []string{"% .1s", "% .1S", "%# .1s", "%# .1S"} {
			s := fmt.Sprintf(f, b)
			// t.Logf(`("%s", %d): "%s"`, f, c.b, s)
			if s != es[i] {
				t.Errorf(`("%s", %#v): want: "%s", got: "%s"`, f, b, es[i], s)
			}
		}
	}
}
