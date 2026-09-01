package kurobbs

import (
	"fmt"
	"time"
)

// MockClient 演示数据源：完整驱动 UI 流程，不发起任何网络请求。
// 供界面开发与离线演示使用；真实登录（⏸）切换 RealClient。
type MockClient struct {
	SignedIn bool
	CurWave  int
	MaxWave  int
}

func NewMockClient() *MockClient {
	return &MockClient{SignedIn: false, CurWave: 206, MaxWave: 240}
}

func (m *MockClient) SendSmsCode(phone string) error {
	if len(phone) != 11 {
		return fmt.Errorf("手机号格式不正确（演示模式同样校验）")
	}
	return nil
}

func (m *MockClient) Login(phone, code string) (*Account, error) {
	if code != "888888" {
		return nil, fmt.Errorf("验证码错误（演示模式固定使用 888888）")
	}
	acc := NewAccount("mock-token", "10086", "漂泊者", maskPhone(phone))
	return acc, nil
}

func (m *MockClient) Mine(account *Account) (*UserInfo, error) {
	return &UserInfo{Name: account.Name, Level: 70}, nil
}

func (m *MockClient) Roles(account *Account) ([]Role, error) {
	return []Role{
		{RoleId: "100252731", ServerId: "1", RoleName: "漂泊者", Level: 80, AreaName: "无音区·一区", IsDefault: true},
	}, nil
}

func (m *MockClient) SignInInit(account *Account) (*SignInInfo, error) {
	now := time.Now()
	return &SignInInfo{
		HadSignIn:   m.SignedIn,
		MonthStart:  now.Format("2006-01"),
		TotalSignIn: 18,
		TodayReward: "星声 ×20",
	}, nil
}

func (m *MockClient) SignIn(account *Account, role Role) error {
	if m.SignedIn {
		return fmt.Errorf("今日已签到（演示模式）")
	}
	m.SignedIn = true
	return nil
}

func (m *MockClient) WidgetData(account *Account) (*WidgetData, error) {
	w := &WidgetData{Chest: 12, Combat: "160/160", Unlock: "83%", BattleExpire: ""}
	w.Energy.Cur = m.CurWave
	w.Energy.Max = m.MaxWave
	full := time.Now().Add(time.Duration(m.MaxWave-m.CurWave) * 6 * time.Minute)
	w.Energy.Full = full
	return w, nil
}

func maskPhone(phone string) string {
	if len(phone) != 11 {
		return phone
	}
	return phone[:3] + "****" + phone[7:]
}
