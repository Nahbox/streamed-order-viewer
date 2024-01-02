package nats_subscriber

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"

	"github.com/nats-io/stan.go"

	"github.com/Nahbox/streamed-order-viewer/service/internal/handler/order"
	"github.com/Nahbox/streamed-order-viewer/service/internal/models"
)

func NatsStreamingSubscribe(ctx context.Context, orderHandler *order.Handler) error {
	// Подключение к серверу NATS
	sc, err := stan.Connect("test-cluster", "subscriber-client", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		return errors.Wrap(err, "stan.Connect")
	}
	defer sc.Close()

	sub, err := sc.Subscribe("test", func(m *stan.Msg) {
		if err = handleMessage(orderHandler, m); err != nil {
			return
		}
	}, stan.DeliverAllAvailable())

	if err != nil {
		return err
	}

	log.Println("nats subscriber is listening...")

	select {
	case <-ctx.Done():
		sub.Unsubscribe()
		sc.Close()
	}

	return nil
}

func handleMessage(orderHandler *order.Handler, m *stan.Msg) error {
	var orderDetails models.OrderDetails
	err := json.Unmarshal(m.Data, &orderDetails)
	if err != nil {
		return err
	}

	if err = orderHandler.DB.AddOrder(&orderDetails); err != nil {
		return err
	}

	byt, err := json.Marshal(orderDetails)
	if err != nil {
		return err
	}
	orderHandler.CacheManager.Cache.SetDefault(orderDetails.Order.OrderUID, byt)

	log.Printf("received a message: Id: %d\n", orderDetails.Order.Id)

	return nil
}
