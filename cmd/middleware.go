package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"runtime"
	"strings"
	"time"

	"matstuff/handlers"
)

const TIMING = false

func RouteLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		realIP := r.Header.Get("X-Real-IP")
		slog.Info("Request", "Method", r.Method, "Path", r.URL.Path, "Who", realIP)
		fmt.Println(r.Method, r.URL.Path)
		start := time.Now()
		next.ServeHTTP(w, r)
		if TIMING {
			fmt.Println("Timing", r.Method, r.URL.Path, time.Since(start).String())
			slog.Info("Response Time", "Method", r.Method, "Path", r.URL.Path, "Time", time.Since(start))
		}
	})
}

func AuthMiddleware(handler *handlers.Auth) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowed := !strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/api/auth/")
			token := r.Header.Get("Authorization")
			if allowed && token == "" {
				next.ServeHTTP(w, r)
				return
			}
			if len(token) < 32 || !strings.HasPrefix(token, "Bearer ") {
				fmt.Println("no token")
				w.WriteHeader(http.StatusForbidden)
				return
			}
			user, ok, err := handler.Check(token)
			if err != nil {
				slog.Error("Errored checking token", "Error", err)
				w.WriteHeader(http.StatusForbidden)
				return
			}
			if !ok {
				ip := r.Header.Get("X-Real-IP")
				slog.Info("Unauthorized", "User", user, "IP", ip)
				w.WriteHeader(http.StatusForbidden)
				return
			}
			ctx := handlers.AddUserToContext(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				c := 0
				for {
					_, file, line, ok := runtime.Caller(c)
					fmt.Println(strings.Repeat(" ", c*2), file, line)
					c++
					if !ok {
						break
					}
				}
				log.Println(err)
				http.Error(w, "Internal Server Error", 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func Middleware(next http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for _, middleware := range middlewares {
		next = middleware(next)
	}
	return next
}
