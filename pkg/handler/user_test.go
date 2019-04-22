package handler

import (
	"log"
	"net/http/httptest"
	"os"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	sample "github.com/garciademarina/verse/cmd/sample-data"
	"github.com/garciademarina/verse/pkg/repository"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
)

func init() {
	TokenAuthHS256 = jwtauth.New("HS256", TokenSecret, nil)
}

func TestFoo(t *testing.T) {

	users := sample.Users
	repoUsers := repository.NewInmemUserRepo(users)
	userHandler := NewUserHandler(repoUsers)
	handler := userHandler.FindById(log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile))

	r := chi.NewRouter()
	r.Use(jwtauth.Verifier(TokenAuthHS256), jwtauth.Authenticator)
	r.Get("/", handler)

	ts := httptest.NewServer(r)
	defer ts.Close()

	if status, resp := testRequest(t, ts, "GET", "/", newAuthHeader(jwt.MapClaims{"user_id": "01D3XZ3ZHCP3KG9VT4FGAD8KDR"}), nil); status != 200 || resp != `{"id":"01D3XZ3ZHCP3KG9VT4FGAD8KDR","name":"Jenny","email":"Jenny@example.com"}` {
		t.Fatalf(resp)
	}
}
