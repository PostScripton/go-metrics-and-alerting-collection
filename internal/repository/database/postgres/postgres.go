package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
)

type postgres struct {
	pool *pgxpool.Pool
}

func ConnectToDB(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	if dsn == "" {
		return nil, fmt.Errorf("no dsn")
	}
	return pgxpool.Connect(ctx, dsn)
}

func Migrate(pool *pgxpool.Pool) {
	dir := "./internal/repository/database/postgres/migrations/"
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Printf("nothing to migrate: %s\n", err)
		return
	}

	for _, file := range files {
		sql, err := os.ReadFile(dir + file.Name())
		if err != nil {
			fmt.Printf("Cannot open file [%s]: %s", file.Name(), err)
			return
		}
		query := string(sql)

		migration := file.Name()[:len(file.Name())-4]
		fmt.Printf("[%s] Migrating...\n", migration)
		_, err = pool.Exec(context.Background(), query)
		if err != nil {
			fmt.Printf("[%s] Migration failed: %s\n", migration, err)
			return
		}
		fmt.Printf("[%s] Migrated!\n", migration)
	}
}

func NewPostgres(pool *pgxpool.Pool) *postgres {
	return &postgres{pool: pool}
}

func (p *postgres) GetCollection() (map[string]metrics.Metrics, error) {
	q := `SELECT id, type, delta, value FROM metrics;`
	rows, err := p.pool.Query(context.Background(), q)
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
	err := p.pool.QueryRow(context.Background(), q, metric.ID, metric.Type).Scan(&m.ID, &m.Type, &m.Delta, &m.Value)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, metrics.ErrNoValue
	}

	return &m, err
}

func (p *postgres) Store(metric metrics.Metrics) error {
	q := `SELECT id, type, delta, value FROM metrics where id = $1 and type = $2;`
	var m metrics.Metrics
	err := p.pool.QueryRow(context.Background(), q, metric.ID, metric.Type).Scan(&m.ID, &m.Type, &m.Delta, &m.Value)

	if errors.Is(err, pgx.ErrNoRows) {
		q = `INSERT INTO metrics (id, type, delta, value) VALUES ($1, $2, $3, $4);`
	} else {
		q = `UPDATE metrics SET delta = $3, value = $4 WHERE id = $1 and type = $2;`
	}
	metrics.Update(&m, &metric)

	if _, err = p.pool.Exec(context.Background(), q, m.ID, m.Type, m.Delta, m.Value); err != nil {
		return err
	}

	return nil
}

func (p *postgres) CleanUp() error {
	q := `TRUNCATE metrics;`
	if _, err := p.pool.Exec(context.Background(), q); err != nil {
		return err
	}
	return nil
}
