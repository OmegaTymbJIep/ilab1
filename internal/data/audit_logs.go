package data

import (
	"encoding/json"
	"github.com/google/uuid"
)

type AuditAction string

const (
	AuditActionLoginSuccess         AuditAction = "login_success"
	AuditActionAccountCreated       AuditAction = "account_created"
	AuditActionAccountDeleted       AuditAction = "account_deleted"
	AuditActionDepositMade          AuditAction = "deposit_made"
	AuditActionWithdrawalMade       AuditAction = "withdrawal_made"
	AuditActionTransferMade         AuditAction = "transfer_made"
	AuditActionExcelReportGenerated AuditAction = "excel_report_generated"
)

type AuditLogs interface {
	CRUDQ[*AuditLog, uuid.UUID]

	WhereCustomerID(customerID uuid.UUID) AuditLogs
	WhereAccountID(accountID uuid.UUID) AuditLogs
	WhereAction(action AuditAction) AuditLogs

	Limit(limit uint64) AuditLogs
	Offset(offset uint64) AuditLogs
	OrderBy(orderBy ...string) AuditLogs

	Count() (uint64, error)
}

type AuditLog struct {
	Entity[uuid.UUID] `structs:"-"`

	CustomerID uuid.UUID       `db:"customer_id"  structs:"customer_id"`
	AccountID  *uuid.UUID      `db:"account_id"   structs:"account_id"`
	Action     AuditAction     `db:"action"       structs:"action"`
	Details    json.RawMessage `db:"details"      structs:"details"`
}
