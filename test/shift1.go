// errorcheck

// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test illegal shifts.
// Issue 1708, illegal cases.
// Does not compile.

package p

func f(x int) int         { return 0 }
func g(x interface{}) int { return 0 }
func h(x float64) int     { return 0 }

// from the spec
var (
	s uint    = 33
	u         = 1.0 << s // ERROR "invalid operation|shift of non-integer operand"
	v float32 = 1 << s   // ERROR "invalid" "as type float32"
)

// non-constant shift expressions
var (
	e1       = g(2.0 << s) // ERROR "invalid" "as type interface"
	f1       = h(2 << s)   // ERROR "invalid" "as type float64"
	g1 int64 = 1.1 << s    // ERROR "truncated"
)

// constant shift expressions
const c uint = 65

var (
	a2 int = 1.0 << c    // ERROR "overflow"
	b2     = 1.0 << c    // ERROR "overflow"
	d2     = f(1.0 << c) // ERROR "overflow"
)

var (
	// issues 4882, 4936.
	a3 = 1.0<<s + 0 // ERROR "invalid operation|shift of non-integer operand"
	// issue 4937
	b3 = 1<<s + 1 + 1.0 // ERROR "invalid operation|shift of non-integer operand"
	// issue 5014
	c3     = complex(1<<s, 0) // ERROR "shift of type float64"
	d3 int = complex(1<<s, 3) // ERROR "cannot use.*as type int" "shift of type float64"
	e3     = real(1 << s)     // ERROR "invalid"
	f3     = imag(1 << s)     // ERROR "invalid"
)

var (
	a4 float64
	b4 int
	c4 = complex(1<<s, a4) // ERROR "shift of type float64"
	d4 = complex(1<<s, b4) // ERROR "invalid"
)
