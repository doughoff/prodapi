package config

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
)

type CustomTracer struct{}

func (t *CustomTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	log.Printf("Start Query: %s \nArgs: %+v\n", data.SQL, data.Args)
	return ctx
}

func (t *CustomTracer) TraceQueryEnd(_ context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	if data.Err != nil {
		fmt.Printf("Query Error: %v\n", data.Err)
	} else {
		fmt.Printf("End Query: %s\n", data.CommandTag)
	}
}

func NewPgxConn() *pgx.Conn {
	ctx := context.Background()
	connConfig, err := pgx.ParseConfig(os.Getenv("DB_URL"))
	//connConfig.Tracer = &CustomTracer{}
	conn, err := pgx.ConnectConfig(ctx, connConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// register custom types
	t, err := conn.LoadType(context.Background(), "status")
	if err != nil {
		panic(err)
	}
	conn.TypeMap().RegisterType(t)

	t, err = conn.LoadType(context.Background(), "_status")
	if err != nil {
		panic(err)
	}
	conn.TypeMap().RegisterType(t)

	t, err = conn.LoadType(context.Background(), "unit")
	if err != nil {
		panic(err)
	}
	conn.TypeMap().RegisterType(t)

	t, err = conn.LoadType(context.Background(), "_unit")
	if err != nil {
		panic(err)
	}
	conn.TypeMap().RegisterType(t)

	return conn
}
