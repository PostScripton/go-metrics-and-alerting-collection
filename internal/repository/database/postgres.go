package database

import (
	"context"
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"github.com/jackc/pgx/v4"
	"time"
)

type postgres struct {
	conn      *pgx.Conn
	Connected bool
}

func NewPostgres(ctx context.Context, addr string) (*postgres, error) {
	conn, err := pgx.Connect(ctx, addr)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err = conn.Ping(ctx); err != nil {
		return nil, err
	}

	return &postgres{conn: conn, Connected: true}, nil
}

func (p *postgres) GetCollection() (map[string]metrics.Metrics, error) {
	q := `SELECT id, type, delta, value FROM metrics;`
	rows, err := p.conn.Query(context.Background(), q)
	if err != nil {
		return nil, err
	}

	metricsCollection := make(map[string]metrics.Metrics)

	for rows.Next() {
		var metric metrics.Metrics
		err = rows.Scan(&metric.ID, &metric.Type, &metric.Delta, &metric.Value)
		if err != nil {
			return nil, err
		}
		metricsCollection[metric.ID] = metric
	}

	return metricsCollection, nil
}

func (p *postgres) StoreCollection(metricsCollection map[string]metrics.Metrics) error {
	q := `TRUNCATE TABLE metrics;`
	p.conn.QueryRow(context.Background(), q)

	for _, m := range metricsCollection {
		if err := p.Store(m); err != nil {
			return err
		}
	}

	return nil
}

func (p *postgres) Get(metric metrics.Metrics) (*metrics.Metrics, error) {
	q := `SELECT id, type, delta, value FROM metrics WHERE id = $1 and type = $2;`

	var m metrics.Metrics
	err := p.conn.QueryRow(context.Background(), q, metric.ID, metric.Type).Scan(&m.ID, &m.Type, &m.Delta, &m.Value)
	if err == pgx.ErrNoRows {
		return nil, metrics.ErrNoValue
	}

	return &m, err
}

func (p *postgres) Store(metric metrics.Metrics) error {
	q := `SELECT id, type, delta, value FROM metrics where id = $1 and type = $2;`
	err := p.conn.QueryRow(context.Background(), q, metric.ID, metric.Type).Scan()
	if err == pgx.ErrNoRows {
		q = `INSERT INTO metrics (id, type, delta, value) VALUES ($1, $2, $3, $4);`
	} else {
		q = `UPDATE metrics SET delta = $3, value = $4 WHERE id = $1 and type = $2;`
	}

	if _, err = p.conn.Exec(context.Background(), q, metric.ID, metric.Type, metric.Delta, metric.Value); err != nil {
		return err
	}

	return nil
}

func (p *postgres) Close(ctx context.Context) error {
	if err := p.conn.Close(ctx); err != nil {
		return err
	}
	return nil
}

func (p *postgres) Ping(ctx context.Context) (context.CancelFunc, error) {
	newCtx, cancel := context.WithTimeout(ctx, 1*time.Second)
	if !p.Connected {
		return cancel, fmt.Errorf("DB is not connected, unable to ping")
	}
	if err := p.conn.Ping(newCtx); err != nil {
		p.Connected = false
		return cancel, err
	}
	return cancel, nil
}
