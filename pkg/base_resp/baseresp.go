package baseresp

import (
	"compete_classes_script/pkg/logger"
	"errors"
)

type BaseResp struct {
	Code    int    `json:"code"`
	Content string `json:"content"`
}

func ErrorResp(err error) *BaseResp {
	resp, ok := ErrRespMp[err]
	if !ok {
		resp = ErrRespMp[ErrServerFail]
	}

	logger.Error(err)

	return resp
}

func SuccessResp() *BaseResp {
	return &BaseResp{
		Code:    200,
		Content: "Success",
	}
}

var (
	ErrServerFail          error = errors.New("server fail")
	ErrLoginWrongPw        error = errors.New("password is wrong")
	ErrLoginUnknownAccount error = errors.New("user does not exist")
	ErrMissToken           error = errors.New("token dose not set")
	ErrExpiredToken        error = errors.New("token expired, please log in again")
	ErrInvalidArgument     error = errors.New("argument is invalid")
	ErrOrderIDNotExist     error = errors.New("order id does not exist")
	ErrInfoUnsuccessOrder  error = errors.New("can not info unsuccess order")
	ErrAuthInvalid         error = errors.New("invalid auth")

	ErrHeuReLogin            error = errors.New("please log in again")
	ErrHeuFullClasses        error = errors.New("classes if full")
	ErrSpecifyPublicNotFound error = errors.New("public specify not found")
	ErrFuzzyNameNotFound     error = errors.New("fuzzy name not found")
	ErrKindNotFound          error = errors.New("nothing in this kind")
	ErrKindNotFoundInWrap    error = errors.New("nothing in this kind of this wrap")

	ErrRespMp map[error]*BaseResp = map[error]*BaseResp{
		ErrServerFail:          {500, ErrServerFail.Error()},
		ErrLoginWrongPw:        {1001, ErrLoginWrongPw.Error()},
		ErrLoginUnknownAccount: {1002, ErrLoginUnknownAccount.Error()},
		ErrMissToken:           {1003, ErrMissToken.Error()},
		ErrExpiredToken:        {1004, ErrExpiredToken.Error()},
		ErrInvalidArgument:     {1005, ErrInvalidArgument.Error()},
		ErrOrderIDNotExist:     {1006, ErrOrderIDNotExist.Error()},
		ErrInfoUnsuccessOrder:  {1007, ErrInfoUnsuccessOrder.Error()},
		ErrAuthInvalid:         {1008, ErrAuthInvalid.Error()},
	}
)
