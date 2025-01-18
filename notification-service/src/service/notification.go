package service

import (
	helper "goodkarma-notification-service/helpers"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MailService struct {
	channel *amqp.Channel
}

func NewMailService(ch *amqp.Channel) MailService {
	return MailService{
		channel: ch,
	}
}

func (m MailService) SendRegistrasiEmailNotification(q amqp.Queue) {
	msgs, err := m.channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatal(err)
	}

	for d := range msgs {
		log.Printf("\033[36mNEW MESSAGE:\033[0m %s", d.Body)

		userData := helper.AssertJsonToUserStruct(d.Body)
		helper.SendRegistrasiEmailNotification(userData)
	}
}

func (m MailService) SendInvoiceEmailNotification(q amqp.Queue) {
	msgs, err := m.channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatal(err)
	}

	for d := range msgs {
		log.Printf("\033[36mNEW MESSAGE:\033[0m %s", d.Body)

		data := helper.AssertJsonToInvoiceStruct(d.Body)
		helper.SendInvoiceEmailNotification(data)
	}
}

func (m MailService) SendGoodslNotification(q amqp.Queue) {
	msgs, err := m.channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatal(err)
	}

	for d := range msgs {
		log.Printf("\033[36mNEW MESSAGE:\033[0m %s", d.Body)

		data := helper.AssertJsonToGoodsStruct(d.Body)
		helper.SendGoodslNotification(data)
	}
}

func (m MailService) SendGoodsArrivalNotification(q amqp.Queue) {
	msgs, err := m.channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatal(err)
	}

	for d := range msgs {
		log.Printf("\033[36mNEW MESSAGE:\033[0m %s", d.Body)

		data := helper.AssertJsonToGoodsStruct(d.Body)
		helper.SendGoodsArrivalNotification(data)
	}
}
