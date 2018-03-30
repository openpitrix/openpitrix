// Copyright confd. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE-confd file.

//go:generate go run gen_tmpl_func_list.go -output tmpl_funcs_zz.go

package libconfd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	pathpkg "path"
	"sort"
	"strconv"
	"strings"
	"time"
)

type TemplateFunc struct {
	FuncMap       map[string]interface{}
	Store         *KVStore
	PGPPrivateKey []byte
}

var _TemplateFunc_initFuncMap func(p *TemplateFunc) = nil

func NewTemplateFunc(store *KVStore, pgpPrivateKey []byte) *TemplateFunc {
	p := &TemplateFunc{
		FuncMap:       map[string]interface{}{},
		Store:         store,
		PGPPrivateKey: pgpPrivateKey,
	}

	if _TemplateFunc_initFuncMap == nil {
		logger.Panic("_TemplateFunc_initFuncMap missing")
	}

	_TemplateFunc_initFuncMap(p)
	return p
}

// ----------------------------------------------------------------------------
// KVStore
// ----------------------------------------------------------------------------

func (p TemplateFunc) Exists(key string) bool {
	return p.Store.Exists(key)
}

func (p TemplateFunc) Ls(filepath string) []string {
	return p.Store.List(filepath)
}

func (p TemplateFunc) Lsdir(filepath string) []string {
	return p.Store.ListDir(filepath)
}

func (p TemplateFunc) Get(key string) (KVPair, error) {
	if v, ok := p.Store.Get(key); ok {
		return v, nil
	}
	return KVPair{}, fmt.Errorf("key not exists")
}

func (p TemplateFunc) Gets(pattern string) ([]KVPair, error) {
	return p.Store.GetAll(pattern)
}

func (p TemplateFunc) Getv(key string, v ...string) (string, error) {
	if v, ok := p.Store.GetValue(key, v...); ok {
		return v, nil
	}
	return "", fmt.Errorf("key not exists")
}

func (p TemplateFunc) Getvs(pattern string) ([]string, error) {
	return p.Store.GetAllValues(pattern)
}

// ----------------------------------------------------------------------------
// Crypt func
// ----------------------------------------------------------------------------

func (p TemplateFunc) Cget(key string) (KVPair, error) {
	if len(p.PGPPrivateKey) == 0 {
		return KVPair{}, fmt.Errorf("PGPPrivateKey is empty")
	}

	kv, err := p.FuncMap["get"].(func(string) (KVPair, error))(key)
	if err != nil {
		return KVPair{}, err
	}

	var b []byte
	b, err = secconfDecode([]byte(kv.Value), bytes.NewBuffer(p.PGPPrivateKey))
	if err != nil {
		return KVPair{}, err
	}

	kv.Value = string(b)
	return kv, nil
}

func (p TemplateFunc) Cgets(pattern string) ([]KVPair, error) {
	if len(p.PGPPrivateKey) == 0 {
		return nil, fmt.Errorf("PGPPrivateKey is empty")
	}

	kvs, err := p.FuncMap["gets"].(func(string) ([]KVPair, error))(pattern)
	if err != nil {
		return nil, err
	}

	for i := range kvs {
		b, err := secconfDecode([]byte(kvs[i].Value), bytes.NewBuffer(p.PGPPrivateKey))
		if err != nil {
			return nil, err
		}
		kvs[i].Value = string(b)
	}
	return kvs, nil
}

func (p TemplateFunc) Cgetv(key string) (string, error) {
	if len(p.PGPPrivateKey) == 0 {
		return "", fmt.Errorf("PGPPrivateKey is empty")
	}

	v, err := p.FuncMap["getv"].(func(string, ...string) (string, error))(key)
	if err != nil {
		return "", err
	}

	var b []byte
	b, err = secconfDecode([]byte(v), bytes.NewBuffer(p.PGPPrivateKey))
	if err != nil {
		return "", err
	}

	return string(b), err
}

func (p TemplateFunc) Cgetvs(pattern string) ([]string, error) {
	if len(p.PGPPrivateKey) == 0 {
		return nil, fmt.Errorf("PGPPrivateKey is empty")
	}

	vs, err := p.FuncMap["getvs"].(func(string) ([]string, error))(pattern)
	if err != nil {
		return nil, err
	}

	for i := range vs {
		b, err := secconfDecode([]byte(vs[i]), bytes.NewBuffer(p.PGPPrivateKey))
		if err != nil {
			return nil, err
		}
		vs[i] = string(b)
	}
	return vs, nil
}

// ----------------------------------------------------------------------------
// util func
// ----------------------------------------------------------------------------

func (_ TemplateFunc) Base(path string) string {
	return pathpkg.Base(path)
}

func (_ TemplateFunc) Split(s, sep string) []string {
	return strings.Split(s, sep)
}

func (_ TemplateFunc) Json(data string) (map[string]interface{}, error) {
	var ret map[string]interface{}
	err := json.Unmarshal([]byte(data), &ret)
	return ret, err
}

func (_ TemplateFunc) JsonArray(data string) ([]interface{}, error) {
	var ret []interface{}
	err := json.Unmarshal([]byte(data), &ret)
	return ret, err
}

func (_ TemplateFunc) Dir(path string) string {
	return pathpkg.Dir(path)
}

// Map creates a key-value map of string -> interface{}
// The i'th is the key and the i+1 is the value
func (_ TemplateFunc) Map(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid map call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("map keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

// getenv retrieves the value of the environment variable named by the key.
// It returns the value, which will the default value if the variable is not present.
// If no default value was given - returns "".
func (_ TemplateFunc) Getenv(key string, defaultValue ...string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}

func (_ TemplateFunc) Join(a []string, sep string) string {
	return strings.Join(a, sep)
}

func (_ TemplateFunc) Datetime() time.Time {
	return time.Now()
}

func (_ TemplateFunc) ToUpper(s string) string {
	return strings.ToUpper(s)
}

func (_ TemplateFunc) ToLower(s string) string {
	return strings.ToLower(s)
}

func (_ TemplateFunc) Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func (_ TemplateFunc) Replace(s, old, new string, n int) string {
	return strings.Replace(s, old, new, n)
}

func (_ TemplateFunc) TrimSuffix(s, suffix string) string {
	return strings.TrimSuffix(s, suffix)
}

func (_ TemplateFunc) LookupIP(data string) []string {
	ips, err := net.LookupIP(data)
	if err != nil {
		return nil
	}
	// "Cast" IPs into strings and sort the array
	ipStrings := make([]string, len(ips))

	for i, ip := range ips {
		ipStrings[i] = ip.String()
	}
	sort.Strings(ipStrings)
	return ipStrings
}

func (_ TemplateFunc) LookupSRV(service, proto, name string) []*net.SRV {
	_, s, err := net.LookupSRV(service, proto, name)
	if err != nil {
		return nil
	}

	sort.Slice(s, func(i, j int) bool {
		str1 := fmt.Sprintf("%s%d%d%d", s[i].Target, s[i].Port, s[i].Priority, s[i].Weight)
		str2 := fmt.Sprintf("%s%d%d%d", s[j].Target, s[j].Port, s[j].Priority, s[j].Weight)
		return str1 < str2
	})
	return s
}

func (_ TemplateFunc) FileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return err == nil
}

func (_ TemplateFunc) Base64Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func (_ TemplateFunc) Base64Decode(data string) (string, error) {
	s, err := base64.StdEncoding.DecodeString(data)
	return string(s), err
}

func (_ TemplateFunc) ParseBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

// reverse returns the array in reversed order
// works with []string and []KVPair
func (_ TemplateFunc) Reverse(values interface{}) interface{} {
	switch values.(type) {
	case []string:
		v := values.([]string)
		for left, right := 0, len(v)-1; left < right; left, right = left+1, right-1 {
			v[left], v[right] = v[right], v[left]
		}
	case []KVPair:
		v := values.([]KVPair)
		for left, right := 0, len(v)-1; left < right; left, right = left+1, right-1 {
			v[left], v[right] = v[right], v[left]
		}
	}
	return values
}

func (_ TemplateFunc) SortKVByLength(values []KVPair) []KVPair {
	sort.Slice(values, func(i, j int) bool {
		return len(values[i].Key) < len(values[j].Key)
	})
	return values
}

func (_ TemplateFunc) SortByLength(values []string) []string {
	sort.Slice(values, func(i, j int) bool {
		return len(values[i]) < len(values[j])
	})
	return values
}

func (_ TemplateFunc) Add(a, b int) int {
	return a + b
}
func (_ TemplateFunc) Sub(a, b int) int {
	return a - b
}
func (_ TemplateFunc) Div(a, b int) int {
	return a / b
}
func (_ TemplateFunc) Mod(a, b int) int {
	return a % b
}
func (_ TemplateFunc) Mul(a, b int) int {
	return a * b
}

// seq creates a sequence of integers. It's named and used as GNU's seq.
// seq takes the first and the last element as arguments. So Seq(3, 5) will generate [3,4,5]
func (_ TemplateFunc) Seq(first, last int) []int {
	var arr []int
	for i := first; i <= last; i++ {
		arr = append(arr, i)
	}
	return arr
}

func (_ TemplateFunc) Atoi(s string) (int, error) {
	return strconv.Atoi(s)
}

// ----------------------------------------------------------------------------
// END
// ----------------------------------------------------------------------------
