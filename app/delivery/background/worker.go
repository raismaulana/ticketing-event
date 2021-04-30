package background

import (
	"log"
	"math/rand"
	"strings"

	"github.com/raismaulana/ticketing-event/app/config"
	"github.com/raismaulana/ticketing-event/app/helper"
	"github.com/raismaulana/ticketing-event/app/usecase"
	"github.com/streadway/amqp"
)

const (
	REMINDER_QUEUE_NAME  = "reminder"
	PROMOTION_QUEUE_NAME = "promotion"
)

type BackgroundTask interface {
	SendReminderPayment()
	SendPromotionEvent()
	ListenerReminderPayment()
	ListenerPromotionEvent()
}

type backgroundTask struct {
	backgroundCase usecase.BackgroundCase
	ch             *amqp.Channel
}

func NewBackgroundTask(backgroundCase usecase.BackgroundCase, conn *amqp.Connection) BackgroundTask {
	ch, err := conn.Channel()
	config.FailOnError(err, "Failed to open a channel")
	return &backgroundTask{
		backgroundCase: backgroundCase,
		ch:             ch,
	}
}
func setupQueue(ch *amqp.Channel, name string) amqp.Queue {
	q, err := ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	config.FailOnError(err, "Failed to declare a queue")
	return q
}

func (b *backgroundTask) SendReminderPayment() {
	log.Println("Reminder start for today")

	q := setupQueue(b.ch, REMINDER_QUEUE_NAME)

	res, err := b.backgroundCase.GetPendingTransaction()
	if err != nil {
		log.Println("worker.go : ", err)
	}
	for _, v := range res {
		err = b.ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(v.Email),
			})
		config.FailOnError(err, "Failed to publish a message")
	}
}

func (b *backgroundTask) SendPromotionEvent() {
	log.Println("Promotion start for today")
	events, users, err := b.backgroundCase.GetPromotionEvent()
	q := setupQueue(b.ch, PROMOTION_QUEUE_NAME)

	if err == nil && len(events) > 0 && len(users) > 0 {
		for _, v := range users {
			i := rand.Intn(len(events))
			event := events[i]
			msg := v.Email + "/" + event.TitleEvent
			err = b.ch.Publish(
				"",     // exchange
				q.Name, // routing key
				false,  // mandatory
				false,  // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(msg),
				})
			config.FailOnError(err, "Failed to publish a message")
		}
	}

}

func (b *backgroundTask) ListenerReminderPayment() {
	defer b.ch.Close()
	q := setupQueue(b.ch, REMINDER_QUEUE_NAME)
	msgs, err := b.ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	config.FailOnError(err, "Failed to register a consumer")

	done := make(chan bool)

	go func() {
		for d := range msgs {
			helper.SendMail(string(d.Body), "Reminder Payment", "v.Eid will be held tommorow, complete your payment now to get the webinar's link.")
			log.Println("send reminder payment to ", string(d.Body))
		}
	}()

	log.Printf(" ListenerReminderPayment is starting in background")
	<-done
}

func (b *backgroundTask) ListenerPromotionEvent() {
	defer b.ch.Close()
	q := setupQueue(b.ch, PROMOTION_QUEUE_NAME)
	msgs, err := b.ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	config.FailOnError(err, "Failed to register a consumer")

	done := make(chan bool)

	go func() {
		for d := range msgs {
			t := strings.Split(string(d.Body), "/")
			helper.SendMail(t[0], "Cool Event Will Be Held", "Join this event "+t[1]+".")

			log.Println("send ", t[1], " to ", t[0])
		}
	}()

	log.Printf(" ListenerPromotionEvent is starting in background")
	<-done
}
