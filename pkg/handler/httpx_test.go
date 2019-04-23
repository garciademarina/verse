package handler

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	sample "github.com/garciademarina/verse/cmd/sample-data"
	"github.com/garciademarina/verse/pkg/account"
	"github.com/go-chi/jwtauth"
)

var (
	TokenAuthHS256 *jwtauth.JWTAuth
	TokenSecret    = []byte("secret")
	logger         *log.Logger
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, header http.Header, body io.Reader) (int, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return 0, ""
	}

	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v[0])
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("%v\n", req)
		t.Fatal(err)
		return 0, ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return 0, ""
	}
	defer resp.Body.Close()
	return resp.StatusCode, string(respBody)
}

func newJwtToken(secret []byte, claims ...jwt.MapClaims) string {
	token := jwt.New(jwt.GetSigningMethod("HS256"))
	if len(claims) > 0 {
		token.Claims = claims[0]
	}
	tokenStr, err := token.SignedString(secret)
	if err != nil {
		log.Fatal(err)
	}
	return tokenStr
}

func newAuthHeader(claims ...jwt.MapClaims) http.Header {
	h := http.Header{}
	h.Set("Authorization", "BEARER "+newJwtToken(TokenSecret, claims...))
	return h
}

func getSampleAccounts() map[account.Num]*account.Account {
	accounts := make(map[account.Num]*account.Account)
	for k2, v2 := range sample.Accounts {
		accounts[k2] = &account.Account{
			Num:     v2.Num,
			UserID:  v2.UserID,
			Name:    v2.Name,
			OpenAt:  v2.OpenAt,
			Balance: v2.Balance,
		}
	}
	return accounts
}
