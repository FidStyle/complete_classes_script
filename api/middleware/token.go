package middleware

import (
	"bytes"
	"compete_classes_script/dao/user"
	baseresp "compete_classes_script/pkg/base_resp"
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func TokenVerify(rtx *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		bts, err := c.GetRawData()
		if err != nil {
			c.JSON(200, baseresp.ErrorResp(err))
			c.Abort()
			return
		}
		var req map[string]interface{}
		if err := json.Unmarshal(bts, &req); err != nil {
			c.JSON(200, baseresp.ErrorResp(err))
			c.Abort()
			return
		}

		token, ok := req["token"]
		if !ok {
			c.JSON(200, baseresp.ErrorResp(baseresp.ErrMissToken))
			c.Abort()
			return
		}
		tk, ok := token.(string)
		if !ok {
			c.JSON(200, baseresp.ErrorResp(baseresp.ErrInvalidArgument))
			c.Abort()
			return
		}

		_, err = user.GetAccountByToken(rtx, tk)
		if err == redis.Nil {
			c.JSON(200, baseresp.ErrorResp(baseresp.ErrExpiredToken))
			c.Abort()
			return
		}
		if err != nil {
			c.JSON(200, baseresp.ErrorResp(err))
			c.Abort()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(bts))
		c.Next()
	}
}
