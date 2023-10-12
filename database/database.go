package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/floetenleague/floetenleague/config"
	"github.com/floetenleague/floetenleague/database/dbgen"
	"github.com/golang-migrate/migrate/v4"
	migratepgx "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zerologadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog/log"
)

//go:embed migrations/*.sql
var migrations embed.FS

type DB struct {
	pool *pgxpool.Pool
}

func (db *DB) Aquire(ctx context.Context) (*Queries, error) {
	conn, err := db.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	return &Queries{
		Queries: dbgen.New(conn),
		Conn:    conn,
	}, nil
}

type Queries struct {
	*dbgen.Queries
	Conn *pgxpool.Conn
}

func (q *Queries) BeginTxFunc(ctx context.Context, f func(*dbgen.Queries) error) error {
	return q.Conn.BeginTxFunc(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable}, func(tx pgx.Tx) error {
		return f(dbgen.New(tx))
	})
}

func (q *Queries) Close() {
	q.Conn.Release()
}

// New creates a db instance.
func New(cfg *config.Config) *DB {
	pgCfg, err := pgxpool.ParseConfig(cfg.DBConnection)
	if err != nil {
		log.Fatal().Err(err).Msg("err parsing conn string")
	}

	stdDB, err := sql.Open("pgx", cfg.DBConnection)
	defer stdDB.Close()
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't open database")
	}
	log.Debug().Msg("Initializing database")

	err = initDB(stdDB, pgCfg.ConnConfig.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't migrate database")
	}

	log.Debug().Msg("Database initialized")
	pgCfg.ConnConfig.Logger = zerologadapter.NewLogger(log.Logger)
	pool, err := pgxpool.ConnectConfig(context.Background(), pgCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't open database")
	}
	return &DB{
		pool: pool,
	}
}

func initDB(db *sql.DB, database string) error {
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}
	driver, err := migratepgx.WithInstance(db, &migratepgx.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("floetenleague", source, database, driver)
	if err != nil {
		return err
	}

	m.Log = &L{}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("up error: %w", err)
	}
	return nil
}

type L struct{}

func (*L) Printf(format string, v ...interface{}) {
	log.Error().Msgf(format, v...)
}
func (*L) Verbose() bool {
	return true
}
