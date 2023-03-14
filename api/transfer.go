package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/guy-ronen/simplebank/db/sqlc"
)

type createTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required, min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required, min=1"`
	Amount        int64  `json:"amount" binding:"required, gt=0"`
	Currency      string `json:"currency" binding:"required, oneof=USD EUR CAD"`
}

func (server *Server) createTransfer(ctx *gin.Context) {

	// createAccountRequest is a struct that contains the request body
	var req createTransferRequest

	// ShouldBindJSON binds the request body into the given struct
	if err := ctx.ShouldBindJSON(&req); err != nil {

		// gin.H is a shortcut for map[string]interface{}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// create a new account in the database
	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	// call the createAccount method in the store
	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if !server.validAccount(ctx, req.FromAccountID, req.Currency) {
		return
	}
	if !server.validAccount(ctx, req.ToAccountID, req.Currency) {
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) bool {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		} else {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return false
		}
	}
	if account.Currency != currency {
		err := fmt.Errorf("account [%d] mismatch: %s, %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}
	return true
}
