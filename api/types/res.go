package types

import baseresp "compete_classes_script/pkg/base_resp"

type CreateOrderResp struct {
	baseresp.BaseResp
}

type LoginResp struct {
	baseresp.BaseResp
	Token string `json:"token"`
}
