// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package infounit

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"regexp"
	"strconv"
	"sync/atomic"
	"unsafe"
)

// BitRate represents a number of bits that are transferred or processed per
// unit of time. BitRate values can be converted to human-readable string
// representations by the standard Printf family functions in the package fmt.
// See the documentation of Format method bellow for details.
type BitRate float64

// Common BitRate values for units with SI and binary prefixes. To convert a
// float value of specific unit to a BitRate, multiply:
//
// 	gbps := 100
// 	fmt.Print(infounit.BitRate(gbps) * infounit.GigabitPerSecond)
const (
	BitPerSecond     BitRate = 1                       // bit/s, bit per second
	KilobitPerSecond         = 1000 * BitPerSecond     // kbit/s, kilobit per second
	MegabitPerSecond         = 1000 * KilobitPerSecond // Mbit/s, megabit per second
	GigabitPerSecond         = 1000 * MegabitPerSecond // Gbit/s, gigabit per second
	TerabitPerSecond         = 1000 * GigabitPerSecond // Tbit/s, terabit per second
	PetabitPerSecond         = 1000 * TerabitPerSecond // Pbit/s, petabit per second
	ExabitPerSecond          = 1000 * PetabitPerSecond // Ebit/s, exabit per second
	KibibitPerSecond         = 1024 * BitPerSecond     // Kibit/s, kibibit per second
	MebibitPerSecond         = 1024 * KibibitPerSecond // Mibit/s, mebibit per second
	GibibitPerSecond         = 1024 * MebibitPerSecond // Gibit/s, gibibit per second
	TebibitPerSecond         = 1024 * GibibitPerSecond // Tibit/s, tebibit per second
	PebibitPerSecond         = 1024 * TebibitPerSecond // Pibit/s, pebibit per second
	ExbibitPerSecond         = 1024 * PebibitPerSecond // Eibit/s, exbibit per second
)

// String returns the human-readable string representing the bit rate using SI
// prefix. This implements the Stringer interface in the package fmt.
func (br BitRate) String() string {
	return fmt.Sprintf("% .1s", br)
}

// GoString returns a string representation of the BitRate value in Go syntax
// format. This implements the GoStringer interface in the package fmt.
func (br BitRate) GoString() string {
	return fmt.Sprintf("BitRate(%s)", strconv.FormatFloat(float64(br), 'f', -1, 64))
}

// EstimateTimeForByteCount(total ByteCount) (time.Duration, error)
// EstimateTimeForBitCount(total BitCount) (time.Duration, error)
// EstimateByteCount(duration time.Duration) ByteCount
// EstimateBitCount(duration time.Duration) BitCount
// BitRateFromByteCount(total ByteCount, duration time.Duration) BitRate
// BitRateFromBitCount(total BitCount, duration time.Duration) BitRate

// IsInf reports whether the bit rate value is an infinity, according to sign.
// If sign > 0, IsInf reports whether the bit rate value is positive infinity.
// If sign < 0, IsInf reports whether the bit rate value is negative infinity.
// If sign == 0, IsInf reports whether the bit rate value is either infinity.
func (br BitRate) IsInf(sign int) bool {
	return math.IsInf(float64(br), sign)
}

// IsNaN reports whether the bit rate value is an IEEE 754 "not-a-number" value.
func (br BitRate) IsNaN() bool {
	return math.IsNaN(float64(br))
}

// Convert converts the bit rate to a float value in the specified unit. If the
// goal is to output or to create a string in a human-readable format,
// fmt.Printf or fmt.Sprintf is preferred.
func (br BitRate) Convert(unit BitRate) float64 {
	return float64(br) / float64(unit)
}

// ConvertRound is the same as Convert except that it returns a value rounded to
// the specified precision. If the goal is to output or to create a string in a
// human-readable format, fmt.Printf or fmt.Sprintf is preferred.
func (br BitRate) ConvertRound(unit BitRate, precision int) float64 {
	p := math.Pow(10, float64(precision))
	v := math.Round(p*float64(br)/float64(unit)) / p
	return v
}

// AtomicLoadBitRate atomically loads *addr. A wrapper function for the
// package sync/atomic.
func AtomicLoadBitRate(addr *BitRate) BitRate {
	return BitRate(math.Float64frombits(atomic.LoadUint64((*uint64)(unsafe.Pointer(addr)))))
}

// AtomicStoreBitRate atomically store val into *addr. A wrapper function for
// the package sync/atomic.
func AtomicStoreBitRate(addr *BitRate, val BitRate) {
	atomic.StoreUint64((*uint64)(unsafe.Pointer(addr)), math.Float64bits(float64(val)))
}

// AtomicSwapBitRate atomically stores new into *addr and returns the previous
// *addr value. A wrapper function for the package sync/atomic.
func AtomicSwapBitRate(addr *BitRate, new BitRate) BitRate {
	return BitRate(math.Float64frombits(atomic.SwapUint64((*uint64)(unsafe.Pointer(addr)), math.Float64bits(float64(new)))))
}

// MarshalBinary encodes the BitRate value into a binary form and returns the
// result. This implements the BinaryMarshaler interface in the
// package encoding.
func (br *BitRate) MarshalBinary() ([]byte, error) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, math.Float64bits(float64(AtomicLoadBitRate(br))))
	return b, nil
}

// UnmarshalBinary decodes the BitRate value from a binary form. This implements
// the BinaryUnmarshaler interface in the package encoding.
func (br *BitRate) UnmarshalBinary(data []byte) error {
	if len(data) != 8 {
		return fmt.Errorf("invalid len: %d", len(data))
	}
	AtomicStoreBitRate(br, BitRate(math.Float64frombits(binary.BigEndian.Uint64(data))))
	return nil
}

// MarshalText encodes the BitRate value into a UTF-8-encoded text and returns
// the result. This implements the TextMarshaler interface in the
// package encoding.
func (br *BitRate) MarshalText() ([]byte, error) {
	v := float64(AtomicLoadBitRate(br))
	return ([]byte)(strconv.FormatFloat(v, 'f', -1, 64) + " bit/s"), nil
}

// UnmarshalText decodes the BitRate value from a UTF-8-encoded text form. This
// implements the TextUnmarshaler interface in the package encoding.
func (br *BitRate) UnmarshalText(text []byte) error {
	var val BitRate
	n, err := fmt.Sscanf(string(text), "%s", &val)
	switch {
	case err != nil:
		return err
	case n != 1:
		return fmt.Errorf("invalid input")
	}
	AtomicStoreBitRate(br, val)
	return nil
}

//
const (
	unitBitRateFull       = unitBitFull
	unitBitRateAbbr       = unitBitRateFull
	unitBitRateAlt        = "bps"
	unitBitRateAbbrSuffix = "/s"
	unitBitRateLongSuffix = "per second"
)

// Format implements the Formatter interface in the package fmt to format
// BitRate values. This gives the ability to format the BitRate values in
// human-readable format using standard Printf family functions in the
// package fmt; fmt.Printf, fmt.Fprintf, fmt.Sprintf, fmt.Errorf, and functions
// derived from them.
//
// For ByteRate type, four custom 'verbs' are implemented:
//
// 	%s, %a	human-readable format with SI prefix
// 	%S, %A	human-readable format with binary prefix
//
// %s and %S use "bit/s" unit suffix; e.g. "Mbit/s", "Gibit/s"
// %a and %A use "bps" unit suffix; e.g. "Mbps", "Gibps"
//
// Width and precision can be specified to all of %s, %S, %a and %A:
//
// 	%s	default width, default precision
// 	%7s	width 7, default precision
// 	%.2s	default width, precision 2
// 	%7.2s	width 7, precision 2
// 	%7.s	width 7, precision 0
//
// The following flags are also available for %s, %S, %a and %A:
//
// 	' '	(space) print a space between digits and unit; e.g. "12.3 Mbit/s"
// 	#	use long unit name; e.g. "kilobits per second", "mebibits per second"
// 	-	pad with spaces on the right rather than the left (left-justify)
// 	0	pad with leading zeros rather than spaces
//
// %v prints in the default format:
//
// 	%v	default format, same as "% .1s"
// 	%#v	GoString(); e.g. "BitRate(1234567.89)"
//
// The following float64 compatible verbs are also supported.
// They print the float values always in bit/s:
//
// 	%b	decimalless scientific notation, e.g. -123456p-78
// 	%e	scientific notation, e.g. -1.234456e+78
// 	%E	scientific notation, e.g. -1.234456E+78
// 	%f	decimal point but no exponent, e.g. 123.456
// 	%F	synonym for %f
// 	%g	%e for large exponents, %f otherwise
// 	%G	%E for large exponents, %F otherwise
// 	%x	hexadecimal notation, e.g. -0x1.23abcp+20
// 	%X	upper-case hexadecimal notation, e.g. -0X1.23ABCP+20
//
// See the package fmt documentation for details.
func (br BitRate) Format(s fmt.State, verb rune) {

	switch verb {

	case 's', 'S', 'a', 'A':
		tFmt := "%"
		if s.Flag(int('-')) {
			tFmt += "-"
		}
		if s.Flag(int('0')) {
			tFmt += "0"
		}
		if wid, ok := s.Width(); ok {
			tFmt += strconv.FormatInt(int64(wid), 10)
		}
		tFmt += "s"
		prec, ok := s.Precision()
		if !ok {
			prec = -1
		}
		full, space := s.Flag(int('#')), s.Flag(int(' '))
		var pfx *prefix
		var uabbr, usuff string
		switch verb {
		case 's':
			pfx, uabbr, usuff = siPrefix, unitBitRateAbbr, unitBitRateAbbrSuffix
		case 'S':
			pfx, uabbr, usuff = binPrefix, unitBitRateAbbr, unitBitRateAbbrSuffix
		case 'a':
			pfx, uabbr, usuff = siPrefix, unitBitRateAlt, ""
		case 'A':
			pfx, uabbr, usuff = binPrefix, unitBitRateAlt, ""
		}
		expr := pfx.formatFloat(float64(br), prec, full, space, uabbr, usuff)
		fmt.Fprintf(s, tFmt, expr)

	case 'v':
		if s.Flag(int('#')) {
			fmt.Fprint(s, br.GoString())
			break
		}
		fmt.Fprint(s, br.String())

	case 'b', 'e', 'E', 'f', 'F', 'g', 'G', 'x', 'X':
		tFmt := "%"
		for _, flag := range []rune{' ', '#', '+', '-', '0'} {
			// fmt.Printf("FLAG[%c]\n", flag)
			if s.Flag(int(flag)) {
				tFmt += string(flag)
				// fmt.Printf("FLAG[%c]\n", flag)
			}
		}
		if wid, ok := s.Width(); ok {
			tFmt += strconv.FormatInt(int64(wid), 10)
		}
		if prec, ok := s.Precision(); ok {
			tFmt += "." + strconv.FormatInt(int64(prec), 10)
		}
		tFmt += fmt.Sprintf("%c", verb)
		// fmt.Printf("T-FMT[%s]\n", tFmt)
		fmt.Fprintf(s, tFmt, float64(br))

	default:
		fmt.Fprintf(s, "%%!%c(BitRate=%f)", verb, float64(br))

	}
}

//
type bitRateScanUnitEnt struct {
	re  *regexp.Regexp
	brs float64
	brb float64
}

var (
	bitRateScanTokenRe []*regexp.Regexp
	bitRateScanUnitRe  []bitRateScanUnitEnt
	bitRateScanUnit3Re []bitRateScanUnitEnt // 3 tokens unit suffix, e.g. "kilobits per second"
)

//
func init() {
	bitRateScanTokenRe = []*regexp.Regexp{
		regexp.MustCompile(`(?i)^(nan|[+-]inf|([0-9]*)(\.[0-9]+)?)([a-z/]*)$`), // 1:num, 2:int, 3:frac, 4:unit
		regexp.MustCompile(`(?i)^([a-z/]+)$`),                                  // 1:unit
		regexp.MustCompile(`(?i)^per$`),
		regexp.MustCompile(`(?i)^sec(ond)?$`),
	}
	ent := func(s string, brs, brb BitRate) bitRateScanUnitEnt {
		return bitRateScanUnitEnt{
			re:  regexp.MustCompile(`(?i)^` + s + `$`),
			brs: float64(brs),
			brb: float64(brb),
		}
	}
	st := "(bps|bit/s)"
	bitRateScanUnitRe = []bitRateScanUnitEnt{
		ent(st, BitPerSecond, BitPerSecond),
		ent("k"+st, KilobitPerSecond, KibibitPerSecond),
		ent("m"+st, MegabitPerSecond, MebibitPerSecond),
		ent("g"+st, GigabitPerSecond, GibibitPerSecond),
		ent("t"+st, TerabitPerSecond, TebibitPerSecond),
		ent("p"+st, PetabitPerSecond, PebibitPerSecond),
		ent("e"+st, ExabitPerSecond, ExbibitPerSecond),
		ent("ki"+st, KibibitPerSecond, KibibitPerSecond),
		ent("mi"+st, MebibitPerSecond, MebibitPerSecond),
		ent("gi"+st, GibibitPerSecond, GibibitPerSecond),
		ent("ti"+st, TebibitPerSecond, TebibitPerSecond),
		ent("pi"+st, PebibitPerSecond, PebibitPerSecond),
		ent("ei"+st, ExbibitPerSecond, ExbibitPerSecond),
	}
	bitRateScanUnit3Re = []bitRateScanUnitEnt{
		ent("bits?", BitPerSecond, BitPerSecond),
		ent("k(ilo)?bits?", KilobitPerSecond, KibibitPerSecond),
		ent("m(ega)?bits?", MegabitPerSecond, MebibitPerSecond),
		ent("g(iga)?bits?", GigabitPerSecond, GibibitPerSecond),
		ent("t(era)?bits?", TerabitPerSecond, TebibitPerSecond),
		ent("p(eta)?bits?", PetabitPerSecond, PebibitPerSecond),
		ent("e(xa)?bits?", ExabitPerSecond, ExbibitPerSecond),
		ent("(ki|kibi)bits?", KibibitPerSecond, KibibitPerSecond),
		ent("(mi|mebi)bits?", MebibitPerSecond, MebibitPerSecond),
		ent("(gi|gibi)bits?", GibibitPerSecond, GibibitPerSecond),
		ent("(ti|tebi)bits?", TebibitPerSecond, TebibitPerSecond),
		ent("(pi|pebi)bits?", PebibitPerSecond, PebibitPerSecond),
		ent("(ei|exbi)bits?", ExbibitPerSecond, ExbibitPerSecond),
	}
}

// Scan implements the Scanner interface in the package fmt to scan BitRate
// values from strings. This allows BitRate values to be scanned from
// human-readable string representations with unit suffixes using the standard
// Scanf family functions in the package fmt; fmt.Scanf, fmt.Fscanf, and
// fmt.Sscanf().
//
// For BitRate type, four custom 'verbs' are implemented:
//
// 	%s, %u	human-readable formats with both SI and binary prefixes
// 	%S, %U	treat SI prefix as binary prefix; 1 kbit/s = 1024 bit/s
//
// Note that, unlike Format, the %s verb can properly scan expressions with
// units using both SI and binary prefixes.
//
// Therefore, it is usually recommended to scan using only the %s verb. The %S
// verb is the same as %s, except that it treats the SI prefix as binary prefix.
// That is, %S scans the expression "100 kbit/s" as 100 Kibit/s (=102400 bit/s).
//
// For verbs %s and %S, unit suffix is mandatory. If the first token consists
// only of digits, it is assumed that the next token is a unit suffix, with one
// space in between. On the other hand, %u and %U do not allow expressions with
// a space between digits and the unit suffix. They always scan only one token.
// They assume that if the token consists only of digits, it is the number of
// bit/s.
//
// The following verbs are compatible with float64 and scans floating point
// values without a unit suffix. If it is clear that there is absolutely no unit
// suffix in the input, the use of these is recommended:
//
// 	%f, %F	floating point representation
//
// See the package fmt documentation for details.
func (br *BitRate) Scan(state fmt.ScanState, verb rune) error {
	// fmt.Printf("**scan[%c]**\n", verb)
	switch verb {

	case 'f', 'F':
		tFmt := "%"
		if wid, ok := state.Width(); ok {
			tFmt += strconv.FormatInt(int64(wid), 10)
		}
		tFmt += string(verb)
		ptr := (*float64)(br)
		n, err := fmt.Fscanf(state, tFmt, ptr)
		switch {
		case err != nil:
			return fmt.Errorf("%%%c: no input: %w", verb, err)
		case n != 1:
			return fmt.Errorf("%%%c: no input", verb)
		}

	case 's', 'S', 'u', 'U':
		token1Bytes, err := state.Token(true, nil)
		switch {
		case err != nil:
			return fmt.Errorf("%%%c: %w", verb, err)
		case len(token1Bytes) < 1:
			return fmt.Errorf("%%%c: no input", verb)
		}
		token1Str := string(token1Bytes)
		token1 := bitRateScanTokenRe[0].FindStringSubmatch(token1Str)
		if token1 == nil {
			return fmt.Errorf("%%%c: invalid expr: %s", verb, token1Str)
		}
		// fmt.Printf("[SCAN] TOKEN1: %+v\n", token1)

		numExpr, unitExpr := token1[1], token1[4]

		if len(numExpr) < 1 {
			return fmt.Errorf("%%%c: invalid expr: %s", verb, token1Str)
		}
		if unitExpr == "" { // no unit suffix within the first token
			switch verb {
			case 'u', 'U':
				// does not read the second token, assumed to be bit/s
				unitExpr = "bit/s"
			case 's', 'S':
				sp, n, err := state.ReadRune() // read only one space
				switch {
				case err != nil:
					return fmt.Errorf("%%%c: no unit suffix: %w", verb, err)
				case n != 1:
					return fmt.Errorf("%%%c: no unit suffix", verb)
				case sp != ' ':
					return fmt.Errorf("%%%c: no space after digits: [%c]", verb, sp)
				}
				token2Bytes, err := state.Token(false, nil)
				switch {
				case err != nil:
					return fmt.Errorf("%%%c: no unit suffix: %w", verb, err)
				case len(token2Bytes) < 1:
					return fmt.Errorf("%%%c: no unit suffix", verb)
				}
				token2Str := string(token2Bytes)
				token2 := bitRateScanTokenRe[1].FindStringSubmatch(token2Str)
				if token2 == nil {
					return fmt.Errorf("%%%c: invalid unit expr: %s", verb, token2Str)
				}
				// fmt.Printf("[SCAN] TOKEN2: %+v\n", token2)
				unitExpr = token2[1]
				if unitExpr == "" {
					return fmt.Errorf("%%%c: no unit suffix", verb)
				}
			}
		}

		// fmt.Printf("[SCAN] FIRST [%s] [%s]\n", numExpr, unitExpr)

		numVal, err := strconv.ParseFloat(numExpr, 64)
		if err != nil {
			return fmt.Errorf("%%%c: invalid expr: %s", verb, numExpr)
		}

		ptr := (*float64)(br)

		for _, unit := range bitRateScanUnitRe {
			if unit.re.MatchString(unitExpr) {
				switch verb {
				case 's', 'u':
					*ptr = numVal * unit.brs
				case 'S', 'U':
					*ptr = numVal * unit.brb
				}
				return nil
			}
		}

		// try 3 tokens units
		// 12.3kilobits per second
		// 12.3 kilobits per second
		eSuf := unitExpr
		for i := 0; i < 2; i++ {
			sp, n, err := state.ReadRune() // read only one space
			switch {
			case err == io.EOF:
				return fmt.Errorf("%%%c: unknown unit: %s", verb, eSuf)
			case err != nil:
				return fmt.Errorf("%%%c: invalid unit suffix: %s: %w", verb, eSuf, err)
			case n != 1:
				return fmt.Errorf("%%%c: unknown unit: %s", verb, eSuf)
			case sp != ' ':
				return fmt.Errorf("%%%c: unknown unit: %s%c", verb, eSuf, sp)
			}
			token34Bytes, err := state.Token(false, nil)
			switch {
			case err == io.EOF:
				return fmt.Errorf("%%%c: unknown unit: %s", verb, eSuf)
			case err != nil:
				return fmt.Errorf("%%%c: invalid unit suffix: %s: %w", verb, eSuf, err)
			case len(token34Bytes) < 1:
				return fmt.Errorf("%%%c: unknown unit: %s", verb, eSuf)
			case !bitRateScanTokenRe[2+i].Match(token34Bytes):
				return fmt.Errorf("%%%c: unknown unit: %s %s", verb, eSuf, string(token34Bytes))
			}
			eSuf += string(sp) + string(token34Bytes)
		}

		for _, unit := range bitRateScanUnit3Re {
			if unit.re.MatchString(unitExpr) {
				switch verb {
				case 's', 'u':
					*ptr = numVal * unit.brs
				case 'S', 'U':
					*ptr = numVal * unit.brb
				}
				return nil
			}
		}
		return fmt.Errorf("%%%c: unknown unit: %s", verb, eSuf)

	default:
		return fmt.Errorf("unknown verb for BitRate: %%%c", verb)

	}
	return nil
}
