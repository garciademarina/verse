package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	repository "github.com/garciademarina/verse/pkg/repository"
)

// Account ...
type Account struct {
	repo repository.AccountRepo
}

// NewAccountHandler explain ...
func NewAccountHandler(repo repository.AccountRepo) Account {
	return Account{
		repo: repo,
	}
}

// BalanceResponse represents api error messages
type BalanceResponse struct {
	Num     string  `json:"num"`
	Balance float64 `json:"balance"`
}

// GetBalance handles GET /balance requests.
func (u *Account) GetBalance(logger *log.Logger) http.HandlerFunc {
	logger.Printf("[handle] AccountRepo.GetBalance\n")
	return func(w http.ResponseWriter, r *http.Request) {

		userID, err := GetJwtValue(r, "user_id")
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, APIError{Type: "api_error", Code: "jwt_user_id_not_found"})
			return
		}

		account, _ := u.repo.FindByUserID(r.Context(), userID)

		respondwithJSON(w, http.StatusOK, &BalanceResponse{
			Num:     account.Num,
			Balance: account.Balance,
		})
	}
}

// TransferMoneyPost ...
type TransferMoneyPost struct {
	UserDestination string
	Amount          float64
}

// TransferMoneyResponse ...
type TransferMoneyResponse struct {
	Success         bool    `json:"success"`
	OriginUser      string  `json:"origin_user"`
	DestinationUser string  `json:"destination_user"`
	Amount          float64 `json:"amount"`
}

// TransferMoney handles POST /transfers requests.
func (u *Account) TransferMoney(logger *log.Logger) http.HandlerFunc {
	logger.Printf("[handle] AccountRepo.TransferMoney\n")
	return func(w http.ResponseWriter, r *http.Request) {

		userOrigin, err := GetJwtValue(r, "user_id")
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, APIError{Type: "api_error", Code: "jwt_user_id_not_found"})
			return
		}
		fmt.Printf("--> %s\n", userOrigin)

		var transfer TransferMoneyPost
		err = decodeTransferMoney(r, &transfer)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, APIError{Type: "api_error", Code: "decode_json_body_failed"})
			return
		}

		err = u.repo.TransferMoney(r.Context(), userOrigin, transfer.UserDestination, transfer.Amount)

		if err != nil {
			switch err {
			case repository.ErrOriginAccountNotFound:
				RespondWithError(w, http.StatusBadRequest, APIError{Type: "api_error", Code: "account_not_found"})
				return
			case repository.ErrBalanceInsufficient:
				RespondWithError(w, http.StatusBadRequest, APIError{Type: "api_error", Code: "balance_insufficient"})
				return
			case repository.ErrDestinationAccountNotFound:
				RespondWithError(w, http.StatusBadRequest, APIError{Type: "api_error", Code: "destination_account_not_found"})
				return
			default:
				RespondWithError(w, http.StatusBadRequest, APIError{Type: "api_error", Code: ""})
			}
			return
		}

		respondwithJSON(w, http.StatusOK, &TransferMoneyResponse{
			Success:         true,
			OriginUser:      userOrigin,
			DestinationUser: transfer.UserDestination,
			Amount:          transfer.Amount,
		})
	}
}

func decodeTransferMoney(r *http.Request, transfer *TransferMoneyPost) error {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(transfer)
	if err != nil {
		return errors.New("Cannot decode TransformMoneyPost from body")
	}

	if transfer.UserDestination == "" {
		return errors.New("UserDestination field not found")
	}
	return nil
}
