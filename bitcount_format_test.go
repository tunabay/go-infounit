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
func TestBitCount_Format_1(t *testing.T) {
	t.Parallel()

	v := infounit.BitCount(987654321)
	tc := []struct {
		f string
		s string
	}{
		{"%b", "111010110111100110100010110001"},
		{"%032b", "00111010110111100110100010110001"},
		{"%32b", "  111010110111100110100010110001"},
		{"%-32b", "111010110111100110100010110001  "},
		{"%#b", "0b111010110111100110100010110001"},
		{"%o", "7267464261"},
		{"%012o", "007267464261"},
		{"%12o", "  7267464261"},
		{"%-12o", "7267464261  "},
		{"%d", "987654321"},
		{"%012d", "000987654321"},
		{"%12d", "   987654321"},
		{"%-12d", "987654321   "},
		{"%x", "3ade68b1"},
		{"%010x", "003ade68b1"},
		{"%#x", "0x3ade68b1"},
		{"%X", "3ADE68B1"},
		{"%010X", "003ADE68B1"},
		{"%#X", "0X3ADE68B1"},
		{"%v", "987.7 Mbit"},
		{"%#v", "BitCount(987654321)"},
		{"%s", "987.654321Mbit"},
		{"% s", "987.654321 Mbit"},
		{"%.1s", "987.7Mbit"},
		{"% .1s", "987.7 Mbit"},
		{"%#s", "987.654321megabits"},
		{"%# s", "987.654321 megabits"},
		{"%#.1s", "987.7megabits"},
		{"%# .1s", "987.7 megabits"},
		{"%.2s", "987.65Mbit"},
		{"% .3s", "987.654 Mbit"},
		{"%# .4s", "987.6543 megabits"},
		{"%12.0s", "     988Mbit"},
		{"%12.s", "     988Mbit"},
		{"%-12.s", "988Mbit     "},
		{"%-12.0s", "988Mbit     "},
		{"%S", "941.900559425354Mibit"},
		{"% S", "941.900559425354 Mibit"},
		{"%.1S", "941.9Mibit"},
		{"% .1S", "941.9 Mibit"},
		{"%#S", "941.900559425354mebibits"},
		{"%# S", "941.900559425354 mebibits"},
		{"%#.1S", "941.9mebibits"},
		{"%# .1S", "941.9 mebibits"},
		{"%.3S", "941.901Mibit"},
		{"% .4S", "941.9006 Mibit"},
		{"%#.5S", "941.90056mebibits"},
		{"%# .6S", "941.900559 mebibits"},
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
func TestBitCount_Format_2(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b uint64
		s string // "% .1s", "% .1S", "%# .1s", "%# .1S"
	}{
		{0, "0 bit, 0 bit, 0 bits, 0 bits"},
		{1, "1 bit, 1 bit, 1 bit, 1 bit"},
		{777, "777 bit, 777 bit, 777 bits, 777 bits"},

		{1000, "1.0 kbit, 1000 bit, 1.0 kilobit, 1000 bits"},
		{1024, "1.0 kbit, 1.0 Kibit, 1.0 kilobits, 1.0 kibibit"},
		{777777, "777.8 kbit, 759.5 Kibit, 777.8 kilobits, 759.5 kibibits"},

		{1000 * 1000, "1.0 Mbit, 976.6 Kibit, 1.0 megabit, 976.6 kibibits"},
		{1024 * 1024, "1.0 Mbit, 1.0 Mibit, 1.0 megabits, 1.0 mebibit"},
		{777777000, "777.8 Mbit, 741.7 Mibit, 777.8 megabits, 741.7 mebibits"},

		{1000 * 1000 * 1000, "1.0 Gbit, 953.7 Mibit, 1.0 gigabit, 953.7 mebibits"},
		{1024 * 1024 * 1024, "1.1 Gbit, 1.0 Gibit, 1.1 gigabits, 1.0 gibibit"},
		{777777000000, "777.8 Gbit, 724.4 Gibit, 777.8 gigabits, 724.4 gibibits"},

		{1000 * 1000 * 1000 * 1000, "1.0 Tbit, 931.3 Gibit, 1.0 terabit, 931.3 gibibits"},
		{1024 * 1024 * 1024 * 1024, "1.1 Tbit, 1.0 Tibit, 1.1 terabits, 1.0 tebibit"},
		{777777000000000, "777.8 Tbit, 707.4 Tibit, 777.8 terabits, 707.4 tebibits"},

		{1000 * 1000 * 1000 * 1000 * 1000, "1.0 Pbit, 909.5 Tibit, 1.0 petabit, 909.5 tebibits"},
		{1024 * 1024 * 1024 * 1024 * 1024, "1.1 Pbit, 1.0 Pibit, 1.1 petabits, 1.0 pebibit"},
		{777777000000000000, "777.8 Pbit, 690.8 Pibit, 777.8 petabits, 690.8 pebibits"},

		{1000 * 1000 * 1000 * 1000 * 1000 * 1000, "1.0 Ebit, 888.2 Pibit, 1.0 exabit, 888.2 pebibits"},
		{1024 * 1024 * 1024 * 1024 * 1024 * 1024, "1.2 Ebit, 1.0 Eibit, 1.2 exabits, 1.0 exbibit"},
		{7777000000000000000, "7.8 Ebit, 6.7 Eibit, 7.8 exabits, 6.7 exbibits"},

		{18446744073709551615, "18.4 Ebit, 16.0 Eibit, 18.4 exabits, 16.0 exbibits"},
	}

	for _, c := range tc {
		b := infounit.BitCount(c.b)
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
