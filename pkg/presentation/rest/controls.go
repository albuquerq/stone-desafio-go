package rest

import (
	"encoding/json"
	goerrors "errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"

	"github.com/albuquerq/stone-desafio-go/pkg/domain/access"
	"github.com/albuquerq/stone-desafio-go/pkg/domain/account"
	"github.com/albuquerq/stone-desafio-go/pkg/domain/errors"
	"github.com/albuquerq/stone-desafio-go/pkg/domain/transfer"
	"github.com/albuquerq/stone-desafio-go/pkg/infraestructure/validator"
	"github.com/albuquerq/stone-desafio-go/pkg/presentation/rest/contextkey"
)

// value type used in token responses.
type tokenResponse struct {
	Token string `json:"token"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) (resp Response) {
	var cred access.Credential

	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		resp.Error = &Error{
			Code:    http.StatusBadRequest,
			Message: "invalid JSON data",
			GoErr:   err,
		}
		return
	}

	err = validator.ValidateAccessCredential(cred)
	if err != nil {
		resp.Error = &Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			GoErr:   err,
		}
		return
	}

	descr, err := h.registry.AccessService().Authenticate(cred)
	if err != nil {
		resp.Error = &Error{
			Code:    http.StatusUnauthorized,
			Message: "access denied",
			GoErr:   err,
		}
		return
	}

	// TODO: implements expiration in exp field.
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"account_id":   descr.AccountID,
			"account_cpf":  descr.CPF,
			"account_name": descr.Name,
		},
	)

	strToken, err := token.SignedString([]byte("my secret"))
	if err != nil {
		resp.Error = &Error{
			Code:    http.StatusInternalServerError,
			Message: "error on make jwt",
			GoErr:   err,
		}
		return
	}

	resp.Value = tokenResponse{Token: strToken}

	return
}

// func (h *Handler) AccountGet(w http.ResponseWriter, r *http.Request) (resp Response) {
// 	panic("not implemented")
// 	return
// }

func (h *Handler) AccountCreate(w http.ResponseWriter, r *http.Request) (resp Response) {

	var acv account.InputValue

	err := json.NewDecoder(r.Body).Decode(&acv)
	if err != nil {
		resp.Error = &Error{
			Code:    http.StatusBadRequest,
			Message: "invalid JSON data",
			GoErr:   err,
		}
		return
	}

	ac := account.Account{
		Name:    acv.Name,
		CPF:     acv.CPF,
		Balance: acv.Balance,
		Secret:  acv.Secret,
	}

	err = h.registry.AccountService().CreateAccount(&ac)
	if err != nil {
		resp.Error = &Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			GoErr:   err,
		}
		return
	}

	// hides secret field using json struct tag omitempty
	ac.Secret = ""
	resp.Value = ac

	w.WriteHeader(http.StatusCreated)

	return
}

func (h *Handler) AccountBalance(w http.ResponseWriter, r *http.Request) (resp Response) {

	accountID, exists := h.getRouteParam(r, "accountID")
	if !exists {
		resp.Error = &Error{
			Code:    http.StatusBadRequest,
			Message: "missing account identifier",
			GoErr:   errors.ErrAccountNotFound,
		}
		return
	}

	balance, err := h.registry.AccountService().AccountBalance(accountID)
	if err != nil {
		resp.Error = &Error{
			Code:    http.StatusInternalServerError,
			Message: "account balance could not be obtained",
			GoErr:   err,
		}

		if goerrors.Is(err, errors.ErrAccountNotFound) {
			resp.Error.Code = http.StatusNotFound
			resp.Error.Message = "account not found"
		}
		return
	}

	resp.Value = balance

	return
}

func (h *Handler) AccountList(w http.ResponseWriter, r *http.Request) (resp Response) {
	var err error

	resp.Value, err = h.registry.AccountService().ListAccounts()
	if err != nil {
		resp.Error = &Error{
			Code:    http.StatusInternalServerError,
			Message: "it was not possible to list the registered accounts",
			GoErr:   err,
		}
		return
	}
	if resp.Value == nil { // returns empty account list.
		resp.Value = []account.Account{}
	}
	return
}

func (h *Handler) AccountTransferList(w http.ResponseWriter, r *http.Request) (resp Response) {
	var err error

	accessAcDescr, exists := r.Context().Value(contextkey.AccountDescription).(access.Description)
	if !exists {
		resp.Error = &Error{
			Code:    http.StatusUnauthorized,
			Message: "missing access token",
			GoErr:   errors.ErrInvalidAccessCredentials,
		}
		return
	}

	resp.Value, err = h.registry.TransferService().ListTransfersFromAccount(accessAcDescr.AccountID)
	if err != nil {
		resp.Error = &Error{
			Code:    http.StatusInternalServerError,
			Message: "account transfers could not be listed",
			GoErr:   err,
		}
		return
	}

	if resp.Value == nil { // returns empty transfers list.
		resp.Value = []transfer.Transfer{}
	}

	return
}

func (h *Handler) TransferCreate(w http.ResponseWriter, r *http.Request) (resp Response) {
	var tv transfer.InputValue

	acOrigin, exists := r.Context().Value(contextkey.AccountDescription).(access.Description)
	if !exists {
		resp.Error = &Error{
			Code:    http.StatusUnauthorized,
			Message: "missing access token",
			GoErr:   errors.ErrTransferNotAllowed,
		}
		return
	}

	err := json.NewDecoder(r.Body).Decode(&tv)
	if err != nil {
		resp.Error = &Error{
			Code:    http.StatusBadRequest,
			Message: "invalid JSON data",
			GoErr:   err,
		}
		return
	}

	tr, err := h.registry.TransferService().Transfer(
		acOrigin.AccountID,
		tv.AccountDestinationID,
		tv.Amount,
	)
	if err != nil {
		resp.Error = &Error{
			Code:    http.StatusInternalServerError,
			Message: "failed to transfer",
			GoErr:   err,
		}

		if goerrors.Is(err, errors.ErrTransferNotAllowed) {
			switch err {

			case errors.ErrTransferAccountOriginNotFound:
				resp.Error.Code = http.StatusNotFound
				resp.Error.Message = "origin acccount not found"

			case errors.ErrTransferAccountDestinationNotFound:
				resp.Error.Code = http.StatusNotFound
				resp.Error.Message = "destination account not found"

			case errors.ErrTransferBetweenSameAccount:
				resp.Error.Code = http.StatusForbidden
				resp.Error.Message = "transfer between the same account not allowed"

			case errors.ErrTransferInsufficientBalance:
				resp.Error.Code = http.StatusForbidden
				resp.Error.Message = "insufficient balance"

			case errors.ErrTransferMissingData:
				resp.Error.Code = http.StatusBadRequest
				resp.Error.Message = "missing data"

			case errors.ErrTransferMissingAmount:
				resp.Error.Code = http.StatusBadRequest
				resp.Error.Message = "missing amount"
			}
		}
		return
	}

	resp.Value = tr

	return
}
