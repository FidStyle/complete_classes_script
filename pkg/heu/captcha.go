package heu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	captchaURL string = "http://api.jfbym.com/api/YmServer/customApi"
)

func GetCaptcha(token string) (uuid string, res string, err error) {
	client := &http.Client{}
	url := "http://jwxk.hrbeu.edu.cn/xsxk/auth/captcha"
	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return
	}

	resp, err := client.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("GetCapthca statusCode = %v", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var data struct {
		Data struct {
			Captcha string `json:"captcha"`
			Uuid    string `json:"uuid"`
		} `json:"data"`
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return
	}

	i := strings.Index(data.Data.Captcha, ",")
	base := data.Data.Captcha[i+1:]

	res, err = getResOfCaptcha(base, token, url)
	if err != nil {
		return
	}

	return data.Data.Uuid, res, nil
}

func getResOfCaptcha(base, token, url string) (string, error) {
	config := map[string]interface{}{}
	config["image"] = base
	config["type"] = "10110"
	config["token"] = token
	configData, _ := json.Marshal(config)
	body := bytes.NewBuffer([]byte(configData))
	resp, err := http.Post(captchaURL, "application/json;charset=utf-8", body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var res struct {
		Data struct {
			Data string `json:"data"`
		} `json:"data"`
	}

	err = json.Unmarshal(data, &res)
	if err != nil {
		return "", err
	}

	return res.Data.Data, nil
}
