package main

import (
	"fmt"
	"time"

	"github.com/scan252/aMC/windows/internal/settings"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// scheduler 每日签到与体力提醒的后台调度。
type scheduler struct {
	kuro              *KurobbsService
	notify            func(title, body string)
	lastSignInDate    string
	lastWaveNotifyAt  time.Time
}

// StartScheduler 启动后台调度（每 30 秒检查一次）。
func StartScheduler(kuro *KurobbsService, notify func(title, body string)) {
	s := &scheduler{kuro: kuro, notify: notify}
	go func() {
		// 启动后先等 15 秒，避开窗口初始化高峰
		time.Sleep(15 * time.Second)
		for {
			s.tick()
			time.Sleep(30 * time.Second)
		}
	}()
}

func (s *scheduler) tick() {
	st, err := settings.Load()
	if err != nil {
		return
	}
	now := time.Now()
	today := now.Format("2006-01-02")

	// 每日自动签到
	if st.SignInAuto && now.Hour() >= st.SignInHour && s.lastSignInDate != today {
		msg, err := s.kuro.SignInNow()
		if err == nil {
			s.lastSignInDate = today
			s.notify("每日签到完成", msg)
		} else if !isNotBound(err) {
			// 未绑定账号不打扰；其它错误明天再试
			s.lastSignInDate = today
			s.notify("自动签到失败", err.Error())
		}
	}

	// 波片回满提醒（每 6 小时最多提醒一次）
	if st.WaveNotify && now.Sub(s.lastWaveNotifyAt) > 6*time.Hour {
		if ov, err := s.kuro.GetOverview(); err == nil && ov.Widget != nil && ov.Widget.Energy.Max > 0 {
			if ov.Widget.Energy.Cur >= ov.Widget.Energy.Max {
				s.lastWaveNotifyAt = now
				s.notify("结晶波片已回满", fmt.Sprintf("%d / %d，记得消耗体力", ov.Widget.Energy.Cur, ov.Widget.Energy.Max))
			}
		}
	}
}

func isNotBound(err error) bool {
	return err != nil && err.Error() == "尚未绑定库街区账号"
}

// SetupTray 创建系统托盘：左键显示/隐藏窗口，右键菜单。
func SetupTray(app *application.App, window *application.WebviewWindow, kuro *KurobbsService) *application.SystemTray {
	tray := app.SystemTray.New()
	tray.SetIcon(buildTrayIcon())
	tray.SetTooltip("aMC Suite · 鸣潮工具站")

	menu := application.NewMenu()
	menu.Add("显示主窗口").OnClick(func(*application.Context) {
		window.Show()
	})
	menu.Add("立即签到").OnClick(func(*application.Context) {
		go func() {
			msg, err := kuro.SignInNow()
			if err != nil {
				app.Event.Emit("app:notify", map[string]string{"title": "签到失败", "body": err.Error()})
				return
			}
			app.Event.Emit("app:notify", map[string]string{"title": "每日签到", "body": msg})
		}()
	})
	autostartItem := menu.Add("开机自启")
	autostartItem.OnClick(func(*application.Context) {
		next := !GetAutostart()
		_ = SetAutostart(next)
		autostartItem.SetLabel(checkLabel(next))
	})
	autostartItem.SetLabel(checkLabel(GetAutostart()))
	menu.AddSeparator()
	menu.Add("退出").OnClick(func(*application.Context) {
		app.Quit()
	})
	tray.SetMenu(menu)
	return tray
}

func checkLabel(on bool) string {
	if on {
		return "开机自启 ✓"
	}
	return "开机自启"
}
