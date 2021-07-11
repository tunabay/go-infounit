// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package infounit

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"sync/atomic"
	"time"
)

// BitCount represents a non-negative bit count. BitCount values can be
// converted to human-readable string representations by the standard Printf
// family functions in the package fmt. See the documentation of Format method
// bellow for details.
//
// Range: 0 bits through 18446744073709551615 bits (=2 EiB)
type BitCount uint64

// Common BitCount values for units with SI and binary prefixes. To convert an
// integer of specific unit to a BitCount, multiply:
//
// 	gbits := 100
// 	fmt.Print(infounit.BitCount(gbits) * infounit.Gigabit)
const (
	Bit     BitCount = 1              // bit
	Kilobit          = 1000 * Bit     // kbit, kilobit
	Megabit          = 1000 * Kilobit // Mbit, megabit
	Gigabit          = 1000 * Megabit // Gbit, gigabit
	Terabit          = 1000 * Gigabit // Tbit, terabit
	Petabit          = 1000 * Terabit // Pbit, petabit
	Exabit           = 1000 * Petabit // Ebit, exabit
	Kibibit          = 1024 * Bit     // KiB, kibibit
	Mebibit          = 1024 * Kibibit // MiB, mebibit
	Gibibit          = 1024 * Mebibit // GiB, gibibit
	Tebibit          = 1024 * Gibibit // TiB, tebibit
	Pebibit          = 1024 * Tebibit // PiB, pebibit
	Exbibit          = 1024 * Pebibit // EiB, exbibit
)

// String returns the human-readable string representing the bit count using SI
// prefix. This implements the Stringer interface in the package fmt.
func (bc BitCount) String() string {
	return fmt.Sprintf("% .1s", bc)
}

// GoString returns a string representation of the BitCount value in Go syntax
// format. This implements the GoStringer interface in the package fmt.
func (bc BitCount) GoString() string {
	return fmt.Sprintf("BitCount(%d)", uint64(bc))
}

// ByteCount returns the value converted to the number of bytes and the
// number of remaining bits.
func (bc BitCount) ByteCount() (ByteCount, BitCount) {
	return ByteCount(uint64(bc) >> 3), bc & 0x7
}

// Convert converts the bit count to a float value in the specified unit. If the
// goal is to output or to create a string in a human-readable format,
// fmt.Printf or fmt.Sprintf are preferred.
func (bc BitCount) Convert(unit BitCount) float64 {
	return float64(bc) / float64(unit)
}

// ConvertRound is the same as Convert except that it returns a value rounded to
// the specified precision. If the goal is to output or to create a string in a
// human-readable format, fmt.Printf or fmt.Sprintf is preferred.
func (bc BitCount) ConvertRound(unit BitCount, precision int) float64 {
	p := math.Pow(10, float64(precision))
	v := math.Round(float64(bc)*p/float64(unit)) / p
	return v
}

// CalcTime calculates the duration it takes to transfer or process the number
// of bits at the specified rate.
func (bc BitCount) CalcTime(rate BitRate) (time.Duration, error) {
	if rate == 0 {
		return 0, ErrDivZeroBitRate
	}
	ns := float64(bc) * float64(time.Second) / float64(rate)
	if ns < float64(math.MinInt64) || float64(math.MaxInt64) < ns {
		return 0, ErrOutOfRange
	}
	return time.Duration(ns), nil
}

// CalcBitRate calculates the bit rate when the number of bits is transferred or
// processed in the specified duration.
func (bc BitCount) CalcBitRate(duration time.Duration) BitRate {
	if duration == 0 {
		if bc == 0 {
			return 0
		}
		return BitRate(math.Inf(+1))
	}
	return BitRate(float64(bc) / duration.Seconds())
}

// AtomicAddBitCount atomically adds delta to *addr and returns the new value.
// A wrapper function for the package sync/atomic.
func AtomicAddBitCount(addr *BitCount, delta BitCount) BitCount {
	return BitCount(atomic.AddUint64((*uint64)(addr), uint64(delta)))
}

// AtomicSubBitCount atomically subtract delta from *addr and returns the new
// value. A wrapper function for the package sync/atomic.
func AtomicSubBitCount(addr *BitCount, delta BitCount) BitCount {
	return BitCount(atomic.AddUint64((*uint64)(addr), ^uint64(delta-1)))
}

/* Does anyone want this?
// AtomicCompareAndSwapBitCount atomically executes the compare-and-swap
// operation for a BitCount value. A wrapper function for the
// package sync/atomic.
func AtomicCompareAndSwapBitCount(addr *BitCount, old, new BitCount) bool {
	return atomic.CompareAndSwapUint64((*uint64)(addr), uint64(old), uint64(new))
}
*/

// AtomicLoadBitCount atomically loads *addr. A wrapper function for the
// package sync/atomic.
func AtomicLoadBitCount(addr *BitCount) BitCount {
	return BitCount(atomic.LoadUint64((*uint64)(addr)))
}

// AtomicStoreBitCount atomically store val into *addr. A wrapper function for
// the package sync/atomic.
func AtomicStoreBitCount(addr *BitCount, val BitCount) {
	atomic.StoreUint64((*uint64)(addr), uint64(val))
}

// AtomicSwapBitCount atomically stores new into *addr and returns the previous
// *addr value. A wrapper function for the package sync/atomic.
func AtomicSwapBitCount(addr *BitCount, new BitCount) BitCount {
	return BitCount(atomic.SwapUint64((*uint64)(addr), uint64(new)))
}

// MarshalBinary encodes the BitCount value into a binary form and returns the
// result. This implements the BinaryMarshaler interface in the
// package encoding.
func (bc *BitCount) MarshalBinary() ([]byte, error) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(AtomicLoadBitCount(bc)))
	return b, nil
}

// UnmarshalBinary decodes the BitCount value from a binary form. This
// implements the BinaryUnmarshaler interface in the package encoding.
func (bc *BitCount) UnmarshalBinary(data []byte) error {
	if len(data) != 8 {
		return fmt.Errorf("invalid len: %d", len(data))
	}
	AtomicStoreBitCount(bc, BitCount(binary.BigEndian.Uint64(data)))
	return nil
}

// MarshalText encodes the BitCount value into a UTF-8-encoded text and returns
// the result. This implements the TextMarshaler interface in the
// package encoding.
func (bc *BitCount) MarshalText() ([]byte, error) {
	v := uint64(AtomicLoadBitCount(bc))
	return ([]byte)(fmt.Sprintf("%d bit", v)), nil
}

// UnmarshalText decodes the BitCount value from a UTF-8-encoded text form. This
// implements the TextUnmarshaler interface in the package encoding.
func (bc *BitCount) UnmarshalText(text []byte) error {
	var val BitCount
	if _, err := fmt.Sscanf(string(text), "%s", &val); err != nil {
		return err
	}
	AtomicStoreBitCount(bc, val)
	return nil
}

// MarshalYAML encodes the BitCount value into a uint64 for a YAML field.
func (bc *BitCount) MarshalYAML() (interface{}, error) {
	return uint64(AtomicLoadBitCount(bc)), nil
}

// UnmarshalYAML decodes the BitCount value from a YAML field.
func (bc *BitCount) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var u64 uint64
	if unmarshal(&u64) == nil {
		AtomicStoreBitCount(bc, BitCount(u64))

		return nil
	}

	var s string
	if unmarshal(&s) == nil {
		v, err := ParseBitCount(s)
		if err != nil {
			return fmt.Errorf("%q: %w: %v", s, ErrMalformedRepresentation, err)
		}
		AtomicStoreBitCount(bc, v)

		return nil
	}

	return fmt.Errorf("%w: unexpected type", ErrMalformedRepresentation)
}

// IsZero returns whether the BitCount value is zero.
func (bc BitCount) IsZero() bool {
	return bc == 0
}

// MarshalJSON encodes the BitCount value into a string for a JSON field.
func (bc *BitCount) MarshalJSON() ([]byte, error) {
	return json.Marshal(AtomicLoadBitCount(bc))
}

// UnmarshalJSON decodes the BitCount value from a JSON field.
func (bc *BitCount) UnmarshalJSON(b []byte) error {
	if string(b) == jsonNULL {
		return nil
	}

	var u64 uint64
	if json.Unmarshal(b, &u64) == nil {
		AtomicStoreBitCount(bc, BitCount(u64))

		return nil
	}

	var s string
	if json.Unmarshal(b, &s) == nil {
		v, err := ParseBitCount(s)
		if err != nil {
			return fmt.Errorf("%q: %w: %v", s, ErrMalformedRepresentation, err)
		}
		AtomicStoreBitCount(bc, v)

		return nil
	}

	return fmt.Errorf("%w: unexpected type", ErrMalformedRepresentation)
}

//
const (
	unitBitFull = "bit"
	unitBitAbbr = unitBitFull
)

// Format implements the Formatter interface in the package fmt to format
// BitCount values. This gives the ability to format BitCount values in
// human-readable format using standard Printf family functions in the
// package fmt; fmt.Printf, fmt.Fprintf, fmt.Sprintf, fmt.Errorf, and functions
// derived from them.
//
// For BitCount type, two custom 'verbs' are implemented:
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
// Regardless of the precision specified, while the unit is bit,
// no decimal parts are printed.
//
// The following flags are also available for both %s and %S:
//
// 	' '	(space) print a space between digits and unit; e.g. "12.3 kbit"
// 	#	use long unit name; e.g. "kilobit", "mebibit"
// 	-	pad with spaces on the right rather than the left (left-justify)
// 	0	pad with leading zeros rather than spaces
//
// %v prints in the default format:
//
// 	%v	default format, same as "% .1s"
// 	%#v	GoString(); e.g. "BitCount(1024)"
//
// The following uint64 compatible verbs are also supported.
// They print the integer values always in bit:
//
// 	%b	base 2
// 	%d	base 10
// 	%o	base 8
// 	%x	base 16, with lower-case letters for a-f
// 	%X	base 16, with upper-case letters for A-F
//
// See the package fmt documentation for details.
func (bc BitCount) Format(s fmt.State, verb rune) {
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
		expr := pfx.formatUint(uint64(bc), prec, full, space, unitBitAbbr, unitBitFull)
		fmt.Fprintf(s, tFmt, expr)

	case 'v':
		if s.Flag(int('#')) {
			fmt.Fprint(s, bc.GoString())
			break
		}
		fmt.Fprint(s, bc.String())

	case 'b', 'd', 'o', 'x', 'X':
		tFmt := "%"
		for _, flag := range " #+-0" {
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
		fmt.Fprintf(s, "%%!%c(BitCount=%d)", verb, uint64(bc))
	}
}

//
type bitCountScanUnitEnt struct {
	re  *regexp.Regexp
	bcs uint64
	bcb uint64
}

var (
	bitCountScanTokenRe []*regexp.Regexp
	bitCountScanUnitRe  []bitCountScanUnitEnt
)

//
func init() {
	ent := func(s string, bcs, bcb BitCount) bitCountScanUnitEnt {
		return bitCountScanUnitEnt{
			re:  regexp.MustCompile(`(?i)^` + s + `$`),
			bcs: uint64(bcs),
			bcb: uint64(bcb),
		}
	}
	bitCountScanUnitRe = []bitCountScanUnitEnt{
		ent("bits?", Bit, Bit),
		ent("k(ilo)?bits?", Kilobit, Kibibit),
		ent("m(ega)?bits?", Megabit, Mebibit),
		ent("g(iga)?bits?", Gigabit, Gibibit),
		ent("t(era)?bits?", Terabit, Tebibit),
		ent("p(eta)?bits?", Petabit, Pebibit),
		ent("e(xa)?bits?", Exabit, Exbibit),
		ent("(ki|kibi)bits?", Kibibit, Kibibit),
		ent("(mi|mebi)bits?", Mebibit, Mebibit),
		ent("(gi|gibi)bits?", Gibibit, Gibibit),
		ent("(ti|tebi)bits?", Tebibit, Tebibit),
		ent("(pi|pebi)bits?", Pebibit, Pebibit),
		ent("(ei|exbi)bits?", Exbibit, Exbibit),
	}
	bitCountScanTokenRe = []*regexp.Regexp{
		regexp.MustCompile(`(?i)^(([0-9]*)(\.[0-9]+)?)([a-z]*)$`), // 1:num, 2:int, 3:frac, 4:unit
		regexp.MustCompile(`(?i)^([a-z]*)$`),                      // 1:unit
	}
}

// Scan implements the Scanner interface in the package fmt to scan BitCount
// values from strings. This allows BitCount values to be scanned from
// human-readable string representations with unit suffixes using the standard
// Scanf family functions in the package fmt; fmt.Scanf, fmt.Fscanf, and
// fmt.Sscanf.
//
// For BitCount type, four custom 'verbs' are implemented:
//
// 	%s, %u	human-readable formats with both SI and binary prefixes
// 	%S, %U	treat SI prefix as binary prefix; 1 kilobits = 1024 bits
//
// Note that, unlike Format, the %s verb can properly scan expressions with
// units using both SI and binary prefixes.
//
// Therefore, it is usually recommended to scan using only the %s verb. The %S
// verb is the same as %s, except that it treats the SI prefix as binary prefix.
// That is, %S scans the expression "100 kbit" as 100 Kibit (=102400 bit).
//
// For verbs %s and %S, unit suffix is mandatory. If the first token consists
// only of digits, it is assumed that the next token is a unit suffix, with one
// space in between. On the other hand, %u and %U do not allow expressions with
// a space between digits and the unit suffix. They always scan only one token.
// They assume that if the token consists only of digits, it is the number of
// bits.
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
func (bc *BitCount) Scan(state fmt.ScanState, verb rune) error {
	// fmt.Printf("**scan[%c]**\n", verb)
	switch verb {
	case 'b', 'd', 'o', 'x', 'X':
		tFmt := "%"
		if wid, ok := state.Width(); ok {
			tFmt += strconv.FormatInt(int64(wid), 10)
		}
		tFmt += string(verb)
		ptr := (*uint64)(bc)
		if _, err := fmt.Fscanf(state, tFmt, ptr); err != nil {
			return fmt.Errorf("%%%c: no input: %w", verb, err)
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
		token1 := bitCountScanTokenRe[0].FindStringSubmatch(token1Str)
		if token1 == nil {
			return fmt.Errorf("%%%c: invalid expr: %s", verb, token1Str)
		}
		// fmt.Printf("[SCAN] TOKEN1: %+v\n", token1)

		numExpr, unitExpr := token1[1], token1[4]
		isInt := 0 < len(token1[2]) && len(token1[3]) < 1

		if len(numExpr) < 1 {
			return fmt.Errorf("%%%c: invalid expr: %s", verb, token1Str)
		}

		if unitExpr == "" { // no unit suffix within the first token
			switch verb {
			case 'u', 'U':
				// does not read the second token, assumed to be bits
				unitExpr = "bit"
			case 's', 'S':
				sp, n, err := state.ReadRune()
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
				token2 := bitCountScanTokenRe[1].FindStringSubmatch(token2Str)
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

		// unit is bit
		if bitCountScanUnitRe[0].re.MatchString(unitExpr) {
			if !isInt {
				return fmt.Errorf("%%%c: non-integer bit count: %s", verb, numExpr)
			}
			numVal, err := strconv.ParseUint(numExpr, 10, 64)
			if err != nil {
				return fmt.Errorf("%%%c: invalid bit count: %s: %w", verb, numExpr, err)
			}
			*ptr = numVal
			return nil
		}

		if isInt { // integer
			numVal, err := strconv.ParseUint(numExpr, 10, 64)
			if err != nil {
				return fmt.Errorf("%%%c: invalid bit count: %s: %w", verb, numExpr, err)
			}
			for _, unit := range bitCountScanUnitRe {
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
			return fmt.Errorf("%%%c: invalid bit count: %s: %w", verb, numExpr, err)
		}
		for _, unit := range bitCountScanUnitRe {
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
		return fmt.Errorf("unknown verb for BitCount: %%%c", verb)
	}
	return nil
}

// ParseBitCount converts a human-readable string representation into a BitCount
// value. The human-readable string is a decimal number with a unit suffix. SI
// and binary prefixes are correctly recognized.
func ParseBitCount(s string) (BitCount, error) {
	var v BitCount
	if _, err := fmt.Sscanf(s, "%s", &v); err != nil {
		return 0, fmt.Errorf("invalid bit count: %s: %w", s, err)
	}
	return v, nil
}

// ParseBitCountBinary is the same as ParseBitCount except that it treats the SI
// prefixes as binary prefixes. That is, it parses "100 kbit" as 100 Kibit
// (=102400 bit).
func ParseBitCountBinary(s string) (BitCount, error) {
	var v BitCount
	if _, err := fmt.Sscanf(s, "%S", &v); err != nil {
		return 0, fmt.Errorf("invalid bit count: %s: %w", s, err)
	}
	return v, nil
}
