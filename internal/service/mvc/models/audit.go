package models

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/omegatymbjiep/ilab1/internal/data"
)

type AuditService struct {
	db data.MainQ
}

func NewAuditService(db data.MainQ) *AuditService {
	return &AuditService{
		db: db,
	}
}

type AuditDetails map[string]interface{}

func (m *AuditService) LogAction(
	customerID uuid.UUID,
	accountID *uuid.UUID,
	action data.AuditAction,
	details AuditDetails,
) error {
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return fmt.Errorf("failed to marshal audit details: %w", err)
	}

	auditLog := &data.AuditLog{
		CustomerID: customerID,
		Action:     action,
		Details:    detailsJSON,
		AccountID:  accountID,
	}

	if err := m.db.AuditLogs().Insert(auditLog); err != nil {
		return fmt.Errorf("failed to insert audit log: %w", err)
	}

	return nil
}

func (m *AuditService) GetUserActivityLogs(
	customerID uuid.UUID,
	limit uint64,
	offset uint64,
) ([]*data.AuditLog, error) {
	logs, err := m.db.AuditLogs().
		WhereCustomerID(customerID).
		OrderBy("created_at DESC").
		Limit(limit).
		Offset(offset).
		Select()

	if err != nil {
		return nil, fmt.Errorf("failed to get activity logs: %w", err)
	}

	return logs, nil
}

func (m *AuditService) GetAccountActivityLogs(
	customerID uuid.UUID,
	accountID uuid.UUID,
	limit uint64,
	offset uint64,
) ([]*data.AuditLog, error) {
	hasAccess, err := m.db.CustomersAccounts().HasAccount(customerID, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to check account access: %w", err)
	}
	if !hasAccess {
		return nil, fmt.Errorf("customer does not have access to this account")
	}

	logs, err := m.db.AuditLogs().
		WhereAccountID(accountID).
		OrderBy("created_at DESC").
		Limit(limit).
		Offset(offset).
		Select()

	if err != nil {
		return nil, fmt.Errorf("failed to get account activity logs: %w", err)
	}

	return logs, nil
}

// GetTotalLogsCount retrieves the total number of logs for a customer
func (m *AuditService) GetTotalLogsCount(customerID uuid.UUID) (uint64, error) {
	var count uint64

	count, err := m.db.AuditLogs().WhereCustomerID(customerID).Count()
	if err != nil {
		return 0, fmt.Errorf("failed to get total logs count: %w", err)
	}

	return count, nil
}

func (m *AuditService) GetAccountDetails(customerID, accountID uuid.UUID) (*data.Account, error) {
	hasAccess, err := m.db.CustomersAccounts().HasAccount(customerID, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to check account access: %w", err)
	}
	if !hasAccess {
		return nil, fmt.Errorf("customer does not have access to this account")
	}

	account := new(data.Account)
	ok, err := m.db.Accounts().WhereID(accountID).Get(account)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("account not found")
	}

	return account, nil
}

func (m *AuditService) logAccountCreated(customerID uuid.UUID, accountID uuid.UUID) error {
	err := m.LogAction(customerID, &accountID, data.AuditActionAccountCreated, nil)
	if err != nil {
		return fmt.Errorf("failed to log audit action: %w", err)
	}

	return nil
}

func (m *AuditService) logAccountDeleted(customerID uuid.UUID, accountID uuid.UUID) error {
	err := m.LogAction(customerID, &accountID, data.AuditActionAccountDeleted, nil)
	if err != nil {
		return fmt.Errorf("failed to log audit action: %w", err)
	}

	return nil
}

func (m *AuditService) logExcelReportGenerated(customerID uuid.UUID, accountID uuid.UUID) error {
	err := m.LogAction(customerID, &accountID, data.AuditActionExcelReportGenerated, nil)
	if err != nil {
		return fmt.Errorf("failed to log audit action: %w", err)
	}

	return nil
}

func (m *AuditService) logDepositMade(customerID uuid.UUID, accountID uuid.UUID, amount uint) error {
	details := AuditDetails{
		"amount": amount,
	}

	err := m.LogAction(customerID, &accountID, data.AuditActionDepositMade, details)
	if err != nil {
		return fmt.Errorf("failed to log audit action: %w", err)
	}

	return nil
}

func (m *AuditService) logWithdrawalMade(customerID uuid.UUID, accountID uuid.UUID, amount uint) error {
	details := AuditDetails{
		"amount": amount,
	}

	err := m.LogAction(customerID, &accountID, data.AuditActionWithdrawalMade, details)
	if err != nil {
		return fmt.Errorf("failed to log audit action: %w", err)
	}

	return nil
}

func (m *AuditService) logTransferMade(
	customerID uuid.UUID,
	fromAccountID uuid.UUID,
	toAccountID uuid.UUID,
	amount uint,
) error {
	details := AuditDetails{
		"from_account": fromAccountID,
		"to_account":   toAccountID,
		"amount":       amount,
	}

	err := m.LogAction(customerID, &fromAccountID, data.AuditActionTransferMade, details)
	if err != nil {
		return fmt.Errorf("failed to log audit action: %w", err)
	}

	return nil
}
