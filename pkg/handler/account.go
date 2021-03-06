package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/garciademarina/verse/pkg/account"
	"github.com/go-chi/chi"
)

// AccountHandler handler struct for account endpoints
type AccountHandler struct {
	repo account.Repository
}

// NewAccountHandler create a new AccountHandler
func NewAccountHandler(repo account.Repository) AccountHandler {
	return AccountHandler{
		repo: repo,
	}
}

// BalanceResponse represents api balance response
type BalanceResponse struct {
	Num     account.Num `json:"num"`
	Balance int64       `json:"balance"`
}

// ListAll handles GET /admin/accounts requests.
func (h *AccountHandler) ListAll(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("[handle] AccountHandler.ListAll\n")

		accounts, _ := h.repo.ListAll(r.Context())

		respondwithJSON(w, http.StatusOK, accounts)
	}
}

// GetBalanceById handles GET /admin/balance/{userID} requests.
func (h *AccountHandler) GetBalanceById(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("[handle] AccountHandler.GetBalanceById\n")
		if userID := chi.URLParam(r, "userID"); userID != "" {
			account, _ := h.repo.FindByUserID(r.Context(), userID)
			respondwithJSON(w, http.StatusOK, &BalanceResponse{
				Num:     account.Num,
				Balance: account.Balance,
			})
		}

	}
}

// GetBalance handles GET /balance requests.
func (h *AccountHandler) GetBalance(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("[handle] AccountHandler.GetBalance\n")

		userID, err := GetJwtValue(r, "user_id")
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, APIError{Type: "api_error", Code: "jwt_user_id_not_found"})
			return
		}

		account, _ := h.repo.FindByUserID(r.Context(), userID)

		respondwithJSON(w, http.StatusOK, &BalanceResponse{
			Num:     account.Num,
			Balance: account.Balance,
		})
	}
}

// TransferMoneyPost represents json post data
type TransferMoneyPost struct {
	DestinationUser string `json:"destination_user"`
	Amount          int64  `json:"amount"`
}

// TransferMoneyPostAdmin represents json post data with user origin id
type TransferMoneyPostAdmin struct {
	OriginUser      string `json:"origin_user"`
	DestinationUser string `json:"destination_user"`
	Amount          int64  `json:"amount"`
}

// TransferMoneyResponse represents /transfers POST response
type TransferMoneyResponse struct {
	OriginUser      string `json:"origin_user"`
	DestinationUser string `json:"destination_user"`
	Amount          int64  `json:"amount"`
}

// TransferMoneyAdmin handles POST /admin/transfers requests.
func (h *AccountHandler) TransferMoneyAdmin(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("[handle] AccountHandler.TransferMoneyAdmin\n")
		var transfer TransferMoneyPostAdmin
		err := decodeTransferMoneyAdmin(r, &transfer)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, APIError{Type: "api_error", Code: "decode_json_body_failed"})
			return
		}

		err = h.repo.TransferMoney(r.Context(), transfer.OriginUser, transfer.DestinationUser, transfer.Amount)
		if err != nil {
			handlerTransferError(w, err)
			return
		}

		respondwithJSON(w, http.StatusOK, &TransferMoneyResponse{
			OriginUser:      transfer.OriginUser,
			DestinationUser: transfer.DestinationUser,
			Amount:          transfer.Amount,
		})
	}
}

// TransferMoney handles POST /transfers requests.
func (h *AccountHandler) TransferMoney(logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("[handle] AccountHandler.TransferMoney\n")

		userOrigin, err := GetJwtValue(r, "user_id")
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, APIError{Type: "api_error", Code: "jwt_user_id_not_found"})
			return
		}

		var transfer TransferMoneyPost
		err = decodeTransferMoney(r, &transfer)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, APIError{Type: "api_error", Code: "decode_json_body_failed"})
			return
		}

		err = h.repo.TransferMoney(r.Context(), userOrigin, transfer.DestinationUser, transfer.Amount)
		if err != nil {
			handlerTransferError(w, err)
			return
		}

		respondwithJSON(w, http.StatusOK, &TransferMoneyResponse{
			OriginUser:      userOrigin,
			DestinationUser: transfer.DestinationUser,
			Amount:          transfer.Amount,
		})
	}
}

func handlerTransferError(w http.ResponseWriter, err error) {
	if err != nil {
		switch err {
		case account.ErrOriginAccountNotFound:
			RespondWithError(w, http.StatusBadRequest, APIError{Type: "api_error", Code: "account_not_found"})
			return
		case account.ErrBalanceInsufficient:
			RespondWithError(w, http.StatusBadRequest, APIError{Type: "api_error", Code: "balance_insufficient"})
			return
		case account.ErrDestinationAccountNotFound:
			RespondWithError(w, http.StatusBadRequest, APIError{Type: "api_error", Code: "destination_account_not_found"})
			return
		default:
			RespondWithError(w, http.StatusBadRequest, APIError{Type: "api_error", Code: ""})
		}
	}
}
func decodeTransferMoney(r *http.Request, transfer *TransferMoneyPost) error {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(transfer)
	if err != nil {
		return errors.New("Cannot decode TransformMoneyPost from body")
	}

	if transfer.DestinationUser == "" {
		return errors.New("DestinationUser field not found")
	}
	return nil
}

func decodeTransferMoneyAdmin(r *http.Request, transfer *TransferMoneyPostAdmin) error {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(transfer)
	if err != nil {
		return errors.New("Cannot decode TransferMoneyPostAdmin from body")
	}

	if transfer.OriginUser == "" {
		return errors.New("OriginUser field cannot be found")
	}
	if transfer.DestinationUser == "" {
		return errors.New("DestinationUser field cannot be found")
	}
	return nil
}
