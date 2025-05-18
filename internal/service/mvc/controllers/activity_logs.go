package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

	"github.com/omegatymbjiep/ilab1/internal/data"
	"github.com/omegatymbjiep/ilab1/internal/service/mvc/models"
	"github.com/omegatymbjiep/ilab1/internal/service/mvc/views"
	"time"
)

const defaultLimitPerPage = 10

type ActivityLogs struct {
	auditService *models.AuditService
}

func NewActivityLogs(auditService *models.AuditService) *ActivityLogs {
	return &ActivityLogs{
		auditService: auditService,
	}
}

func (c *ActivityLogs) UserActivityPage(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := uint64(defaultLimitPerPage)
	if limitStr != "" {
		parsedLimit, err := strconv.ParseUint(limitStr, 10, 64)
		if err == nil && parsedLimit > 0 && parsedLimit <= 50 {
			limit = parsedLimit
		}
	}

	var offset uint64
	if offsetStr != "" {
		if parsedOffset, err := strconv.ParseUint(offsetStr, 10, 64); err == nil {
			offset = parsedOffset
		}
	}

	customerID := CustomerID(r)
	logs, err := c.auditService.GetUserActivityLogs(customerID, limit+1, offset)
	if err != nil {
		InternalError(w, r, fmt.Errorf("failed to get activity logs: %w", err))
		return
	}

	hasMore := false
	if len(logs) > int(limit) {
		hasMore = true
		logs = logs[:limit]
	}

	totalCount, err := c.auditService.GetTotalLogsCount(customerID)
	if err != nil {
		Log(r).WithError(err).Warn("failed to get total logs count")
	}

	currentPage := offset/limit + 1
	totalPages := uint64(0)
	if totalCount > 0 {
		totalPages = (totalCount + limit - 1) / limit // Ceiling division
	}

	viewData := &views.ActivityLogs{
		Logs: formatLogsForDisplay(logs),
		Pagination: views.Pagination{
			CurrentPage: currentPage,
			TotalPages:  totalPages,
			Limit:       limit,
			Offset:      offset,
			HasMore:     hasMore,
			TotalItems:  totalCount,
		},
	}

	if err := Templates(r).ExecuteTemplate(w, views.ActivityLogsTemplateName, viewData); err != nil {
		Log(r).WithError(err).Error("failed to execute template")
		ape.RenderErr(w, problems.InternalError())
		return
	}
}

func formatLogsForDisplay(logs []*data.AuditLog) []views.FormattedLog {
	result := make([]views.FormattedLog, len(logs))
	for i, log := range logs {
		// Parse details JSON
		var details map[string]interface{}
		if err := json.Unmarshal(log.Details, &details); err != nil {
			details = map[string]interface{}{"error": "Could not parse details"}
		}

		result[i] = views.FormattedLog{
			ID:        log.ID.String(),
			Action:    formatActionForDisplay(log.Action),
			AccountID: log.AccountID.String(),
			Details:   details,
			CreatedAt: log.CreatedAt.Format(time.Stamp),
		}
	}
	return result
}

// Helper to format action enum for display
func formatActionForDisplay(action data.AuditAction) string {
	switch action {
	case data.AuditActionLoginSuccess:
		return "Login"
	case data.AuditActionAccountCreated:
		return "Account Created"
	case data.AuditActionAccountDeleted:
		return "Account Deleted"
	case data.AuditActionDepositMade:
		return "Deposit"
	case data.AuditActionWithdrawalMade:
		return "Withdrawal"
	case data.AuditActionTransferMade:
		return "Transfer"
	case data.AuditActionExcelReportGenerated:
		return "Excel Report Generated"
	default:
		return string(action)
	}
}
