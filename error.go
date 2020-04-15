// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package infounit

import (
	"fmt"
)

// ErrOutOfRange is the error thrown when the result exceeds the range.
var ErrOutOfRange = fmt.Errorf("out of range")
