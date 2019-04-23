package handler

import (
	"io/ioutil"
	"log"
	"net/http/httptest"
	"strings"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/garciademarina/verse/pkg/repository"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
)

func init() {
	TokenAuthHS256 = jwtauth.New("HS256", TokenSecret, nil)
	logger = log.New(ioutil.Discard, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func TestListAll(t *testing.T) {
	repoAccount := repository.NewInmemAccountRepo(getSampleAccounts())
	accountHandler := NewAccountHandler(repoAccount)
	handler := accountHandler.ListAll(logger)

	r := chi.NewRouter()
	r.Get("/", handler)

	ts := httptest.NewServer(r)
	defer ts.Close()

	expected := `[{"num":"62AC2","userID":"01D3XZ89NFJZ9QT2DHVD462AC2","name":"Rainbow account","openAt":"0001-01-01T00:00:00Z","balance":10000},{"num":"D5DP9","userID":"01D3XZ7CN92AKS9HAPSZ4D5DP9","name":"Billy account","openAt":"0001-01-01T00:00:00Z","balance":10000},{"num":"D8KDR","userID":"01D3XZ3ZHCP3KG9VT4FGAD8KDR","name":"Jenny account","openAt":"0001-01-01T00:00:00Z","balance":10000},{"num":"VE9Z2","userID":"01D3XZ8JXHTDA6XY05EVJVE9Z2","name":"Bjorn account","openAt":"0001-01-01T00:00:00Z","balance":10000}]`

	if status, resp := testRequest(t, ts, "GET", "/?key=admin", nil, nil); status != 200 || resp != expected {
		t.Fatalf(resp)
	}

}
func TestAdminGetBalance(t *testing.T) {
	repoAccount := repository.NewInmemAccountRepo(getSampleAccounts())
	accountHandler := NewAccountHandler(repoAccount)
	handler := accountHandler.GetBalanceById(logger)

	r := chi.NewRouter()
	r.Get("/admin/balance/{userID}", handler)

	ts := httptest.NewServer(r)
	defer ts.Close()

	expected := `{"num":"D8KDR","balance":10000}`
	if status, resp := testRequest(t, ts, "GET", "/admin/balance/01D3XZ3ZHCP3KG9VT4FGAD8KDR?key=admin", nil, nil); status != 200 || resp != expected {
		t.Fatalf(resp)
	}
}

func TestBalance(t *testing.T) {

	repoAccount := repository.NewInmemAccountRepo(getSampleAccounts())
	accountHandler := NewAccountHandler(repoAccount)
	handler := accountHandler.GetBalance(logger)

	r := chi.NewRouter()
	r.Use(jwtauth.Verifier(TokenAuthHS256), jwtauth.Authenticator)
	r.Get("/", handler)

	ts := httptest.NewServer(r)
	defer ts.Close()

	expected := `{"num":"D8KDR","balance":10000}`
	if status, resp := testRequest(t, ts, "GET", "/", newAuthHeader(jwt.MapClaims{"user_id": "01D3XZ3ZHCP3KG9VT4FGAD8KDR"}), nil); status != 200 || resp != expected {
		t.Fatalf(resp)
	}
}

func TestAdminTransfer(t *testing.T) {

	repoAccount := repository.NewInmemAccountRepo(getSampleAccounts())
	accountHandler := NewAccountHandler(repoAccount)
	handler := accountHandler.TransferMoneyAdmin(logger)
	handlerBalance := accountHandler.GetBalanceById(logger)

	r := chi.NewRouter()
	r.Post("/admin/transfers", handler)
	r.Get("/admin/balance/{userID}", handlerBalance)

	ts := httptest.NewServer(r)
	defer ts.Close()

	payload := ` { "UserOrigin": "01D3XZ3ZHCP3KG9VT4FGAD8KDR", "UserDestination": "01D3XZ7CN92AKS9HAPSZ4D5DP9", "Amount": 1040 }`
	expected := `{"origin_user":"01D3XZ3ZHCP3KG9VT4FGAD8KDR","destination_user":"01D3XZ7CN92AKS9HAPSZ4D5DP9","amount":1040}`
	if status, resp := testRequest(t, ts, "POST", "/admin/transfers?key=admin", nil, strings.NewReader(payload)); status != 200 || resp != expected {
		t.Fatalf("%v\n", resp)
	}

	expected = `{"num":"D5DP9","balance":11040}`
	if status, resp := testRequest(t, ts, "GET", "/admin/balance/01D3XZ7CN92AKS9HAPSZ4D5DP9?key=admin", nil, nil); status != 200 || resp != expected {
		t.Fatalf("%s\n", resp)
	}

}
func TestAdminTransferBalanceInsufficient(t *testing.T) {

	repoAccount := repository.NewInmemAccountRepo(getSampleAccounts())
	accountHandler := NewAccountHandler(repoAccount)
	handler := accountHandler.TransferMoneyAdmin(logger)

	r := chi.NewRouter()
	r.Post("/admin/transfers", handler)

	ts := httptest.NewServer(r)
	defer ts.Close()

	payload := ` { "UserOrigin": "01D3XZ3ZHCP3KG9VT4FGAD8KDR", "UserDestination": "01D3XZ7CN92AKS9HAPSZ4D5DP9", "Amount": 9990000 }`
	expected := `{"type":"api_error","code":"balance_insufficient","message":""}`
	if status, resp := testRequest(t, ts, "POST", "/admin/transfers?key=admin", nil, strings.NewReader(payload)); status != 400 || resp != expected {
		t.Fatalf("%v\n", resp)
	}
}
func TestTransferNotUserIdInJwt(t *testing.T) {

	repoAccount := repository.NewInmemAccountRepo(getSampleAccounts())
	accountHandler := NewAccountHandler(repoAccount)
	handler := accountHandler.TransferMoney(logger)

	r := chi.NewRouter()
	r.Use(jwtauth.Verifier(TokenAuthHS256), jwtauth.Authenticator)
	r.Post("/transfers", handler)

	ts := httptest.NewServer(r)
	defer ts.Close()

	payload := ` { "UserDestination": "01D3XZ7CN92AKS9HAPSZ4D5DP9", "Amount": 9990000 }`
	expected := `{"type":"api_error","code":"jwt_user_id_not_found","message":""}`
	if status, resp := testRequest(t, ts, "POST", "/transfers", newAuthHeader(jwt.MapClaims{"no_valid_user_id": "01D3XZ3ZHCP3KG9VT4FGAD8KDR"}), strings.NewReader(payload)); status != 400 || resp != expected {
		t.Fatalf("%v\n", resp)
	}
}

func TestTransfer(t *testing.T) {

	repoAccount := repository.NewInmemAccountRepo(getSampleAccounts())
	accountHandler := NewAccountHandler(repoAccount)
	handler := accountHandler.TransferMoney(logger)
	handlerBalance := accountHandler.GetBalance(logger)

	r := chi.NewRouter()
	r.Use(jwtauth.Verifier(TokenAuthHS256), jwtauth.Authenticator)
	r.Post("/transfer", handler)
	r.Get("/b", handlerBalance)

	ts := httptest.NewServer(r)
	defer ts.Close()

	expected2 := `{"num":"D5DP9","balance":10000}`
	if status, resp := testRequest(t, ts, "GET", "/b", newAuthHeader(jwt.MapClaims{"user_id": "01D3XZ7CN92AKS9HAPSZ4D5DP9"}), nil); status != 200 || resp != expected2 {
		t.Fatalf("%s\n", resp)
	}

	payload := ` { "UserDestination": "01D3XZ7CN92AKS9HAPSZ4D5DP9", "Amount": 1040 }`
	expected := `{"origin_user":"01D3XZ3ZHCP3KG9VT4FGAD8KDR","destination_user":"01D3XZ7CN92AKS9HAPSZ4D5DP9","amount":1040}`
	if status, resp := testRequest(t, ts, "POST", "/transfer", newAuthHeader(jwt.MapClaims{"user_id": "01D3XZ3ZHCP3KG9VT4FGAD8KDR"}), strings.NewReader(payload)); status != 200 || resp != expected {
		t.Fatalf("%v\n", resp)
	}

	expected = `{"num":"D5DP9","balance":11040}`
	if status, resp := testRequest(t, ts, "GET", "/b", newAuthHeader(jwt.MapClaims{"user_id": "01D3XZ7CN92AKS9HAPSZ4D5DP9"}), nil); status != 200 || resp != expected {
		t.Fatalf("%s\n", resp)
	}

}
