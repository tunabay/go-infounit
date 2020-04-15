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
func TestBitCount_Scan_1(t *testing.T) {
	t.Parallel()

	tc := []struct {
		src string
		fmt string
		bc  infounit.BitCount
		es  string
	}{
		{"0bit", "%s", infounit.Bit * 0, ""},
		{"0 bit", "%s", infounit.Bit * 0, ""},
		{"0", "%u", infounit.Bit * 0, ""},
		{"0", "%U", infounit.Bit * 0, ""},
		{"18446744073709551615 bit", "%s", infounit.Bit * 18446744073709551615, ""},
		{"18446744073709551615 bit", "%S", infounit.Bit * 18446744073709551615, ""},
		{"18446744073709551615", "%u", infounit.Bit * 18446744073709551615, ""},
		{"18446744073709551615", "%U", infounit.Bit * 18446744073709551615, ""},

		{"110", "%u", infounit.Bit * 110, ""},
		{"110", "%U", infounit.Bit * 110, ""},
		{"110BIT", "%u", infounit.Bit * 110, ""},
		{"110BIT", "%U", infounit.Bit * 110, ""},
		{"110bit", "%u", infounit.Bit * 110, ""},
		{"110bit", "%U", infounit.Bit * 110, ""},
		{"111BIT", "%s", infounit.Bit * 111, ""},
		{"111BIT", "%S", infounit.Bit * 111, ""},
		{"111bit", "%s", infounit.Bit * 111, ""},
		{"111bit", "%S", infounit.Bit * 111, ""},
		{"112 BIT", "%s", infounit.Bit * 112, ""},
		{"112 BIT", "%S", infounit.Bit * 112, ""},
		{"112 BIT", "%u", infounit.Bit * 112, ""}, // %u does not read "BIT"
		{"112 BIT", "%U", infounit.Bit * 112, ""}, // %U does not read "BIT"
		{"113bit", "%s", infounit.Bit * 113, ""},
		{"113bIt", "%s", infounit.Bit * 113, ""},
		{"113bit", "%S", infounit.Bit * 113, ""},
		{"113bit", "%u", infounit.Bit * 113, ""},
		{"113bit", "%U", infounit.Bit * 113, ""},
		{"114bits", "%s", infounit.Bit * 114, ""},
		{"114BiTS", "%s", infounit.Bit * 114, ""},
		{"114bits", "%S", infounit.Bit * 114, ""},
		{"115 bit", "%s", infounit.Bit * 115, ""},
		{"116 bits", "%s", infounit.Bit * 116, ""},
		{"77000000000000 bits", "%s", infounit.Terabit * 77, ""},

		{"210kBit", "%s", infounit.Kilobit * 210, ""},
		{"210Kbit", "%s", infounit.Kilobit * 210, ""},
		{"210kbit", "%S", infounit.Kibibit * 210, ""},
		{"210kbit", "%u", infounit.Kilobit * 210, ""},
		{"210kbit", "%U", infounit.Kibibit * 210, ""},
		{"211 kbit", "%s", infounit.Kilobit * 211, ""},
		{"211 kBit", "%S", infounit.Kibibit * 211, ""},
		{"211 kbit", "%u", infounit.Bit * 211, ""}, // %u does not read "kB"
		{"211 kbit", "%U", infounit.Bit * 211, ""}, // %U does not read "kB"
		{"212kilobit", "%s", infounit.Kilobit * 212, ""},
		{"213KiLoBITS", "%s", infounit.Kilobit * 213, ""},
		{"214 kilobit", "%s", infounit.Kilobit * 214, ""},
		{"215 kilobits", "%s", infounit.Kilobit * 215, ""},
		{"220.5kbit", "%s", infounit.Bit * 220500, ""},
		{".75 kbit", "%s", infounit.Bit * 750, ""},
		{"00.777 kilobits", "%s", infounit.Bit * 777, ""},

		{"310Mbit", "%s", infounit.Megabit * 310, ""},
		{"310Mbit", "%S", infounit.Mebibit * 310, ""},
		{"310Mbit", "%u", infounit.Megabit * 310, ""},
		{"310Mbit", "%U", infounit.Mebibit * 310, ""},
		{"311 mbit", "%s", infounit.Megabit * 311, ""},
		{"311 Mbit", "%S", infounit.Mebibit * 311, ""},
		{"312megabits", "%s", infounit.Megabit * 312, ""},
		{"312 megabits", "%S", infounit.Mebibit * 312, ""},
		{"320.25 Mbit", "%s", infounit.Kilobit * 320250, ""},
		{"0.000567 megabits", "%s", infounit.Bit * 567, ""},
		{"30000.0megabit", "%s", infounit.Gigabit * 30, ""},

		{"410GBIT", "%s", infounit.Gigabit * 410, ""},
		{"410gbit", "%S", infounit.Gibibit * 410, ""},
		{"410GBit", "%u", infounit.Gigabit * 410, ""},
		{"410Gbit", "%U", infounit.Gibibit * 410, ""},
		{"411 Gbit", "%s", infounit.Gigabit * 411, ""},
		{"411 gBit", "%S", infounit.Gibibit * 411, ""},
		{"412gigaBIT", "%s", infounit.Gigabit * 412, ""},
		{"412 gigaBIts", "%S", infounit.Gibibit * 412, ""},
		{"420.001 GBit", "%s", infounit.Megabit * 420001, ""},
		{"0.00000001GBit", "%s", infounit.Bit * 10, ""},

		{"510TBIT", "%s", infounit.Terabit * 510, ""},
		{"510tbit", "%S", infounit.Tebibit * 510, ""},
		{"510Tbit", "%u", infounit.Terabit * 510, ""},
		{"510Tbit", "%U", infounit.Tebibit * 510, ""},
		{"511 terabits", "%s", infounit.Terabit * 511, ""},
		{"511 terabits", "%S", infounit.Tebibit * 511, ""},
		{"0.000000012 Tbit", "%s", infounit.Kilobit * 12, ""},

		{"610Pbit", "%s", infounit.Petabit * 610, ""},
		{"610pbit", "%S", infounit.Pebibit * 610, ""},
		{"611 petabit", "%s", infounit.Petabit * 611, ""},
		{"611 petaBITs", "%S", infounit.Pebibit * 611, ""},
		{"18446 Pbit", "%s", infounit.Petabit * 18446, ""},
		{"18.446 Ebit", "%s", infounit.BitCount(float64(infounit.Exabit) * 18.446), ""},

		{"11Ebit", "%s", infounit.Exabit * 11, ""},
		{"11ebit", "%S", infounit.Exbibit * 11, ""},
		{"18.2 Ebit", "%s", infounit.Petabit * 18200, ""},
		{"15.5 Ebit", "%S", infounit.Exbibit*15 + infounit.Pebibit*512, ""},

		{"260kibit", "%s", infounit.Kibibit * 260, ""},
		{"260Kibit", "%S", infounit.Kibibit * 260, ""},
		{"260Kibit", "%u", infounit.Kibibit * 260, ""},
		{"260KiBit", "%U", infounit.Kibibit * 260, ""},
		{"261 KiBit", "%s", infounit.Kibibit * 261, ""},
		{"261 KiBit", "%S", infounit.Kibibit * 261, ""},
		{"262 kibibit", "%s", infounit.Kibibit * 262, ""},
		{"263 kiBIBIts", "%s", infounit.Kibibit * 263, ""},

		{"360Mibit", "%s", infounit.Mebibit * 360, ""},
		{"360mibit", "%S", infounit.Mebibit * 360, ""},
		{"360MiBit", "%u", infounit.Mebibit * 360, ""},
		{"360MiBit", "%U", infounit.Mebibit * 360, ""},
		{"361 MiBit", "%s", infounit.Mebibit * 361, ""},
		{"361 MiBit", "%S", infounit.Mebibit * 361, ""},
		{"362 mebibit", "%s", infounit.Mebibit * 362, ""},
		{"363 mebibits", "%s", infounit.Mebibit * 363, ""},

		{"460GiBit", "%s", infounit.Gibibit * 460, ""},
		{"460gibit", "%S", infounit.Gibibit * 460, ""},
		{"460GiBit", "%u", infounit.Gibibit * 460, ""},
		{"460GiBit", "%U", infounit.Gibibit * 460, ""},
		{"461 GiBit", "%s", infounit.Gibibit * 461, ""},
		{"461 GiBit", "%S", infounit.Gibibit * 461, ""},
		{"462 gibibit", "%s", infounit.Gibibit * 462, ""},
		{"463 gibibits", "%s", infounit.Gibibit * 463, ""},

		{"560TiBit", "%s", infounit.Tebibit * 560, ""},
		{"560tibit", "%S", infounit.Tebibit * 560, ""},
		{"560TiBit", "%u", infounit.Tebibit * 560, ""},
		{"560TiBit", "%U", infounit.Tebibit * 560, ""},
		{"561 TiBit", "%s", infounit.Tebibit * 561, ""},
		{"561 TiBit", "%S", infounit.Tebibit * 561, ""},
		{"562 tebibit", "%s", infounit.Tebibit * 562, ""},
		{"563 tebibits", "%s", infounit.Tebibit * 563, ""},

		{"660PiBit", "%s", infounit.Pebibit * 660, ""},
		{"660pibit", "%S", infounit.Pebibit * 660, ""},
		{"660PiBit", "%u", infounit.Pebibit * 660, ""},
		{"660PiBit", "%U", infounit.Pebibit * 660, ""},
		{"661 PiBit", "%s", infounit.Pebibit * 661, ""},
		{"661 PiBit", "%S", infounit.Pebibit * 661, ""},
		{"662 pebibit", "%s", infounit.Pebibit * 662, ""},
		{"663 pebibits", "%s", infounit.Pebibit * 663, ""},

		{"10EiBit", "%s", infounit.Exbibit * 10, ""},
		{"10eibit", "%S", infounit.Exbibit * 10, ""},
		{"12.5 EiBit", "%s", infounit.Exbibit*12 + infounit.Pebibit*512, ""},
		{"13.5 EiBit", "%S", infounit.Exbibit*13 + infounit.Pebibit*512, ""},

		{"0.00001 kbit", "%s", 0, ""},
		{"999", "%t", 0, "unknown verb for BitCount: %t"},
		{"", "%s", 0, "%s: no input"},
		{"many bits", "%s", 0, "%s: invalid expr: many"},
		{"0.1bit", "%s", 0, "%s: non-integer bit count: 0.1"},
		{"0.1bit", "%S", 0, "%S: non-integer bit count: 0.1"},
		{"0.1bit", "%u", 0, "%u: non-integer bit count: 0.1"},
		{"0.1bit", "%U", 0, "%U: non-integer bit count: 0.1"},
		{"0.12 bit", "%s", 0, "%s: non-integer bit count: 0.12"},
		{"0.12 bit", "%S", 0, "%S: non-integer bit count: 0.12"},
		{"0.123bits", "%s", 0, "%s: non-integer bit count: 0.123"},
		{"0.123bits", "%S", 0, "%S: non-integer bit count: 0.123"},
		{"999", "%s", 0, "%s: no unit suffix: EOF"},
		{"999", "%S", 0, "%S: no unit suffix: EOF"},
		{"999  Gbit", "%s", 0, "%s: no unit suffix"},
		{"999 666", "%s", 0, "%s: invalid unit expr: 666"},
		{"+9999", "%s", 0, "%s: invalid expr: +9999"},
		{"999-666", "%s", 0, "%s: invalid expr: 999-666"},
		{"999 megabits", "%u", infounit.Bit * 999, ""}, // %u does not read "megabits"
		{"999 megabits", "%U", infounit.Bit * 999, ""}, // %U does not read "megabits"
		{"1.21jigowatts", "%s", 0, "%s: unknown unit: jigowatts"},
		{"1.21jigowatts", "%S", 0, "%S: unknown unit: jigowatts"},
		{"1.21jigowatts", "%u", 0, "%u: unknown unit: jigowatts"},
		{"1.21jigowatts", "%U", 0, "%U: unknown unit: jigowatts"},
		{"1.21 jigowatts", "%s", 0, "%s: unknown unit: jigowatts"},
		{"1.21 jigowatts", "%S", 0, "%S: unknown unit: jigowatts"},
		{"1.21 jigowatts", "%u", 0, "%u: non-integer bit count: 1.21"}, // %u does not read "jigowatts"
		{"1.21 jigowatts", "%U", 0, "%U: non-integer bit count: 1.21"}, // %U does not read "jigowatts"

		{"0", "%d", infounit.Bit * 0, ""},
		{"1", "%d", infounit.Bit * 1, ""},
		{"18446744073709551615", "%d", infounit.Bit * 18446744073709551615, ""},
		{"0", "%x", infounit.Bit * 0, ""},
		{"ffffffffffffffff", "%x", infounit.Bit * 0xffffffffffffffff, ""},
		{"DeadBeef", "%X", infounit.Bit * 0xdeadbeef, ""},
		{"11110000", "%b", infounit.Bit * 0xf0, ""},
		{"1111111111111111111111111111111111111111111111111111111111111111", "%b", infounit.Bit * 0xffffffffffffffff, ""},
		{"77777777", "%o", infounit.Bit * 0xffffff, ""},
	}

	for _, c := range tc {
		var bc infounit.BitCount
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
