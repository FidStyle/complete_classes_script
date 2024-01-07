package types

import (
	"compete_classes_script/dao/order"
	baseresp "compete_classes_script/pkg/base_resp"
)

type CreateOrderReq struct {
	// unencrypted password
	Pw                  string
	Account             string
	A                   float64
	B                   float64
	C                   float64
	D                   float64
	E                   float64
	F                   float64
	A0                  float64
	PublicRandom        int      `json:"public_random"`
	ProfessionalRandom  int      `json:"profeesional_random"`
	SpecifyPublic       []string `json:"specify_public"`
	SpecifyProfessional []string `json:"specify_professional"`
	B_n                 float64  `json:"B_n"`
	F_n                 float64  `json:"F_n"`
	A0_n                float64  `json:"A0_n"`
	Token               string   `json:"token"`
}

type GetOrderByCreaterReq struct {
	Creater string
	Finish  bool
	Limit   int
	Offset  int
}

type GetOrderByCreaterResp struct {
	baseresp.BaseResp
	Orders []*order.Order `json:"orders"`
}

type LoginReq struct {
	Account string
	Pw      string
}

type LoginResp struct {
	baseresp.BaseResp
	Token string `json:"token"`
}
