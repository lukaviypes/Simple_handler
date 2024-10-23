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

var Logger *logger.Logmid

type Request struct {
	Nums []int `json:"Nums"`
}
type Respons struct {
	Ans int `json:"Res"`
}

func LogMIddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		Logger.Err = nil
		Logger.Status = http.StatusOK
		next.ServeHTTP(w, r)
		if Logger.Err != nil {
			Logger.Logger.LogAttrs(
				r.Context(),
				slog.LevelError,
				"Something malicious happened",
				slog.Any("URL", r.URL),
				slog.Int("duration", int(time.Since(start).Microseconds())),
				slog.Int("status", Logger.Status),
			)
			return

		}
		Logger.Logger.LogAttrs(
			r.Context(),
			slog.LevelInfo,
			"Succes",
			slog.Any("URL", r.URL),
			slog.Int("duration", int(time.Since(start).Microseconds())),
			slog.Int("status", Logger.Status),
		)

	})

}

func tokenAuthMIddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if Logger.Err != nil {
				Logger.Status = http.StatusUnauthorized
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)

		}()
		authHeader := r.Header.Get("Authorization")
		token := strings.TrimPrefix(authHeader, "Bearer ")

		if authHeader == "" {
			Logger.Err = fmt.Errorf("")

			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			Logger.Err = fmt.Errorf("")

			return
		}

		if token != os.Getenv("AUTH_BEARER_TOKEN") {
			Logger.Err = fmt.Errorf("")

			return
		}

	})
}

func formhandler(w http.ResponseWriter, r *http.Request) {

	req := Request{}
	resp := Respons{}
	var byteresp []byte

	defer func() {

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(Logger.Status)
		w.Write(byteresp)
	}()

	if r.Method != http.MethodGet {
		Logger.Status = http.StatusMethodNotAllowed
		Logger.Err = fmt.Errorf("")
		return
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		Logger.Status = http.StatusBadRequest

		return
	}
	err = json.Unmarshal(data, &req)
	if err != nil {
		Logger.Status = http.StatusBadRequest
		Logger.Err = fmt.Errorf("")
		return
	}
	for _, num := range req.Nums {
		resp.Ans += num
	}

	byteresp, err = json.Marshal(resp)
	if err != nil {
		Logger.Status = http.StatusInternalServerError
		Logger.Err = fmt.Errorf("")
		return
	}

}

func main() {
	Logger = logger.New()
	mux := http.NewServeMux()

	mux.Handle("/form", http.HandlerFunc(formhandler))

	fmt.Println("Server is listening")
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	http.ListenAndServe(port, LogMIddleware(tokenAuthMIddleware(mux)))

}
