package sample

import (
	"time"

	"github.com/garciademarina/verse/pkg/models"
)

var Users = map[string]*models.User{
	"01D3XZ3ZHCP3KG9VT4FGAD8KDR": &models.User{
		ID:    "01D3XZ3ZHCP3KG9VT4FGAD8KDR",
		Name:  "Jenny",
		Email: "Jenny@example.com",
	},
	"01D3XZ7CN92AKS9HAPSZ4D5DP9": &models.User{
		ID:    "01D3XZ7CN92AKS9HAPSZ4D5DP9",
		Name:  "Billy",
		Email: "Billy@example.com",
	},
	"01D3XZ89NFJZ9QT2DHVD462AC2": &models.User{
		ID:    "01D3XZ89NFJZ9QT2DHVD462AC2",
		Name:  "Rainbow",
		Email: "Rainbow@example.com",
	},
	"01D3XZ8JXHTDA6XY05EVJVE9Z2": &models.User{
		ID:    "01D3XZ8JXHTDA6XY05EVJVE9Z2",
		Name:  "Bjorn",
		Email: "Bjorn@example.com",
	},
}
var t1 time.Time

func init() {
	t1, _ = time.Parse(
		time.RFC3339,
		"2019-04-23T22:08:41+00:00")
}

var Accounts = map[string]*models.Account{
	"D8KDR": &models.Account{
		Num:     "D8KDR",
		UserID:  "01D3XZ3ZHCP3KG9VT4FGAD8KDR",
		Name:    "Jenny account",
		OpenAt:  t1,
		Balance: 10000,
	},
	"D5DP9": &models.Account{
		Num:     "D5DP9",
		UserID:  "01D3XZ7CN92AKS9HAPSZ4D5DP9",
		Name:    "Billy account",
		OpenAt:  t1,
		Balance: 10000,
	},
	"62AC2": &models.Account{
		Num:     "62AC2",
		UserID:  "01D3XZ89NFJZ9QT2DHVD462AC2",
		Name:    "Rainbow account",
		OpenAt:  t1,
		Balance: 10000,
	},
	"VE9Z2": &models.Account{
		Num:     "VE9Z2",
		UserID:  "01D3XZ8JXHTDA6XY05EVJVE9Z2",
		Name:    "Bjorn account",
		OpenAt:  t1,
		Balance: 10000,
	},
}