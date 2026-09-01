package main

import (
	"fmt"

	"github.com/scan252/aMC/windows/internal/kurobbs"
)

// KurobbsService 账号域对前端的入口。
// 模式说明：mock（默认，演示数据）/ real（真实库街区接口，⏸ 待真人测试短信登录）。
type KurobbsService struct {
	client kurobbs.Client
}

// NewKurobbsService 默认演示模式。
func NewKurobbsService() *KurobbsService {
	return &KurobbsService{client: kurobbs.NewMockClient()}
}

// KurobbsStatus 登录状态概览。
type KurobbsStatus struct {
	Bound   bool   `json:"bound"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Mode    string `json:"mode"`
	UserID  string `json:"userId"`
}

// KurobbsOverview 账号中心页聚合数据。
type KurobbsOverview struct {
	Status    KurobbsStatus       `json:"status"`
	Roles     []kurobbs.Role      `json:"roles"`
	SignIn    *kurobbs.SignInInfo `json:"signIn"`
	Widget    *kurobbs.WidgetData `json:"widget"`
}

func (s *KurobbsService) store() (*kurobbs.Store, error) {
	return kurobbs.DefaultStore()
}

// GetStatus 当前绑定状态。
func (s *KurobbsService) GetStatus() (*KurobbsStatus, error) {
	st, err := s.store()
	if err != nil {
		return nil, err
	}
	acc, err := st.Load()
	if err != nil {
		return nil, err
	}
	status := &KurobbsStatus{Mode: "mock"}
	if acc != nil {
		status.Bound = true
		status.Name = acc.Name
		status.Phone = acc.Phone
		status.UserID = acc.UserID
	}
	return status, nil
}

// SendSms 向手机号发送验证码。
func (s *KurobbsService) SendSms(phone string) error {
	return s.client.SendSmsCode(phone)
}

// Login 使用手机号 + 验证码绑定库街区账号。
func (s *KurobbsService) Login(phone, code string) (*KurobbsStatus, error) {
	acc, err := s.client.Login(phone, code)
	if err != nil {
		return nil, err
	}
	st, err := s.store()
	if err != nil {
		return nil, err
	}
	if err := st.Save(acc); err != nil {
		return nil, err
	}
	return s.GetStatus()
}

// Logout 解绑账号（清除本地凭证）。
func (s *KurobbsService) Logout() error {
	st, err := s.store()
	if err != nil {
		return err
	}
	return st.Clear()
}

// GetOverview 账号中心聚合数据（未绑定时仅返回状态）。
func (s *KurobbsService) GetOverview() (*KurobbsOverview, error) {
	status, err := s.GetStatus()
	if err != nil {
		return nil, err
	}
	ov := &KurobbsOverview{Status: *status}
	if !status.Bound {
		return ov, nil
	}
	st, _ := s.store()
	acc, _ := st.Load()
	if acc == nil {
		return ov, nil
	}

	if roles, err := s.client.Roles(acc); err == nil {
		ov.Roles = roles
	}
	if info, err := s.client.SignInInit(acc); err == nil {
		ov.SignIn = info
	}
	if w, err := s.client.WidgetData(acc); err == nil {
		ov.Widget = w
	}
	return ov, nil
}

// SignInNow 对默认角色执行今日签到。
func (s *KurobbsService) SignInNow() (string, error) {
	st, err := s.store()
	if err != nil {
		return "", err
	}
	acc, err := st.Load()
	if err != nil {
		return "", err
	}
	if acc == nil {
		return "", fmt.Errorf("尚未绑定库街区账号")
	}
	roles, err := s.client.Roles(acc)
	if err != nil {
		return "", err
	}
	if len(roles) == 0 {
		return "", fmt.Errorf("账号未绑定任何鸣潮角色")
	}
	for _, role := range roles {
		if role.IsDefault || len(roles) == 1 {
			if err := s.client.SignIn(acc, role); err != nil {
				return "", err
			}
			return fmt.Sprintf("角色 %s 签到成功", role.RoleName), nil
		}
	}
	if err := s.client.SignIn(acc, roles[0]); err != nil {
		return "", err
	}
	return fmt.Sprintf("角色 %s 签到成功", roles[0].RoleName), nil
}
