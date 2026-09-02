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

func (m *MockClient) ForumList(account *Account, limit int) ([]ForumPost, error) {
	posts := []ForumPost{
		{ID: "p1", Title: "「潮汐回响」版本活动预告", Summary: "新版本限时活动即将开启，参与可获得星声与养成材料。", Time: "2026-09-01", URL: "https://www.kurobbs.com/mc"},
		{ID: "p2", Title: "唤取「今汐」概率UP公告", Summary: "角色活动唤取「汐潮逐浪」概率提升开启时间说明。", Time: "2026-08-28", URL: "https://www.kurobbs.com/mc"},
		{ID: "p3", Title: "维护更新说明 2.7.1", Summary: "修复了若干已知问题，优化了部分体验。", Time: "2026-08-25", URL: "https://www.kurobbs.com/mc"},
		{ID: "p4", Title: "深塔/sol3塔层奖励调整说明", Summary: "对逆境深塔奖励内容进行例行轮换。", Time: "2026-08-20", URL: "https://www.kurobbs.com/mc"},
	}
	if limit > 0 && len(posts) > limit {
		posts = posts[:limit]
	}
	return posts, nil
}

func (m *MockClient) RedeemCodes() ([]RedeemCode, error) {
	return []RedeemCode{
		{Code: "WUMINGYOU", Reward: "星声 ×50", Source: "版本直播", Date: "2026-08-13", Expired: false},
		{Code: "MINGCHAO2026", Reward: "贝币 ×10000", Source: "官方社区", Date: "2026-08-01", Expired: true},
	}, nil
}
