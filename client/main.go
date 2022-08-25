package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

type Time struct {
	Time    string `json:"time"`
	Message string `json:"message"`
}

func main() {
	a := app.New()
	w := a.NewWindow("GET REST API")
	w.Resize(fyne.NewSize(300, 200))
	w.SetMaster()

	w.SetContent(widget.NewButton("GET", func() {
		res, _ := http.Get("http://localhost:3000/time")
		defer res.Body.Close()

		body, _ := io.ReadAll(res.Body)
		var time Time
		if err := json.Unmarshal(body, &time); err != nil {
			log.Fatal(err)
		}

		resw := a.NewWindow("RESPONSE")
		resw.Resize(fyne.NewSize(200, 100))
		resw.SetContent(widget.NewLabel(time.Time + " " + time.Message))
		//		resw.SetContent(widget.NewLabel(string(body)))
		resw.Show()
	}))

	w.ShowAndRun()
}
