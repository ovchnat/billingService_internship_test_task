package repository

import (
	"billingService/internal/entity"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

//goland:noinspection ALL
const accByIdQuery = `  SELECT id
	FROM accounts
	WHERE user_id = $1 LIMIT 1
`

//goland:noinspection ALL
const getIdBalanceQuery = `  SELECT id, curr_amount
	FROM accounts
	WHERE user_id = $1 LIMIT 1
`

//goland:noinspection ALL

const logTransactionQuery = ` INSERT INTO 
	transactions_log(account_id_from, account_id_to, transaction_sum, status, event_type, created_at, updated_at)
	VALUES ((SELECT id FROM accounts WHERE user_id = $1), 
	        (SELECT id FROM accounts WHERE user_id = $2), 
	        $3, $4, $5, current_timestamp, current_timestamp)
`

//goland:noinspection ALL
const getTransactionQuery = `SELECT account_id_from, account_id_to, transaction_sum, status, event_type, created_at
	FROM transactions_log
	WHERE account_id_from = (SELECT id FROM accounts WHERE user_id = $1)
	AND created_at = (select created_at from transactions_log order by created_at desc limit 1)
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
const decreasePendingAmountQuery = ` UPDATE accounts
	SET pending_amount = pending_amount - $2,
	    last_updated = current_timestamp
	WHERE user_id = $1
	RETURNING user_id, curr_amount, pending_amount
`

//goland:noinspection ALL
const createAccQuery = ` INSERT INTO 
	accounts(user_id, curr_amount, pending_amount, last_updated)
	VALUES ($1, 0, 0, current_timestamp)
	RETURNING id
`

//goland:noinspection ALL
const logServiceOrderQuery = ` INSERT INTO
	service_log(account_id, invoice, service_id, order_id, status, created_at, updated_at)
	VALUES ((SELECT id FROM accounts WHERE user_id = $1), $2, $3, $4, $5, current_timestamp, current_timestamp)
	RETURNING account_id, service_id, invoice, status, created_at
`

//goland:noinspection ALL
const changeServiceStatusQuery = ` UPDATE service_log
	SET status = $5, 
	updated_at = current_timestamp
	WHERE account_id = (SELECT id FROM accounts WHERE user_id = $1) 
	AND order_id = $2 
	AND service_id = $3
	AND invoice = $4
`

//goland:noinspection ALL
const reserveAmountQuery = ` UPDATE accounts
	SET pending_amount = pending_amount + $2,
	    last_updated = current_timestamp
	WHERE user_id = $1
`

//goland:noinspection ALL
const getLastServiceQuery = ` SELECT account_id, service_id, order_id, invoice, status, created_at, updated_at
	FROM service_log
	WHERE account_id = (SELECT id FROM accounts WHERE user_id = $1)
	AND order_id = $2 
	AND service_id = $3
	AND invoice = $4
`

//goland:noinspection ALL
const getLastServiceStatusQuery = ` SELECT status
	FROM service_log
	WHERE account_id = (SELECT id FROM accounts WHERE user_id = $1)
	AND order_id = $2 
	AND service_id = $3
	AND invoice = $4
`

type AccPostgres struct {
	db *sqlx.DB
}

func NewAccPostgres(db *sqlx.DB) *AccPostgres {
	return &AccPostgres{db: db}
}

func (r *AccPostgres) GetBalance(userid entity.GetBalanceRequest, ctx *gin.Context) (entity.GetBalanceResponse, error) {
	var balanceRes entity.GetBalanceResponse

	fail := func(err error) (entity.GetBalanceResponse, error) {
		return balanceRes, fmt.Errorf("GetBalance: %v", err)
	}

	if userid.UserId <= 0 {
		err := errors.New("illegal user ID")
		return fail(err)
	}

	query := fmt.Sprintf(
		"SELECT ac.curr_amount, ac.pending_amount FROM accounts ac " +
			"WHERE user_id = $1")
	row := r.db.QueryRowContext(ctx, query, userid.UserId)

	if err := row.Scan(
		&balanceRes.Balance,
		&balanceRes.Pending,
	); err != nil {
		return entity.GetBalanceResponse{}, err
	}
	return entity.GetBalanceResponse{Balance: balanceRes.Balance, Pending: balanceRes.Pending}, nil
}

func (r *AccPostgres) DepositMoney(depositReq entity.UpdateBalanceRequest, ctx *gin.Context) (entity.UpdateBalanceDepositResponse, error) {
	var depositResponse entity.UpdateBalanceDepositResponse

	fail := func(err error) (entity.UpdateBalanceDepositResponse, error) {
		return depositResponse, fmt.Errorf("DepositMoney: %v", err)
	}

	if depositReq.Sum <= 0 {
		err := errors.New("can't add negative or zero funds")
		return fail(err)
	}

	if depositReq.UserId <= 0 {
		err := errors.New("illegal user ID")
		return fail(err)
	}

	var exists int64

	if err := r.db.QueryRow(accByIdQuery, depositReq.UserId).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			rows := r.db.QueryRow(createAccQuery, depositReq.UserId)
			if err := rows.Scan(
				&depositResponse.AccountId,
			); err != nil {
				return depositResponse, err
			}
			logrus.Print("created new account ", depositReq.UserId, " in database")
		}
	}

	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		logrus.Errorf("")
	}

	defer func() {
		if err := recover(); err != nil {
			rb := tx.Rollback()
			if rb != nil {
				errMsg := errors.New("rollback error")
				_, err := fail(errMsg)
				if err != nil {
					return
				}

			}
		}
	}()

	_, err = tx.ExecContext(ctx, addFundsQuery, depositReq.UserId, depositReq.Sum)

	if err != nil {
		return fail(err)
	}

	_, err = tx.ExecContext(ctx, logTransactionQuery, depositReq.UserId, depositReq.UserId, depositReq.Sum, "Completed", "Deposit")

	if err != nil {
		return fail(err)
	}

	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	var holder int

	rows := r.db.QueryRow(getTransactionQuery, depositReq.UserId)
	if err := rows.Scan(
		&depositResponse.AccountId,
		&holder,
		&depositResponse.Sum,
		&depositResponse.Status,
		&depositResponse.EventType,
		&depositResponse.CreatedAt,
	); err != nil {
		return depositResponse, err
	}

	return entity.UpdateBalanceDepositResponse{
		AccountId: depositResponse.AccountId,
		Sum:       depositResponse.Sum,
		Status:    depositResponse.Status,
		EventType: depositResponse.EventType,
		CreatedAt: depositResponse.CreatedAt,
	}, nil
}

func (r *AccPostgres) WithdrawMoney(withdrawReq entity.UpdateBalanceRequest, ctx *gin.Context) (entity.UpdateBalanceWithdrawResponse, error) {
	var withdrawResponse entity.UpdateBalanceWithdrawResponse

	fail := func(err error) (entity.UpdateBalanceWithdrawResponse, error) {
		return withdrawResponse, fmt.Errorf("WithdrawMoney: %v", err)
	}

	if withdrawReq.Sum < 0 {
		err := errors.New("can't withdraw negative funds")
		return fail(err)
	}

	if withdrawReq.UserId <= 0 {
		err := errors.New("illegal user ID")
		return fail(err)
	}

	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		logrus.Errorf("")
	}

	defer func() {
		if err := recover(); err != nil {
			rb := tx.Rollback()
			if rb != nil {
				errMsg := errors.New("rollback error")
				_, err := fail(errMsg)
				if err != nil {
					return
				}

			}
		}
	}()

	var idBalanceHolder struct {
		Id      int64
		Balance int64
	}

	if err = tx.QueryRowContext(ctx, getIdBalanceQuery, withdrawReq.UserId).Scan(&idBalanceHolder.Id, &idBalanceHolder.Balance); err != nil {
		if err == sql.ErrNoRows {
			logrus.Errorf("no account with that user id: add a new one by depositing money")
			return fail(err)
		}
		return fail(err)
	}
	if idBalanceHolder.Balance < withdrawReq.Sum {
		err = errors.New("not enough funds")
		logrus.Errorf("not enough funds on the account")
		return fail(err)
	}
	_, err = tx.ExecContext(ctx, withdrawFundsQuery, withdrawReq.UserId, withdrawReq.Sum)

	if err != nil {
		return fail(err)
	}

	_, err = tx.ExecContext(ctx, logTransactionQuery, withdrawReq.UserId, withdrawReq.UserId, withdrawReq.Sum, "Completed", "Withdrawal")

	if err != nil {
		return fail(err)
	}

	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	var holder int
	rows := r.db.QueryRow(getTransactionQuery, withdrawReq.UserId)
	if err := rows.Scan(
		&withdrawResponse.AccountId,
		&holder,
		&withdrawResponse.Sum,
		&withdrawResponse.Status,
		&withdrawResponse.EventType,
		&withdrawResponse.CreatedAt,
	); err != nil {
		return withdrawResponse, err
	}
	logrus.Print("found acc ", withdrawReq.UserId, " in database, withdrew ", withdrawReq.Sum, " funds")
	return entity.UpdateBalanceWithdrawResponse{
		AccountId: withdrawResponse.AccountId,
		Sum:       withdrawResponse.Sum,
		Status:    withdrawResponse.Status,
		EventType: withdrawResponse.EventType,
		CreatedAt: withdrawResponse.CreatedAt,
	}, nil
}

func (r *AccPostgres) Transfer(transferReq entity.TransferRequest, ctx *gin.Context) (entity.TransferResponse, error) {
	var transferRes entity.TransferResponse

	fail := func(err error) (entity.TransferResponse, error) {
		return transferRes, fmt.Errorf("TransferMoney: %v", err)
	}

	if transferReq.Sum < 0 {
		err := errors.New("can't transfer negative amount")
		return fail(err)
	}

	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		logrus.Errorf("")
	}

	defer func() {
		if err := recover(); err != nil {
			rb := tx.Rollback()
			if rb != nil {
				errMsg := errors.New("rollback error")
				_, err := fail(errMsg)
				if err != nil {
					return
				}

			}
		}
	}()

	var idBalanceHolder struct {
		Id      int64
		Balance int64
	}

	if err = tx.QueryRowContext(ctx, accByIdQuery, transferReq.ReceiverId).Scan(&idBalanceHolder.Id); err != nil {
		if err == sql.ErrNoRows {
			logrus.Errorf("no account with that receiver id: add a new one by depositing money")
			return fail(err)
		}
		return fail(err)
	}

	if err = tx.QueryRowContext(ctx, getIdBalanceQuery, transferReq.SenderId).Scan(&idBalanceHolder.Id, &idBalanceHolder.Balance); err != nil {
		if err == sql.ErrNoRows {
			logrus.Errorf("no account with that sender id: add a new one by depositing money.")
			return fail(err)
		}
		return fail(err)
	}
	if idBalanceHolder.Balance < transferReq.Sum {
		err = errors.New("not enough funds to transfer")
		return fail(err)
	}

	_, err = tx.ExecContext(ctx, withdrawFundsQuery, transferReq.SenderId, transferReq.Sum)

	if err != nil {
		return fail(err)
	}

	_, err = tx.ExecContext(ctx, addFundsQuery, transferReq.ReceiverId, transferReq.Sum)

	if err != nil {
		return fail(err)
	}

	_, err = tx.ExecContext(ctx, logTransactionQuery, transferReq.SenderId, transferReq.ReceiverId, transferReq.Sum, "Completed", "Withdrawn-Transfer")

	if err != nil {
		return fail(err)
	}

	_, err = tx.ExecContext(ctx, logTransactionQuery, transferReq.ReceiverId, transferReq.SenderId, transferReq.Sum, "Completed", "Received-Transfer")

	if err != nil {
		return fail(err)
	}

	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	rows := r.db.QueryRow(getTransactionQuery, transferReq.SenderId)
	if err := rows.Scan(
		&transferRes.AccountIdFrom,
		&transferRes.AccountIdTo,
		&transferRes.Amount,
		&transferRes.Status,
		&transferRes.EventType,
		&transferRes.Timecode,
	); err != nil {
		return transferRes, err
	}

	return entity.TransferResponse{
		AccountIdTo:   transferRes.AccountIdTo,
		AccountIdFrom: transferRes.AccountIdFrom,
		Amount:        transferRes.Amount,
		Status:        transferRes.Status,
		EventType:     transferRes.EventType,
		Timecode:      transferRes.Timecode,
	}, nil
}

func (r *AccPostgres) ReserveServiceFee(reserveSerFeeReq entity.ReserveServiceFeeRequest, ctx *gin.Context) (entity.ReserveServiceFeeResponse, error) {
	var reserveRes entity.ReserveServiceFeeResponse

	fail := func(err error) (entity.ReserveServiceFeeResponse, error) {
		return reserveRes, fmt.Errorf("ReserveServiceFee: %v", err)
	}

	if reserveSerFeeReq.Fee < 0 {
		err := errors.New("can't reserve negative sum")
		return fail(err)
	}

	if reserveSerFeeReq.ServiceId < 0 {
		err := errors.New("illegal service ID")
		return fail(err)
	}

	if reserveSerFeeReq.OrderId < 0 {
		err := errors.New("illegal order ID")
		return fail(err)
	}

	if reserveSerFeeReq.UserId <= 0 {
		err := errors.New("illegal user ID")
		return fail(err)
	}

	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		logrus.Errorf("")
	}

	// Defer a rollback in case anything fails.
	defer func() {
		if err := recover(); err != nil {
			rb := tx.Rollback()
			if rb != nil {
				errMsg := errors.New("rollback error")
				_, err := fail(errMsg)
				if err != nil {
					return
				}

			}
		}
	}()

	var exists int64

	if err = tx.QueryRowContext(ctx, accByIdQuery, reserveSerFeeReq.UserId).Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			logrus.Errorf("no account with that user id: add a new one by depositing money")
			return fail(err)
		}
		return fail(err)
	}

	_, err = tx.ExecContext(ctx, logServiceOrderQuery, reserveSerFeeReq.UserId, reserveSerFeeReq.Fee, reserveSerFeeReq.ServiceId,
		reserveSerFeeReq.OrderId, "Pending")

	if err != nil {
		return fail(err)
	}

	_, err = tx.ExecContext(ctx, reserveAmountQuery, reserveSerFeeReq.UserId, reserveSerFeeReq.Fee)

	if err != nil {
		return fail(err)
	}

	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	logServiceOrderRes := r.db.QueryRowContext(ctx, getLastServiceQuery, reserveSerFeeReq.UserId,
		reserveSerFeeReq.OrderId, reserveSerFeeReq.ServiceId, reserveSerFeeReq.Fee)
	if err := logServiceOrderRes.Scan(
		&reserveRes.AccountId,
		&reserveRes.ServiceId,
		&reserveRes.OrderId,
		&reserveRes.Invoice,
		&reserveRes.Status,
		&reserveRes.CreatedAt,
		&reserveRes.UpdatedAt,
	); err != nil {
		return reserveRes, err
	}
	return reserveRes, nil
}

func (r *AccPostgres) ApproveServiceFee(approveSerFeeReq entity.StatusServiceFeeRequest, ctx *gin.Context) (entity.StatusServiceFeeResponse, error) {
	var approvalServiceFeeResponse entity.StatusServiceFeeResponse

	fail := func(err error) (entity.StatusServiceFeeResponse, error) {
		return approvalServiceFeeResponse, fmt.Errorf("ApproveServiceFee: %v", err)
	}

	if approveSerFeeReq.Fee < 0 {
		err := errors.New("can't withdraw negative funds")
		return fail(err)
	}

	if approveSerFeeReq.UserId <= 0 {
		err := errors.New("illegal user ID")
		return fail(err)
	}

	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		logrus.Errorf("")
	}

	defer func() {
		if err := recover(); err != nil {
			rb := tx.Rollback()
			if rb != nil {
				errMsg := errors.New("rollback error")
				_, err := fail(errMsg)
				if err != nil {
					return
				}

			}
		}
	}()

	var idBalanceHolder struct {
		Id      int64
		Balance int64
	}

	if err = tx.QueryRowContext(ctx, getIdBalanceQuery, approveSerFeeReq.UserId).Scan(
		&idBalanceHolder.Id,
		&idBalanceHolder.Balance,
	); err != nil {
		if err == sql.ErrNoRows {
			logrus.Errorf("No account with that user id. Add a new one by depositing money.")
			return fail(err)
		}
		return fail(err)
	}

	if idBalanceHolder.Balance < approveSerFeeReq.Fee {
		err = errors.New("not enough funds")
		return fail(err)
	}

	var status string

	if err = tx.QueryRowContext(ctx, getLastServiceStatusQuery, approveSerFeeReq.UserId, approveSerFeeReq.OrderId,
		approveSerFeeReq.ServiceId, approveSerFeeReq.Fee).Scan(&status); err != nil {
		if err == sql.ErrNoRows {
			logrus.Errorf("No account with that user id. Add a new one by depositing money.")
			return fail(err)
		}
		return fail(err)
	} else {
		if status == "Approved" {
			err = errors.New("this fee has already been approved")
			return fail(err)
		}
	}

	_, err = tx.ExecContext(ctx, changeServiceStatusQuery, approveSerFeeReq.UserId, approveSerFeeReq.OrderId,
		approveSerFeeReq.ServiceId, approveSerFeeReq.Fee, "Approved")

	if err != nil {
		return fail(err)
	}

	_, err = tx.ExecContext(ctx, withdrawFundsQuery, approveSerFeeReq.UserId, approveSerFeeReq.Fee)

	if err != nil {
		return fail(err)
	}

	_, err = tx.ExecContext(ctx, decreasePendingAmountQuery, approveSerFeeReq.UserId, approveSerFeeReq.Fee)

	if err != nil {
		return fail(err)
	}

	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	logServiceOrderRes := r.db.QueryRowContext(ctx, getLastServiceQuery, approveSerFeeReq.UserId, approveSerFeeReq.OrderId,
		approveSerFeeReq.ServiceId, approveSerFeeReq.Fee)
	if err := logServiceOrderRes.Scan(
		&approvalServiceFeeResponse.AccountId,
		&approvalServiceFeeResponse.ServiceId,
		&approvalServiceFeeResponse.OrderId,
		&approvalServiceFeeResponse.Invoice,
		&approvalServiceFeeResponse.Status,
		&approvalServiceFeeResponse.CreatedAt,
		&approvalServiceFeeResponse.UpdatedAt,
	); err != nil {
		return approvalServiceFeeResponse, err
	}
	return approvalServiceFeeResponse, nil
}

func (r *AccPostgres) FailedServiceFee(failedServiceFeeReq entity.StatusServiceFeeRequest, ctx *gin.Context) (entity.StatusServiceFeeResponse, error) {
	var failedServiceFee entity.StatusServiceFeeResponse

	fail := func(err error) (entity.StatusServiceFeeResponse, error) {
		return failedServiceFee, fmt.Errorf("FailedServiceFee: %v", err)
	}

	if failedServiceFeeReq.UserId <= 0 {
		err := errors.New("illegal user ID")
		return fail(err)
	}

	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		logrus.Errorf("")
	}

	defer func() {
		if err := recover(); err != nil {
			rb := tx.Rollback()
			if rb != nil {
				errMsg := errors.New("rollback error")
				_, err := fail(errMsg)
				if err != nil {
					return
				}

			}
		}
	}()

	var idBalance struct {
		Id  int64
		Bal int64
	}

	if err = tx.QueryRowContext(ctx, getIdBalanceQuery, failedServiceFeeReq.UserId).Scan(
		&idBalance.Id,
		&idBalance.Bal,
	); err != nil {
		if err == sql.ErrNoRows {
			logrus.Errorf("no account with that user-id: create a new one by depositing money")
			err = errors.New("no account with that user-id")
			return fail(err)
		}
		return fail(err)
	}

	var status string

	if err = tx.QueryRowContext(ctx, getLastServiceStatusQuery, failedServiceFeeReq.UserId, failedServiceFeeReq.OrderId,
		failedServiceFeeReq.ServiceId, failedServiceFeeReq.Fee).Scan(&status); err != nil {
		if err == sql.ErrNoRows {
			logrus.Errorf("no service log with that parameters")
			err = errors.New("no service log with that parameters")
			return fail(err)
		}
		return fail(err)
	} else {
		if status == "Approved" || status == "Cancelled" {
			err = fmt.Errorf("this fee has already been %s", status)
			return fail(err)
		}
	}

	_, err = tx.ExecContext(ctx, changeServiceStatusQuery, failedServiceFeeReq.UserId, failedServiceFeeReq.OrderId,
		failedServiceFeeReq.ServiceId, failedServiceFeeReq.Fee, "Cancelled")

	if err != nil {
		return fail(err)
	}

	_, err = tx.ExecContext(ctx, decreasePendingAmountQuery, failedServiceFeeReq.UserId, failedServiceFeeReq.Fee)

	if err != nil {
		return fail(err)
	}

	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	logServiceOrderRes := r.db.QueryRowContext(ctx, getLastServiceQuery, failedServiceFeeReq.UserId, failedServiceFeeReq.OrderId,
		failedServiceFeeReq.ServiceId, failedServiceFeeReq.Fee)
	if err := logServiceOrderRes.Scan(
		&failedServiceFee.AccountId,
		&failedServiceFee.ServiceId,
		&failedServiceFee.OrderId,
		&failedServiceFee.Invoice,
		&failedServiceFee.Status,
		&failedServiceFee.CreatedAt,
		&failedServiceFee.UpdatedAt,
	); err != nil {
		return failedServiceFee, err
	}
	return failedServiceFee, nil
}
