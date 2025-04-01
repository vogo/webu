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

package httpresp

import (
	"encoding/json"
	"net/http"

	"github.com/vogo/logger"
	"github.com/vogo/webu/cerror"
)

type ResponseBody struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func ResponseData(w http.ResponseWriter, req *http.Request, code int, data interface{}) {
	WriteResp(w, req, code, "", data)
}

func ResponseCodeData(w http.ResponseWriter, req *http.Request, code int, msg string, data interface{}) {
	WriteResp(w, req, code, msg, data)
}

func ResponseOK(w http.ResponseWriter, req *http.Request) {
	WriteResp(w, req, cerror.CodeOK, "", "ok")
}

func ResponseSuccess(w http.ResponseWriter, req *http.Request, data interface{}) {
	WriteResp(w, req, cerror.CodeOK, "", data)
}

func ResponseCodeError(w http.ResponseWriter, req *http.Request, code int, err error) {
	ResponseCodeMsg(w, req, code, err.Error())
}

func ErrorResponse(w http.ResponseWriter, req *http.Request, err error) {
	if c, ok := err.(cerror.StatusState); ok {
		w.WriteHeader(c.Status())
	}

	code := cerror.CodeUnknownErr

	if c, ok := err.(cerror.Coder); ok {
		code = c.Code()
	}

	ResponseCodeMsg(w, req, code, err.Error())
}

func ResponseBadMsg(w http.ResponseWriter, req *http.Request, msg string) {
	ResponseCodeMsg(w, req, cerror.CodeBadRequestErr, msg)
}

func ResponseBadError(w http.ResponseWriter, req *http.Request, err error) {
	logger.Errorf("bad request: %v", err)
	ResponseBadMsg(w, req, err.Error())
}

func ResponseCodeMsg(w http.ResponseWriter, req *http.Request, code int, msg string) {
	WriteResp(w, req, code, msg, nil)
}

func WriteResp(w http.ResponseWriter, req *http.Request, code int, msg string, data interface{}) {
	resp := ResponseBody{
		Code: code,
		Msg:  msg,
		Data: data,
	}

	jsonBytes, err := json.Marshal(resp)
	if err != nil {
		logger.Infof("json marshal error: %+v", err)

		_, _ = w.Write([]byte("internal server error"))

		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jsonBytes)
}
