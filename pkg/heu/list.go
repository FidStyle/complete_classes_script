package heu

import (
	"bytes"
	baseresp "compete_classes_script/pkg/base_resp"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var KindMap map[string]string = map[string]string{"f": "17", "a0": "18", "d": "15", "e": "16", "c": "14", "b": "13", "a": "12"}

type list struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Rows []struct {
			// classID
			JXBID string `json:"JXBID"`
			// className
			KCM string `json:"KCM"`
			// score
			XF string `json:"XF"`
			// capacity of class
			KRL int `json:"KRL"`
			// selected number of people
			YXRS int `json:"YXRS"`
			// pre selected number of people
			DYZYRS    int    `json:"DYZYRS"`
			SecretVal string `json:"secretVal"`
		} `json:"rows"`
	} `json:"data"`
}

func GetList(kind string, token, batchid string) (rows *list, err error) {
	client := &http.Client{}
	url := "http://jwxk.hrbeu.edu.cn/xsxk/elective/clazz/list"
	request, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(fmt.Sprintf(`{"teachingClassType":"XGKC","pageNumber":1,"pageSize":300,"orderBy":"","campus":"01","SFCT":"0","XGXKLB": "%v"}`, KindMap[kind]))))
	if err != nil {
		return
	}

	request.Header.Add("Content-Type", "application/json;charset=UTF-8")
	request.Header.Add("Authorization", token)
	request.Header.Add("batchId", batchid)

	resp, err := client.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("GetList resp.StatusCode = %v", resp.StatusCode)
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &rows)
	if err != nil {
		return
	}

	if rows.Code != 200 {
		err = fmt.Errorf("GetList Rows.Code = %v Msg: %v", rows.Code, rows.Msg)
		return
	} else if rows.Code == 401 && strings.Contains(rows.Msg, "重新登录") {
		err = baseresp.ErrHeuReLogin
		return
	}

	return
}

func GetListByFuzzyName(fuzzyName string, token string, batchid string) (rows *list, err error) {
	client := &http.Client{}
	url := "http://jwxk.hrbeu.edu.cn/xsxk/elective/clazz/list"
	request, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(fmt.Sprintf(`{"teachingClassType":"XGKC","pageNumber":1,"pageSize":300,"orderBy":"","campus":"01","SFCT":"0", "KEY":"%v"}`, fuzzyName))))
	if err != nil {
		return
	}

	request.Header.Add("Content-Type", "application/json;charset=UTF-8")
	request.Header.Add("Authorization", token)
	request.Header.Add("batchId", batchid)

	resp, err := client.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("GetList resp.StatusCode = %v", resp.StatusCode)
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &rows)
	if err != nil {
		return
	}

	if rows.Code != 200 {
		err = fmt.Errorf("GetList Rows.Code = %v Msg: %v", rows.Code, rows.Msg)
		return
	} else if rows.Code == 401 && strings.Contains(rows.Msg, "重新登录") {
		err = baseresp.ErrHeuReLogin
		return
	}

	return
}

type professionalList struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Rows []struct {
			TcList []struct {
				// classID
				JXBID string `json:"JXBID"`
				// className
				KCM string `json:"KCM"`
				// score
				XF string `json:"XF"`
				// capacity of class
				KRL int `json:"KRL"`
				// selected number of people
				YXRS int `json:"YXRS"`
				// pre selected number of people
				DYZYRS    int    `json:"DYZYRS"`
				SecretVal string `json:"secretVal"`
			} `json:"tcList"`
		} `json:"rows"`
	} `json:"data"`
}

func GetListProfessional(token, batchid string) (rows *professionalList, err error) {
	client := &http.Client{}
	url := "http://jwxk.hrbeu.edu.cn/xsxk/elective/clazz/list"
	request, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(fmt.Sprintln(`{"teachingClassType":"TJKC","pageNumber":1,"pageSize":300,"orderBy":"","campus":"01","SFCT":"0"}`))))
	if err != nil {
		return
	}

	request.Header.Add("Content-Type", "application/json;charset=UTF-8")
	request.Header.Add("Authorization", token)
	request.Header.Add("batchId", batchid)

	resp, err := client.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("GetList resp.StatusCode = %v", resp.StatusCode)
		return
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &rows)
	if err != nil {
		return
	}

	if rows.Code != 200 {
		err = fmt.Errorf("GetList Rows.Code = %v Msg: %v", rows.Code, rows.Msg)
		return
	} else if rows.Code == 401 && strings.Contains(rows.Msg, "重新登录") {
		err = baseresp.ErrHeuReLogin
		return
	}

	return
}

type allList struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Rows []struct {
			// classID
			JXBID string `json:"JXBID"`
			// className
			KCM string `json:"KCM"`
			// score
			XF        string `json:"XF"`
			SecretVal string `json:"secretVal"`
			// classKind
			XGXKLB string `json:"XGXKLB"`
		} `json:"rows"`
	} `json:"data"`
}
