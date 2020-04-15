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
func TestByteCount_Format_1(t *testing.T) {
	t.Parallel()

	v := infounit.ByteCount(987654321)
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
		{"%v", "987.7 MB"},
		{"%#v", "ByteCount(987654321)"},
		{"%s", "987.654321MB"},
		{"% s", "987.654321 MB"},
		{"%.1s", "987.7MB"},
		{"% .1s", "987.7 MB"},
		{"%#s", "987.654321megabytes"},
		{"%# s", "987.654321 megabytes"},
		{"%#.1s", "987.7megabytes"},
		{"%# .1s", "987.7 megabytes"},
		{"%.2s", "987.65MB"},
		{"% .3s", "987.654 MB"},
		{"%# .4s", "987.6543 megabytes"},
		{"%10.0s", "     988MB"},
		{"%10.s", "     988MB"},
		{"%-12.s", "988MB       "},
		{"%-12.0s", "988MB       "},
		{"%S", "941.900559425354MiB"},
		{"% S", "941.900559425354 MiB"},
		{"%.1S", "941.9MiB"},
		{"% .1S", "941.9 MiB"},
		{"%#S", "941.900559425354mebibytes"},
		{"%# S", "941.900559425354 mebibytes"},
		{"%#.1S", "941.9mebibytes"},
		{"%# .1S", "941.9 mebibytes"},
		{"%.3S", "941.901MiB"},
		{"% .4S", "941.9006 MiB"},
		{"%#.5S", "941.90056mebibytes"},
		{"%# .6S", "941.900559 mebibytes"},
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
func TestByteCount_Format_2(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b uint64
		s string // "% .1s", "% .1S", "%# .1s", "%# .1S"
	}{
		{0, "0 B, 0 B, 0 bytes, 0 bytes"},
		{1, "1 B, 1 B, 1 byte, 1 byte"},
		{777, "777 B, 777 B, 777 bytes, 777 bytes"},

		{1000, "1.0 kB, 1000 B, 1.0 kilobyte, 1000 bytes"},
		{1024, "1.0 kB, 1.0 KiB, 1.0 kilobytes, 1.0 kibibyte"},
		{777777, "777.8 kB, 759.5 KiB, 777.8 kilobytes, 759.5 kibibytes"},

		{1000 * 1000, "1.0 MB, 976.6 KiB, 1.0 megabyte, 976.6 kibibytes"},
		{1024 * 1024, "1.0 MB, 1.0 MiB, 1.0 megabytes, 1.0 mebibyte"},
		{777777000, "777.8 MB, 741.7 MiB, 777.8 megabytes, 741.7 mebibytes"},

		{1000 * 1000 * 1000, "1.0 GB, 953.7 MiB, 1.0 gigabyte, 953.7 mebibytes"},
		{1024 * 1024 * 1024, "1.1 GB, 1.0 GiB, 1.1 gigabytes, 1.0 gibibyte"},
		{777777000000, "777.8 GB, 724.4 GiB, 777.8 gigabytes, 724.4 gibibytes"},

		{1000 * 1000 * 1000 * 1000, "1.0 TB, 931.3 GiB, 1.0 terabyte, 931.3 gibibytes"},
		{1024 * 1024 * 1024 * 1024, "1.1 TB, 1.0 TiB, 1.1 terabytes, 1.0 tebibyte"},
		{777777000000000, "777.8 TB, 707.4 TiB, 777.8 terabytes, 707.4 tebibytes"},

		{1000 * 1000 * 1000 * 1000 * 1000, "1.0 PB, 909.5 TiB, 1.0 petabyte, 909.5 tebibytes"},
		{1024 * 1024 * 1024 * 1024 * 1024, "1.1 PB, 1.0 PiB, 1.1 petabytes, 1.0 pebibyte"},
		{777777000000000000, "777.8 PB, 690.8 PiB, 777.8 petabytes, 690.8 pebibytes"},

		{1000 * 1000 * 1000 * 1000 * 1000 * 1000, "1.0 EB, 888.2 PiB, 1.0 exabyte, 888.2 pebibytes"},
		{1024 * 1024 * 1024 * 1024 * 1024 * 1024, "1.2 EB, 1.0 EiB, 1.2 exabytes, 1.0 exbibyte"},
		{7777000000000000000, "7.8 EB, 6.7 EiB, 7.8 exabytes, 6.7 exbibytes"},

		{18446744073709551615, "18.4 EB, 16.0 EiB, 18.4 exabytes, 16.0 exbibytes"},
	}

	for _, c := range tc {
		b := infounit.ByteCount(c.b)
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
