package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"

	sample "github.com/garciademarina/verse/cmd/sample-data"
	account "github.com/garciademarina/verse/pkg/account/inmem"
	"github.com/garciademarina/verse/pkg/handler"
	puser "github.com/garciademarina/verse/pkg/user"
	user "github.com/garciademarina/verse/pkg/user/inmem"
)

// Server serves http requests.
type Server struct {
	Config
	Logger *log.Logger
	Router chi.Router
}

// Config represents server configuration
type Config struct {
	Port int
	Env  string
}

// TokenAuth contains auth token
var TokenAuth *jwtauth.JWTAuth

// AdminKey use for /admin/ routes
var AdminKey string

// NewConfig creates new config.
func NewConfig(
	port int,
	env string,
) Config {
	return Config{
		Port: port,
		Env:  env,
	}
}

// NewServer creates new Server.
func NewServer(config Config, logger *log.Logger) *Server {
	r := chi.NewRouter()
	TokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
	AdminKey = "admin"
	registerRoutes(r, config, logger)

	return &Server{
		Config: config,
		Logger: logger,
		Router: r,
	}
}

func registerRoutes(r chi.Router, config Config, logger *log.Logger) {

	r.Use(middleware.RequestID)

	users := sample.Users
	repoUsers := user.NewInmemoryRepository(users)
	userHandler := handler.NewUserHandler(repoUsers)

	accounts := sample.Accounts
	repoAccounts := account.NewInmemoryRepository(accounts)
	accountHandler := handler.NewAccountHandler(repoAccounts)

	// Protected routes
	r.Group(func(r chi.Router) {
		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(TokenAuth))

		// Handle valid / invalid tokens.
		r.Use(authenticator(repoUsers))

		r.Get("/user", userHandler.FindById(logger))
		r.Get("/balance", accountHandler.GetBalance(logger))
		r.Post("/transfers", accountHandler.TransferMoney(logger))

	})

	// Admin Protected routes
	r.Group(func(r chi.Router) {

		// Handle valid / invalid admin key.
		r.Use(AuthenticatorAdmin(AdminKey))

		r.Get("/admin/accounts", accountHandler.ListAll(logger))
		r.Get("/admin/balance/{userID}", accountHandler.GetBalanceById(logger))
		r.Post("/admin/transfers", accountHandler.TransferMoney(logger))

	})
}

func AuthenticatorAdmin(key string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			queryValues := r.URL.Query()
			if queryValues.Get("key") != AdminKey {
				e := handler.APIError{Type: "authentication_error", Message: "Admin key invalid"}
				handler.RespondWithError(w, http.StatusUnauthorized, e)
				return
			}

			// Token is authenticated, pass it through
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func authenticator(userRepo puser.Repository) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			token, _, err := jwtauth.FromContext(r.Context())

			if err != nil {
				e := handler.APIError{Type: "authentication_error", Message: err.Error()}
				handler.RespondWithError(w, http.StatusUnauthorized, e)
				return
			}

			if token == nil || !token.Valid {
				e := handler.APIError{Type: "authentication_error", Message: "Authentication token invalid"}
				handler.RespondWithError(w, http.StatusUnauthorized, e)
				return
			}

			// check if user exist ...
			userID, _ := handler.GetJwtValue(r, "user_id")
			user, err := userRepo.FindById(r.Context(), string(userID))
			if err != nil {
				http.Error(w, "User does not exist", 404)
				return
			}

			if user.ID != userID {
				http.Error(w, "User does not exist", 404)
				return
			}

			// Token is authenticated, pass it through
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// Run runs the server.
func (s *Server) Run(ctx context.Context) error {
	s.Logger.Printf("Server started. PORT:%d\n", s.Port)

	http.Handle("/", s.Router)

	return http.ListenAndServe(fmt.Sprintf(":%v", s.Config.Port), s.Router)
}

// ServeHTTP serve just one request.
func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.Router.ServeHTTP(w, req)
}
