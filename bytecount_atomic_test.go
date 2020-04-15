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
func TestAtomicAddByteCount_1(t *testing.T) {
	t.Parallel()

	var wg sync.WaitGroup
	bc := infounit.Megabyte
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			infounit.AtomicAddByteCount(&bc, 100)
		}()
	}
	wg.Wait()

	exbc := infounit.Megabyte + infounit.ByteCount(100*10000)
	if bc != exbc {
		t.Errorf(`want: %s, got: %s`, exbc, bc)
	}
}

//
func TestAtomicAddByteCount_2(t *testing.T) {
	t.Parallel()

	bc := infounit.Byte * 30000
	v := infounit.AtomicAddByteCount(&bc, 999)

	exbc := infounit.Byte * 30999
	if bc != exbc {
		t.Errorf(`bc: want: %s, got: %s`, exbc, bc)
	}
	if v != exbc {
		t.Errorf(`v: want: %s, got: %s`, exbc, v)
	}
}

//
func TestAtomicSubByteCount_1(t *testing.T) {
	t.Parallel()

	var wg sync.WaitGroup
	bc := infounit.Megabyte
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			infounit.AtomicSubByteCount(&bc, 10)
		}()
	}
	wg.Wait()

	exbc := infounit.Megabyte - infounit.ByteCount(10*10000)
	if bc != exbc {
		t.Errorf(`want: %s, got: %s`, exbc, bc)
	}
}

//
func TestAtomicSubByteCount_2(t *testing.T) {
	t.Parallel()

	bc := infounit.Byte * 30000
	v := infounit.AtomicSubByteCount(&bc, 999)

	exbc := infounit.Byte * 29001
	if bc != exbc {
		t.Errorf(`bc: want: %s, got: %s`, exbc, bc)
	}
	if v != exbc {
		t.Errorf(`v: want: %s, got: %s`, exbc, v)
	}
}

//
func TestAtomicLoadByteCount_1(t *testing.T) {
	t.Parallel()

	bc := infounit.Byte * 30000
	v := infounit.AtomicLoadByteCount(&bc)

	exbc := infounit.Byte * 30000
	if v != exbc {
		t.Errorf(`want: %s, got: %s`, exbc, bc)
	}
}

//
func TestAtomicStoreByteCount_1(t *testing.T) {
	t.Parallel()

	bc := infounit.Byte * 30000
	infounit.AtomicStoreByteCount(&bc, 12345)

	exbc := infounit.Byte * 12345
	if bc != exbc {
		t.Errorf(`want: %s, got: %s`, exbc, bc)
	}
}
