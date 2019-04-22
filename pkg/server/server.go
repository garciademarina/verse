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
	"github.com/garciademarina/verse/pkg/handler"
	"github.com/garciademarina/verse/pkg/repository"
)

// Server serves http requests.
type Server struct {
	Config
	Logger *log.Logger
	Router chi.Router
}

// Config ...
type Config struct {
	Port              int
	Env               string
	BasicAuthUser     string
	BasicAuthPassword string
}

// TokenAuth contains auth token
var TokenAuth *jwtauth.JWTAuth

// NewConfig creates new config.
func NewConfig(
	port int,
	env string,
	basicAuthUser string,
	basicAuthPassword string,
) Config {
	return Config{
		Port:              port,
		Env:               env,
		BasicAuthUser:     basicAuthUser,
		BasicAuthPassword: basicAuthPassword,
	}
}

// NewServer creates new Server.
func NewServer(config Config, logger *log.Logger) *Server {
	r := chi.NewRouter()
	TokenAuth = jwtauth.New("HS256", []byte("secret"), nil)
	registerRoutes(r, config, logger)

	return &Server{
		Config: config,
		Logger: logger,
		Router: r,
	}
}

func registerRoutes(r chi.Router, config Config, logger *log.Logger) {

	r.Use(middleware.RequestID)

	// Protected routes
	r.Group(func(r chi.Router) {
		users := sample.Users
		repoUsers := repository.NewInmemUserRepo(users)
		// userHandler := handler.NewUserHandler(repoUsers)

		accounts := sample.Accounts
		repoAccounts := repository.NewInmemAccountRepo(accounts)
		accountHandler := handler.NewAccountHandler(repoAccounts)

		// Seek, verify and validate JWT tokens
		r.Use(jwtauth.Verifier(TokenAuth))

		// Handle valid / invalid tokens.
		r.Use(authenticator(repoUsers))

		r.Get("/balance", accountHandler.GetBalance(logger))
		r.Post("/transfers", accountHandler.TransferMoney(logger))

	})
}

func authenticator(userRepo repository.UserRepo) func(next http.Handler) http.Handler {
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

// A completely separate router for posts routes
// r.Get("/", pHandler.ListAll)
// r.Get("/{id:[0-9]+}", pHandler.FindById)
// r.Post("/", pHandler.Create)
// r.Put("/{id:[0-9]+}", pHandler.Update)
// r.Delete("/{id:[0-9]+}", pHandler.Delete)

// logger.Printf("Register route /user\n")
// r.Get("/user", userHandler.FindById(logger))

// logger.Printf("Register route /balance\n")
// r.Get("/balance", accountHandler.FindByUserID(logger))

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
