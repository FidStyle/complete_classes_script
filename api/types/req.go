package types

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
	Token  string `json:"token"`
	Finish bool
	Limit  int
	Offset int
	Info   string
}

type UpdateOrderInfoByIDReq struct {
	ID    int
	Info  bool
	Token string
}

type DeleteOrderByIDReq struct {
	ID    int
	Token string
}

type LoginReq struct {
	Account string `json:"account"`
	Pw      string `json:"pw"`
}
