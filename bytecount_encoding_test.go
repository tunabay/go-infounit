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

func TestByteCount_MarshalYAML(t *testing.T) {
	var (
		x = infounit.ByteCount(99991111)
		y = infounit.ByteCount(999992222)
		z = infounit.ByteCount(9999993333)
	)
	v := &struct {
		Val      infounit.ByteCount
		Ptr      *infounit.ByteCount
		PtrNil   *infounit.ByteCount
		ValSlice []infounit.ByteCount
		PtrSlice []*infounit.ByteCount
		Renamed  infounit.ByteCount  `yaml:"xyzBC"`
		ZeroVal  infounit.ByteCount  `yaml:"zeroVal,omitempty"`
		ZeroPtr  *infounit.ByteCount `yaml:",omitempty"`
	}{
		Val:      infounit.ByteCount(1111),
		Ptr:      &x,
		ValSlice: []infounit.ByteCount{777111, 777222, 777333},
		PtrSlice: []*infounit.ByteCount{&y, &z},
		Renamed:  66666666,
		ZeroVal:  infounit.ByteCount(0),
		ZeroPtr:  nil,
	}

	expected := strings.Join([]string{
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
		"xyzBC: 66666666",
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

func TestByteCount_UnmarshalYAML(t *testing.T) {
	v := struct {
		Val      infounit.ByteCount
		Ptr      *infounit.ByteCount
		PtrNil   *infounit.ByteCount
		ValSlice []infounit.ByteCount
		PtrSlice []*infounit.ByteCount
		Renamed  infounit.ByteCount `yaml:"xyzBC"`
		VarExprs []infounit.ByteCount
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
		"xyzBC: 66666666",
		"varexprs:",
		`- "123 kilobytes"`,
		`- "345 MB"`,
		`- "67.8GB"`,
		"",
	}, "\n")

	if err := yaml.UnmarshalStrict(([]byte)(yamlSrc), &v); err != nil {
		t.Errorf("yaml.Unmarshal() failed: %v", err)
	}
	if v.Val != 1111 {
		t.Errorf("Val: unexpected value: got: %v, want: 1111 B", v.Val)
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
	if v.Renamed != 66666666 {
		t.Errorf("Renamed: unexpected value: got: %d, want: 66666666", v.Renamed)
	}
	switch {
	case len(v.VarExprs) != 3:
		t.Errorf("VarExprs: unexpected length: got: %d, want: 3", len(v.VarExprs))
	case v.VarExprs[0] != infounit.Kilobyte*123:
		t.Errorf("VarExprs[0]: unexpected value: got: %d, want: %d", v.VarExprs[0], infounit.Kilobyte*123)
	case v.VarExprs[1] != infounit.Megabyte*345:
		t.Errorf("VarExprs[1]: unexpected value: got: %d, want: %d", v.VarExprs[1], infounit.Megabyte*345)
	case v.VarExprs[2] != infounit.Gigabyte/10*678:
		t.Errorf("VarExprs[2]: unexpected value: got: %d, want: %d", v.VarExprs[2], infounit.Gigabyte/10*678)
	}

	// t.Logf("%+v\n", v)
}
