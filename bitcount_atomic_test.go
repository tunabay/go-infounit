// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package infounit_test

import (
	"sync"
	"testing"

	"github.com/tunabay/go-infounit"
)

//
func TestAtomicAddBitCount_1(t *testing.T) {
	t.Parallel()

	var wg sync.WaitGroup
	bc := infounit.Megabit
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			infounit.AtomicAddBitCount(&bc, 100)
		}()
	}
	wg.Wait()

	exbc := infounit.Megabit + infounit.BitCount(100*10000)
	if bc != exbc {
		t.Errorf(`want: %s, got: %s`, exbc, bc)
	}
}

//
func TestAtomicAddBitCount_2(t *testing.T) {
	t.Parallel()

	bc := infounit.Bit * 30000
	v := infounit.AtomicAddBitCount(&bc, 999)

	exbc := infounit.Bit * 30999
	if bc != exbc {
		t.Errorf(`want: %s, got: %s`, exbc, bc)
	}
	if v != exbc {
		t.Errorf(`want: %s, got: %s`, exbc, v)
	}
}

//
func TestAtomicSubBitCount_1(t *testing.T) {
	t.Parallel()

	var wg sync.WaitGroup
	bc := infounit.Megabit
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			infounit.AtomicSubBitCount(&bc, 10)
		}()
	}
	wg.Wait()

	exbc := infounit.Megabit - infounit.BitCount(10*10000)
	if bc != exbc {
		t.Errorf(`want: %s, got: %s`, exbc, bc)
	}
}

//
func TestAtomicSubBitCount_2(t *testing.T) {
	t.Parallel()

	bc := infounit.Bit * 30000
	v := infounit.AtomicSubBitCount(&bc, 999)

	exbc := infounit.Bit * 29001
	if bc != exbc {
		t.Errorf(`want: %s, got: %s`, exbc, bc)
	}
	if v != exbc {
		t.Errorf(`want: %s, got: %s`, exbc, v)
	}
}

//
func TestAtomicLoadBitCount_1(t *testing.T) {
	t.Parallel()

	bc := infounit.Bit * 30000
	v := infounit.AtomicLoadBitCount(&bc)

	exbc := infounit.Bit * 30000
	if v != exbc {
		t.Errorf(`want: %s, got: %s`, exbc, bc)
	}
}

//
func TestAtomicStoreBitCount_1(t *testing.T) {
	t.Parallel()

	bc := infounit.Bit * 30000
	infounit.AtomicStoreBitCount(&bc, 12345)

	exbc := infounit.Bit * 12345
	if bc != exbc {
		t.Errorf(`want: %s, got: %s`, exbc, bc)
	}
}
