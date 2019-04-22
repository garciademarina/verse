package handler

import (
	"log"
	"net/http/httptest"
	"os"
	"strings"
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

func TestBalance(t *testing.T) {
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	users := sample.Users
	_ = repository.NewInmemUserRepo(users)

	accounts := sample.Accounts
	repoAccount := repository.NewInmemAccountRepo(accounts)
	accountHandler := NewAccountHandler(repoAccount)
	handler := accountHandler.GetBalance(logger)

	r := chi.NewRouter()
	r.Use(jwtauth.Verifier(TokenAuthHS256), jwtauth.Authenticator)
	r.Get("/", handler)

	ts := httptest.NewServer(r)
	defer ts.Close()

	expected := `{"num":"D8KDR","balance":100}`
	if status, resp := testRequest(t, ts, "GET", "/", newAuthHeader(jwt.MapClaims{"user_id": "01D3XZ3ZHCP3KG9VT4FGAD8KDR"}), nil); status != 200 || resp != expected {
		t.Fatalf(resp)
	}
}

func TestTransfer(t *testing.T) {
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	users := sample.Users
	_ = repository.NewInmemUserRepo(users)

	accounts := sample.Accounts
	repoAccount := repository.NewInmemAccountRepo(accounts)
	accountHandler := NewAccountHandler(repoAccount)
	handler := accountHandler.TransferMoney(logger)
	handlerBalance := accountHandler.GetBalance(logger)

	r := chi.NewRouter()
	r.Use(jwtauth.Verifier(TokenAuthHS256), jwtauth.Authenticator)
	r.Get("/", handler)
	r.Get("/b", handlerBalance)

	ts := httptest.NewServer(r)
	defer ts.Close()

	payload := ` { "UserDestination": "01D3XZ7CN92AKS9HAPSZ4D5DP9", "Amount": 10.4 }`
	expected := `{"success":true,"origin_user":"01D3XZ3ZHCP3KG9VT4FGAD8KDR","destination_user":"01D3XZ7CN92AKS9HAPSZ4D5DP9","amount":10.4}`
	if status, resp := testRequest(t, ts, "GET", "/", newAuthHeader(jwt.MapClaims{"user_id": "01D3XZ3ZHCP3KG9VT4FGAD8KDR"}), strings.NewReader(payload)); status != 200 || resp != expected {
		t.Fatalf("%v\n", resp)
	}

	expected = `{"num":"D5DP9","balance":110.4}`
	if status, resp := testRequest(t, ts, "GET", "/b", newAuthHeader(jwt.MapClaims{"user_id": "01D3XZ7CN92AKS9HAPSZ4D5DP9"}), strings.NewReader(payload)); status != 200 || resp != expected {
		t.Fatalf("%s\n", resp)
	}

}
