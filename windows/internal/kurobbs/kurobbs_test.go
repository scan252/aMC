package kurobbs

import (
	"path/filepath"
	"testing"
	"time"
)

func TestStoreRoundTrip(t *testing.T) {
	dir := t.TempDir()
	s := &Store{Path: filepath.Join(dir, "kurobbs.json")}

	acc, err := s.Load()
	if err != nil {
		t.Fatal(err)
	}
	if acc != nil {
		t.Fatal("未登录时应返回 nil")
	}

	want := NewAccount("tok", "u1", "测试", "138****1234")
	if err := s.Save(want); err != nil {
		t.Fatal(err)
	}
	got, err := s.Load()
	if err != nil {
		t.Fatal(err)
	}
	if got.Token != "tok" || got.UserID != "u1" || got.DevCode == "" {
		t.Fatalf("回读不符: %+v", got)
	}
	if got.LoggedIn.After(time.Now()) {
		t.Fatal("LoggedIn 不应在未来")
	}

	if err := s.Clear(); err != nil {
		t.Fatal(err)
	}
	acc, err = s.Load()
	if err != nil || acc != nil {
		t.Fatalf("清除后应返回 nil: %v %v", acc, err)
	}
}

func TestMockLoginFlow(t *testing.T) {
	m := NewMockClient()
	if err := m.SendSmsCode("13800138000"); err != nil {
		t.Fatal(err)
	}
	if err := m.SendSmsCode("123"); err == nil {
		t.Fatal("非法手机号应报错")
	}
	if _, err := m.Login("13800138000", "000000"); err == nil {
		t.Fatal("错误验证码应失败")
	}
	acc, err := m.Login("13800138000", "888888")
	if err != nil {
		t.Fatal(err)
	}
	if acc.Phone != "138****8000" {
		t.Fatalf("手机号应脱敏: %s", acc.Phone)
	}
	roles, err := m.Roles(acc)
	if err != nil || len(roles) != 1 {
		t.Fatalf("演示角色列表异常: %v %v", roles, err)
	}
	if err := m.SignIn(acc, roles[0]); err != nil {
		t.Fatal(err)
	}
	if err := m.SignIn(acc, roles[0]); err == nil {
		t.Fatal("重复签到应报错")
	}
	w, err := m.WidgetData(acc)
	if err != nil || w.Energy.Max != 240 {
		t.Fatalf("小组件数据异常: %+v %v", w, err)
	}
}

func TestMockSignInInit(t *testing.T) {
	m := NewMockClient()
	acc, _ := m.Login("13800138000", "888888")
	info, err := m.SignInInit(acc)
	if err != nil {
		t.Fatal(err)
	}
	if info.HadSignIn {
		t.Fatal("演示模式初始应未签到")
	}
	_ = time.Now()
}
