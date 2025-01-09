package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rickj1ang/RRS/internal/data"
	"github.com/rickj1ang/RRS/internal/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		mu.Lock()

		if _, found := clients[ip]; !found {
			clients[ip] = &client{
				limiter: rate.NewLimiter(4, 8),
			}
		}

		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			app.rateLimitExceededResponse(w, r)
			return
		}

		mu.Unlock()
		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidCredentialsResponse(w, r)
			return
		}

		token := headerParts[1]

		v := validator.New()

		if data.ValidateTokenPlaintext(v, token); !v.Valid() {
			app.invalidCredentialsResponse(w, r)
			return
		}

		id, err := app.models.Tokens.GetIdByToken(token)
		if err != nil {
			switch {
			case errors.Is(err, redis.Nil):
				app.invalidAuthenticationResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		user, err := app.models.Users.Get("_id", id)
		if err != nil {
			switch {
			case errors.Is(err, mongo.ErrNoDocuments):
				app.invalidAuthenticationResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		r = app.contextSetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireNormalUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user := app.contextGetUser(r)

		if user.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)
			return
		}

		if user.Level == -1 {
			app.userBanedResponse(w, r)
			return
		}

		if user.Level == 0 {
			app.inactiveAccountResponse(w, r)
			return
		}

		id, err := readIDFromReq(r)
		if err != nil {
			app.badRequestResponse(w, r, err)
		}

		if id != primitive.NilObjectID {
			if !slices.Contains(user.Records, id) {
				app.untouchableResponse(w, r)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireGod(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		if user.Level != 3 {
			app.notLordResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
