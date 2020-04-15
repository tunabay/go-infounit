// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package infounit_test

import (
	"fmt"
	"testing"

	"github.com/tunabay/go-infounit"
)

//
func TestByteCount_Scan_1(t *testing.T) {
	t.Parallel()

	tc := []struct {
		src string
		fmt string
		bc  infounit.ByteCount
		es  string
	}{
		{"0B", "%s", infounit.Byte * 0, ""},
		{"0 B", "%s", infounit.Byte * 0, ""},
		{"0", "%u", infounit.Byte * 0, ""},
		{"0", "%U", infounit.Byte * 0, ""},
		{"18446744073709551615 B", "%s", infounit.Byte * 18446744073709551615, ""},
		{"18446744073709551615 B", "%S", infounit.Byte * 18446744073709551615, ""},
		{"18446744073709551615", "%u", infounit.Byte * 18446744073709551615, ""},
		{"18446744073709551615", "%U", infounit.Byte * 18446744073709551615, ""},

		{"110", "%u", infounit.Byte * 110, ""},
		{"110", "%U", infounit.Byte * 110, ""},
		{"110B", "%u", infounit.Byte * 110, ""},
		{"110B", "%U", infounit.Byte * 110, ""},
		{"110b", "%u", infounit.Byte * 110, ""},
		{"110b", "%U", infounit.Byte * 110, ""},
		{"111B", "%s", infounit.Byte * 111, ""},
		{"111B", "%S", infounit.Byte * 111, ""},
		{"111b", "%s", infounit.Byte * 111, ""},
		{"111b", "%S", infounit.Byte * 111, ""},
		{"112 B", "%s", infounit.Byte * 112, ""},
		{"112 B", "%S", infounit.Byte * 112, ""},
		{"112 B", "%u", infounit.Byte * 112, ""}, // %u does not read "B"
		{"112 B", "%U", infounit.Byte * 112, ""}, // %U does not read "B"
		{"113byte", "%s", infounit.Byte * 113, ""},
		{"113bYtE", "%s", infounit.Byte * 113, ""},
		{"113byte", "%S", infounit.Byte * 113, ""},
		{"113byte", "%u", infounit.Byte * 113, ""},
		{"113byte", "%U", infounit.Byte * 113, ""},
		{"114bytes", "%s", infounit.Byte * 114, ""},
		{"114ByTeS", "%s", infounit.Byte * 114, ""},
		{"114bytes", "%S", infounit.Byte * 114, ""},
		{"115 byte", "%s", infounit.Byte * 115, ""},
		{"116 bytes", "%s", infounit.Byte * 116, ""},
		{"77000000000000 bytes", "%s", infounit.Terabyte * 77, ""},

		{"210kB", "%s", infounit.Kilobyte * 210, ""},
		{"210Kb", "%s", infounit.Kilobyte * 210, ""},
		{"210kB", "%S", infounit.Kibibyte * 210, ""},
		{"210kB", "%u", infounit.Kilobyte * 210, ""},
		{"210kB", "%U", infounit.Kibibyte * 210, ""},
		{"211 kB", "%s", infounit.Kilobyte * 211, ""},
		{"211 kB", "%S", infounit.Kibibyte * 211, ""},
		{"211 kB", "%u", infounit.Byte * 211, ""}, // %u does not read "kB"
		{"211 kB", "%U", infounit.Byte * 211, ""}, // %U does not read "kB"
		{"212kilobyte", "%s", infounit.Kilobyte * 212, ""},
		{"213KiLoBYTES", "%s", infounit.Kilobyte * 213, ""},
		{"214 kilobyte", "%s", infounit.Kilobyte * 214, ""},
		{"215 kilobytes", "%s", infounit.Kilobyte * 215, ""},
		{"220.5kB", "%s", infounit.Byte * 220500, ""},
		{".75 kB", "%s", infounit.Byte * 750, ""},
		{"00.777 kilobytes", "%s", infounit.Byte * 777, ""},

		{"310MB", "%s", infounit.Megabyte * 310, ""},
		{"310MB", "%S", infounit.Mebibyte * 310, ""},
		{"310MB", "%u", infounit.Megabyte * 310, ""},
		{"310MB", "%U", infounit.Mebibyte * 310, ""},
		{"311 mb", "%s", infounit.Megabyte * 311, ""},
		{"311 MB", "%S", infounit.Mebibyte * 311, ""},
		{"312megabytes", "%s", infounit.Megabyte * 312, ""},
		{"312 megabytes", "%S", infounit.Mebibyte * 312, ""},
		{"320.25 MB", "%s", infounit.Kilobyte * 320250, ""},
		{"0.000567 megabytes", "%s", infounit.Byte * 567, ""},
		{"30000.0megabyte", "%s", infounit.Gigabyte * 30, ""},

		{"410GB", "%s", infounit.Gigabyte * 410, ""},
		{"410gb", "%S", infounit.Gibibyte * 410, ""},
		{"410GB", "%u", infounit.Gigabyte * 410, ""},
		{"410Gb", "%U", infounit.Gibibyte * 410, ""},
		{"411 Gb", "%s", infounit.Gigabyte * 411, ""},
		{"411 gB", "%S", infounit.Gibibyte * 411, ""},
		{"412gigaBYTE", "%s", infounit.Gigabyte * 412, ""},
		{"412 gigaBYtes", "%S", infounit.Gibibyte * 412, ""},
		{"420.001 GB", "%s", infounit.Megabyte * 420001, ""},
		{"0.00000001GB", "%s", infounit.Byte * 10, ""},

		{"510TB", "%s", infounit.Terabyte * 510, ""},
		{"510tb", "%S", infounit.Tebibyte * 510, ""},
		{"510TB", "%u", infounit.Terabyte * 510, ""},
		{"510TB", "%U", infounit.Tebibyte * 510, ""},
		{"511 terabytes", "%s", infounit.Terabyte * 511, ""},
		{"511 terabytes", "%S", infounit.Tebibyte * 511, ""},
		{"0.000000012 TB", "%s", infounit.Kilobyte * 12, ""},

		{"610PB", "%s", infounit.Petabyte * 610, ""},
		{"610pb", "%S", infounit.Pebibyte * 610, ""},
		{"611 petabyte", "%s", infounit.Petabyte * 611, ""},
		{"611 petaBYTes", "%S", infounit.Pebibyte * 611, ""},
		{"18446 PB", "%s", infounit.Petabyte * 18446, ""},
		{"18.446 EB", "%s", infounit.ByteCount(float64(infounit.Exabyte) * 18.446), ""},

		{"11EB", "%s", infounit.Exabyte * 11, ""},
		{"11eb", "%S", infounit.Exbibyte * 11, ""},
		{"18.2 EB", "%s", infounit.Petabyte * 18200, ""},
		{"15.5 EB", "%S", infounit.Exbibyte*15 + infounit.Pebibyte*512, ""},

		{"260kib", "%s", infounit.Kibibyte * 260, ""},
		{"260KiB", "%S", infounit.Kibibyte * 260, ""},
		{"260KiB", "%u", infounit.Kibibyte * 260, ""},
		{"260KiB", "%U", infounit.Kibibyte * 260, ""},
		{"261 KiB", "%s", infounit.Kibibyte * 261, ""},
		{"261 KiB", "%S", infounit.Kibibyte * 261, ""},
		{"262 kibibyte", "%s", infounit.Kibibyte * 262, ""},
		{"263 kiBIBYtes", "%s", infounit.Kibibyte * 263, ""},

		{"360MiB", "%s", infounit.Mebibyte * 360, ""},
		{"360mib", "%S", infounit.Mebibyte * 360, ""},
		{"360MiB", "%u", infounit.Mebibyte * 360, ""},
		{"360MiB", "%U", infounit.Mebibyte * 360, ""},
		{"361 MiB", "%s", infounit.Mebibyte * 361, ""},
		{"361 MiB", "%S", infounit.Mebibyte * 361, ""},
		{"362 mebibyte", "%s", infounit.Mebibyte * 362, ""},
		{"363 mebibytes", "%s", infounit.Mebibyte * 363, ""},

		{"460GiB", "%s", infounit.Gibibyte * 460, ""},
		{"460gib", "%S", infounit.Gibibyte * 460, ""},
		{"460GiB", "%u", infounit.Gibibyte * 460, ""},
		{"460GiB", "%U", infounit.Gibibyte * 460, ""},
		{"461 GiB", "%s", infounit.Gibibyte * 461, ""},
		{"461 GiB", "%S", infounit.Gibibyte * 461, ""},
		{"462 gibibyte", "%s", infounit.Gibibyte * 462, ""},
		{"463 gibibytes", "%s", infounit.Gibibyte * 463, ""},

		{"560TiB", "%s", infounit.Tebibyte * 560, ""},
		{"560tib", "%S", infounit.Tebibyte * 560, ""},
		{"560TiB", "%u", infounit.Tebibyte * 560, ""},
		{"560TiB", "%U", infounit.Tebibyte * 560, ""},
		{"561 TiB", "%s", infounit.Tebibyte * 561, ""},
		{"561 TiB", "%S", infounit.Tebibyte * 561, ""},
		{"562 tibibyte", "%s", infounit.Tebibyte * 562, ""},
		{"563 tibibytes", "%s", infounit.Tebibyte * 563, ""},

		{"660PiB", "%s", infounit.Pebibyte * 660, ""},
		{"660pib", "%S", infounit.Pebibyte * 660, ""},
		{"660PiB", "%u", infounit.Pebibyte * 660, ""},
		{"660PiB", "%U", infounit.Pebibyte * 660, ""},
		{"661 PiB", "%s", infounit.Pebibyte * 661, ""},
		{"661 PiB", "%S", infounit.Pebibyte * 661, ""},
		{"662 pebibyte", "%s", infounit.Pebibyte * 662, ""},
		{"663 pebibytes", "%s", infounit.Pebibyte * 663, ""},

		{"10EiB", "%s", infounit.Exbibyte * 10, ""},
		{"10eib", "%S", infounit.Exbibyte * 10, ""},
		{"12.5 EiB", "%s", infounit.Exbibyte*12 + infounit.Pebibyte*512, ""},
		{"13.5 EiB", "%S", infounit.Exbibyte*13 + infounit.Pebibyte*512, ""},

		{"0.00001 kB", "%s", 0, ""},
		{"999", "%t", 0, "unknown verb for ByteCount: %t"},
		{"", "%s", 0, "%s: no input"},
		{"many bytes", "%s", 0, "%s: invalid expr: many"},
		{"0.1B", "%s", 0, "%s: non-integer byte count: 0.1"},
		{"0.1B", "%S", 0, "%S: non-integer byte count: 0.1"},
		{"0.1B", "%u", 0, "%u: non-integer byte count: 0.1"},
		{"0.1B", "%U", 0, "%U: non-integer byte count: 0.1"},
		{"0.12 b", "%s", 0, "%s: non-integer byte count: 0.12"},
		{"0.12 b", "%S", 0, "%S: non-integer byte count: 0.12"},
		{"0.123bytes", "%s", 0, "%s: non-integer byte count: 0.123"},
		{"0.123bytes", "%S", 0, "%S: non-integer byte count: 0.123"},
		{"999", "%s", 0, "%s: no unit suffix: EOF"},
		{"999", "%S", 0, "%S: no unit suffix: EOF"},
		{"999  GB", "%s", 0, "%s: no unit suffix"},
		{"999 666", "%s", 0, "%s: invalid unit expr: 666"},
		{"+9999", "%s", 0, "%s: invalid expr: +9999"},
		{"999-666", "%s", 0, "%s: invalid expr: 999-666"},
		{"999 megabytes", "%u", infounit.Byte * 999, ""}, // %u does not read "megabytes"
		{"999 megabytes", "%U", infounit.Byte * 999, ""}, // %U does not read "megabytes"
		{"1.21jigowatts", "%s", 0, "%s: unknown unit: jigowatts"},
		{"1.21jigowatts", "%S", 0, "%S: unknown unit: jigowatts"},
		{"1.21jigowatts", "%u", 0, "%u: unknown unit: jigowatts"},
		{"1.21jigowatts", "%U", 0, "%U: unknown unit: jigowatts"},
		{"1.21 jigowatts", "%s", 0, "%s: unknown unit: jigowatts"},
		{"1.21 jigowatts", "%S", 0, "%S: unknown unit: jigowatts"},
		{"1.21 jigowatts", "%u", 0, "%u: non-integer byte count: 1.21"}, // %u does not read "jigowatts"
		{"1.21 jigowatts", "%U", 0, "%U: non-integer byte count: 1.21"}, // %U does not read "jigowatts"

		{"0", "%d", infounit.Byte * 0, ""},
		{"1", "%d", infounit.Byte * 1, ""},
		{"18446744073709551615", "%d", infounit.Byte * 18446744073709551615, ""},
		{"0", "%x", infounit.Byte * 0, ""},
		{"ffffffffffffffff", "%x", infounit.Byte * 0xffffffffffffffff, ""},
		{"DeadBeef", "%X", infounit.Byte * 0xdeadbeef, ""},
		{"11110000", "%b", infounit.Byte * 0xf0, ""},
		{"1111111111111111111111111111111111111111111111111111111111111111", "%b", infounit.Byte * 0xffffffffffffffff, ""},
		{"77777777", "%o", infounit.Byte * 0xffffff, ""},
	}

	for _, c := range tc {
		var bc infounit.ByteCount
		n, err := fmt.Sscanf(c.src, c.fmt, &bc)
		switch c.es {
		case "": // expecting no error
			switch {
			case err != nil:
				t.Errorf("src='%s', fmt='%s': %s", c.src, c.fmt, err)
				continue
			case n != 1:
				t.Errorf("src='%s', fmt='%s': n(%d) != 1", c.src, c.fmt, n)
				continue
			case bc != c.bc:
				t.Errorf("src='%s', fmt='%s': want: %#v, got: %#v", c.src, c.fmt, c.bc, bc)
				continue
			}
			// t.Logf("src='%s', fmt='%s': OK: %v", c.src, c.fmt, bc)
		default: // expecting error
			switch {
			case err == nil:
				t.Errorf("src='%s', fmt='%s': error expected: got: %v", c.src, c.fmt, bc)
				continue
			case err.Error() != c.es:
				t.Errorf("src='%s', fmt='%s': error want: %s, got: %s", c.src, c.fmt, c.es, err.Error())
				continue
			}
			// t.Logf("src='%s', fmt='%s': OK: %s", c.src, c.fmt, err.Error())
		}
	}
}
