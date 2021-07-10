// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package infounit_test

import (
	"bytes"
	"encoding/hex"
	"math"
	"strings"
	"testing"

	"github.com/tunabay/go-infounit"
	"gopkg.in/yaml.v2"
)

//
func TestBitRate_MarshalBinary_1(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b   float64
		hex string
	}{
		{0.0, "0000000000000000"},
		{1.0, "3ff0000000000000"},
		{-1.0, "bff0000000000000"},
		{123456.78, "40fe240c7ae147ae"},
		{math.NaN(), "7ff8000000000001"},
		{math.Inf(+1), "7ff0000000000000"},
		{math.Inf(-1), "fff0000000000000"},
		{math.MaxFloat64, "7fefffffffffffff"},
		{math.SmallestNonzeroFloat64, "0000000000000001"},
	}

	for _, c := range tc {
		br := infounit.BitRate(c.b)
		bin, err := br.MarshalBinary()
		if err != nil {
			t.Error(err)
		}
		exbin, err := hex.DecodeString(c.hex)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(bin, exbin) {
			t.Errorf(`%v: want: %x, got: %x`, br, exbin, bin)
		}
		// t.Logf(`%v: %x`, br, bin)
	}
}

//
func TestBitRate_UnmarshalBinary_1(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b   float64
		hex string
	}{
		{0.0, "0000000000000000"},
		{1.0, "3ff0000000000000"},
		{-1.0, "bff0000000000000"},
		{123456.78, "40fe240c7ae147ae"},
		{math.NaN(), "7ff8000000000001"},
		{math.Inf(+1), "7ff0000000000000"},
		{math.Inf(-1), "fff0000000000000"},
		{math.MaxFloat64, "7fefffffffffffff"},
		{math.SmallestNonzeroFloat64, "0000000000000001"},
	}

	for _, c := range tc {
		var br infounit.BitRate
		bin, err := hex.DecodeString(c.hex)
		if err != nil {
			t.Fatal(err)
		}
		if err := br.UnmarshalBinary(bin); err != nil {
			t.Error(err)
		}

		var ok bool
		var exbr infounit.BitRate
		switch {
		case math.IsNaN(c.b):
			exbr = infounit.BitRate(math.NaN())
			ok = math.IsNaN(float64(br))
		case math.IsInf(c.b, +1):
			exbr = infounit.BitRate(math.Inf(+1))
			ok = math.IsInf(float64(br), +1)
		case math.IsInf(c.b, -1):
			exbr = infounit.BitRate(math.Inf(-1))
			ok = math.IsInf(float64(br), -1)
		default:
			exbr = infounit.BitRate(c.b)
			ok = br == exbr
		}
		if !ok {
			t.Errorf(`%x: want: %v, got: %v`, c.hex, exbr, br)
		}
		// t.Logf(`%s: %v`, c.hex, br)
	}
}

//
func TestBitRate_MarshalText_1(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b   float64
		txt string
	}{
		{0, "0 bit/s"},
		{1, "1 bit/s"},
		{987654.321, "987654.321 bit/s"},
		{0.00987654321, "0.00987654321 bit/s"},
		{99999999999.9999, "99999999999.9999 bit/s"},
	}

	for _, c := range tc {
		br := infounit.BitRate(c.b)
		txtb, err := br.MarshalText()
		if err != nil {
			t.Error(err)
		}
		txt := string(txtb)
		if txt != c.txt {
			t.Errorf(`%f: want: "%s", got: "%s"`, br, c.txt, txt)
		}
	}
}

//
func TestBitRate_UnmarshalText_1(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b   float64
		txt string
	}{
		{0, "0 bit/s"},
		{1, "1 bit/s"},
		{987654.321, "987654.321 bit/s"},
		{0.00987654321, "0.00987654321 bit/s"},
		{99999999999.9999, "99999999999.9999 bit/s"},
	}

	for _, c := range tc {
		var br infounit.BitRate
		exbr := infounit.BitRate(c.b)
		if err := br.UnmarshalText(([]byte)(c.txt)); err != nil {
			t.Errorf(`%s: want: %s, got: %s`, c.txt, exbr, err)
		}
		if br != exbr {
			t.Errorf(`%s: want: %s, got: %s`, c.txt, exbr, br)
		}
	}
}

func TestBitRate_MarshalYAML(t *testing.T) {
	var (
		x = infounit.BitRate(999.99999)
		y = infounit.BitRate(888.88888)
		z = infounit.BitRate(777.77777)
	)
	v := &struct {
		Val      infounit.BitRate
		Ptr      *infounit.BitRate
		PtrNil   *infounit.BitRate
		ValSlice []infounit.BitRate
		PtrSlice []*infounit.BitRate
		Renamed  infounit.BitRate  `yaml:"xyzBC"`
		ZeroVal  infounit.BitRate  `yaml:"zeroVal,omitempty"`
		ZeroPtr  *infounit.BitRate `yaml:",omitempty"`
	}{
		Val:      infounit.BitRate(666.666),
		Ptr:      &x,
		ValSlice: []infounit.BitRate{555.111, 555.222, 555.333},
		PtrSlice: []*infounit.BitRate{&y, &z},
		Renamed:  444.4444,
		ZeroVal:  infounit.BitRate(0),
		ZeroPtr:  nil,
	}

	expected := strings.Join([]string{
		"val: 666.666",
		"ptr: 999.99999",
		"ptrnil: null",
		"valslice:",
		"- 555.111",
		"- 555.222",
		"- 555.333",
		"ptrslice:",
		"- 888.88888",
		"- 777.77777",
		"xyzBC: 444.4444",
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

func TestBitRate_UnmarshalYAML(t *testing.T) {
	v := struct {
		Val      infounit.BitRate
		Ptr      *infounit.BitRate
		PtrNil   *infounit.BitRate
		ValSlice []infounit.BitRate
		PtrSlice []*infounit.BitRate
		Renamed  infounit.BitRate `yaml:"xyzBC"`
		VarExprs []infounit.BitRate
	}{}

	yamlSrc := strings.Join([]string{
		"val: 9999.99999",
		"ptr: 8888.88888",
		"ptrnil: null",
		"valslice:",
		"- 777.111",
		"- 777.222",
		"- 777.333",
		"ptrslice:",
		"- 66666.2222",
		"- 66666.3333",
		"xyzBC: 5555555.555",
		"varexprs:",
		`- "12345.678 kilobits per second"`,
		`- "345 Mbit/s"`,
		`- "67.8Gbit/s"`,
		"",
	}, "\n")

	if err := yaml.UnmarshalStrict(([]byte)(yamlSrc), &v); err != nil {
		t.Errorf("yaml.Unmarshal() failed: %v", err)
	}
	if v.Val != 9999.99999 {
		t.Errorf("Val: unexpected value: got: %v, want: 9999.99999 bit/s", v.Val)
	}

	switch {
	case v.Ptr == nil:
		t.Errorf("Ptr: unexpected value: got: <nil>, want: %v", 8888.88888)
	case *v.Ptr != 8888.88888:
		t.Errorf("Ptr: unexpected value: got: %v, want: %v", *v.Ptr, 8888.88888)
	}
	if v.PtrNil != nil {
		t.Errorf("PtrNil: unexpected value: got: %v, want: <nil>", *v.PtrNil)
	}
	switch {
	case len(v.ValSlice) != 3:
		t.Errorf("ValSlice: unexpected length: got: %d, want: 3", len(v.ValSlice))
	case v.ValSlice[0] != 777.111:
		t.Errorf("ValSlice[0]: unexpected value: got: %v, want: 777.111", v.ValSlice[0])
	case v.ValSlice[1] != 777.222:
		t.Errorf("ValSlice[1]: unexpected value: got: %v, want: 777.222", v.ValSlice[1])
	case v.ValSlice[2] != 777.333:
		t.Errorf("ValSlice[2]: unexpected value: got: %d, want: 777.333", v.ValSlice[2])
	}
	switch {
	case len(v.PtrSlice) != 2:
		t.Errorf("PtrSlice: unexpected length: got: %d, want: 2", len(v.PtrSlice))
	case *v.PtrSlice[0] != 66666.2222:
		t.Errorf("PtrSlice[0]: unexpected value: got: %v, want: 66666.2222", *v.PtrSlice[0])
	case *v.PtrSlice[1] != 66666.3333:
		t.Errorf("PtrSlice[1]: unexpected value: got: %v, want: 66666.3333", *v.PtrSlice[1])
	}
	if v.Renamed != 5555555.555 {
		t.Errorf("Renamed: unexpected value: got: %v, want: 5555555.555", v.Renamed)
	}
	switch {
	case len(v.VarExprs) != 3:
		t.Errorf("VarExprs: unexpected length: got: %v, want: 3", len(v.VarExprs))
	case v.VarExprs[0] != infounit.BitRate(12345678):
		t.Errorf("VarExprs[0]: unexpected value: got: %v, want: %v", v.VarExprs[0], infounit.BitRate(12345678))
	case v.VarExprs[1] != infounit.MegabitPerSecond*345:
		t.Errorf("VarExprs[1]: unexpected value: got: %v, want: %v", v.VarExprs[1], infounit.MegabitPerSecond*345)
	case v.VarExprs[2] != infounit.GigabitPerSecond/10*678:
		t.Errorf("VarExprs[2]: unexpected value: got: %v, want: %v", v.VarExprs[2], infounit.GigabitPerSecond/10*678)
	}
}
