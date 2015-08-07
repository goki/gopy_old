// Copyright 2015 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package funcs

type Func func()

type S1 struct {
	F Func
}

type S2 struct {
	F func()
}
