package data

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rickj1ang/RRS/internal/mailer"
)

type NotifyQue struct {
	Conn *amqp.Connection
}

const DELAYED_QUEUE = "work.later"
const DESTINATION_QUEUE = "work.now"
const DELAY_TIME = 5000

func (n NotifyQue) Publish(address string) {
	conn, err := amqp.Dial("amqps://oqwhsvhx:pbtoNpE6D7Xiwcns1W2V-za6R4ZMNzbh@gerbil.rmq.cloudamqp.com/oqwhsvhx")
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	defer ch.Close()

	args := amqp.Table{
		"x-dead-letter-exchange":    "",
		"x-dead-letter-routing-key": DESTINATION_QUEUE,
		"x-message-ttl":             5000,
	}

	q, err := ch.QueueDeclare(DELAYED_QUEUE, false, false, false, false, args)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(address),
		})
}

func (n NotifyQue) Subscribe() {
	ch, err := n.Conn.Channel()
	if err != nil {
		panic(err)
	}

	defer ch.Close()
	q, err := ch.QueueDeclare(
		DESTINATION_QUEUE, // name
		false,             // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		panic(err)
	}
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		panic(err)
	}

	var forever chan struct{}
	go func() {
		for d := range msgs {
			mailer.Send(string(d.Body))
			d.Ack(false)
		}
	}()
	<-forever
}
