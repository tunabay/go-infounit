// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package infounit_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/tunabay/go-infounit"
)

//
func TestByteCount_MarshalBinary_1(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b   uint64
		hex string
	}{
		{0, "0000000000000000"},
		{1, "0000000000000001"},
		{987654321, "000000003ADE68B1"},
		{18446744073709551614, "FFFFFFFFFFFFFFFE"},
		{18446744073709551615, "FFFFFFFFFFFFFFFF"},
	}

	for _, c := range tc {
		bc := infounit.ByteCount(c.b)
		bin, err := bc.MarshalBinary()
		if err != nil {
			t.Error(err)
		}
		exbin, err := hex.DecodeString(c.hex)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(bin, exbin) {
			t.Errorf(`%s: want: %x, got: %x`, bc, exbin, bin)
		}
	}
}

//
func TestByteCount_UnmarshalBinary_1(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b   uint64
		hex string
	}{
		{0, "0000000000000000"},
		{1, "0000000000000001"},
		{987654321, "000000003ADE68B1"},
		{18446744073709551614, "FFFFFFFFFFFFFFFE"},
		{18446744073709551615, "FFFFFFFFFFFFFFFF"},
	}

	for _, c := range tc {
		var bc infounit.ByteCount
		bin, err := hex.DecodeString(c.hex)
		if err != nil {
			t.Fatal(err)
		}
		if err := bc.UnmarshalBinary(bin); err != nil {
			t.Error(err)
		}
		if exbc := infounit.ByteCount(c.b); bc != exbc {
			t.Errorf(`%x: want: %d, got: %d`, c.hex, exbc, bc)
		}
	}
}

//
func TestByteCount_MarshalText_1(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b   uint64
		txt string
	}{
		{0, "0 B"},
		{1, "1 B"},
		{987654321, "987654321 B"},
		{18446744073709551614, "18446744073709551614 B"},
		{18446744073709551615, "18446744073709551615 B"},
	}

	for _, c := range tc {
		bc := infounit.ByteCount(c.b)
		txtb, err := bc.MarshalText()
		if err != nil {
			t.Error(err)
		}
		txt := string(txtb)
		if txt != c.txt {
			t.Errorf(`%d: want: "%s", got: "%s"`, bc, c.txt, txt)
		}
	}
}

//
func TestByteCount_UnmarshalText_1(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b   uint64
		txt string
	}{
		{0, "0 B"},
		{1, "1 B"},
		{987654321, "987654321 B"},
		{18446744073709551614, "18446744073709551614 B"},
		{18446744073709551615, "18446744073709551615 B"},
	}

	for _, c := range tc {
		var bc infounit.ByteCount
		if err := bc.UnmarshalText(([]byte)(c.txt)); err != nil {
			t.Error(err)
		}
		if exbc := infounit.ByteCount(c.b); bc != exbc {
			t.Errorf(`%s: want: %d, got: %d`, c.txt, exbc, bc)
		}
	}
}
