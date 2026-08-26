package apiRsp

import (
	"net/http"
	"testing"
)

func TestApiRsp(t *testing.T) {
	testCases := []struct {
		name       string
		response   IApiRsp
		wantStatus int
		wantCode   int
		wantMsg    string
		wantData   any
	}{
		{
			name:       "ok",
			response:   ApiRsp{}.Ok(),
			wantStatus: http.StatusOK,
			wantCode:   SUCCESS,
			wantMsg:    Msg_OperationSuccessful,
		},
		{
			name:       "ok with data",
			response:   ApiRsp{}.OkWithData("payload"),
			wantStatus: http.StatusOK,
			wantCode:   SUCCESS,
			wantMsg:    Msg_OperationSuccessful,
			wantData:   "payload",
		},
		{
			name:       "formatted failure",
			response:   ApiRsp{}.FailWithMsg("failed %d", 7),
			wantStatus: http.StatusBadRequest,
			wantCode:   ERROR,
			wantMsg:    "failed 7",
		},
		{
			name:       "unauthorized",
			response:   ApiRsp{}.NoAuth("denied"),
			wantStatus: http.StatusUnauthorized,
			wantCode:   ERROR,
			wantMsg:    "denied",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if got := testCase.response.HttpStatus(); got != testCase.wantStatus {
				t.Fatalf("HttpStatus() = %d", got)
			}
			if got := testCase.response.Code(); got != testCase.wantCode {
				t.Fatalf("Code() = %d", got)
			}
			if got := testCase.response.Msg(); got != testCase.wantMsg {
				t.Fatalf("Msg() = %q", got)
			}
			if got := testCase.response.Data(); got != testCase.wantData {
				t.Fatalf("Data() = %#v", got)
			}
		})
	}
}

func TestApiRspError(t *testing.T) {
	response := ApiRsp{}.Info(http.StatusTeapot, 99, nil, "problem").(*ApiRsp)
	if got := response.Error(); got != "WantCode: 99, msg: problem" {
		t.Fatalf("Error() = %q", got)
	}
}
