package main

import (
	"encoding/json"
	"fmt"
	logger "handlerserver/CustomLogger"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

var Logger *slog.Logger

type Request struct {
	Nums []int `json:"Nums"`
}
type Respons struct {
	Ans int `json:"Res"`
}

func tokenAuthMIddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			if err != nil {
				Logger.LogAttrs(
					r.Context(),
					slog.LevelError,
					"Authorization failed",
					slog.Any("URL", r.URL),

					slog.Int("status", http.StatusUnauthorized),
				)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			Logger.LogAttrs(
				r.Context(),
				slog.LevelInfo,
				"Authorization passed",
				slog.Any("URL", r.URL),

				slog.Int("status", http.StatusUnauthorized),
			)
			next.ServeHTTP(w, r)

		}()
		authHeader := r.Header.Get("Authorization")
		token := strings.TrimPrefix(authHeader, "Bearer ")

		if authHeader == "" {
			err = http.ErrAbortHandler
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			err = http.ErrAbortHandler
			return
		}

		if token != os.Getenv("AUTH_BEARER_TOKEN") {
			err = http.ErrAbortHandler
			return
		}

	})
}

func formhandler(w http.ResponseWriter, r *http.Request) {

	req := Request{}
	resp := Respons{}
	var byteresp []byte
	statuscode := http.StatusOK

	start := time.Now()
	var err error

	defer func() {
		if err != nil {
			Logger.LogAttrs(
				r.Context(),
				slog.LevelError,
				"Something malicious happened",
				slog.Any("URL", r.URL),
				slog.Int("duration", int(time.Since(start).Microseconds())),
				slog.Int("status", statuscode),
			)

		} else {
			Logger.LogAttrs(
				r.Context(),
				slog.LevelInfo,
				"Succes",
				slog.Any("URL", r.URL),
				slog.Int("duration", int(time.Since(start).Microseconds())),
				slog.Int("status", statuscode),
			)

		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statuscode)
		w.Write(byteresp)
	}()

	if r.Method != http.MethodGet {
		statuscode = http.StatusMethodNotAllowed

		return
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		statuscode = http.StatusBadRequest

		return
	}
	err = json.Unmarshal(data, &req)
	if err != nil {
		statuscode = http.StatusBadRequest

		return
	}
	for _, num := range req.Nums {
		resp.Ans += num
	}

	byteresp, err = json.Marshal(resp)
	if err != nil {
		statuscode = http.StatusInternalServerError
		return
	}

}

func main() {
	Logger = logger.New()
	mux := http.NewServeMux()

	mux.Handle("/form", tokenAuthMIddleware(http.HandlerFunc(formhandler)))

	fmt.Println("Server is listening")
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	http.ListenAndServe(port, mux)

}
