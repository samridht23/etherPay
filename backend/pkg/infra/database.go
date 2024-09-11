package infra

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Pool struct {
	Pool *pgxpool.Pool
}

func InitializePool() (*Pool, error) {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("failed initializing database pool: %v", err)
	}
	zap.S().Info("Database pool initialized")

	p := &Pool{
		Pool: pool,
	}

	if err := p.PingDatabase(); err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	return p, nil
}

func (p *Pool) PingDatabase() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt to ping the database
	if err := p.Pool.Ping(ctx); err != nil {
		zap.S().Errorf("Failed to connect to the PostgreSQL database: %v", err)
		return err
	}
	zap.S().Info("Successfully connected to the PostgreSQL database")
	return nil
}

func (p *Pool) New() *pgxpool.Pool {
	return p.Pool
}

func (p *Pool) Close() {
	if p.Pool != nil {
		p.Pool.Close()
		zap.S().Info("Database pool closed")
	}
}
