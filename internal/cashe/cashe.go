package cashe

import (
	"log"
	"main/internal/db"
	"main/internal/parser"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Cash struct {
	m    map[string]parser.Order
	mydb *pgxpool.Pool
}

func (c *Cash) Init(p *pgxpool.Pool) {
	c.mydb = p
	c.m = make(map[string]parser.Order)
	orders, err := db.LoadOrders(p)
	if err != nil {
		log.Printf("Faled to load cashe from DB\n")
		return
	}
	for _, i := range orders {
		c.m[i.OrderUid] = i
	}
}

func (c *Cash) Insert(O parser.Order) {
	c.m[O.OrderUid] = O
	err := db.Insert(&O, c.mydb)
	if err != nil {
		log.Print(err)
	} else {
		log.Printf("add order with uuid %s\n", O.OrderUid)
	}
}
func (c *Cash) Get(uuid string) (parser.Order, error) {
	O, isex := c.m[uuid]
	if !isex {
		O, err := db.LoadOrder(c.mydb, uuid)
		if err != nil {
			return parser.Order{}, err
		}
		c.Insert(O)
		return O, nil
	}
	return O, nil
}
