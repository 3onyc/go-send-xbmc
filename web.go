package sendxbmc

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type ApiResponse struct {
	Error *string `json:"error"`
}

func NewApiResponse(err error) *ApiResponse {
	if err != nil {
		errStr := err.Error()
		return &ApiResponse{
			Error: &errStr,
		}
	} else {
		return &ApiResponse{
			Error: nil,
		}
	}

}

type Webserver struct {
	api *XbmcApi
}

func NewWebserver(a *XbmcApi) *Webserver {
	return &Webserver{
		api: a,
	}
}

func (s Webserver) handlePost(r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	us, ok := r.Form["url"]
	if !ok {
		return errors.New("No URL")
	}

	u := strings.TrimSpace(us[0])
	if u == "" {
		return errors.New("No URL")
	}

	url, err := Resolve(u)
	if err != nil {
		return err
	}

	return s.api.Play(url)
}

func (s Webserver) Root(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	showErr := ""

	if r.Method == "POST" {
		if err := s.handlePost(r); err != nil {
			showErr = err.Error()
		}
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

func (s Webserver) Api(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not Allowed (Allowed: POST)", 405)
		return
	}

	resp := NewApiResponse(s.handlePost(r))
	enc, err := json.Marshal(resp)

	if err != nil {
		http.Error(w, "Something went wrong", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(enc)
}

func (s Webserver) Listen(listen string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.Root)
	mux.HandleFunc("/api", s.Api)

	log.Fatal(http.ListenAndServe(listen, mux))
}
