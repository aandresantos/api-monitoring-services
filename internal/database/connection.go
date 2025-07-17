package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBClient struct {
	Conn *pgxpool.Pool
}

var DB *pgxpool.Pool

func ConnectDB() (*DBClient, error) {
	databaseUrl := os.Getenv("DATABASE_URL")

	if databaseUrl == "" {
		return nil, fmt.Errorf("DATABASE_URL não está definida")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config, errConfig := pgxpool.ParseConfig(databaseUrl)

	if errConfig != nil {
		return nil, fmt.Errorf("erro ao fazer o parse da conexão: %w", errConfig)
	}

	pool, errPool := pgxpool.NewWithConfig(ctx, config)

	if errPool != nil {
		return nil, fmt.Errorf("erro ao conectar no banco: %w", errPool)
	}

	if errPing := pool.Ping(ctx); errPing != nil {
		return nil, fmt.Errorf("erro ao pingar no db: %w", errPing)
	}

	DB = pool
	fmt.Println(("DB conectado"))

	return &DBClient{Conn: DB}, nil
}

func CloseConnectionDB() {
	if DB != nil {
		DB.Close()

		fmt.Println(("DB desconectado"))
	}
}
