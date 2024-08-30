package service

import (
	"context"

	audit "github.com/wilfridterry/audit-log/pkg/domain"
)

type AuditClient interface {
	SendLogRequest(ctx context.Context, req audit.LogItem) error
}
