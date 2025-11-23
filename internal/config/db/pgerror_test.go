package db

import (
	"errors"
	"testing"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

func TestNewPostgresErrorClassifier(t *testing.T) {
	classifier := NewPostgresErrorClassifier()
	require.NotNil(t, classifier)
}

func TestClassify(t *testing.T) {
	classifier := NewPostgresErrorClassifier()

	tests := []struct {
		name  string
		value error
		want  PGErrorClassification
	}{
		{
			name: "ConnectionException",
			value: &pgconn.PgError{
				Code: pgerrcode.ConnectionException,
			},
			want: Retriable,
		},
		{
			name: "ConnectionDoesNotExist",
			value: &pgconn.PgError{
				Code: pgerrcode.ConnectionDoesNotExist,
			},
			want: Retriable,
		},
		{
			name: "ConnectionFailure",
			value: &pgconn.PgError{
				Code: pgerrcode.ConnectionFailure,
			},
			want: Retriable,
		},

		{
			name: "TransactionRollback",
			value: &pgconn.PgError{
				Code: pgerrcode.TransactionRollback,
			},
			want: NonRetriable,
		},
		{
			name: "CaseNotFound",
			value: &pgconn.PgError{
				Code: pgerrcode.CaseNotFound,
			},
			want: NonRetriable,
		},
		{
			name: "random code",
			value: &pgconn.PgError{
				Code: "qwertyuiioooo",
			},
			want: NonRetriable,
		},
		{
			name:  "NilError",
			value: nil,
			want:  NonRetriable,
		},
		{
			name:  "NonPgError",
			value: errors.New("regular error"),
			want:  NonRetriable,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := classifier.Classify(tt.value)
			require.Equal(t, res, tt.want)
		})
	}
}
