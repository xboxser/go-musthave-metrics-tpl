package db

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// ErrorClassification тип для классификации ошибок
type PGErrorClassification int

const (
	// NonRetriable - операцию не следует повторять
	NonRetriable PGErrorClassification = iota

	// Retriable - операцию можно повторить
	Retriable
)

// PostgresErrorClassifier классификатор ошибок PostgreSQL
type PostgresErrorClassifier struct{}

func NewPostgresErrorClassifier() *PostgresErrorClassifier {
	return &PostgresErrorClassifier{}
}

// Classify классифицирует ошибку и возвращает PGErrorClassification
func (c *PostgresErrorClassifier) Classify(err error) PGErrorClassification {
	if err == nil {
		return NonRetriable
	}

	// Проверяем и конвертируем в pgconn.PgError, если это возможно
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return ClassifyPgError(pgErr)
	}

	// По умолчанию считаем ошибку неповторяемой
	return NonRetriable
}

func ClassifyPgError(pgErr *pgconn.PgError) PGErrorClassification {
	// Коды ошибок PostgreSQL: https://www.postgresql.org/docs/current/errcodes-appendix.html
	if pgErr == nil {
		return NonRetriable
	}
	switch pgErr.Code {
	// Класс 08 - Ошибки соединения
	case pgerrcode.ConnectionException,
		pgerrcode.ConnectionDoesNotExist,
		pgerrcode.ConnectionFailure:
		return Retriable
	}

	// По умолчанию считаем ошибку неповторяемой
	return NonRetriable
}
