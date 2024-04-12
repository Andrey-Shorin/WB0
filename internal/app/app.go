package app

import (
	"encoding/json"
	"log"
	"main/internal/cashe"
	"main/internal/config"
	"main/internal/db"
	"main/internal/parser"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

type App struct {
	cash cashe.Cash
	mydb *pgxpool.Pool
	nc   *nats.Conn
	sc   stan.Conn
	sub  stan.Subscription
}

func (a *App) Handler(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Query().Get("uuid")
	if uuid == "" {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}
	order, err := a.cash.Get(uuid)
	if err != nil {
		log.Print(err)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	str, err := json.MarshalIndent(order, "", "  ")
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(str)

}
func (a *App) Init(conf config.Config) {
	var err error
	a.mydb = db.ConnectDB(conf.DbURL)
	a.cash.Init(a.mydb)
	a.nc, err = nats.Connect(conf.NatsIP)
	if err != nil {
		a.mydb.Close()
		log.Fatal(err)
	}
	a.sc, err = stan.Connect(
		conf.ClasterName, conf.ClientID,
		stan.NatsConn(a.nc), stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}))
	if err != nil {
		a.mydb.Close()
		a.nc.Close()
		log.Fatal(err)
	}
	a.sub, err = a.sc.QueueSubscribe("test", "rrr", a.handler)
	if err != nil {
		a.mydb.Close()
		a.sc.Close()
		a.nc.Close()
		log.Fatal(err)
	}

}
func (a *App) Close() {
	a.sub.Unsubscribe()
	a.mydb.Close()
	a.sc.Close()
	a.nc.Close()

}
func (a *App) handler(msg *stan.Msg) {
	//fmt.Print(msg)f
	data, err := parser.Parser(msg.Data)
	if err != nil {
		return
	}
	a.cash.Insert(*data)

}
