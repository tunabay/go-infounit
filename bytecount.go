// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package infounit

import (
	"encoding/binary"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"sync/atomic"
	"time"
)

// ByteCount represents a non-negative byte count. ByteCount values can be
// converted to human-readable string representations by the standard Printf
// family functions in the package fmt. See the documentation of Format method
// bellow for details.
//
// Range: 0 bytes through 18446744073709551615bytes (=16 EiB)
type ByteCount uint64

// Common ByteCount values for units with SI and binary prefixes. To convert an
// integer of specific unit to a ByteCount, multiply:
//
// 	gb := 100
// 	fmt.Print(infounit.ByteCount(gb) * infounit.Gigabyte)
const (
	Byte     ByteCount = 1               // B, byte
	Kilobyte           = 1000 * Byte     // kB, kilobyte
	Megabyte           = 1000 * Kilobyte // MB, megabyte
	Gigabyte           = 1000 * Megabyte // GB, gigabyte
	Terabyte           = 1000 * Gigabyte // TB, terabyte
	Petabyte           = 1000 * Terabyte // PB, petabyte
	Exabyte            = 1000 * Petabyte // EB, exabyte
	Kibibyte           = 1024 * Byte     // KiB, kibibyte
	Mebibyte           = 1024 * Kibibyte // MiB, mebibyte
	Gibibyte           = 1024 * Mebibyte // GiB, gibibyte
	Tebibyte           = 1024 * Gibibyte // TiB, tebibyte
	Pebibyte           = 1024 * Tebibyte // PiB, pebibyte
	Exbibyte           = 1024 * Pebibyte // EiB, exbibyte
)

// String returns the human-readable string representing the byte count using SI
// prefix. This implements the Stringer interface in the package fmt.
func (bc ByteCount) String() string {
	return fmt.Sprintf("% .1s", bc)
}

// GoString returns a string representation of the ByteCount value in Go syntax
// format. This implements the GoStringer interface in the package fmt.
func (bc ByteCount) GoString() string {
	return fmt.Sprintf("ByteCount(%d)", uint64(bc))
}

// BitCount returns the value converted to the number of bits. If the number of
// bits is too large, an ErrOutOfRange will be returned.
func (bc ByteCount) BitCount() (BitCount, error) {
	if bc&(0b_111<<61) != 0 {
		return BitCount(0), ErrOutOfRange
	}
	return BitCount(uint64(bc) << 3), nil
}

// Convert converts the byte count to a float value in the specified unit. If
// the goal is to output or to create a string in a human-readable format,
// fmt.Printf or fmt.Sprintf is preferred.
func (bc ByteCount) Convert(unit ByteCount) float64 {
	return float64(bc) / float64(unit)
}

// ConvertRound is the same as Convert except that it returns a value rounded to
// the specified precision. If the goal is to output or to create a string in a
// human-readable format, fmt.Printf or fmt.Sprintf is preferred.
func (bc ByteCount) ConvertRound(unit ByteCount, precision int) float64 {
	p := math.Pow(10, float64(precision))
	v := math.Round(p*float64(bc)/float64(unit)) / p
	return v
}

// CalcTime calculates the duration it takes to transfer or process the number
// of bytes at the specified rate.
func (bc ByteCount) CalcTime(rate BitRate) (time.Duration, error) {
	if rate == 0 {
		return 0, ErrDivZeroBitRate
	}
	ns := float64(bc) * 8 * float64(time.Second) / float64(rate)
	if ns < float64(math.MinInt64) || float64(math.MaxInt64) < ns {
		return 0, ErrOutOfRange
	}
	return time.Duration(ns), nil
}

// CalcBitRate calculates the bit rate when the number of bytes is transferred
// or processed in the specified duration.
func (bc ByteCount) CalcBitRate(duration time.Duration) BitRate {
	if duration == 0 {
		if bc == 0 {
			return 0
		}
		return BitRate(math.Inf(+1))
	}
	return BitRate(float64(bc) * 8 / duration.Seconds())
}

// AtomicAddByteCount atomically adds delta to *addr and returns the new value.
// A wrapper function for the package sync/atomic.
func AtomicAddByteCount(addr *ByteCount, delta ByteCount) ByteCount {
	return ByteCount(atomic.AddUint64((*uint64)(addr), uint64(delta)))
}

// AtomicSubByteCount atomically subtract delta from *addr and returns the new
// value. A wrapper function for the package sync/atomic.
func AtomicSubByteCount(addr *ByteCount, delta ByteCount) ByteCount {
	return ByteCount(atomic.AddUint64((*uint64)(addr), ^uint64(delta-1)))
}

/* Does anyone want this?
// AtomicCompareAndSwapByteCount atomically executes the compare-and-swap
// operation for a ByteCount value. A wrapper function for the
// package sync/atomic.
func AtomicCompareAndSwapByteCount(addr *ByteCount, old, new ByteCount) bool {
	return atomic.CompareAndSwapUint64((*uint64)(addr), uint64(old), uint64(new))
}
*/

// AtomicLoadByteCount atomically loads *addr. A wrapper function for the
// package sync/atomic.
func AtomicLoadByteCount(addr *ByteCount) ByteCount {
	return ByteCount(atomic.LoadUint64((*uint64)(addr)))
}

// AtomicStoreByteCount atomically store val into *addr. A wrapper function for
// the package sync/atomic.
func AtomicStoreByteCount(addr *ByteCount, val ByteCount) {
	atomic.StoreUint64((*uint64)(addr), uint64(val))
}

// AtomicSwapByteCount atomically stores new into *addr and returns the previous
// *addr value. A wrapper function for the package sync/atomic.
func AtomicSwapByteCount(addr *ByteCount, new ByteCount) ByteCount {
	return ByteCount(atomic.SwapUint64((*uint64)(addr), uint64(new)))
}

// MarshalBinary encodes the ByteCount value into a binary form and returns the
// result. This implements the BinaryMarshaler interface in the
// package encoding.
func (bc *ByteCount) MarshalBinary() ([]byte, error) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(AtomicLoadByteCount(bc)))
	return b, nil
}

// UnmarshalBinary decodes the ByteCount value from a binary form. This
// implements the BinaryUnmarshaler interface in the package encoding.
func (bc *ByteCount) UnmarshalBinary(data []byte) error {
	if len(data) != 8 {
		return fmt.Errorf("invalid len: %d", len(data))
	}
	AtomicStoreByteCount(bc, ByteCount(binary.BigEndian.Uint64(data)))
	return nil
}

// MarshalText encodes the ByteCount value into a UTF-8-encoded text and returns
// the result. This implements the TextMarshaler interface in the
// package encoding.
func (bc *ByteCount) MarshalText() ([]byte, error) {
	return ([]byte)(fmt.Sprintf("%d B", uint64(AtomicLoadByteCount(bc)))), nil
}

// UnmarshalText decodes the ByteCount value from a UTF-8-encoded text form.
// This implements the TextUnmarshaler interface in the package encoding.
func (bc *ByteCount) UnmarshalText(text []byte) error {
	var val ByteCount
	n, err := fmt.Sscanf(string(text), "%s", &val)
	switch {
	case err != nil:
		return err
	case n != 1:
		return fmt.Errorf("invalid input")
	}
	AtomicStoreByteCount(bc, val)
	return nil
}

//
const (
	unitByteFull = "byte"
	unitByteAbbr = "B"
)

// Format implements the Formatter interface in the package fmt to format
// ByteCount values. This gives the ability to format ByteCount values in
// human-readable format using standard Printf family functions in the
// package fmt; fmt.Printf, fmt.Fprintf, fmt.Sprintf, fmt.Errorf, and functions
// derived from them.
//
// For ByteCount type, two custom 'verbs' are implemented:
//
// 	%s	human-readable format with SI prefix
// 	%S	human-readable format with binary prefix
//
// Width and precision can be specified to both %s and %S:
//
// 	%s	default width, default precision
// 	%7s	width 7, default precision
// 	%.2s	default width, precision 2
// 	%7.2s	width 7, precision 2
// 	%7.s	width 7, precision 0
//
// Regardless of the precision specified, while the unit is byte,
// no decimal parts are printed.
//
// The following flags are also available for both %s and %S:
//
// 	' '	(space) print a space between digits and unit; e.g. "12.3 kB"
// 	#	use long unit name; e.g. "kilobyte", "mebibyte"
// 	-	pad with spaces on the right rather than the left (left-justify)
// 	0	pad with leading zeros rather than spaces
//
// %v prints in the default format:
//
// 	%v	default format, same as "% .1s"
// 	%#v	GoString(); e.g. "ByteCount(1024)"
//
// The following uint64 compatible verbs are also supported.
// They print the integer values always in byte:
//
// 	%b	base 2
// 	%d	base 10
// 	%o	base 8
// 	%x	base 16, with lower-case letters for a-f
// 	%X	base 16, with upper-case letters for A-F
//
// See the package fmt documentation for details.
func (bc ByteCount) Format(s fmt.State, verb rune) {

	switch verb {

	case 's', 'S':
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
		switch verb {
		case 's':
			pfx = siPrefix
		case 'S':
			pfx = binPrefix
		}
		expr := pfx.formatUint(uint64(bc), prec, full, space, unitByteAbbr, unitByteFull)
		fmt.Fprintf(s, tFmt, expr)

	case 'v':
		if s.Flag(int('#')) {
			fmt.Fprint(s, bc.GoString())
			break
		}
		fmt.Fprint(s, bc.String())

	case 'b', 'd', 'o', 'x', 'X':
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
		tFmt += string(verb)
		// fmt.Printf("T-FMT[%s]\n", tFmt)
		fmt.Fprintf(s, tFmt, uint64(bc))

	default:
		fmt.Fprintf(s, "%%!%c(ByteCount=%d)", verb, uint64(bc))

	}
}

//
type byteCountScanUnitEnt struct {
	re  *regexp.Regexp
	bcs uint64
	bcb uint64
}

var (
	byteCountScanTokenRe []*regexp.Regexp
	byteCountScanUnitRe  []byteCountScanUnitEnt
)

//
func init() {
	ent := func(s string, bcs, bcb ByteCount) byteCountScanUnitEnt {
		return byteCountScanUnitEnt{
			re:  regexp.MustCompile(`(?i)^` + s + `$`),
			bcs: uint64(bcs),
			bcb: uint64(bcb),
		}
	}
	byteCountScanUnitRe = []byteCountScanUnitEnt{
		ent("b(ytes?)?", Byte, Byte),
		ent("kb|kilobytes?", Kilobyte, Kibibyte),
		ent("mb|megabytes?", Megabyte, Mebibyte),
		ent("gb|gigabytes?", Gigabyte, Gibibyte),
		ent("tb|terabytes?", Terabyte, Tebibyte),
		ent("pb|petabytes?", Petabyte, Pebibyte),
		ent("eb|exabytes?", Exabyte, Exbibyte),
		ent("kib|kibibytes?", Kibibyte, Kibibyte),
		ent("mib|mebibytes?", Mebibyte, Mebibyte),
		ent("gib|gibibytes?", Gibibyte, Gibibyte),
		ent("tib|tebibytes?", Tebibyte, Tebibyte),
		ent("pib|pebibytes?", Pebibyte, Pebibyte),
		ent("eib|exbibytes?", Exbibyte, Exbibyte),
	}
	byteCountScanTokenRe = []*regexp.Regexp{
		regexp.MustCompile(`(?i)^(([0-9]*)(\.[0-9]+)?)([a-z]*)$`), // 1:num, 2:int, 3:frac, 4:unit
		regexp.MustCompile(`(?i)^([a-z]*)$`),                      // 1:unit
	}
}

// Scan implements the Scanner interface in the package fmt to scan ByteCount
// values from strings. This allows ByteCount values to be scanned from
// human-readable string representations with unit suffixes using the standard
// Scanf family functions in the package fmt; fmt.Scanf, fmt.Fscanf, and
// fmt.Sscanf.
//
// For ByteCount type, four custom 'verbs' are implemented:
//
// 	%s, %u	human-readable formats with both SI and binary prefixes
// 	%S, %U	treat SI prefix as binary prefix; 1 kilobyte = 1024 bytes
//
// Note that, unlike Format, the %s verb can properly scan expressions with
// units using both SI and binary prefixes.
//
// Therefore, it is usually recommended to scan using only the %s verb. The %S
// verb is the same as %s, except that it treats the SI prefix as binary prefix.
// That is, %S scans the expression "100 kB" as 100 KiB (=102400 B).
//
// For verbs %s and %S, unit suffix is mandatory. If the first token consists
// only of digits, it is assumed that the next token is a unit suffix, with one
// space in between. On the other hand, %u and %U do not allow expressions with
// a space between digits and the unit suffix. They always scan only one token.
// They assume that if the token consists only of digits, it is the number of
// bytes.
//
// The following verbs are compatible with uint64 and scans integers without a
// unit suffix. If it is clear that there is absolutely no unit suffix in the
// input, the use of these is recommended:
//
// 	%b	base 2
// 	%o	base 8
// 	%d	base 10
// 	%x, %X	base 16
//
// See the package fmt documentation for details.
func (bc *ByteCount) Scan(state fmt.ScanState, verb rune) error {
	// fmt.Printf("**scan[%c]**\n", verb)
	switch verb {

	case 'b', 'd', 'o', 'x', 'X':
		tFmt := "%"
		if wid, ok := state.Width(); ok {
			tFmt += strconv.FormatInt(int64(wid), 10)
		}
		tFmt += string(verb)
		ptr := (*uint64)(bc)
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
		token1 := byteCountScanTokenRe[0].FindStringSubmatch(token1Str)
		if token1 == nil {
			return fmt.Errorf("%%%c: invalid expr: %s", verb, token1Str)
		}
		// fmt.Printf("[SCAN] TOKEN1: %+v\n", token1)

		numExpr := token1[1]
		isInt := 0 < len(token1[2]) && len(token1[3]) < 1
		unitExpr := token1[4]

		if len(numExpr) < 1 {
			return fmt.Errorf("%%%c: invalid expr: %s", verb, token1Str)
		}

		if unitExpr == "" { // no unit suffix within the first token
			switch verb {
			case 'u', 'U':
				// does not read the second token, assumed to be bytes
				unitExpr = "b"
			case 's', 'S':
				sp, n, err := state.ReadRune()
				if err != nil {
					return fmt.Errorf("%%%c: no unit suffix: %w", verb, err)
				}
				if n != 1 {
					return fmt.Errorf("%%%c: no unit suffix", verb)
				}
				if sp != ' ' {
					return fmt.Errorf("%%%c: no space after digits: [%c]", verb, sp)
				}
				token2Bytes, err := state.Token(false, nil)
				if err != nil {
					return fmt.Errorf("%%%c: no unit suffix: %w", verb, err)
				}
				if len(token2Bytes) < 1 {
					return fmt.Errorf("%%%c: no unit suffix", verb)
				}
				token2Str := string(token2Bytes)
				token2 := byteCountScanTokenRe[1].FindStringSubmatch(token2Str)
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

		// fmt.Printf("[SCAN] LAST [%s] [%s]\n", numExpr, unitExpr)

		ptr := (*uint64)(bc)

		// unit is byte
		if byteCountScanUnitRe[0].re.MatchString(unitExpr) {
			if !isInt {
				return fmt.Errorf("%%%c: non-integer byte count: %s", verb, numExpr)
			}
			numVal, err := strconv.ParseUint(numExpr, 10, 64)
			if err != nil {
				return fmt.Errorf("%%%c: invalid byte count: %s: %s", verb, numExpr, err)
			}
			*ptr = numVal
			return nil
		}

		if isInt { // integer
			numVal, err := strconv.ParseUint(numExpr, 10, 64)
			if err != nil {
				return fmt.Errorf("%%%c: invalid byte count: %s: %s", verb, numExpr, err)
			}
			for _, unit := range byteCountScanUnitRe {
				if unit.re.MatchString(unitExpr) {
					switch verb {
					case 's', 'u':
						*ptr = numVal * unit.bcs
					case 'S', 'U':
						*ptr = numVal * unit.bcb
					}
					return nil
				}
			}
			return fmt.Errorf("%%%c: unknown unit: %s", verb, unitExpr)
		}

		// float
		numVal, err := strconv.ParseFloat(numExpr, 64)
		if err != nil {
			return fmt.Errorf("%%%c: invalid byte count: %s: %s", verb, numExpr, err)
		}
		for _, unit := range byteCountScanUnitRe {
			if unit.re.MatchString(unitExpr) {
				switch verb {
				case 's', 'u':
					*ptr = uint64(math.Round(numVal * float64(unit.bcs)))
				case 'S', 'U':
					*ptr = uint64(math.Round(numVal * float64(unit.bcb)))
				}
				return nil
			}
		}
		return fmt.Errorf("%%%c: unknown unit: %s", verb, unitExpr)

	default:
		return fmt.Errorf("unknown verb for ByteCount: %%%c", verb)

	}
	return nil
}

// ParseByteCount converts a human-readable string representation into a
// ByteCount value. The human-readable string is a decimal number with a unit
// suffix. SI and binary prefixes are correctly recognized.
func ParseByteCount(s string) (ByteCount, error) {
	var v ByteCount
	if _, err := fmt.Sscanf(s, "%s", &v); err != nil {
		return 0, fmt.Errorf("invalid byte count: %s", s)
	}
	return v, nil
}

// ParseByteCountBinary is the same as ParseByteCount except that it treats the
// SI prefixes as binary prefixes. That is, it parses "100 kB" as 100 KiB
// (=102400 B).
func ParseByteCountBinary(s string) (ByteCount, error) {
	var v ByteCount
	if _, err := fmt.Sscanf(s, "%S", &v); err != nil {
		return 0, fmt.Errorf("invalid byte count: %s", s)
	}
	return v, nil
}
