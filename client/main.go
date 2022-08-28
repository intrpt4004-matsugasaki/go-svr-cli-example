package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type SignUpRes struct {
	Message string `json:"message"`
}

type AuthRes struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

type TimeRes struct {
	Time    string `json:"time"`
	Message string `json:"message"`
}

type UserAgentRes struct {
	UserAgent string `json:"user-agent"`
	Message   string `json:"message"`
}

type ReverseRes struct {
	Text    string `json:"text"`
	Message string `json:"message"`
}

type PostCreateRes struct {
	Message string `json:"message"`
}

type PostAllRes struct {
	Posts   []PostSubRes `json:"posts"`
	Message string       `json:"message"`
}

type PostSubRes struct {
	Id   string `json:"id"`
	Post string `json:"post"`
}

func main() {
	app := app.NewWithID("go-server-client-sample")
	window := app.NewWindow("Client Program")
	window.Resize(fyne.NewSize(450, 0))
	window.SetMaster()

	window.SetMainMenu(fyne.NewMainMenu(
		fyne.NewMenu("Account", fyne.NewMenuItem("Sign Up", func() {
			signupWindow := app.NewWindow("POST /signup")
			signupWindow.Resize(fyne.NewSize(450, 0))

			idEntry := widget.NewEntry()
			idEntry.SetPlaceHolder("Enter your id ...")
			pwEntry := widget.NewEntry()
			pwEntry.SetPlaceHolder("Enter your password ...")

			signUpForm := &widget.Form{
				Items: []*widget.FormItem{
					{Text: "ID", Widget: idEntry},
					{Text: "PW", Widget: pwEntry},
				},
			}

			statusLabel := widget.NewLabel("")

			signUpContent := container.NewGridWithColumns(
				3,
				statusLabel,
				widget.NewLabel(""),
				widget.NewButton("Sign Up", func() {
					params := url.Values{}
					params.Add("id", idEntry.Text)
					params.Add("pw", pwEntry.Text)
					resp, err := http.PostForm("http://localhost:3000/signup", params)
					if err != nil {
						statusLabel.SetText("connect failed.")
						return
					}
					defer resp.Body.Close()

					body, err := io.ReadAll(resp.Body)
					if err != nil {
						statusLabel.SetText("read failed.")
						return
					}

					var res AuthRes
					if err := json.Unmarshal(body, &res); err != nil {
						statusLabel.SetText("unmarshal failed.")
						return
					}

					if resp.StatusCode == http.StatusOK {
						fyne.CurrentApp().SendNotification(&fyne.Notification{
							Title:   "Client Program",
							Content: "Hello " + idEntry.Text + "!  Your account is created!",
						})
						signupWindow.Hide()
					} else {
						statusLabel.SetText(res.Message)
					}
				}),
			)

			signupWindow.SetContent(container.NewVBox(
				signUpForm,
				signUpContent,
			))

			signupWindow.Show()
		})),
	))

	idEntry := widget.NewEntry()
	idEntry.SetPlaceHolder("Enter your id ...")
	pwEntry := widget.NewEntry()
	pwEntry.SetPlaceHolder("Enter your password ...")

	loginForm := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "ID", Widget: idEntry},
			{Text: "PW", Widget: pwEntry},
		},
	}

	tokenLabel := widget.NewLabel("")

	sepContent := widget.NewSeparator()
	sepContent.Hide()

	bodyArea := widget.NewMultiLineEntry()
	bodyContent := container.NewMax(
		bodyArea,
	)
	bodyContent.Hide()

	reqContent := container.NewGridWithColumns(
		3,
		widget.NewButton("GET /time", func() {
			resp, err := http.Get("http://localhost:3000/time?token=" + tokenLabel.Text)
			if err != nil {
				bodyArea.SetText("connect failed.")
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				bodyArea.SetText("read failed.")
				return
			}

			var res TimeRes
			if err := json.Unmarshal(body, &res); err != nil {
				bodyArea.SetText("unmarshal failed.")
				return
			}

			if resp.StatusCode == http.StatusOK {
				bodyArea.SetText(res.Time)
			} else {
				bodyArea.SetText(res.Message)
			}
		}),

		widget.NewButton("GET /user-agent", func() {
			resp, err := http.Get("http://localhost:3000/user-agent?token=" + tokenLabel.Text)
			if err != nil {
				bodyArea.SetText("connect failed.")
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				bodyArea.SetText("read failed.")
				return
			}

			var res UserAgentRes
			if err := json.Unmarshal(body, &res); err != nil {
				bodyArea.SetText("unmarshal failed.")
				return
			}

			if resp.StatusCode == http.StatusOK {
				bodyArea.SetText(res.UserAgent)
			} else {
				bodyArea.SetText(res.Message)
			}
		}),

		widget.NewButton("POST /reverse", func() {
			params := url.Values{}
			params.Add("text", bodyArea.Text)
			resp, err := http.PostForm("http://localhost:3000/reverse?token="+tokenLabel.Text, params)
			if err != nil {
				bodyArea.SetText("connect failed.")
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				bodyArea.SetText("read failed.")
				return
			}

			var res ReverseRes
			if err := json.Unmarshal(body, &res); err != nil {
				bodyArea.SetText("unmarshal failed.")
				return
			}

			if resp.StatusCode == http.StatusOK {
				bodyArea.SetText(res.Text)
			} else {
				bodyArea.SetText(res.Message)
			}
		}),

		widget.NewButton("POST /post/create", func() {
			params := url.Values{}
			params.Add("post", bodyArea.Text)
			resp, err := http.PostForm("http://localhost:3000/post/create?token="+tokenLabel.Text, params)
			if err != nil {
				bodyArea.SetText("connect failed.")
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				bodyArea.SetText("read failed.")
				return
			}

			var res PostCreateRes
			if err := json.Unmarshal(body, &res); err != nil {
				bodyArea.SetText("unmarshal failed.")
				return
			}

			bodyArea.SetText(res.Message)
		}),

		widget.NewButton("GET /post/all", func() {
			resp, err := http.Get("http://localhost:3000/post/all?token=" + tokenLabel.Text)
			if err != nil {
				bodyArea.SetText("connect failed.")
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				bodyArea.SetText("read failed.")
				return
			}

			var res PostAllRes
			if err := json.Unmarshal(body, &res); err != nil {
				bodyArea.SetText("unmarshal failed.")
				return
			}

			if resp.StatusCode == http.StatusOK {
				bodyArea.SetText("")
				for i := range res.Posts {
					post := res.Posts[len(res.Posts)-1-i]
					bodyArea.SetText(bodyArea.Text + "@" + post.Id + ": " + post.Post + "\n")
				}
			} else {
				bodyArea.SetText(res.Message)
			}
		}),
	)
	reqContent.Hide()

	authContent := container.NewGridWithColumns(
		3,
		tokenLabel,
		widget.NewLabel(""),
		widget.NewButton("Auth", func() {
			params := url.Values{}
			params.Add("id", idEntry.Text)
			params.Add("pw", pwEntry.Text)
			resp, err := http.PostForm("http://localhost:3000/auth", params)
			if err != nil {
				tokenLabel.SetText("connect failed.")
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				tokenLabel.SetText("read failed.")
				return
			}

			var res AuthRes
			if err := json.Unmarshal(body, &res); err != nil {
				tokenLabel.SetText("unmarshal failed.")
				return
			}

			if resp.StatusCode == http.StatusOK {
				tokenLabel.SetText(res.Token)

				sepContent.Show()
				reqContent.Show()
				bodyContent.Show()

				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Title:   "Client Program",
					Content: "Hello " + idEntry.Text + "!  your token is " + tokenLabel.Text + ".",
				})
			} else {
				tokenLabel.SetText(res.Message)

				sepContent.Hide()
				reqContent.Hide()
				bodyContent.Hide()
			}
		}),
	)

	window.SetContent(container.NewVBox(
		loginForm,
		authContent,
		sepContent,
		reqContent,
		bodyContent,
	))

	window.ShowAndRun()
}
