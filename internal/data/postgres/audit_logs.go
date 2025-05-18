package postgres

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"gitlab.com/distributed_lab/kit/pgdb"

	"github.com/omegatymbjiep/ilab1/internal/data"
)

const (
	auditLogsTableName = "audit_logs"
	customerIDColumn   = "customer_id"
	accountIDColumn    = "account_id"
	actionColumn       = "action"
)

type auditLogsQ struct {
	*crudQ[*data.AuditLog, uuid.UUID]
}

func NewAuditLogsQ(db *pgdb.DB) data.AuditLogs {
	return &auditLogsQ{
		newCRUDQ[*data.AuditLog, uuid.UUID](db, auditLogsTableName),
	}
}

func (q *auditLogsQ) WhereCustomerID(customerID uuid.UUID) data.AuditLogs {
	q.sel = q.sel.Where(sq.Eq{customerIDColumn: customerID})
	return q
}

func (q *auditLogsQ) WhereAccountID(accountID uuid.UUID) data.AuditLogs {
	q.sel = q.sel.Where(sq.Eq{accountIDColumn: accountID})
	return q
}

func (q *auditLogsQ) WhereAction(action data.AuditAction) data.AuditLogs {
	q.sel = q.sel.Where(sq.Eq{actionColumn: action})
	return q
}

func (q *auditLogsQ) Limit(limit uint64) data.AuditLogs {
	q.sel = q.sel.Limit(limit)
	return q
}

func (q *auditLogsQ) Offset(offset uint64) data.AuditLogs {
	q.sel = q.sel.Offset(offset)
	return q
}

func (q *auditLogsQ) OrderBy(orderBy ...string) data.AuditLogs {
	q.sel = q.sel.OrderBy(orderBy...)
	return q
}
