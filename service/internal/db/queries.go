package db

const (
	orderInstertQuery = `
			INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, 
			                    delivery_service, shardkey, sm_id, date_created, oof_shard) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id;
			`
	deliveryInstertQuery = `
			INSERT INTO deliveries (order_id, name, phone, zip, city, address, region, email ) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id; 
			`
	paymentInstertQuery = `
			INSERT INTO payments (order_id, transaction,  request_id, currency, provider, amount, payment_dt, 
			                    bank, delivery_cost, goods_total, custom_fee) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id;
			`
	itemInstertQuery = `
			INSERT INTO items (order_id, chrt_id, track_number, price, rid, name, sale, 
			                    size, total_price, nm_id, brand, status) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id;
			`

	orderSelectQuery    = `SELECT * FROM orders WHERE id = $1;`
	deliverySelectQuery = `SELECT * FROM deliveries WHERE order_id = $1;`
	paymentSelectQuery  = `SELECT * FROM payments WHERE order_id = $1;`
	itemSelecttQuery    = `SELECT * FROM items WHERE order_id = $1;`

	AllOrdersSelectQuery = `SELECT * FROM orders;`
)
