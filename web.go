package sendxbmc

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Webserver struct {
	api *XbmcApi
}

func NewWebserver(a *XbmcApi) *Webserver {
	return &Webserver{
		api: a,
	}
}

func (s Webserver) handleRootPost(r *http.Request) string {
	if err := r.ParseForm(); err != nil {
		return err.Error()
	}

	us, ok := r.Form["url"]
	if !ok {
		return "No URL"
	}

	u := strings.TrimSpace(us[0])
	if u == "" {
		return "No URL"
	}

	url, err := Resolve(u)
	if err != nil {
		return err.Error()
	}

	if err := s.api.Play(url); err != nil {
		return err.Error()
	}

	return ""
}
func (s Webserver) Root(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	showErr := ""

	if r.Method == "POST" {
		showErr = s.handleRootPost(r)
	}

	fmt.Fprintf(w, `
		<html>
			<head>
				<style>
					body, html {
						background-color: #eeeeee;
						font-size: 18px;
					}
					.error {
						text-align: center;
						color: red;
						font-weight: bold;
					}
					.wrapper {
						margin-left: auto;
						margin-right: auto;
						max-width: 30rem;
						background-color: #fdfdfd;
						padding: 2rem;
						box-shadow: 3px 3px 4px #888888;
					}
					.branding {
						text-align: center;
					}
					.play {
						font-size: 4rem;
						width: 100%%;
					}

					input {
						width:100%%;
						font-size: 2rem;
					}

				</style>
			</head>
			<body>
				<div class="wrapper">
					<h2 class="branding">Play XBMC</h2>
					<form method="POST" action="/">
						<div class="error">%s</div>
						<p>
							<input type="text" name="url" id="url" placeholder="URL..." />
						</p>
						<p>
							<button class="play" type="submit">Play</button>
						</p>
					</form>
				</div>
			<body>
		</html>
	`, showErr)
}

func (s Webserver) Listen(listen string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.Root)

	log.Fatal(http.ListenAndServe(listen, mux))
}
