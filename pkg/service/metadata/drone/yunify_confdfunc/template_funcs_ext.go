// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

// copy from https://github.com/yunify/confd/blob/master/resource/template/template_funcs_ext.go

package yunify_confdfunc

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"

	"openpitrix.io/openpitrix/pkg/libconfd"
)

//Part of the following func is copied from spf13/hugo and do some refactor.

var timeType = reflect.TypeOf((*time.Time)(nil)).Elem()

// DoArithmetic performs arithmetic operations (+,-,*,/) using reflection to
// determine the type of the two terms.
// This func will auto convert string and uint to int64/float64, then apply  operations,
// return float64, or int64
func DoArithmetic(a, b interface{}, op rune) (interface{}, error) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	var ai, bi int64
	var af, bf float64

	var err error
	if av.Kind() == reflect.String {
		av, err = stringToNumber(av)
		if err != nil {
			return nil, err
		}
	}
	if bv.Kind() == reflect.String {
		bv, err = stringToNumber(bv)
		if err != nil {
			return nil, err
		}
	}

	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ai = av.Int()
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			bi = bv.Int()
		case reflect.Float32, reflect.Float64:
			af = float64(ai) // may overflow
			ai = 0
			bf = bv.Float()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			bi = int64(bv.Uint()) // may overflow
		default:
			return nil, errors.New("Can't apply the operator to the values")
		}
	case reflect.Float32, reflect.Float64:
		af = av.Float()
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			bf = float64(bv.Int()) // may overflow
		case reflect.Float32, reflect.Float64:
			bf = bv.Float()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			bf = float64(bv.Uint()) // may overflow
		default:
			return nil, errors.New("Can't apply the operator to the values")
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		ai = int64(av.Uint()) // may overflow
		switch bv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			bi = bv.Int()
		case reflect.Float32, reflect.Float64:
			af = float64(ai) // may overflow
			ai = 0
			bf = bv.Float()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			bi = int64(bv.Uint()) // may overflow
		default:
			return nil, errors.New("Can't apply the operator to the values")
		}
	default:
		return nil, errors.New("Can't apply the operator to the values")
	}

	switch op {
	case '+':
		if af != 0 || bf != 0 {
			return af + bf, nil
		} else if ai != 0 || bi != 0 {
			return ai + bi, nil
		}
		return 0, nil
	case '-':
		if af != 0 || bf != 0 {
			return af - bf, nil
		} else if ai != 0 || bi != 0 {
			return ai - bi, nil
		}
		return 0, nil
	case '*':
		if af != 0 || bf != 0 {
			return af * bf, nil
		}
		if ai != 0 || bi != 0 {
			return ai * bi, nil
		}
		return 0, nil
	case '/':
		if bf != 0 {
			return af / bf, nil
		} else if bi != 0 {
			return ai / bi, nil
		}
		return nil, errors.New("Can't divide the value by 0")
	default:
		return nil, errors.New("There is no such an operation")
	}
}

func stringToNumber(value reflect.Value) (reflect.Value, error) {
	var result reflect.Value
	str := value.String()
	if isFloat(str) {
		vf, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("Can't apply the operator to the value [%s] ,err [%s] ", str, err.Error())
		}
		result = reflect.ValueOf(vf)
	} else {
		vi, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("Can't apply the operator to the value [%s] ,err [%s] ", str, err.Error())
		}
		result = reflect.ValueOf(vi)
	}
	return result, nil
}

func isFloat(value string) bool {
	return strings.Index(value, ".") >= 0
}

// eq returns the boolean truth of arg1 == arg2.
func eq(x, y interface{}) bool {
	normalize := func(v interface{}) interface{} {
		vv := reflect.ValueOf(v)
		nv, err := stringToNumber(vv)
		if err == nil {
			vv = nv
		}
		switch vv.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return float64(vv.Int()) //may overflow
		case reflect.Float32, reflect.Float64:
			return vv.Float()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return float64(vv.Uint()) //may overflow
		default:
			return v
		}
	}
	x = normalize(x)
	y = normalize(y)
	return reflect.DeepEqual(x, y)
}

// ne returns the boolean truth of arg1 != arg2.
func ne(x, y interface{}) bool {
	return !eq(x, y)
}

// ge returns the boolean truth of arg1 >= arg2.
func ge(a, b interface{}) bool {
	left, right := compareGetFloat(a, b)
	return left >= right
}

// gt returns the boolean truth of arg1 > arg2.
func gt(a, b interface{}) bool {
	left, right := compareGetFloat(a, b)
	return left > right
}

// le returns the boolean truth of arg1 <= arg2.
func le(a, b interface{}) bool {
	left, right := compareGetFloat(a, b)
	return left <= right
}

// lt returns the boolean truth of arg1 < arg2.
func lt(a, b interface{}) bool {
	left, right := compareGetFloat(a, b)
	return left < right
}

func compareGetFloat(a interface{}, b interface{}) (float64, float64) {
	var left, right float64
	var leftStr, rightStr *string
	av := reflect.ValueOf(a)

	switch av.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		left = float64(av.Len())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		left = float64(av.Int())
	case reflect.Float32, reflect.Float64:
		left = av.Float()
	case reflect.String:
		var err error
		left, err = strconv.ParseFloat(av.String(), 64)
		if err != nil {
			str := av.String()
			leftStr = &str
		}
	case reflect.Struct:
		switch av.Type() {
		case timeType:
			left = float64(toTimeUnix(av))
		}
	}

	bv := reflect.ValueOf(b)

	switch bv.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		right = float64(bv.Len())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		right = float64(bv.Int())
	case reflect.Float32, reflect.Float64:
		right = bv.Float()
	case reflect.String:
		var err error
		right, err = strconv.ParseFloat(bv.String(), 64)
		if err != nil {
			str := bv.String()
			rightStr = &str
		}
	case reflect.Struct:
		switch bv.Type() {
		case timeType:
			right = float64(toTimeUnix(bv))
		}
	}

	switch {
	case leftStr == nil || rightStr == nil:
	case *leftStr < *rightStr:
		return 0, 1
	case *leftStr > *rightStr:
		return 1, 0
	default:
		return 0, 0
	}

	return left, right
}

// mod returns a % b.
func mod(a, b interface{}) (int64, error) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)
	var err error
	if av.Kind() == reflect.String {
		av, err = stringToNumber(av)
		if err != nil {
			return 0, err
		}
	}
	if bv.Kind() == reflect.String {
		bv, err = stringToNumber(bv)
		if err != nil {
			return 0, err
		}
	}

	var ai, bi int64

	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ai = av.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		ai = int64(av.Uint()) //may overflow
	default:
		return 0, errors.New("Modulo operator can't be used with non integer value")
	}

	switch bv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		bi = bv.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		bi = int64(bv.Uint()) //may overflow
	default:
		return 0, errors.New("Modulo operator can't be used with non integer value")
	}

	if bi == 0 {
		return 0, errors.New("The number can't be divided by zero at modulo operation")
	}

	return ai % bi, nil
}

// max returns the larger of a or b
func max(a, b interface{}) (float64, error) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)
	var err error
	if av.Kind() == reflect.String {
		av, err = stringToNumber(av)
		if err != nil {
			return 0, err
		}
	}
	if bv.Kind() == reflect.String {
		bv, err = stringToNumber(bv)
		if err != nil {
			return 0, err
		}
	}

	var af, bf float64

	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		af = float64(av.Int()) //may overflow
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		af = float64(av.Uint()) //may overflow
	case reflect.Float64, reflect.Float32:
		af = av.Float()
	default:
		return 0, errors.New("Max operator can't apply to the values")
	}

	switch bv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		bf = float64(bv.Int()) //may overflow
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		bf = float64(bv.Uint()) //may overflow
	case reflect.Float64, reflect.Float32:
		bf = bv.Float()
	default:
		return 0, errors.New("Max operator can't apply to the values")
	}

	return math.Max(af, bf), nil
}

// min returns the smaller of a or b
func min(a, b interface{}) (float64, error) {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)
	var err error
	if av.Kind() == reflect.String {
		av, err = stringToNumber(av)
		if err != nil {
			return 0, err
		}
	}
	if bv.Kind() == reflect.String {
		bv, err = stringToNumber(bv)
		if err != nil {
			return 0, err
		}
	}

	var af, bf float64

	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		af = float64(av.Int()) //may overflow
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		af = float64(av.Uint()) //may overflow
	case reflect.Float64, reflect.Float32:
		af = av.Float()
	default:
		return 0, errors.New("Max operator can't apply to the values")
	}

	switch bv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		bf = float64(bv.Int()) //may overflow
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		bf = float64(bv.Uint()) //may overflow
	case reflect.Float64, reflect.Float32:
		bf = bv.Float()
	default:
		return 0, errors.New("Max operator can't apply to the values")
	}

	return math.Min(af, bf), nil
}

func toTimeUnix(v reflect.Value) int64 {
	if v.Kind() == reflect.Interface {
		return toTimeUnix(v.Elem())
	}
	if v.Type() != timeType {
		panic("coding error: argument must be time.Time type reflect Value")
	}
	return v.MethodByName("Unix").Call([]reflect.Value{})[0].Int()
}

var kvType = reflect.TypeOf(libconfd.KVPair{}).Kind()

func Filter(regex string, c interface{}) ([]interface{}, error) {
	cv := reflect.ValueOf(c)

	switch cv.Kind() {
	case reflect.Array, reflect.Slice:
		result := make([]interface{}, 0, cv.Len())
		for i := 0; i < cv.Len(); i++ {
			v := cv.Index(i)
			if v.Kind() == reflect.Interface {
				v = reflect.ValueOf(v.Interface())
			}
			if v.Kind() == reflect.String {
				matched, err := regexp.MatchString(regex, v.String())
				if err != nil {
					return nil, err
				}
				if matched {
					result = append(result, v.String())
				}
			} else if v.Kind() == kvType {
				kv := v.Interface().(libconfd.KVPair)
				matched, err := regexp.MatchString(regex, kv.Value)
				if err != nil {
					return nil, err
				}
				if matched {
					result = append(result, kv)
				}
			}
		}
		return result, nil
	default:
		return nil, errors.New("filter only support slice or array.")
	}
}

func ToJson(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func ToYaml(v interface{}) (string, error) {
	b, err := yaml.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func Base64Encode(v interface{}) (string, error) {
	var input []byte
	s, ok := v.(string)
	if ok {
		input = []byte(s)
	} else {
		b, ok := v.([]byte)
		if ok {
			input = b
		}
	}
	if input != nil {
		return base64.StdEncoding.EncodeToString(input), nil
	}
	return "", fmt.Errorf("unsupported type %s", reflect.ValueOf(v).Kind().String())
}

func Base64Decode(v interface{}) (string, error) {
	var input string
	s, ok := v.(string)
	if ok {
		input = s
	} else {
		return "", fmt.Errorf("unsupported type %s", reflect.ValueOf(v).Kind().String())
	}
	r, err := base64.StdEncoding.DecodeString(input)
	return string(r), err
}
