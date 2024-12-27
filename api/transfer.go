package api

import (
	"fmt"
	"net/http"

	db "github.com/BariqDev/ias-bank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type createTransferRequest struct {
	FromAccountId int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountId   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	isFromCurrencyMatched := server.validateAccount(ctx, req.FromAccountId, req.Currency)

	if !isFromCurrencyMatched {
		return
	}
	isToCurrencyMatched := server.validateAccount(ctx, req.ToAccountId, req.Currency)
	if !isToCurrencyMatched {
		return
	}

	args := db.TransferTxParams{
		FromAccountID: pgtype.Int8{Int64: req.FromAccountId, Valid: true},
		ToAccountID:   pgtype.Int8{Int64: req.ToAccountId, Valid: true},
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, args)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validateAccount(ctx *gin.Context, accountID int64, currency string) bool {

	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {

		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return false

	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch %s vs %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return false
	}
	return true

}
