package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	kuroSvc := NewKurobbsService()
	app := application.New(application.Options{
		Name:        "aMC Suite",
		Description: "Wuthering Waves Companion for Windows",
		Services: []application.Service{
			application.NewService(&AppService{}),
			application.NewService(&GachaService{}),
			application.NewService(kuroSvc),
			application.NewService(&SettingsService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
	})

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "aMC Suite",
		Width:            1360,
		Height:           860,
		MinWidth:         1080,
		MinHeight:        700,
		BackgroundColour: application.NewRGB(5, 5, 7),
		URL:              "/",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
	})

	kuro := kuroSvc
	SetupTray(app, window, kuro)
	StartScheduler(kuro, func(title, body string) {
		app.Event.Emit("app:notify", map[string]string{"title": title, "body": body})
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
