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

package fastutil

import (
	"strconv"

	"github.com/valyala/fasthttp"
	"github.com/vogo/webu/cerror"
)

func RequireQueryArg(ctx *fasthttp.RequestCtx, arg string) ([]byte, error) {
	val := ctx.QueryArgs().Peek(arg)
	if len(val) == 0 {
		return nil, cerror.ErrArgRequired
	}

	return val, nil
}

func RequireQueryString(ctx *fasthttp.RequestCtx, arg string) (string, error) {
	val := ctx.QueryArgs().Peek(arg)
	if len(val) == 0 {
		return "", cerror.ErrArgRequired
	}

	return string(val), nil
}

func RequireQueryInt(ctx *fasthttp.RequestCtx, arg string) (int, error) {
	val, err := RequireQueryArg(ctx, arg)
	if err != nil {
		return 0, err
	}

	i, err := strconv.Atoi(string(val))
	if err != nil {
		return 0, cerror.ErrValueInvalid
	}

	return i, nil
}
