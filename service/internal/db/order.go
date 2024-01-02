package db

import (
	"fmt"

	"github.com/Nahbox/streamed-order-viewer/service/internal/models"
)

var ErrNoMatch = fmt.Errorf("no matching record")

func (db *Database) GetAllOrders() ([]models.OrderDetails, error) {
	var orders []models.Order

	ordersRows, err := db.conn.Query(AllOrdersSelectQuery)
	if err != nil {
		return nil, err
	}
	for ordersRows.Next() {
		var order models.Order
		err := ordersRows.Scan(&order.Id, &order.OrderUID,
			&order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature, &order.CustomerID,
			&order.DeliveryService, &order.ShardKey, &order.SmID, &order.DateCreated, &order.OofShard)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	OrdersDetails := make([]models.OrderDetails, 0, len(orders))

	for _, order := range orders {
		var OrderDetails models.OrderDetails

		delivery, payment, items, err := db.GetOrderDataById(order.Id)
		if err != nil {
			return nil, err
		}
		OrderDetails.Order = order
		OrderDetails.Delivery = *delivery
		OrderDetails.Payment = *payment
		OrderDetails.Items = items

		OrdersDetails = append(OrdersDetails, OrderDetails)
	}

	return OrdersDetails, nil
}

func (db *Database) GetOrderDataById(orderId int) (*models.Delivery, *models.Payment, []models.Item, error) {
	var delivery models.Delivery
	var payment models.Payment
	var items []models.Item

	tx, err := db.conn.Begin()
	if err != nil {
		return nil, nil, nil, err
	}
	defer tx.Rollback()

	err = tx.QueryRow(deliverySelectQuery, orderId).
		Scan(&delivery.Id, &delivery.OrderId, &delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Address,
			&delivery.Region, &delivery.Email)
	if err != nil {
		return nil, nil, nil, err
	}

	err = tx.QueryRow(paymentSelectQuery, orderId).
		Scan(&payment.Id, &payment.OrderId, &payment.Transaction, &payment.RequestID, &payment.Currency, &payment.Provider,
			&payment.Amount, &payment.PaymentDt, &payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee)
	if err != nil {
		return nil, nil, nil, err
	}

	itemsRows, err := tx.Query(itemSelecttQuery, orderId)
	if err != nil {
		return nil, nil, nil, err
	}
	for itemsRows.Next() {
		var item models.Item
		err := itemsRows.Scan(&item.Id, &item.OrderId, &item.ChrtID, &item.TrackNumber, &item.Price, &item.RID,
			&item.Name, &item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
		if err != nil {
			return nil, nil, nil, err
		}
		items = append(items, item)
	}

	return &delivery, &payment, items, tx.Commit()
}

func (db *Database) AddOrder(orderDetails *models.OrderDetails) error {
	tx, err := db.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var id int
	err = tx.QueryRow(orderInstertQuery, orderDetails.Order.OrderUID, orderDetails.Order.TrackNumber,
		orderDetails.Order.Entry, orderDetails.Order.Locale, orderDetails.Order.InternalSignature,
		orderDetails.Order.CustomerID, orderDetails.Order.DeliveryService, orderDetails.Order.ShardKey,
		orderDetails.Order.SmID, orderDetails.Order.DateCreated, orderDetails.Order.OofShard).Scan(&id)
	if err != nil {
		return err
	}
	orderDetails.Order.Id = id

	err = tx.QueryRow(deliveryInstertQuery, orderDetails.Order.Id, orderDetails.Delivery.Name, orderDetails.Delivery.Phone,
		orderDetails.Delivery.Zip, orderDetails.Delivery.City, orderDetails.Delivery.Address,
		orderDetails.Delivery.Region, orderDetails.Delivery.Email).Scan(&id)
	if err != nil {
		return err
	}
	orderDetails.Delivery.Id = id

	err = tx.QueryRow(paymentInstertQuery, orderDetails.Order.Id, orderDetails.Payment.Transaction, orderDetails.Payment.RequestID,
		orderDetails.Payment.Currency, orderDetails.Payment.Provider, orderDetails.Payment.Amount,
		orderDetails.Payment.PaymentDt, orderDetails.Payment.Bank, orderDetails.Payment.DeliveryCost,
		orderDetails.Payment.GoodsTotal, orderDetails.Payment.CustomFee).Scan(&id)
	if err != nil {
		return err
	}
	orderDetails.Payment.Id = id

	for i := 0; i < len(orderDetails.Items); i++ {
		err = tx.QueryRow(itemInstertQuery, orderDetails.Order.Id, orderDetails.Items[i].ChrtID, orderDetails.Items[i].TrackNumber,
			orderDetails.Items[i].Price, orderDetails.Items[i].RID, orderDetails.Items[i].Name, orderDetails.Items[i].Sale,
			orderDetails.Items[i].Size, orderDetails.Items[i].TotalPrice, orderDetails.Items[i].NmID,
			orderDetails.Items[i].Brand, orderDetails.Items[i].Status).Scan(&id)
		if err != nil {
			return err
		}
		orderDetails.Items[i].Id = id
	}

	return tx.Commit()
}

func (db *Database) GetOrderById(Id int) (*models.OrderDetails, error) {
	tx, err := db.conn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	orderDetails := models.OrderDetails{}

	err = tx.QueryRow(orderSelectQuery, Id).
		Scan(&orderDetails.Order.Id, &orderDetails.Order.OrderUID, &orderDetails.Order.TrackNumber,
			&orderDetails.Order.Entry, &orderDetails.Order.Locale, &orderDetails.Order.InternalSignature,
			&orderDetails.Order.CustomerID, &orderDetails.Order.DeliveryService, &orderDetails.Order.ShardKey,
			&orderDetails.Order.SmID, &orderDetails.Order.DateCreated, &orderDetails.Order.OofShard)
	if err != nil {
		return nil, err
	}

	err = tx.QueryRow(deliverySelectQuery, orderDetails.Order.Id).
		Scan(&orderDetails.Delivery.Id, &orderDetails.Delivery.OrderId, &orderDetails.Delivery.Name, &orderDetails.Delivery.Phone,
			&orderDetails.Delivery.Zip, &orderDetails.Delivery.City, &orderDetails.Delivery.Address,
			&orderDetails.Delivery.Region, &orderDetails.Delivery.Email)
	if err != nil {
		return nil, err
	}

	err = tx.QueryRow(paymentSelectQuery, orderDetails.Order.Id).
		Scan(&orderDetails.Payment.Id, &orderDetails.Delivery.OrderId, &orderDetails.Payment.Transaction, &orderDetails.Payment.RequestID,
			&orderDetails.Payment.Currency, &orderDetails.Payment.Provider, &orderDetails.Payment.Amount,
			&orderDetails.Payment.PaymentDt, &orderDetails.Payment.Bank, &orderDetails.Payment.DeliveryCost,
			&orderDetails.Payment.GoodsTotal, &orderDetails.Payment.CustomFee)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(orderDetails.Items); i++ {
		err = tx.QueryRow(itemSelecttQuery, orderDetails.Order.Id).
			Scan(&orderDetails.Items[i].Id, &orderDetails.Delivery.OrderId, &orderDetails.Items[i].ChrtID, &orderDetails.Items[i].TrackNumber,
				&orderDetails.Items[i].Price, &orderDetails.Items[i].RID, &orderDetails.Items[i].Name,
				&orderDetails.Items[i].Sale, &orderDetails.Items[i].Size, &orderDetails.Items[i].TotalPrice,
				&orderDetails.Items[i].NmID, &orderDetails.Items[i].Brand, &orderDetails.Items[i].Status)
		if err != nil {
			return nil, err
		}
	}

	return &orderDetails, tx.Commit()
}
