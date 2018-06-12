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

	HookKeyAdjuster func(key string) (realKey string)
}

var _TemplateFunc_initFuncMap func(p *TemplateFunc) = nil

func NewTemplateFunc(
	store *KVStore, pgpPrivateKey []byte,
	hookKeyAdjuster func(key string) (realKey string),
) *TemplateFunc {
	p := &TemplateFunc{
		FuncMap:         map[string]interface{}{},
		Store:           store,
		PGPPrivateKey:   pgpPrivateKey,
		HookKeyAdjuster: hookKeyAdjuster,
	}

	if _TemplateFunc_initFuncMap == nil {
		GetLogger().Panic("_TemplateFunc_initFuncMap missing")
	}

	_TemplateFunc_initFuncMap(p)
	return p
}

// ----------------------------------------------------------------------------
// KVStore
// ----------------------------------------------------------------------------

func (p TemplateFunc) Exists(key string) bool {
	if p.HookKeyAdjuster != nil {
		key = p.HookKeyAdjuster(key)
	}
	return p.Store.Exists(key)
}

func (p TemplateFunc) Ls(filepath string) []string {
	if p.HookKeyAdjuster != nil {
		filepath = p.HookKeyAdjuster(filepath)
	}
	return p.Store.List(filepath)
}

func (p TemplateFunc) Lsdir(filepath string) []string {
	if p.HookKeyAdjuster != nil {
		filepath = p.HookKeyAdjuster(filepath)
	}
	return p.Store.ListDir(filepath)
}

func (p TemplateFunc) Get(key string) KVPair {
	if p.HookKeyAdjuster != nil {
		key = p.HookKeyAdjuster(key)
	}
	v, ok := p.Store.Get(key)
	if !ok {
		GetLogger().Error("key not exits:", key)
		return KVPair{}
	}
	return v
}

func (p TemplateFunc) Gets(pattern string) []KVPair {
	if p.HookKeyAdjuster != nil {
		pattern = p.HookKeyAdjuster(pattern)
	}
	v, err := p.Store.GetAll(pattern)
	if err != nil {
		GetLogger().Error(err)
		return nil
	}
	return v
}

func (p TemplateFunc) Getv(key string, v ...string) string {
	if p.HookKeyAdjuster != nil {
		key = p.HookKeyAdjuster(key)
	}

	value, ok := p.Store.GetValue(key, v...)
	if !ok {
		GetLogger().Error("key not exits:", key)
		return ""
	}
	return value
}

func (p TemplateFunc) Getvs(pattern string) []string {
	if p.HookKeyAdjuster != nil {
		pattern = p.HookKeyAdjuster(pattern)
	}
	v, err := p.Store.GetAllValues(pattern)
	if err != nil {
		GetLogger().Error(err)
		return nil
	}
	return v
}

// ----------------------------------------------------------------------------
// Crypt func
// ----------------------------------------------------------------------------

func (p TemplateFunc) Cget(key string) KVPair {
	if len(p.PGPPrivateKey) == 0 {
		err := fmt.Errorf("PGPPrivateKey is empty")
		GetLogger().Error(err)
		return KVPair{}
	}

	if p.HookKeyAdjuster != nil {
		key = p.HookKeyAdjuster(key)
	}
	kv, ok := p.Store.Get(key)
	if !ok {
		return KVPair{}
	}

	var b []byte
	b, err := secconfDecode([]byte(kv.Value), bytes.NewBuffer(p.PGPPrivateKey))
	if err != nil {
		GetLogger().Error(err)
		return KVPair{}
	}

	kv.Value = string(b)
	return kv
}

func (p TemplateFunc) Cgets(pattern string) []KVPair {
	if len(p.PGPPrivateKey) == 0 {
		err := fmt.Errorf("PGPPrivateKey is empty")
		GetLogger().Error(err)
		return nil
	}

	if p.HookKeyAdjuster != nil {
		pattern = p.HookKeyAdjuster(pattern)
	}
	kvs, err := p.Store.GetAll(pattern)
	if err != nil {
		GetLogger().Error(err)
		return nil
	}

	for i := range kvs {
		b, err := secconfDecode([]byte(kvs[i].Value), bytes.NewBuffer(p.PGPPrivateKey))
		if err != nil {
			GetLogger().Error(err)
			return nil
		}
		kvs[i].Value = string(b)
	}
	return kvs
}

func (p TemplateFunc) Cgetv(key string) string {
	if len(p.PGPPrivateKey) == 0 {
		err := fmt.Errorf("PGPPrivateKey is empty")
		GetLogger().Error(err)
		return ""
	}

	if p.HookKeyAdjuster != nil {
		key = p.HookKeyAdjuster(key)
	}
	v, ok := p.Store.GetValue(key)
	if !ok {
		return ""
	}

	var b []byte
	b, err := secconfDecode([]byte(v), bytes.NewBuffer(p.PGPPrivateKey))
	if err != nil {
		GetLogger().Error(err)
		return ""
	}

	return string(b)
}

func (p TemplateFunc) Cgetvs(pattern string) []string {
	if len(p.PGPPrivateKey) == 0 {
		err := fmt.Errorf("PGPPrivateKey is empty")
		GetLogger().Error(err)
		return nil
	}

	if p.HookKeyAdjuster != nil {
		pattern = p.HookKeyAdjuster(pattern)
	}
	vs, err := p.Store.GetAllValues(pattern)
	if err != nil {
		GetLogger().Error(err)
		return nil
	}

	for i := range vs {
		b, err := secconfDecode([]byte(vs[i]), bytes.NewBuffer(p.PGPPrivateKey))
		if err != nil {
			GetLogger().Error(err)
			return nil
		}
		vs[i] = string(b)
	}
	return vs
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

func (_ TemplateFunc) Json(data string) map[string]interface{} {
	var ret map[string]interface{}
	err := json.Unmarshal([]byte(data), &ret)
	if err != nil {
		GetLogger().Error(err)
		return nil
	}
	return ret
}

func (_ TemplateFunc) JsonArray(data string) []interface{} {
	var ret []interface{}
	err := json.Unmarshal([]byte(data), &ret)
	if err != nil {
		GetLogger().Error(err)
		return nil
	}
	return ret
}

func (_ TemplateFunc) Dir(path string) string {
	return pathpkg.Dir(path)
}

// Map creates a key-value map of string -> interface{}
// The i'th is the key and the i+1 is the value
func (_ TemplateFunc) Map(values ...interface{}) map[string]interface{} {
	if len(values)%2 != 0 {
		GetLogger().Error(errors.New("invalid map call"))
		return nil
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			GetLogger().Error(errors.New("map keys must be strings"))
			return nil
		}
		dict[key] = values[i+1]
	}
	return dict
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
		GetLogger().Error(err)
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

func (_ TemplateFunc) LookupIPV6(data string) []string {
	var addresses []string
	for _, ip := range (TemplateFunc{}).LookupIP(data) {
		if strings.Contains(ip, ":") {
			addresses = append(addresses, ip)
		}
	}
	return addresses
}

func (_ TemplateFunc) LookupIPV4(data string) []string {
	var addresses []string
	for _, ip := range (TemplateFunc{}).LookupIP(data) {
		if strings.Contains(ip, ".") {
			addresses = append(addresses, ip)
		}
	}
	return addresses
}

func (_ TemplateFunc) LookupSRV(service, proto, name string) []*net.SRV {
	_, s, err := net.LookupSRV(service, proto, name)
	if err != nil {
		GetLogger().Error(err)
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

func (_ TemplateFunc) Base64Decode(data string) string {
	s, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		GetLogger().Error(err)
		return ""
	}
	return string(s)
}

func (_ TemplateFunc) ParseBool(s string) bool {
	v, err := strconv.ParseBool(s)
	if err != nil {
		GetLogger().Error(err)
		return false
	}
	return v
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

func (_ TemplateFunc) Atoi(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		GetLogger().Error(err)
		return 0
	}
	return v
}

// ----------------------------------------------------------------------------
// END
// ----------------------------------------------------------------------------
