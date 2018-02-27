// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package db

import (
	"github.com/gocraft/dbr"
)

// And creates AND from a list of conditions
func And(cond ...dbr.Builder) dbr.Builder {
	return dbr.And(cond...)
}

// Or creates OR from a list of conditions
func Or(cond ...dbr.Builder) dbr.Builder {
	return dbr.Or(cond...)
}

// Eq is `=`.
// When value is nil, it will be translated to `IS NULL`.
// When value is a slice, it will be translated to `IN`.
// Otherwise it will be translated to `=`.
func Eq(column string, value interface{}) dbr.Builder {
	return dbr.Eq(column, value)
}

// Neq is `!=`.
// When value is nil, it will be translated to `IS NOT NULL`.
// When value is a slice, it will be translated to `NOT IN`.
// Otherwise it will be translated to `!=`.
func Neq(column string, value interface{}) dbr.Builder {
	return dbr.Neq(column, value)
}

// Gt is `>`.
func Gt(column string, value interface{}) dbr.Builder {
	return dbr.Gt(column, value)
}

// Gte is '>='.
func Gte(column string, value interface{}) dbr.Builder {
	return dbr.Gte(column, value)
}

// Lt is '<'.
func Lt(column string, value interface{}) dbr.Builder {
	return dbr.Lt(column, value)
}

// Lte is `<=`.
func Lte(column string, value interface{}) dbr.Builder {
	return dbr.Lte(column, value)
}
