package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "aMC Suite",
		Description: "Wuthering Waves Companion for Windows",
		Services: []application.Service{
			application.NewService(&AppService{}),
			application.NewService(&GachaService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
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

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
