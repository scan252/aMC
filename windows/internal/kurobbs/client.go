package kurobbs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// APIBase 库街区 APP 端接口基地址。
const APIBase = "https://api.kurobbs.com"

const requestTimeout = 15 * time.Second

// RealClient 真实 HTTP 客户端。⚠️ 接口字段以 Kuro-API-Collection 为准，
// 首次真实联调（⏸）时如遇字段变更仅需调整本文件。
type RealClient struct {
	HTTP *http.Client
}

func NewRealClient() *RealClient {
	return &RealClient{HTTP: &http.Client{Timeout: requestTimeout}}
}

// post 统一 POST：表单编码 + 常规 APP 请求头。
func (c *RealClient) post(path string, form map[string]string, account *Account) ([]byte, error) {
	vals := make([]byte, 0, 128)
	buf := bytes.NewBuffer(vals)
	first := true
	for k, v := range form {
		if !first {
			buf.WriteByte('&')
		}
		first = false
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(v)
	}
	req, err := http.NewRequest(http.MethodPost, APIBase+path, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("source", "h5")
	if account != nil {
		req.Header.Set("token", account.Token)
		req.Header.Set("userId", account.UserID)
		req.Header.Set("devCode", account.DevCode)
		req.Header.Set("did", account.Did)
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求 %s 失败: %w", path, err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var base struct {
		Code    int             `json:"code"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(raw, &base); err != nil {
		return nil, fmt.Errorf("响应解析失败: %w", err)
	}
	if base.Code != 200 && base.Message != "success" {
		return nil, fmt.Errorf("库街区接口错误: code=%d, message=%s", base.Code, base.Message)
	}
	return base.Data, nil
}

func (c *RealClient) SendSmsCode(phone string) error {
	_, err := c.post("/user/getSmsCode?captcha=", map[string]string{"mobile": phone}, nil)
	return err
}

func (c *RealClient) Login(phone, code string) (*Account, error) {
	data, err := c.post("/user/sdkLogin", map[string]string{
		"phone":    phone,
		"password": code,
	}, nil)
	if err != nil {
		return nil, err
	}
	var info struct {
		Token    string `json:"token"`
		UserId   string `json:"userId"`
		UserName string `json:"userName"`
		Mobile   string `json:"mobile"`
	}
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, err
	}
	return NewAccount(info.Token, info.UserId, info.UserName, info.Mobile), nil
}

func (c *RealClient) Mine(account *Account) (*UserInfo, error) {
	data, err := c.post("/user/mineV2", map[string]string{}, account)
	if err != nil {
		return nil, err
	}
	var info UserInfo
	// mineV2 的 data.almond？此处取常用字段，真实联调时校正（⏸）
	_ = json.Unmarshal(data, &info)
	return &info, nil
}

func (c *RealClient) Roles(account *Account) ([]Role, error) {
	data, err := c.post("/gamer/role/list", map[string]string{"gameId": GameID}, account)
	if err != nil {
		return nil, err
	}
	var roles []Role
	if err := json.Unmarshal(data, &roles); err != nil {
		return nil, err
	}
	return roles, nil
}

func (c *RealClient) SignInInit(account *Account) (*SignInInfo, error) {
	data, err := c.post("/encourage/signIn/initSignInV2", map[string]string{"gameId": GameID}, account)
	if err != nil {
		return nil, err
	}
	var info SignInInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

func (c *RealClient) SignIn(account *Account, role Role) error {
	_, err := c.post("/encourage/signIn/v2", map[string]string{
		"gameId":   GameID,
		"serverId": role.ServerId,
		"roleId":   role.RoleId,
		"userId":   account.UserID,
	}, account)
	return err
}

func (c *RealClient) WidgetData(account *Account) (*WidgetData, error) {
	data, err := c.post("/gamer/widget/game3/getData", map[string]string{"gameId": GameID}, account)
	if err != nil {
		return nil, err
	}
	var w WidgetData
	if err := json.Unmarshal(data, &w); err != nil {
		return nil, err
	}
	return &w, nil
}
