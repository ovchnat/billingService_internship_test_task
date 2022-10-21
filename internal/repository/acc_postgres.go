package repository

import (
	"billingService/internal/entity"
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

//goland:noinspection ALL
const accByIDQuery = ` SELECT * 
	FROM accounts
	WHERE user_id = $1 LIMIT 1
`

//goland:noinspection ALL
const addFundsQuery = ` UPDATE accounts
	SET curr_amount = curr_amount + $2,
	    last_updated = current_timestamp
	WHERE user_id = $1
	RETURNING user_id, curr_amount, pending_amount
`

//goland:noinspection ALL
const withdrawFundsQuery = ` UPDATE accounts
	SET curr_amount = curr_amount - $2,
	    last_updated = current_timestamp
	WHERE user_id = $1
	RETURNING user_id, curr_amount, pending_amount
`

//goland:noinspection ALL
const createAcc = ` INSERT INTO 
	accounts(user_id, curr_amount, pending_amount, last_updated)
	VALUES ($1, $2, 0, current_timestamp)
	RETURNING user_id, curr_amount, pending_amount
`

type AccPostgres struct {
	db *sqlx.DB
}

func NewAccPostgres(db *sqlx.DB) *AccPostgres {
	return &AccPostgres{db: db}
}

func (r *AccPostgres) GetBalance(userid entity.GetBalanceRequest) (entity.GetBalanceResponse, error) {
	var balanceRes entity.GetBalanceResponse
	query := fmt.Sprintf(
		"SELECT ac.curr_amount, ac.pending_amount FROM accounts ac " +
			"WHERE user_id = $1")
	row := r.db.QueryRow(query, userid.UserId)
	// logrus.Print("Queried with", query, userid.UserId)
	if err := row.Scan(
		&balanceRes.Balance,
		&balanceRes.Pending,
	); err != nil {
		return entity.GetBalanceResponse{}, err
	}
	return entity.GetBalanceResponse{Balance: balanceRes.Balance, Pending: balanceRes.Pending}, nil
}

func (r *AccPostgres) DepositMoney(depositReq entity.UpdateBalanceRequest) (entity.UpdateBalanceResponse, error) {
	var depositRes entity.UpdateBalanceResponse

	accExists, err := r.db.Query(accByIDQuery, depositReq.UserId)
	if err != nil {
		logrus.Fatalf(err.Error())
	}
	if accExists.Next() {
		rows := r.db.QueryRow(addFundsQuery, depositReq.UserId, depositReq.Sum)
		if err := rows.Scan(
			&depositRes.UserId,
			&depositRes.Balance,
			&depositRes.Pending,
		); err != nil {
			return depositRes, err
		}
		logrus.Print("Found acc ", depositReq.UserId, "in database, added ", depositReq.Sum, " funds")
	} else {
		rows := r.db.QueryRow(createAcc, depositReq.UserId, depositReq.Sum)
		if err := rows.Scan(
			&depositRes.UserId,
			&depositRes.Balance,
			&depositRes.Pending,
		); err != nil {
			return depositRes, err
		}
		logrus.Print("Created new acc", depositReq.UserId, "in database, added ", depositReq.Sum, " funds")
	}
	return entity.UpdateBalanceResponse{
		UserId:  depositRes.UserId,
		Balance: depositRes.Balance,
		Pending: depositRes.Pending,
	}, nil
}

func (r *AccPostgres) WithdrawMoney(depositReq entity.UpdateBalanceRequest) (entity.UpdateBalanceResponse, error) {
	var depositRes entity.UpdateBalanceResponse

	accExists, err := r.db.Query(accByIDQuery, depositReq.UserId)
	if err != nil {
		logrus.Fatalf(err.Error())
	}

	if accExists.Next() {
		rows := r.db.QueryRow(withdrawFundsQuery, depositReq.UserId, depositReq.Sum)
		if err := rows.Scan(
			&depositRes.UserId,
			&depositRes.Balance,
			&depositRes.Pending,
		); err != nil {
			return depositRes, err
		}
		logrus.Print("Found acc ", depositReq.UserId, "in database, withdrew ", depositReq.Sum, " funds")
	} else {
		logrus.Print("No such account in database: create one with new deposit. ")
	}
	return entity.UpdateBalanceResponse{
		UserId:  depositRes.UserId,
		Balance: depositRes.Balance,
		Pending: depositRes.Pending,
	}, nil
}

func (r *AccPostgres) Transfer(transferReq entity.TransferRequest) (entity.TransferResponse, error) {
	var transferRes entity.TransferResponse

	return transferRes, nil
}

func (r *AccPostgres) ReserveServiceFee(reserveSerFeeReq entity.ReserveServiceFeeRequest, ctx context.Context) (entity.ReserveServiceFeeResponse, error) {
	var reserveRes entity.ReserveServiceFeeResponse

	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		logrus.Errorf("")
	}

	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	var exists bool
	if err = tx.QueryRowContext(ctx, accByIDQuery, reserveSerFeeReq.UserId).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			logrus.Errorf("No account with that user id. Add a new one by depositing money.")
		}
	}

	_, err = tx.ExecContext(ctx, "")
	if err != nil {
		logrus.Fatalf("")
	}
	return reserveRes, nil
}
