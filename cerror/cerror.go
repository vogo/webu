/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cerror

import "net/http"

const (
	CodeOK         = 0
	CodeUnknownErr = 10
	CodeAuthErr    = 20
	CodeRequestErr = 100
)

var (
	ErrBadRequest      = NewStatusCodeError(http.StatusBadRequest, 100, "forbidden")
	ErrNotFound        = NewStatusCodeError(http.StatusNotFound, 101, "not found")
	ErrArgRequired     = NewStatusCodeError(http.StatusBadRequest, 102, "arg required")
	ErrValueInvalid    = NewStatusCodeError(http.StatusBadRequest, 103, "value invalid")
	ErrUnauthenticated = NewStatusCodeError(http.StatusUnauthorized, 20, "unauthenticated")
	ErrUnauthorized    = NewStatusCodeError(http.StatusUnauthorized, 21, "unauthorized")
	ErrForbidden       = NewStatusCodeError(http.StatusForbidden, 22, "forbidden")
)

type Coder interface {
	Code() int
}

type StatusState interface {
	Status() int
}

type CodeError interface {
	error
	Coder
}
type StatusCodeError interface {
	CodeError
	StatusState
}

func NewCodeError(code int, err string) CodeError {
	return &codeError{c: code, m: err}
}

type codeError struct {
	c int
	m string
}

func (e *codeError) Code() int {
	return e.c
}

func (e *codeError) Error() string {
	return e.m
}

func NewStatusCodeError(status, code int, err string) CodeError {
	return &statusCodeError{s: status, c: code, m: err}
}

type statusCodeError struct {
	s int
	c int
	m string
}

func (e *statusCodeError) Code() int {
	return e.c
}

func (e *statusCodeError) Status() int {
	return e.s
}

func (e *statusCodeError) Error() string {
	return e.m
}
