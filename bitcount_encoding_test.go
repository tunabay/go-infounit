// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package infounit_test

import (
	"bytes"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/tunabay/go-infounit"
	"gopkg.in/yaml.v2"
)

//
func TestBitCount_MarshalBinary_1(t *testing.T) {
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
		bc := infounit.BitCount(c.b)
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
func TestBitCount_UnmarshalBinary_1(t *testing.T) {
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
		var bc infounit.BitCount
		bin, err := hex.DecodeString(c.hex)
		if err != nil {
			t.Fatal(err)
		}
		if err := bc.UnmarshalBinary(bin); err != nil {
			t.Error(err)
		}
		if exbc := infounit.BitCount(c.b); bc != exbc {
			t.Errorf(`%x: want: %d, got: %d`, c.hex, exbc, bc)
		}
	}
}

//
func TestBitCount_MarshalText_1(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b   uint64
		txt string
	}{
		{0, "0 bit"},
		{1, "1 bit"},
		{987654321, "987654321 bit"},
		{18446744073709551614, "18446744073709551614 bit"},
		{18446744073709551615, "18446744073709551615 bit"},
	}

	for _, c := range tc {
		bc := infounit.BitCount(c.b)
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
func TestBitCount_UnmarshalText_1(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b   uint64
		txt string
	}{
		{0, "0 bit"},
		{1, "1 bit"},
		{987654321, "987654321 bit"},
		{18446744073709551614, "18446744073709551614 bit"},
		{18446744073709551615, "18446744073709551615 bit"},
	}

	for _, c := range tc {
		var bc infounit.BitCount
		if err := bc.UnmarshalText(([]byte)(c.txt)); err != nil {
			t.Error(err)
		}
		if exbc := infounit.BitCount(c.b); bc != exbc {
			t.Errorf(`%s: want: %d, got: %d`, c.txt, exbc, bc)
		}
	}
}

func TestBitCount_MarshalYAML(t *testing.T) {
	var (
		x = infounit.BitCount(9991111)
		y = infounit.BitCount(99992222)
		z = infounit.BitCount(999993333)
	)
	v := &struct {
		Val      infounit.BitCount
		Ptr      *infounit.BitCount
		PtrNil   *infounit.BitCount
		ValSlice []infounit.BitCount
		PtrSlice []*infounit.BitCount
		Renamed  infounit.BitCount  `yaml:"xyzBC"`
		ZeroVal  infounit.BitCount  `yaml:"zeroVal,omitempty"`
		ZeroPtr  *infounit.BitCount `yaml:",omitempty"`
	}{
		Val:      infounit.BitCount(1111),
		Ptr:      &x,
		ValSlice: []infounit.BitCount{777111, 777222, 777333},
		PtrSlice: []*infounit.BitCount{&y, &z},
		Renamed:  88888888,
		ZeroVal:  infounit.BitCount(0),
		ZeroPtr:  nil,
	}

	expected := strings.Join([]string{
		"val: 1111",
		"ptr: 9991111",
		"ptrnil: null",
		"valslice:",
		"- 777111",
		"- 777222",
		"- 777333",
		"ptrslice:",
		"- 99992222",
		"- 999993333",
		"xyzBC: 88888888",
		"",
	}, "\n")

	yamlBytes, err := yaml.Marshal(v)
	if err != nil {
		t.Errorf("yaml.Marshal() failed: %v", err)
	}
	yaml := string(yamlBytes)

	if yaml != expected {
		t.Errorf("yaml.Marshal() unexpected result: %q", yaml)
	}
}

func TestBitCount_UnmarshalYAML(t *testing.T) {
	v := struct {
		Val      infounit.BitCount
		Ptr      *infounit.BitCount
		PtrNil   *infounit.BitCount
		ValSlice []infounit.BitCount
		PtrSlice []*infounit.BitCount
		Renamed  infounit.BitCount `yaml:"xyzBC"`
		VarExprs []infounit.BitCount
	}{}

	yamlSrc := strings.Join([]string{
		"val: 1111",
		"ptr: 99991111",
		"ptrnil: null",
		"valslice:",
		"- 777111",
		"- 777222",
		"- 777333",
		"ptrslice:",
		"- 999992222",
		"- 9999993333",
		"xyzBC: 88888888",
		"varexprs:",
		`- "123 kilobits"`,
		`- "345 Mbit"`,
		`- "67.8Gbit"`,
		"",
	}, "\n")

	if err := yaml.UnmarshalStrict(([]byte)(yamlSrc), &v); err != nil {
		t.Errorf("yaml.Unmarshal() failed: %v", err)
	}
	if v.Val != 1111 {
		t.Errorf("Val: unexpected value: got: %v, want: 1111 bit", v.Val)
	}
	switch {
	case v.Ptr == nil:
		t.Errorf("Ptr: unexpected value: got: <nil>, want: %d", 99991111)
	case *v.Ptr != 99991111:
		t.Errorf("Ptr: unexpected value: got: %v, want: %d", *v.Ptr, 99991111)
	}
	if v.PtrNil != nil {
		t.Errorf("PtrNil: unexpected value: got: %v, want: <nil>", *v.PtrNil)
	}
	switch {
	case len(v.ValSlice) != 3:
		t.Errorf("ValSlice: unexpected length: got: %d, want: 3", len(v.ValSlice))
	case v.ValSlice[0] != 777111:
		t.Errorf("ValSlice[0]: unexpected value: got: %d, want: 777111", v.ValSlice[0])
	case v.ValSlice[1] != 777222:
		t.Errorf("ValSlice[1]: unexpected value: got: %d, want: 777222", v.ValSlice[1])
	case v.ValSlice[2] != 777333:
		t.Errorf("ValSlice[2]: unexpected value: got: %d, want: 777333", v.ValSlice[2])
	}
	switch {
	case len(v.PtrSlice) != 2:
		t.Errorf("PtrSlice: unexpected length: got: %d, want: 2", len(v.PtrSlice))
	case *v.PtrSlice[0] != 999992222:
		t.Errorf("PtrSlice[0]: unexpected value: got: %d, want: 9999922222", *v.PtrSlice[0])
	case *v.PtrSlice[1] != 9999993333:
		t.Errorf("PtrSlice[1]: unexpected value: got: %d, want: 99999933333", *v.PtrSlice[1])
	}
	if v.Renamed != 88888888 {
		t.Errorf("Renamed: unexpected value: got: %d, want: 88888888", v.Renamed)
	}
	switch {
	case len(v.VarExprs) != 3:
		t.Errorf("VarExprs: unexpected length: got: %d, want: 3", len(v.VarExprs))
	case v.VarExprs[0] != infounit.Kilobit*123:
		t.Errorf("VarExprs[0]: unexpected value: got: %d, want: %d", v.VarExprs[0], infounit.Kilobit*123)
	case v.VarExprs[1] != infounit.Megabit*345:
		t.Errorf("VarExprs[1]: unexpected value: got: %d, want: %d", v.VarExprs[1], infounit.Megabit*345)
	case v.VarExprs[2] != infounit.Gigabit/10*678:
		t.Errorf("VarExprs[2]: unexpected value: got: %d, want: %d", v.VarExprs[2], infounit.Gigabit/10*678)
	}

	// t.Logf("%+v\n", v)
}
