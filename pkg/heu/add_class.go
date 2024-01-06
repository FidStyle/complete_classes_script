package heu

import (
	"bytes"
	baseresp "compete_classes_script/pkg/base_resp"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

func AddClass(classid, seval, token, batchid string, choices ...string) (ok bool, err error) {
	tp := "XGKC"
	if len(choices) != 0 {
		if choices[0] == "professional" {
			tp = "TJKC"
		}
	}
	formval := url.Values{}
	formval.Set("clazzType", tp)
	formval.Set("clazzId", classid)
	formval.Set("secretVal", seval)

	formDataStr := formval.Encode()
	formDataBytes := []byte(formDataStr)
	formBytesReader := bytes.NewReader(formDataBytes)

	client := &http.Client{}
	url := "http://jwxk.hrbeu.edu.cn/xsxk/elective/clazz/add"
	request, err := http.NewRequest("POST", url, formBytesReader)
	if err != nil {
		return false, err
	}
	request.Header.Add("Authorization", token)
	request.Header.Add("BatchId", batchid)

	resp, err := client.Do(request)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()
	resd, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	type Result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	var res Result

	err = json.Unmarshal(resd, &res)
	if err != nil {
		return false, err
	}
	if res.Code == 200 && res.Msg == "操作成功" {
		return true, nil
	} else if res.Code == 500 && strings.Contains(res.Msg, "已选满5门") {
		return false, baseresp.ErrHeuFullClasses
	} else if res.Code == 401 && strings.Contains(res.Msg, "重新登录") {
		return false, baseresp.ErrHeuReLogin
	} else if res.Code == 403 && strings.Contains(res.Msg, "请求过快") {
		return false, nil
	} else if res.Code == 500 && strings.Contains(res.Msg, "已在选课结果中") {
		return false, nil
	} else if res.Code == 500 && strings.Contains(res.Msg, "该教学班不可选") {
		return false, nil
	} else if res.Code == 500 && strings.Contains(res.Msg, "不能重复选课") {
		return false, nil
	} else if res.Code != 200 {
		return false, fmt.Errorf("AddClass Faild because of %v", res.Msg)
	}

	return false, nil
}

func AddClassVolunteer(classid, seval, token, batchid string, choices ...string) (ok bool, err error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	tp := "XGKC"
	if len(choices) != 0 {
		if choices[0] == "professional" {
			tp = "TJKC"
		}
	}
	_ = writer.WriteField("clazzType", tp)
	_ = writer.WriteField("clazzId", classid)
	_ = writer.WriteField("secretVal", seval)
	_ = writer.WriteField("chooseVolunteer", "1")
	err = writer.Close()
	if err != nil {
		return
	}

	client := &http.Client{}
	url := "http://jwxk.hrbeu.edu.cn/xsxk/elective/clazz/add"
	request, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return false, err
	}
	request.Header.Add("Authorization", token)
	request.Header.Add("BatchId", batchid)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(request)
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()
	resd, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	type Result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}

	var res Result

	err = json.Unmarshal(resd, &res)
	if err != nil {
		return false, err
	}

	if res.Code == 200 && res.Msg == "操作成功" {
		return true, nil
	} else if res.Code == 500 && strings.Contains(res.Msg, "已选满5门") {
		return false, baseresp.ErrHeuFullClasses
	} else if res.Code == 401 && strings.Contains(res.Msg, "重新登录") {
		return false, baseresp.ErrHeuReLogin
	} else if res.Code == 403 && strings.Contains(res.Msg, "请求过快") {
		return false, nil
	} else if res.Code == 500 && strings.Contains(res.Msg, "已在选课结果中") {
		return false, nil
	} else if res.Code == 500 && strings.Contains(res.Msg, "该教学班不可选") {
		return false, nil
	} else if res.Code == 500 && strings.Contains(res.Msg, "不能重复选课") {
		return false, nil
	} else if res.Code != 200 {
		return false, fmt.Errorf("AddClassVolunteer Faild because of %v with code: %v", res.Msg, res.Code)
	}

	return false, nil
}
