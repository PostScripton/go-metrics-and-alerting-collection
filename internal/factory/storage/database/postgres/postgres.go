package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/PostScripton/go-metrics-and-alerting-collection/internal/metrics"
)

type Postgres struct {
	pool *pgxpool.Pool
}

var Migrated = false

func NewPostgres(dsn string) (*Postgres, error) {
	if dsn == "" {
		return nil, fmt.Errorf("no dsn")
	}

	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	postgres := &Postgres{pool: pool}
	postgres.Migrate()

	return postgres, nil
}

func (p *Postgres) Migrate() {
	if Migrated {
		return
	}

	dir := "./internal/factory/storage/database/postgres/migrations/"
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Debug().Err(err).Msg("Nothing to migrate")
		return
	}

	for _, file := range files {
		sql, err := os.ReadFile(dir + file.Name())
		if err != nil {
			log.Warn().Err(err).Msgf("Cannot open file [%s]", file.Name())
			return
		}
		query := string(sql)

		migration := file.Name()[:len(file.Name())-4]
		log.Info().Msgf("[%s] Migrating...", migration)
		_, err = p.pool.Exec(context.Background(), query)
		if err != nil {
			log.Warn().Err(err).Msgf("[%s] Migration failed", migration)
			return
		}
		log.Info().Msgf("[%s] Migrated!", migration)
	}

	Migrated = true
}

func (p *Postgres) GetCollection() (map[string]metrics.Metrics, error) {
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

func (p *Postgres) StoreCollection(metricsCollection map[string]metrics.Metrics) error {
	ctx := context.Background()
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, m := range metricsCollection {
		if err := p.Store(m); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (p *Postgres) Get(metric metrics.Metrics) (*metrics.Metrics, error) {
	q := `SELECT id, type, delta, value FROM metrics WHERE id = $1 and type = $2;`

	var m metrics.Metrics
	err := p.pool.QueryRow(context.Background(), q, metric.ID, metric.Type).Scan(&m.ID, &m.Type, &m.Delta, &m.Value)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, metrics.ErrNoValue
	}

	return &m, err
}

func (p *Postgres) Store(metric metrics.Metrics) error {
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

func (p *Postgres) CleanUp() error {
	q := `TRUNCATE metrics;`
	if _, err := p.pool.Exec(context.Background(), q); err != nil {
		return err
	}
	return nil
}

func (p *Postgres) Ping(ctx context.Context) error {
	return p.pool.Ping(ctx)
}

func (p *Postgres) Close() {
	p.pool.Close()
}
