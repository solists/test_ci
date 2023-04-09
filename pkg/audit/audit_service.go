package audit

import (
	"context"

	"github.com/jmoiron/sqlx"
	jsoniter "github.com/json-iterator/go"
	"github.com/solists/test_ci/pkg/logger"
)

//go:generate mockgen -source=${GOFILE} -destination=mock/mock_${GOFILE}
type Service interface {
	Log(log *Log)
}

type AuditService struct {
	db *sqlx.DB
}

func NewAuditService(db *sqlx.DB) *AuditService {
	return &AuditService{db: db}
}

type Log struct {
	UserID    uint64      `db:"user_id"`
	Data      interface{} `db:"data" sql:"type:jsonb"`
	Operation string      `db:"operation"`
	Response  interface{} `db:"response" sql:"type:jsonb"`
	Status    *uint64     `db:"status"`
}

func (a *AuditService) Log(log *Log) {
	go func() {
		req, err := jsoniter.Marshal(log.Data)
		if err != nil {
			logger.Errorf("auditLog marshal: %v", err)
			return
		}
		log.Data = req

		if log.Response != nil {
			resp, err := jsoniter.Marshal(log.Response)
			if err != nil {
				logger.Errorf("auditLog marshal: %v", err)
				return
			}
			log.Response = resp
		}

		if _, err = a.db.NamedExecContext(context.Background(), query, log); err != nil {
			logger.Errorf("auditLog insert: %v", err)
			return
		}
	}()
}

const query = `INSERT INTO audit_log (user_id, data, operation, response, status) 
					VALUES (:user_id, :data, :operation, :response, :status)`
