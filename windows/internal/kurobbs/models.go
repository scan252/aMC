// Package kurobbs 是库街区（KuroBBS）账号域的基础底座：
// 短信验证码登录、每日签到、小组件数据（体力/资源）、角色绑定列表的接口封装，
// 以及本地凭证存储。接口路径基于社区整理的 Kuro-API-Collection，
// 真实联调标记为 ⏸ 待真人测试（需要真实手机号接收验证码）。
package kurobbs

import "time"

// GameID 鸣潮在库街区体系中的 gameId。
const GameID = "3"

// Account 登录后的账号信息。
type Account struct {
	Token    string    `json:"token"`
	UserID   string    `json:"userId"`
	Name     string    `json:"name"`
	Phone    string    `json:"phone"`
	DevCode  string    `json:"devCode"`
	Did      string    `json:"did"`
	LoggedIn time.Time `json:"loggedIn"`
}

// Role 绑定的游戏角色。
type Role struct {
	RoleId    string `json:"roleId"`
	ServerId  string `json:"serverId"`
	RoleName  string `json:"roleName"`
	Level     int    `json:"level"`
	AreaName  string `json:"areaName"`
	IsDefault bool   `json:"isDefault"`
}

// SignInInfo 每日签到状态。
type SignInInfo struct {
	HadSignIn    bool   `json:"hadSignIn"`
	MonthStart   string `json:"monthStart"`
	TotalSignIn  int    `json:"totalSignIn"`
	TodayReward  string `json:"todayReward"`
	SelectedDays []bool `json:"selectedDays"`
}

// WidgetData 小组件数据（体力、宝箱等）。
type WidgetData struct {
	Energy struct {
		Cur  int `json:"cur"`
		Max  int `json:"max"`
		Full time.Time
	} `json:"energy"`
	Chest        int    `json:"chest"`
	Combat       string `json:"combat"`
	Unlock       string `json:"unlock"`
	BattleExpire string `json:"battleExpire"`
}

// Client 库街区接口抽象（真实 HTTP 客户端与演示实现共用）。
type Client interface {
	// SendSmsCode 向手机号发送验证码。
	SendSmsCode(phone string) error
	// Login 使用手机号 + 短信验证码登录。
	Login(phone, code string) (*Account, error)
	// Mine 获取个人信息（校验 token 有效性）。
	Mine(account *Account) (*UserInfo, error)
	// Roles 获取绑定的鸣潮角色列表。
	Roles(account *Account) ([]Role, error)
	// SignInInit 查询签到面板状态。
	SignInInit(account *Account) (*SignInInfo, error)
	// SignIn 执行每日签到。
	SignIn(account *Account, role Role) error
	// WidgetData 拉取小组件（体力等）数据。
	WidgetData(account *Account) (*WidgetData, error)
	// ForumList 拉取官方资讯帖子列表。
	ForumList(account *Account, limit int) ([]ForumPost, error)
	// RedeemCodes 获取兑换码列表。
	RedeemCodes() ([]RedeemCode, error)
}

// UserInfo 个人信息。
type UserInfo struct {
	Name  string `json:"name"`
	Level int    `json:"level"`
}

// ForumPost 资讯帖子。
type ForumPost struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Time    string `json:"time"`
	URL     string `json:"url"`
}

// RedeemCode 兑换码条目（种子 + 社区渠道，结构先行）。
type RedeemCode struct {
	Code    string `json:"code"`
	Reward  string `json:"reward"`
	Source  string `json:"source"`
	Date    string `json:"date"`
	Expired bool   `json:"expired"`
}
