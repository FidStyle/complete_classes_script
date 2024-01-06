package heu

import (
	"bytes"
	"compete_classes_script/pkg/logger"
	"crypto/aes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

const AesKey = "MWMqg2tPcDkxcm11"

func Login(account, password, token string) (loginToken string, batchID string, err error) {
	uuid, res, err := GetCaptcha(token)
	if err != nil {
		return
	}

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("loginname", account)
	_ = writer.WriteField("password", AesEncryptECB([]byte(password), []byte(AesKey)))
	_ = writer.WriteField("captcha", res)
	_ = writer.WriteField("uuid", uuid)
	err = writer.Close()
	if err != nil {
		return
	}

	client := &http.Client{}
	url := "http://jwxk.hrbeu.edu.cn/xsxk/auth/login"
	request, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("Login StatusCode = %v", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var data struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Token   string `json:"token"`
			Student struct {
				ElectiveBatchList []struct {
					Code string `json:"code"`
				} `json:"electiveBatchList"`
			} `json:"student"`
		} `json:"data"`
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return
	}

	if data.Code != 200 {
		err = fmt.Errorf("Login Data Code != 200 Msg: %v", data.Msg)
		return
	}

	if len(data.Data.Student.ElectiveBatchList) == 0 {
		err = fmt.Errorf("len(data.Data.Student.ElectiveBatchList) == 0")
		return
	}
	return data.Data.Token, data.Data.Student.ElectiveBatchList[0].Code, nil
}

func LoginUntilSuccess(account, password, token string) (loginToken string, batchID string) {
	for {
		var err error
		loginToken, batchID, err = Login(account, password, token)
		if err != nil {
			logger.Errorf("Login Failed Because of %v", err.Error())
		} else {
			return
		}

		time.Sleep(3 * time.Minute)
	}
}

func PKCS7Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func AesEncryptECB(origData []byte, key []byte) string {
	cipher, _ := aes.NewCipher(key)
	plain := PKCS7Padding(origData, aes.BlockSize)
	encrypted := make([]byte, len(plain))
	for bs, be := 0, cipher.BlockSize(); bs <= len(origData); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}

	return base64.StdEncoding.EncodeToString(encrypted)
}
