package response

type Result struct {
	Code       int  `json:"code"`
	Successful bool `json:"successful"`
}
