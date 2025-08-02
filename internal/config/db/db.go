package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	models "metrics/internal/model"
	"time"

	"github.com/jackc/pgx/v5"
)

type DB struct {
	conn *pgx.Conn
}

func NewDB(ctx context.Context, connStr string) (*DB, error) {
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}
	return &DB{
		conn: conn,
	}, nil
}

func (db *DB) Ping() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := db.conn.Ping(ctx)
	return err == nil
}

func (db *DB) ReadAll() ([]models.Metrics, error) {
	var metrics []models.Metrics
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := db.conn.Query(ctx, "SELECT id, type, delta, value FROM metrics")
	if err != nil {
		log.Printf("Ошибка выполнения запроса: %v", err)
		return metrics, err
	}
	defer rows.Close()

	for rows.Next() {
		var m models.Metrics
		var delta sql.NullInt64
		var value sql.NullFloat64
		err := rows.Scan(&m.ID, &m.MType, &delta, &value)
		if err != nil {
			log.Printf("Ошибка при сканировании строки: %v", err)
			return metrics, err
		}
		// Устанавливаем Delta только если тип counter и значение не NULL
		if m.MType == "counter" && delta.Valid {
			m.Delta = &delta.Int64
		}

		// Устанавливаем Value только если тип gauge и значение не NULL
		if m.MType == "gauge" && value.Valid {
			m.Value = &value.Float64
		}

		metrics = append(metrics, m)
	}
	return metrics, nil
}

func (db *DB) SaveAll(metrics []models.Metrics) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	for _, m := range metrics {
		_, err := db.conn.Exec(ctx, `
            INSERT INTO metrics (id, type, delta, value)
            VALUES (@id, @type, @delta, @value)
            ON CONFLICT (id) 
            DO UPDATE SET
                type = EXCLUDED.type,
                delta = EXCLUDED.delta,
                value = EXCLUDED.value;
        `, pgx.NamedArgs{
			"id":    m.ID,
			"type":  m.MType,
			"delta": m.Delta,
			"value": m.Value,
		})
		if err != nil {
			return fmt.Errorf("ошибка при сохранении метрики %s: %w", m.ID, err)
		}
	}
	return nil
}
