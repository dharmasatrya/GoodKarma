package main

import (
	"goodkarma-notification-service/config"
	"goodkarma-notification-service/src/service"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	conn, mbChannel := config.InitMbChannel()
	defer conn.Close()

	regist_queue := config.InitMbQueue(mbChannel, os.Getenv("EMAIL_QUEUE_NAME"))
	invoice_xendit_queue := config.InitMbQueue(mbChannel, os.Getenv("INVOICE_QUEUE_NAME"))
	send_goods_queue := config.InitMbQueue(mbChannel, os.Getenv("SEND_GOODS_QUEUE_NAME"))
	arrival_queue := config.InitMbQueue(mbChannel, os.Getenv("ARRIVAL_QUEUE_NAME"))

	mailService := service.NewMailService(mbChannel)

	go mailService.SendRegistrasiEmailNotification(regist_queue)
	go mailService.SendInvoiceEmailNotification(invoice_xendit_queue)
	go mailService.SendGoodslNotification(send_goods_queue)
	go mailService.SendGoodsArrivalNotification(arrival_queue)

	var forever chan struct{}
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
