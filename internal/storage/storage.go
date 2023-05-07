package storage

import (
	"context"
	"database/sql"
)

type QueryContext interface {
	Exec(query string, params ...interface{}) (sql.Result, error)
	Query(query string, params ...interface{}) (*sql.Rows, error)
}

type Transaction interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	QueryContext
}

type Transactioner interface {
	Transaction
	Commit() error
	Rollback() error
}

type SQL interface {
	BeginTx(context.Context) (Transactioner, error)
	Transaction
}

type KeyValue interface {
	PutMethodMetric(context.Context, *PutMethodMetricInput) error
	GetMethodMetric(context.Context, *GetMethodMetricInput) (*MethodMetricRecord, error)
}
