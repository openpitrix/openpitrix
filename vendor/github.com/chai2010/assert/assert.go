// Copyright 2014 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.9

/*
Package assert provides assert helper functions for testing package.

Example:

	package somepkg_test

	import (
		. "github.com/chai2010/assert"
	)

	func TestAssert(t *testing.T) {
		Assert(t, 1 == 1)
		Assert(t, 1 == 1, "message1", "message2")
	}

	func TestAssertf(t *testing.T) {
		Assertf(t, 1 == 1, "%v:%v", "message1", "message2")
	}

See failed test output (assert_failed_test.go):

	go test -assert.failed

Report bugs to <chaishushan@gmail.com>.

Thanks!
*/
package assert

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"reflect"
	"regexp"
	"testing"
)

func Assert(tb testing.TB, condition bool, args ...interface{}) {
	tb.Helper()
	if !condition {
		if msg := fmt.Sprint(args...); msg != "" {
			tb.Fatalf("Assert failed, %s", msg)
		} else {
			tb.Fatalf("Assert failed")
		}
	}
}

func Assertf(tb testing.TB, condition bool, format string, a ...interface{}) {
	tb.Helper()
	if !condition {
		if msg := fmt.Sprintf(format, a...); msg != "" {
			tb.Fatalf("tAssert failed, %s", msg)
		} else {
			tb.Fatalf("tAssert failed")
		}
	}
}

func AssertNil(tb testing.TB, p interface{}, args ...interface{}) {
	tb.Helper()
	if p != nil {
		if msg := fmt.Sprint(args...); msg != "" {
			if err, ok := p.(error); ok && err != nil {
				tb.Fatalf("AssertNil failed, err = %v, %s", err, msg)
			} else {
				tb.Fatalf("AssertNil failed, %s", msg)
			}
		} else {
			if err, ok := p.(error); ok && err != nil {
				tb.Fatalf("AssertNil failed, err = %v", err)
			} else {
				tb.Fatalf("AssertNil failed")
			}
		}
	}
}

func AssertNotNil(tb testing.TB, p interface{}, args ...interface{}) {
	tb.Helper()
	if p == nil {
		if msg := fmt.Sprint(args...); msg != "" {
			tb.Fatalf("AssertNotNil failed, %s", msg)
		} else {
			tb.Fatalf("AssertNotNil failed")
		}
	}
}

func AssertTrue(tb testing.TB, condition bool, args ...interface{}) {
	tb.Helper()
	if condition != true {
		if msg := fmt.Sprint(args...); msg != "" {
			tb.Fatalf("AssertTrue failed, %s", msg)
		} else {
			tb.Fatalf("AssertTrue failed")
		}
	}
}

func AssertFalse(tb testing.TB, condition bool, args ...interface{}) {
	tb.Helper()
	if condition != false {
		if msg := fmt.Sprint(args...); msg != "" {
			tb.Fatalf("AssertFalse failed, %s", msg)
		} else {
			tb.Fatalf("AssertFalse failed")
		}
	}
}

func AssertEqual(tb testing.TB, expected, got interface{}, args ...interface{}) {
	tb.Helper()
	// reflect.DeepEqual is failed for `int == int64?`
	if fmt.Sprintf("%v", expected) != fmt.Sprintf("%v", got) {
		if msg := fmt.Sprint(args...); msg != "" {
			tb.Fatalf("AssertEqual failed, expected = %v, got = %v, %s", expected, got, msg)
		} else {
			tb.Fatalf("AssertEqual failed, expected = %v, got = %v", expected, got)
		}
	}
}

func AssertNotEqual(tb testing.TB, expected, got interface{}, args ...interface{}) {
	tb.Helper()

	// reflect.DeepEqual is failed for `int == int64?`
	if fmt.Sprintf("%v", expected) == fmt.Sprintf("%v", got) {
		if msg := fmt.Sprint(args...); msg != "" {
			tb.Fatalf("AssertNotEqual failed, expected = %v, got = %v, %s", expected, got, msg)
		} else {
			tb.Fatalf("AssertNotEqual failed, expected = %v, got = %v", expected, got)
		}
	}
}

func AssertNear(tb testing.TB, expected, got, abs float64, args ...interface{}) {
	tb.Helper()
	if math.Abs(expected-got) > abs {
		if msg := fmt.Sprint(args...); msg != "" {
			tb.Fatalf("AssertNear failed, expected = %v, got = %v, abs = %v, %s", expected, got, abs, msg)
		} else {
			tb.Fatalf("AssertNear failed, expected = %v, got = %v, abs = %v", expected, got, abs)
		}
	}
}

func AssertBetween(tb testing.TB, min, max, val float64, args ...interface{}) {
	tb.Helper()
	if val < min || max < val {
		if msg := fmt.Sprint(args...); msg != "" {
			tb.Fatalf("AssertBetween failed, min = %v, max = %v, val = %v, %s", min, max, val, msg)
		} else {
			tb.Fatalf("AssertBetween failed, min = %v, max = %v, val = %v", min, max, val)
		}
	}
}

func AssertNotBetween(tb testing.TB, min, max, val float64, args ...interface{}) {
	tb.Helper()
	if min <= val && val <= max {
		if msg := fmt.Sprint(args...); msg != "" {
			tb.Fatalf("AssertNotBetween failed, min = %v, max = %v, val = %v, %s", min, max, val, msg)
		} else {
			tb.Fatalf("AssertNotBetween failed, min = %v, max = %v, val = %v", min, max, val)
		}
	}
}

func AssertMatch(tb testing.TB, expectedPattern string, got []byte, args ...interface{}) {
	tb.Helper()
	if matched, err := regexp.Match(expectedPattern, got); err != nil || !matched {
		if err != nil {
			if msg := fmt.Sprint(args...); msg != "" {
				tb.Fatalf("AssertMatch failed, expected = %q, got = %v, err = %v, %s", expectedPattern, got, err, msg)
			} else {
				tb.Fatalf("AssertMatch failed, expected = %q, got = %v, err = %v", expectedPattern, got, err)
			}
		} else {
			if msg := fmt.Sprint(args...); msg != "" {
				tb.Fatalf("AssertMatch failed, expected = %q, got = %v, %s", expectedPattern, got, msg)
			} else {
				tb.Fatalf("AssertMatch failed, expected = %q, got = %v", expectedPattern, got)
			}
		}
	}
}

func AssertMatchString(tb testing.TB, expectedPattern, got string, args ...interface{}) {
	tb.Helper()
	if matched, err := regexp.MatchString(expectedPattern, got); err != nil || !matched {
		if err != nil {
			if msg := fmt.Sprint(args...); msg != "" {
				tb.Fatalf("AssertMatchString failed, expected = %q, got = %v, err = %v, %s", expectedPattern, got, err, msg)
			} else {
				tb.Fatalf("AssertMatchString failed, expected = %q, got = %v, err = %v", expectedPattern, got, err)
			}
		} else {
			if msg := fmt.Sprint(args...); msg != "" {
				tb.Fatalf("AssertMatchString failed, expected = %q, got = %v, %s", expectedPattern, got, msg)
			} else {
				tb.Fatalf("AssertMatchString failed, expected = %q, got = %v", expectedPattern, got)
			}
		}
	}
}

func AssertSliceContain(tb testing.TB, slice, val interface{}, args ...interface{}) {
	tb.Helper()
	sliceVal := reflect.ValueOf(slice)
	if sliceVal.Kind() != reflect.Slice {
		tb.Fatalf("AssertSliceContain called with non-slice value of type %T", slice)
	}
	var contained bool
	for i := 0; i < sliceVal.Len(); i++ {
		if reflect.DeepEqual(sliceVal.Index(i).Interface(), val) {
			contained = true
			break
		}
	}
	if !contained {
		if msg := fmt.Sprint(args...); msg != "" {
			tb.Fatalf("AssertSliceContain failed, slice = %v, val = %v, %s", slice, val, msg)
		} else {
			tb.Fatalf("AssertSliceContain failed, slice = %v, val = %v", slice, val)
		}
	}
}

func AssertSliceNotContain(tb testing.TB, slice, val interface{}, args ...interface{}) {
	tb.Helper()
	sliceVal := reflect.ValueOf(slice)
	if sliceVal.Kind() != reflect.Slice {
		tb.Fatalf("AssertSliceNotContain called with non-slice value of type %T", slice)
	}
	var contained bool
	for i := 0; i < sliceVal.Len(); i++ {
		if reflect.DeepEqual(sliceVal.Index(i).Interface(), val) {
			contained = true
			break
		}
	}
	if contained {
		if msg := fmt.Sprint(args...); msg != "" {
			tb.Fatalf("AssertSliceNotContain failed, slice = %v, val = %v, %s", slice, val, msg)
		} else {
			tb.Fatalf("AssertSliceNotContain failed, slice = %v, val = %v", slice, val)
		}
	}
}

func AssertMapEqual(tb testing.TB, expected, got interface{}, args ...interface{}) {
	tb.Helper()
	expectedMap := reflect.ValueOf(expected)
	if expectedMap.Kind() != reflect.Map {
		tb.Fatalf("AssertMapEqual called with non-map expected value of type %T", expected)
	}
	gotMap := reflect.ValueOf(got)
	if gotMap.Kind() != reflect.Map {
		tb.Fatalf("AssertMapEqual called with non-map got value of type %T", got)
	}

	if a, b := expectedMap.Len(), gotMap.Len(); a != b {
		if msg := fmt.Sprint(args...); msg != "" {
			tb.Fatalf("AssertMapEqual failed, len(expected) = %d, len(got) = %d, %s", a, b, msg)
		} else {
			tb.Fatalf("AssertMapEqual failed, len(expected) = %d, len(got) = %d", a, b)
		}
		return
	}

	for _, key := range expectedMap.MapKeys() {
		expectedVal := expectedMap.MapIndex(key).Interface()
		gotVal := gotMap.MapIndex(key).Interface()

		if fmt.Sprintf("%v", expectedVal) != fmt.Sprintf("%v", gotVal) {
			if msg := fmt.Sprint(args...); msg != "" {
				tb.Fatalf(
					"AssertMapEqual failed, key = %v, expected = %v, got = %v, %s",
					key.Interface(), expectedVal, gotVal, msg,
				)
			} else {
				tb.Fatalf(
					"AssertMapEqual failed, key = %v, expected = %v, got = %v",
					key.Interface(), expectedVal, gotVal,
				)
			}
			return
		}
	}
}

func AssertMapContain(tb testing.TB, m, key, val interface{}, args ...interface{}) {
	tb.Helper()
	mapVal := reflect.ValueOf(m)
	if mapVal.Kind() != reflect.Map {
		tb.Fatalf("AssertMapContain called with non-map value of type %T", m)
	}
	elemVal := mapVal.MapIndex(reflect.ValueOf(key))
	if !elemVal.IsValid() || !reflect.DeepEqual(elemVal.Interface(), val) {
		if msg := fmt.Sprint(args...); msg != "" {
			tb.Fatalf("AssertMapContain failed, map = %v, key = %v, val = %v, %s", m, key, val, msg)
		} else {
			tb.Fatalf("AssertMapContain failed, map = %v, key = %v, val = %v", m, key, val)
		}
	}
}

func AssertMapContainKey(tb testing.TB, m, key interface{}, args ...interface{}) {
	tb.Helper()
	mapVal := reflect.ValueOf(m)
	if mapVal.Kind() != reflect.Map {
		tb.Fatalf("AssertMapContainKey called with non-map value of type %T", m)
	}
	elemVal := mapVal.MapIndex(reflect.ValueOf(key))
	if !elemVal.IsValid() {
		if msg := fmt.Sprint(args...); msg != "" {
			tb.Fatalf("AssertMapContainKey failed, map = %v, key = %v, %s", m, key, msg)
		} else {
			tb.Fatalf("AssertMapContainKey failed, map = %v, key = %v", m, key)
		}
	}
}

func AssertMapContainVal(tb testing.TB, m, val interface{}, args ...interface{}) {
	tb.Helper()
	mapVal := reflect.ValueOf(m)
	if mapVal.Kind() != reflect.Map {
		tb.Fatalf("AssertMapContainVal called with non-map value of type %T", m)
	}
	var contained bool
	for _, key := range mapVal.MapKeys() {
		elemVal := mapVal.MapIndex(key)
		if elemVal.IsValid() && reflect.DeepEqual(elemVal.Interface(), val) {
			contained = true
			break
		}
	}
	if !contained {
		if msg := fmt.Sprint(args...); msg != "" {
			tb.Fatalf("AssertMapContainVal failed, map = %v, val = %v, %s", m, val, msg)
		} else {
			tb.Fatalf("AssertMapContainVal failed, map = %v, val = %v", m, val)
		}
	}
}

func AssertMapNotContain(tb testing.TB, m, key, val interface{}, args ...interface{}) {
	tb.Helper()
	mapVal := reflect.ValueOf(m)
	if mapVal.Kind() != reflect.Map {
		tb.Fatalf("AssertMapNotContain called with non-map value of type %T", m)
	}
	elemVal := mapVal.MapIndex(reflect.ValueOf(key))
	if elemVal.IsValid() && reflect.DeepEqual(elemVal.Interface(), val) {
		if msg := fmt.Sprint(args...); msg != "" {
			tb.Fatalf("AssertMapNotContain failed, map = %v, key = %v, val = %v, %s", m, key, val, msg)
		} else {
			tb.Fatalf("AssertMapNotContain failed, map = %v, key = %v, val = %v", m, key, val)
		}
	}
}

func AssertMapNotContainKey(tb testing.TB, m, key interface{}, args ...interface{}) {
	tb.Helper()
	mapVal := reflect.ValueOf(m)
	if mapVal.Kind() != reflect.Map {
		tb.Fatalf("AssertMapNotContainKey called with non-map value of type %T", m)
	}
	elemVal := mapVal.MapIndex(reflect.ValueOf(key))
	if elemVal.IsValid() {
		if msg := fmt.Sprint(args...); msg != "" {
			tb.Fatalf("AssertMapNotContainKey failed, map = %v, key = %v, %s", m, key, msg)
		} else {
			tb.Fatalf("AssertMapNotContainKey failed, map = %v, key = %v", m, key)
		}
	}
}

func AssertMapNotContainVal(tb testing.TB, m, val interface{}, args ...interface{}) {
	tb.Helper()
	mapVal := reflect.ValueOf(m)
	if mapVal.Kind() != reflect.Map {
		tb.Fatalf("AssertMapNotContainVal called with non-map value of type %T", m)
	}
	var contained bool
	for _, key := range mapVal.MapKeys() {
		elemVal := mapVal.MapIndex(key)
		if elemVal.IsValid() && reflect.DeepEqual(elemVal.Interface(), val) {
			contained = true
			break
		}
	}
	if contained {
		if msg := fmt.Sprint(args...); msg != "" {
			tb.Fatalf("AssertMapNotContainVal failed, map = %v, val = %v, %s", m, val, msg)
		} else {
			tb.Fatalf("AssertMapNotContainVal failed, map = %v, val = %v", m, val)
		}
	}
}

func AssertZero(tb testing.TB, val interface{}, args ...interface{}) {
	tb.Helper()
	if !reflect.DeepEqual(reflect.Zero(reflect.TypeOf(val)).Interface(), val) {
		if msg := fmt.Sprint(args...); msg != "" {
			tb.Fatalf("AssertZero failed, val = %v, %s", val, msg)
		} else {
			tb.Fatalf("AssertZero failed, val = %v", val)
		}
	}
}

func AssertNotZero(tb testing.TB, val interface{}, args ...interface{}) {
	tb.Helper()
	if reflect.DeepEqual(reflect.Zero(reflect.TypeOf(val)).Interface(), val) {
		if msg := fmt.Sprint(args...); msg != "" {
			tb.Fatalf("AssertNotZero failed, val = %v, %s", val, msg)
		} else {
			tb.Fatalf("AssertNotZero failed, val = %v", val)
		}
	}
}

func AssertFileExists(tb testing.TB, path string, args ...interface{}) {
	tb.Helper()
	if _, err := os.Stat(path); err != nil {
		if msg := fmt.Sprint(args...); msg != "" {
			if err != nil {
				tb.Fatalf("AssertFileExists failed, path = %v, err = %v, %s", path, err, msg)
			} else {
				tb.Fatalf("AssertFileExists failed, path = %v, %s", path, msg)
			}
		} else {
			if err != nil {
				tb.Fatalf("AssertFileExists failed, path = %v, err = %v", path, err)
			} else {
				tb.Fatalf("AssertFileExists failed, path = %v", path)
			}
		}
	}
}

func AssertFileNotExists(tb testing.TB, path string, args ...interface{}) {
	tb.Helper()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		if msg := fmt.Sprint(args...); msg != "" {
			if err != nil {
				tb.Fatalf("AssertFileNotExists failed, path = %v, err = %v, %s", path, err, msg)
			} else {
				tb.Fatalf("AssertFileNotExists failed, path = %v, %s", path, msg)
			}
		} else {
			if err != nil {
				tb.Fatalf("AssertFileNotExists failed, path = %v, err = %v", path, err)
			} else {
				tb.Fatalf("AssertFileNotExists failed, path = %v", path)
			}
		}
	}
}

func AssertImplements(tb testing.TB, interfaceObj, obj interface{}, args ...interface{}) {
	tb.Helper()
	if !reflect.TypeOf(obj).Implements(reflect.TypeOf(interfaceObj).Elem()) {
		if msg := fmt.Sprint(args...); msg != "" {
			tb.Fatalf("AssertImplements failed, interface = %T, obj = %T, %s", interfaceObj, obj, msg)
		} else {
			tb.Fatalf("AssertImplements failed, interface = %T, obj = %T", interfaceObj, obj)
		}
	}
}

func AssertSameType(tb testing.TB, expectedType interface{}, obj interface{}, args ...interface{}) {
	tb.Helper()
	if !reflect.DeepEqual(reflect.TypeOf(obj), reflect.TypeOf(expectedType)) {
		if msg := fmt.Sprint(args...); msg != "" {
			tb.Fatalf("AssertSameType failed, expected = %T, obj = %T, %s", expectedType, obj, msg)
		} else {
			tb.Fatalf("AssertSameType failed, expected = %T, obj = %T", expectedType, obj)
		}
	}
}

func AssertPanic(tb testing.TB, f func(), args ...interface{}) {
	tb.Helper()

	panicVal := func() (panicVal interface{}) {
		defer func() {
			panicVal = recover()
		}()
		f()
		return
	}()

	if panicVal == nil {
		if msg := fmt.Sprint(args...); msg != "" {
			tb.Fatalf("AssertPanic failed, %s", msg)
		} else {
			tb.Fatalf("AssertPanic failed")
		}
	}
}

func AssertNotPanic(tb testing.TB, f func(), args ...interface{}) {
	tb.Helper()

	panicVal := func() (panicVal interface{}) {
		defer func() {
			panicVal = recover()
		}()
		f()
		return
	}()

	if panicVal != nil {
		if msg := fmt.Sprint(args...); msg != "" {
			tb.Fatalf("AssertNotPanic failed, panic = %v, %s", panicVal, msg)
		} else {
			tb.Fatalf("AssertNotPanic failed, panic = %v", panicVal)
		}
	}
}

func AssertImageEqual(tb testing.TB, expected, got image.Image, maxDelta color.Color, args ...interface{}) {
	tb.Helper()
	if equal, pos := tImageEqual(expected, got, maxDelta); !equal {
		if msg := fmt.Sprint(args...); msg != "" {
			tb.Fatalf("AssertImageEqual failed, pos = %v, expected = %#v, got = %#v, max = %#v, %s",
				pos, expected.At(pos.X, pos.Y), got.At(pos.X, pos.Y),
				maxDelta, msg,
			)
		} else {
			tb.Fatalf("AssertImageEqual failed, pos = %v, expected = %#v, got = %#v, max = %#v",
				pos, expected.At(pos.X, pos.Y), got.At(pos.X, pos.Y),
				maxDelta,
			)
		}
	}
}

func tImageEqual(m0, m1 image.Image, maxDelta color.Color) (ok bool, failedPixelPos image.Point) {
	b := m0.Bounds()
	maxDelta_R, maxDelta_G, maxDelta_B, maxDelta_A := maxDelta.RGBA()

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c0 := m0.At(x, y)
			c1 := m1.At(x, y)
			r0, g0, b0, a0 := c0.RGBA()
			r1, g1, b1, a1 := c1.RGBA()
			if tDeltaUint32(r0, r1) > maxDelta_R {
				return false, image.Pt(x, y)
			}
			if tDeltaUint32(g0, g1) > maxDelta_G {
				return false, image.Pt(x, y)
			}
			if tDeltaUint32(b0, b1) > maxDelta_B {
				return false, image.Pt(x, y)
			}
			if tDeltaUint32(a0, a1) > maxDelta_A {
				return false, image.Pt(x, y)
			}
		}
	}
	return true, image.Pt(0, 0)
}

func tDeltaUint32(a, b uint32) uint32 {
	if a >= b {
		return a - b
	}
	return b - a
}
