package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"net/http"
	"strconv"
	"time"

	"git.ultraware.nl/martin/graceful"
	"github.com/PotatoesFall/ots/app"
	"github.com/PotatoesFall/ots/inf/repo"
	"github.com/PotatoesFall/ots/inf/router"

	_ "github.com/lib/pq"
)

func main() {
	portStr := flag.String(`port`, `80`, `the port to run on, default 80.`)
	dbString := flag.String(`db`, `host=localhost port=5432 user=postgres password=P@ssw0rd sslmode=disable`, `database connection string`)
	flag.Parse()

	_, err := strconv.Atoi(*portStr)
	if err != nil {
		flag.CommandLine.Usage()
	}

	db, err := sql.Open(`postgres`, *dbString)
	if err != nil {
		panic(err)
	}

	repo := repo.New(db)

	app := app.New(repo)

	r := router.Make(app)

	srv := &http.Server{Addr: `:` + *portStr, Handler: r}
	go func() {
		err := srv.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	graceful.OnShutdown(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if srv.Shutdown(ctx) != nil {
			_ = srv.Close()
		}
	})
	graceful.ShutdownOnSignal(0)
	select {}
}
