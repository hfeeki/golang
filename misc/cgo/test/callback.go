// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cgotest

/*
void callback(void *f);
void callGoFoo(void);
*/
import "C"

import (
	"./backdoor"
	"runtime"
	"testing"
	"unsafe"
)

// nestedCall calls into C, back into Go, and finally to f.
func nestedCall(f func()) {
	// NOTE: Depends on representation of f.
	// callback(x) calls goCallback(x)
	C.callback(*(*unsafe.Pointer)(unsafe.Pointer(&f)))
}

//export goCallback
func goCallback(p unsafe.Pointer) {
	(*(*func())(unsafe.Pointer(&p)))()
}

func testCallback(t *testing.T) {
	var x = false
	nestedCall(func() { x = true })
	if !x {
		t.Fatal("nestedCall did not call func")
	}
}

func testCallbackGC(t *testing.T) {
	nestedCall(runtime.GC)
}

var lockedOSThread = backdoor.LockedOSThread

func testCallbackPanic(t *testing.T) {
	// Make sure panic during callback unwinds properly.
	if lockedOSThread() {
		t.Fatal("locked OS thread on entry to TestCallbackPanic")
	}
	defer func() {
		s := recover()
		if s == nil {
			t.Fatal("did not panic")
		}
		if s.(string) != "callback panic" {
			t.Fatal("wrong panic:", s)
		}
		if lockedOSThread() {
			t.Fatal("locked OS thread on exit from TestCallbackPanic")
		}
	}()
	nestedCall(func() { panic("callback panic") })
	panic("nestedCall returned")
}

func testCallbackPanicLoop(t *testing.T) {
	// Make sure we don't blow out m->g0 stack.
	for i := 0; i < 100000; i++ {
		testCallbackPanic(t)
	}
}

func testCallbackPanicLocked(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if !lockedOSThread() {
		t.Fatal("runtime.LockOSThread didn't")
	}
	defer func() {
		s := recover()
		if s == nil {
			t.Fatal("did not panic")
		}
		if s.(string) != "callback panic" {
			t.Fatal("wrong panic:", s)
		}
		if !lockedOSThread() {
			t.Fatal("lost lock on OS thread after panic")
		}
	}()
	nestedCall(func() { panic("callback panic") })
	panic("nestedCall returned")
}

// Callback with zero arguments used to make the stack misaligned,
// which broke the garbage collector and other things.
func testZeroArgCallback(t *testing.T) {
	defer func() {
		s := recover()
		if s != nil {
			t.Fatal("panic during callback:", s)
		}
	}()
	C.callGoFoo()
}

//export goFoo
func goFoo() {
	x := 1
	for i := 0; i < 10000; i++ {
		// variadic call mallocs + writes to
		variadic(x, x, x)
		if x != 1 {
			panic("bad x")
		}
	}
}

func variadic(x ...interface{}) {}

func testBlocking(t *testing.T) {
	c := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			c <- <-c
		}
	}()
	nestedCall(func() {
		for i := 0; i < 10; i++ {
			c <- i
			if j := <-c; j != i {
				t.Errorf("out of sync %d != %d", j, i)
			}
		}
	})
}
