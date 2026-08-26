// Package apiRsp 提供 HTTP API 的统一业务响应类型。
package apiRsp

import (
	"fmt"
	"net/http"
)

const (
	InternalServerError = -1
	SUCCESS             = 0
	ERROR               = 7
	PARAMS              = 400
)

const (
	Msg_Unknown             = "Unknown Error"
	Msg_OperationSuccessful = "Operation Successful"
	Msg_OperationFailed     = "Operation Failed"
	Msg_Success             = "Success"
	Msg_Params              = "Params Error"
)

type IApiRsp interface {
	HttpStatus() int
	Code() int
	Msg() string
	Data() any
}

type ApiRsp struct {
	status int    `json:"status"`
	code   int    `json:"code"`
	msg    string `json:"msg"`
	data   any    `json:"data"`
}

func (r ApiRsp) Error() string {
	return fmt.Sprintf("WantCode: %d, msg: %s", r.code, r.msg)
}

func (r ApiRsp) HttpStatus() int {
	return r.status
}

func (r ApiRsp) Code() int {
	return r.code
}

func (r ApiRsp) Msg() string {
	return r.msg
}

func (r ApiRsp) Data() any {
	return r.data
}

func (r ApiRsp) Info(status int, code int, data any, msg string, values ...any) IApiRsp {
	return &ApiRsp{
		status: status,
		code:   code,
		data:   data,
		msg:    fmt.Sprintf(msg, values...),
	}
}

func (r ApiRsp) Ok() IApiRsp {
	return r.OkWithMsg(Msg_OperationSuccessful)
}

func (r ApiRsp) OkWithMsg(msg string, values ...any) IApiRsp {
	return r.OkWithDet(nil, msg, values...)
}

func (r ApiRsp) OkWithData(data any) IApiRsp {
	return r.OkWithDet(data, Msg_OperationSuccessful)
}

func (r ApiRsp) OkWithDet(data any, msg string, values ...any) IApiRsp {
	return r.Info(http.StatusOK, SUCCESS, data, msg, values...)
}

func (r ApiRsp) Fail() IApiRsp {
	return r.FailWithDet(nil, Msg_OperationFailed)
}

func (r ApiRsp) FailWithMsg(msg string, values ...any) IApiRsp {
	return r.FailWithDet(nil, msg, values...)
}

func (r ApiRsp) FailWithDet(data any, msg string, values ...any) IApiRsp {
	return r.Info(http.StatusBadRequest, ERROR, data, msg, values...)
}

func (r ApiRsp) NoAuth(msg string, values ...any) IApiRsp {
	return r.Info(http.StatusUnauthorized, ERROR, nil, msg, values...)
}
