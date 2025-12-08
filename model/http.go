package model

// UserInfoData 私信通用户信息 data 字段
type UserInfoData struct {
	AccountNo        string                   `json:"account_no,omitempty"`
	AccoutNo         string                   `json:"accout_no,omitempty"`
	AvatarUrl        string                   `json:"avatar_url,omitempty"`
	BUserId          string                   `json:"b_user_id,omitempty"`
	CUserId          string                   `json:"c_user_id,omitempty"`
	CsProviderId     string                   `json:"cs_provider_id,omitempty"`
	Email            string                   `json:"email,omitempty"`
	IsPrimaryAcc     bool                     `json:"is_primary_acc,omitempty"`
	Name             string                   `json:"name,omitempty"`
	Permissions      []string                 `json:"permissions,omitempty"`
	PrimaryAccountNo string                   `json:"primary_account_no,omitempty"`
	Roles            []map[string]interface{} `json:"roles,omitempty"`
	SellerImage      string                   `json:"seller_image,omitempty"`
	UserId           string                   `json:"user_id,omitempty"`
}

// UserInfo 私信通用户信息
type UserInfo struct {
	Code    int          `json:"code,omitempty"`
	Data    UserInfoData `json:"data,omitempty"`
	Msg     string       `json:"msg,omitempty"`
	Success bool         `json:"success,omitempty"`
}

// FlowUserInfoDataContent Flow 用户信息 data 字段 - flow_user 字段
type FlowUserInfoDataContent struct {
	AccountNo    string `json:"account_no,omitempty"`
	AvatarUrl    string `json:"avatar_url,omitempty"`
	ContactWay   string `json:"contact_way,omitempty"`
	CsProviderId string `json:"cs_provider_id,omitempty"`
	CurrLoad     int    `json:"curr_load,omitempty"`
	MaxLoad      int    `json:"max_load,omitempty"`
	Name         string `json:"name,omitempty"`
	Status       string `json:"status,omitempty"`
}

// FlowUserInfoData Flow 用户信息 data 字段
type FlowUserInfoData struct {
	FlowUser FlowUserInfoDataContent `json:"flow_user,omitempty"`
	Success  bool                    `json:"success,omitempty"`
}

// FlowUserInfo Flow 用户信息
type FlowUserInfo struct {
	Code    int              `json:"code,omitempty"`
	Data    FlowUserInfoData `json:"data,omitempty"`
	Msg     string           `json:"msg,omitempty"`
	Success bool             `json:"success,omitempty"`
}
