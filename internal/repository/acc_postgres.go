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
const accIdExists = ` SELECT EXISTS (SELECT * 
	FROM accounts
	WHERE user_id = $1 LIMIT 1)
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
	RETURNING id, user_id, curr_amount, pending_amount
`

//goland:noinspection ALL
const logServiceOrder = ` INSERT INTO
	service_log(account_id, invoice, service_id, order_id, status, created_at)
	VALUES ((SELECT id FROM accounts WHERE user_id = $1), $2, $3, $4, 'Pending', current_timestamp)
	RETURNING account_id, service_id, invoice, status, created_at
`

//goland:noinspection ALL
const reserveAmount = ` UPDATE accounts
	SET pending_amount = pending_amount + $2,
	    last_updated = current_timestamp
	WHERE user_id = $1
`

//goland:noinspection ALL
const getLastService = ` SELECT account_id, service_id, order_id, invoice, status, created_at
	FROM service_log
	WHERE account_id = (SELECT id FROM accounts WHERE user_id = $1)
	AND order_id = $2 
	AND service_id = $3
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
			&depositRes.AccountId,
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

func (r *AccPostgres) WithdrawMoney(withdrawReq entity.UpdateBalanceRequest) (entity.UpdateBalanceResponse, error) {
	var depositRes entity.UpdateBalanceResponse

	accExists, err := r.db.Query(accByIDQuery, withdrawReq.UserId)
	if err != nil {
		logrus.Fatalf(err.Error())
	}

	if accExists.Next() {
		rows := r.db.QueryRow(withdrawFundsQuery, withdrawReq.UserId, withdrawReq.Sum)
		if err := rows.Scan(
			&depositRes.UserId,
			&depositRes.Balance,
			&depositRes.Pending,
		); err != nil {
			return depositRes, err
		}
		logrus.Print("Found acc ", withdrawReq.UserId, "in database, withdrew ", withdrawReq.Sum, " funds")
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

	fail := func(err error) (entity.ReserveServiceFeeResponse, error) {
		return reserveRes, fmt.Errorf("ReserveServiceFee: %v", err)
	}

	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		logrus.Errorf("")
	}

	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	var exists bool

	if err = tx.QueryRowContext(ctx, accIdExists, reserveSerFeeReq.UserId).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			logrus.Errorf("No account with that user id. Add a new one by depositing money.")
			return fail(err)
		}
		return fail(err)
	}

	_, err = tx.ExecContext(ctx, logServiceOrder, reserveSerFeeReq.UserId, reserveSerFeeReq.Fee, reserveSerFeeReq.ServiceId,
		reserveSerFeeReq.OrderId)

	if err != nil {
		return fail(err)
	}

	_, err = tx.ExecContext(ctx, reserveAmount, reserveSerFeeReq.UserId, reserveSerFeeReq.Fee)

	if err != nil {
		return fail(err)
	}

	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	logServiceOrderRes := r.db.QueryRowContext(ctx, getLastService, reserveSerFeeReq.UserId, reserveSerFeeReq.OrderId, reserveSerFeeReq.ServiceId)
	if err := logServiceOrderRes.Scan(
		&reserveRes.UserId,
		&reserveRes.ServiceId,
		&reserveRes.OrderId,
		&reserveRes.Invoice,
		&reserveRes.Status,
		&reserveRes.CreatedAt,
	); err != nil {
		return reserveRes, err
	}
	return reserveRes, nil
}
