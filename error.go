// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package infounit

import (
	"errors"
)

// ErrOutOfRange is the error thrown when the result exceeds the range.
var ErrOutOfRange = errors.New("out of range")

// ErrDivZeroBitRate is the error thrown when trying to divide by zero bit rate.
var ErrDivZeroBitRate = errors.New("division by zero bit rate")
