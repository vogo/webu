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
	"encoding/json"
	"github.com/vogo/webu/httputil"
	"html/template"

	"github.com/vogo/logger"
	"github.com/vogo/webu/cerror"

	"github.com/valyala/fasthttp"
)

func ResponseData(ctx *fasthttp.RequestCtx, code int, data interface{}) {
	WriteResp(ctx, code, "", data)
}

func ResponseCodeData(ctx *fasthttp.RequestCtx, code int, msg string, data interface{}) {
	WriteResp(ctx, code, msg, data)
}

func ResponseOK(ctx *fasthttp.RequestCtx) {
	WriteResp(ctx, cerror.CodeOK, "", "ok")
}

func ResponseSuccess(ctx *fasthttp.RequestCtx, data interface{}) {
	WriteResp(ctx, cerror.CodeOK, "", data)
}

func ResponseCodeError(ctx *fasthttp.RequestCtx, code int, err error) {
	ResponseCodeMsg(ctx, code, err.Error())
}

func ErrorResponse(ctx *fasthttp.RequestCtx, err error) {
	if c, ok := err.(cerror.StatusState); ok {
		ctx.SetStatusCode(c.Status())
	}

	code := cerror.CodeUnknownErr

	if c, ok := err.(cerror.Coder); ok {
		code = c.Code()
	}

	ResponseCodeMsg(ctx, code, err.Error())
}

func ResponseBadMsg(ctx *fasthttp.RequestCtx, msg string) {
	ResponseCodeMsg(ctx, cerror.CodeBadRequestErr, msg)
}

func ResponseBadError(ctx *fasthttp.RequestCtx, err error) {
	logger.Errorf("bad request: %v", err)
	ResponseBadMsg(ctx, err.Error())
}

func ResponseCodeMsg(ctx *fasthttp.RequestCtx, code int, msg string) {
	WriteResp(ctx, code, msg, nil)
}

func WriteResp(ctx *fasthttp.RequestCtx, code int, msg string, data interface{}) {
	resp := httputil.ResponseBody{
		Code: code,
		Msg:  msg,
		Data: data,
	}

	jsonBytes, err := json.Marshal(resp)
	if err != nil {
		logger.Infof("json marshal error: %+v", err)

		_, _ = ctx.Write([]byte("internal server error"))

		return
	}

	ctx.SetContentType("application/json")
	_, _ = ctx.Write(jsonBytes)
}

func ResponseTemplate(ctx *fasthttp.RequestCtx, tpl *template.Template, data interface{}) {
	ctx.SetContentType("text/html")
	err := tpl.Execute(ctx.Response.BodyWriter(), data)
	if err != nil {
		logger.Fatalf("template format error: %v", err)
	}
}
