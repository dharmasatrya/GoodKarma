package service

import (
	"fmt"
	pb "goodkarma-notification-service/pb"
	"time"

	"goodkarma-notification-service/utils"
)

type NotificationService struct {
	pb.UnimplementedNotificationServiceServer
}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

func SendRegistrasiEmailNotification(req *pb.RegistrationData) (*pb.WebResponse, error) {
	to := req.GetEmail()
	confirmationLink := req.GetLink()

	subject := "Register GoodKarma Successfully"
	content := fmt.Sprintf("Your register in our website is success! Please use this link %v to activate your account!", confirmationLink)
	utils.SendEmailNotification(to, subject, content)

	webResponse := pb.WebResponse{
		Message: "Success send registration email!",
	}

	return &webResponse, nil
}

func SendInvoiceEmailNotification(req *pb.InvoiceData) (*pb.WebResponse, error) {
	// Send email notification
	to := req.GetEmail()
	subject := fmt.Sprintf("Payment for Top Up %v is Successful", req.GetName())

	// Format the content to include detailed payment and booking information
	content := fmt.Sprintf(`
	Dear %s,

	Thank you for making your payment with us! Your booking has been successfully processed. Here are the details:

	- Payment Created Date: %v

	Payment Details:
	- Payment Status: %s
	- Invoice ID: %s
	- Description: %s
	- Payment Link: %s
	- Ammount: Rp. %v

	Please save this email for your reference. If you have any questions or require assistance, feel free to reach out to us.

	Best regards,  
	Your Booking Team
	`,
		req.GetName(),
		time.Now().Format("January 2, 2006, 3:04 PM"),
		req.GetStatus(),
		req.GetInvoiceId(),
		req.GetDescription(),
		req.GetLink(),
		req.GetAmmount())

	// Send the email
	utils.SendEmailNotification(to, subject, content)

	webResponse := pb.WebResponse{
		Message: "Success send registration email!",
	}

	return &webResponse, nil
}
