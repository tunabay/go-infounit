// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package infounit_test

import (
	"testing"

	"github.com/tunabay/go-infounit"
)

//
func TestAtomicLoadBitRate_1(t *testing.T) {
	t.Parallel()

	br := infounit.BitPerSecond * 30000.987
	v := infounit.AtomicLoadBitRate(&br)

	exbr := infounit.BitPerSecond * 30000.987
	if v != exbr {
		t.Errorf(`want: %s, got: %s`, exbr, br)
	}
}

//
func TestAtomicStoreBitRate_1(t *testing.T) {
	t.Parallel()

	br := infounit.BitPerSecond * 30000
	infounit.AtomicStoreBitRate(&br, 12345.67)

	exbr := infounit.BitPerSecond * 12345.67
	if br != exbr {
		t.Errorf(`want: %s, got: %s`, exbr, br)
	}
}
