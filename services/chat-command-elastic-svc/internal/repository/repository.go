package repository

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
)

type DB interface {
	AwaitSchemaAgreement(ctx context.Context) error
	Bind(stmt string, b func(q *gocql.QueryInfo) ([]interface{}, error)) *gocql.Query
	Close()
	Closed() bool
	ContextQuery(ctx context.Context, stmt string, names []string) *gocqlx.Queryx
	ExecStmt(stmt string) error
	ExecuteBatch(batch *gocql.Batch) error
	ExecuteBatchCAS(batch *gocql.Batch, dest ...interface{}) (applied bool, iter *gocql.Iter, err error)
	KeyspaceMetadata(keyspace string) (*gocql.KeyspaceMetadata, error)
	MapExecuteBatchCAS(batch *gocql.Batch, dest map[string]interface{}) (applied bool, iter *gocql.Iter, err error)
	NewBatch(typ gocql.BatchType) *gocql.Batch
	Query(stmt string, names []string) *gocqlx.Queryx
	SetConsistency(cons gocql.Consistency)
	SetPageSize(n int)
	SetPrefetch(p float64)
	SetTrace(trace gocql.Tracer)
}
