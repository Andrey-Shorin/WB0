package db

import (
	"context"
	"fmt"
	"log"
	"main/internal/parser"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(dbURL string) *pgxpool.Pool {
	dbpool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal(err)

	}

	err = dbpool.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return dbpool

}

// func LoadCashe(dbpool *pgxpool.Pool) ([]parser.Order, error) {
// 	cashe, err := loadOrders(dbpool)

// }
func loadDelivery(dbpool *pgxpool.Pool, uuid string) (parser.Delivery, error) {
	query := `SELECT name, phone, zip, city, address, 
	region,email FROM delivery WHERE order_uid = $1 LIMIT 1`
	rows, err := dbpool.Query(context.Background(), query, uuid)
	if err != nil {
		fmt.Println(err)
		return parser.Delivery{}, err
	}
	var i parser.Delivery
	for rows.Next() {
		err = rows.Scan(
			&i.Name, &i.Phone, &i.Zip, &i.City, &i.Address,
			&i.Region, &i.Email)
	}
	if err != nil {
		fmt.Println(err)
		return parser.Delivery{}, err
	}
	return i, nil

}
func loadPayment(dbpool *pgxpool.Pool, uuid string) (parser.Payment, error) {
	query := `SELECT transaction, request_id, currency, provider, amount, 
	payment_dt,bank,delivery_cost,goods_total,custom_fee FROM payment WHERE order_uid = $1 LIMIT 1;`
	rows, err := dbpool.Query(context.Background(), query, uuid)
	if err != nil {
		fmt.Println(err)
		return parser.Payment{}, err
	}
	var i parser.Payment
	for rows.Next() {
		err = rows.Scan(
			&i.Transaction, &i.RequestId, &i.Currency, &i.Provider, &i.Amount,
			&i.PaymentDt, &i.Bank, &i.DeliveryCost, &i.GoodsTotal, &i.CustomFee)
	}
	if err != nil {
		fmt.Println(err)
		return parser.Payment{}, err
	}
	return i, nil

}
func loadItems(dbpool *pgxpool.Pool, uuid string) ([]parser.Item, error) {
	query := `SELECT i.chrt_id, i.track_number,
	i.price, i.rid, i.name, i.sale, i.size, i.total_price, i.nm_id,
	i.brand, i.status 
	FROM orders_item oi
	JOIN item i ON oi.item_id = i.id
	WHERE oi.order_uid = $1;`
	rows, err := dbpool.Query(context.Background(), query, uuid)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var Item []parser.Item
	for rows.Next() {
		var i parser.Item
		err = rows.Scan(
			&i.ChrtId, &i.TrackNumber, &i.Price, &i.Rid, &i.Name,
			&i.Sale, &i.Size, &i.TotalPrice, &i.NmId, &i.Brand, &i.Status)
		Item = append(Item, i)

	}
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return Item, nil

}
func LoadOrder(dbpool *pgxpool.Pool, uuid string) (parser.Order, error) {
	query := `SELECT order_uid, track_number, entry, locale, internal_signature, 
	customer_id,delivery_service,shardkey,sm_id,date_created,oof_shard FROM orders WHERE order_uid = $1 LIMIT 1`
	rows, err := dbpool.Query(context.Background(), query, uuid)
	if err != nil {
		fmt.Println(err)
		return parser.Order{}, err
	}

	var Order parser.Order
	if !rows.Next() {
		return parser.Order{}, fmt.Errorf("can't find uuid %s", uuid)
	}
	err = rows.Scan(
		&Order.OrderUid, &Order.TrackNumber, &Order.Entry, &Order.Locale, &Order.InternalSignature,
		&Order.CustomerId, &Order.DeliveryService, &Order.Shardkey, &Order.SmId, &Order.DateCreated, &Order.OofShard,
	)
	fmt.Print("LLLLL")

	if err != nil {
		fmt.Println(err)
		return parser.Order{}, err
	}
	Order.Delivery, err = loadDelivery(dbpool, uuid)
	if err != nil {
		fmt.Println(err)
		return parser.Order{}, err
	}
	Order.Payment, err = loadPayment(dbpool, uuid)
	if err != nil {
		fmt.Println(err)
		return parser.Order{}, err
	}
	Order.Items, err = loadItems(dbpool, uuid)
	if err != nil {
		fmt.Println(err)
		return parser.Order{}, err
	}
	return Order, nil
}
func LoadOrders(dbpool *pgxpool.Pool) ([]parser.Order, error) {
	var cashe []parser.Order
	query := `SELECT order_uid, track_number, entry, locale, internal_signature, 
	customer_id,delivery_service,shardkey,sm_id,date_created, oof_shard FROM orders LIMIT 10`
	rows, err := dbpool.Query(context.Background(), query)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for rows.Next() {
		var i parser.Order

		err := rows.Scan(
			&i.OrderUid, &i.TrackNumber, &i.Entry, &i.Locale, &i.InternalSignature,
			&i.CustomerId, &i.DeliveryService, &i.Shardkey, &i.SmId, &i.DateCreated, &i.OofShard,
		)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		i.Delivery, err = loadDelivery(dbpool, i.OrderUid)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		i.Payment, err = loadPayment(dbpool, i.OrderUid)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		i.Items, err = loadItems(dbpool, i.OrderUid)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		cashe = append(cashe, i)
	}
	return cashe, nil

}

func Insert(order *parser.Order, dbpool *pgxpool.Pool) error {
	Tx, err := dbpool.Begin(context.Background())
	defer func() {
		if err == nil {
			err = Tx.Commit(context.Background())
		} else {
			Tx.Rollback(context.Background())
		}
	}()
	err = insertOrder(Tx, order)
	if err != nil {
		return err
	}
	err = insertPayment(Tx, &order.Payment, order.OrderUid)
	if err != nil {
		return err
	}
	err = insertDelivery(Tx, &order.Delivery, order.OrderUid)
	if err != nil {
		return err
	}

	for _, i := range order.Items {
		err = insertItem(Tx, &i, order.OrderUid)
		if err != nil {
			return err
		}
	}
	return err

}
func insertOrder(Tx pgx.Tx, order *parser.Order) error {
	query := `
	INSERT INTO orders (
						order_uid, track_number, entry, internal_signature,
						locale, customer_id, delivery_service, shardkey,
						sm_id, date_created, oof_shard
						)
	VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
			)
`

	_, err := Tx.Exec(context.Background(),
		query,
		order.OrderUid, order.TrackNumber, order.Entry, order.InternalSignature,
		order.Locale, order.CustomerId, order.DeliveryService, order.Shardkey,
		order.SmId, order.DateCreated, order.OofShard,
	)

	return err
}
func insertPayment(Tx pgx.Tx, payment *parser.Payment, orderId string) error {
	query := `INSERT INTO payment (
		transaction, currency, provider, amount,
		payment_dt, bank, delivery_cost, goods_total,
		custom_fee, order_uid
		)

VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
`

	_, err := Tx.Exec(context.Background(), query,
		payment.Transaction, payment.Currency, payment.Provider, payment.Amount,
		payment.PaymentDt, payment.Bank, payment.DeliveryCost, payment.GoodsTotal,
		payment.CustomFee, orderId,
	)

	return err
}
func insertDelivery(Tx pgx.Tx, delivery *parser.Delivery, orderId string) error {
	query := `INSERT INTO delivery (
		name, phone, zip, city,
		address, region, email, order_uid
		)

VALUES (
	  $1, $2, $3, $4,
	  $5, $6, $7, $8
	  )
`

	_, err := Tx.Exec(context.Background(), query,
		delivery.Name, delivery.Phone, delivery.Zip, delivery.City,
		delivery.Address, delivery.Region, delivery.Email, orderId,
	)

	return err
}
func insertItem(Tx pgx.Tx, item *parser.Item, orderId string) error {
	query := `	
	WITH inserted_item AS (
  		INSERT INTO item(
							chrt_id, track_number, price, rid,
							 name, sale, size, total_price,
							 nm_id, brand, status
							) 
		VALUES (
				$1, $2, $3, $4,
				$5, $6, $7, $8,
				$9, $10, $11
				) 
		RETURNING id
	)
	INSERT INTO orders_item (item_id, order_uid)
	SELECT id, $12 FROM inserted_item;
`

	_, err := Tx.Exec(context.Background(), query,
		item.ChrtId, item.TrackNumber, item.Price, item.Rid,
		item.Name, item.Sale, item.Size, item.TotalPrice,
		item.NmId, item.Brand, item.Status, orderId,
	)

	return err
}
