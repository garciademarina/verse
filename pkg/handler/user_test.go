package handler

import (
	"io/ioutil"
	"log"
	"net/http/httptest"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	sample "github.com/garciademarina/verse/cmd/sample-data"
	user "github.com/garciademarina/verse/pkg/user/inmem"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
)

func init() {
	TokenAuthHS256 = jwtauth.New("HS256", TokenSecret, nil)
	logger = log.New(ioutil.Discard, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func TestFindById(t *testing.T) {

	users := sample.Users
	repoUsers := user.NewInmemoryRepository(users)
	userHandler := NewUserHandler(repoUsers)
	handler := userHandler.FindById(logger)

	r := chi.NewRouter()
	r.Use(jwtauth.Verifier(TokenAuthHS256), jwtauth.Authenticator)
	r.Get("/", handler)

	ts := httptest.NewServer(r)
	defer ts.Close()

	if status, resp := testRequest(t, ts, "GET", "/", newAuthHeader(jwt.MapClaims{"user_id": "01D3XZ3ZHCP3KG9VT4FGAD8KDR"}), nil); status != 200 || resp != `{"id":"01D3XZ3ZHCP3KG9VT4FGAD8KDR","name":"Jenny","email":"Jenny@example.com"}` {
		t.Fatalf(resp)
	}
}
