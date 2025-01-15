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
		Message: fmt.Sprintf("Success send registration to email: %v!", req.GetEmail()),
	}

	return &webResponse, nil
}

func SendInvoiceEmailNotification(req *pb.InvoiceData) (*pb.WebResponse, error) {
	// Kirim notifikasi email
	to := req.GetEmail()
	subject := fmt.Sprintf("Pembayaran Donasi GoodKarma", req.GetName())

	// Format konten email dengan informasi pembayaran dan pemesanan secara rinci
	content := fmt.Sprintf(`
	Yth. %s,

	Terima kasih telah melakukan transaksi pembayaran kepada kami! Berikut adalah detailnya:

	- Tanggal Pembuatan Transaksi: %v

	Detail Pembayaran:
	- Status Pembayaran: %s
	- ID Faktur: %s
	- Deskripsi: %s
	- Tautan Pembayaran: %s
	- Jumlah: Rp. %v

	Silakan simpan email ini sebagai referensi Anda. Jika Anda memiliki pertanyaan atau membutuhkan bantuan, jangan ragu untuk menghubungi kami.

	Terima Kasih,
	Tim GoodKarma
	`,
		req.GetName(),
		time.Now().Format("2 Januari 2006, 15:04 WIB"),
		req.GetStatus(),
		req.GetInvoiceId(),
		req.GetDescription(),
		req.GetLink(),
		req.GetAmmount())

	// Kirim email
	utils.SendEmailNotification(to, subject, content)

	webResponse := pb.WebResponse{
		Message: fmt.Sprintf("Berhasil mengirim faktur ke email :%v!", req.GetEmail()),
	}

	return &webResponse, nil
}

func SendGoodslNotification(req *pb.SendGoodsData) (*pb.WebResponse, error) {
	// Send email notification
	to := req.GetEmail()
	subject := fmt.Sprintf("Pengiriman Barang GoodKarma", req.GetName())

	// Format the content to include detailed payment and booking information
	content := fmt.Sprintf(`
	Halo %s,

	Terima kasih atas donasi anda, barang anda dengan detail:
	
	- Tanggal Pengiriman: %v

	Barang:
	- Status Pengiriman: %s
	- Alamat: %s
	- Ammount: Rp. %v

	Sedang dalam pengiriman, kami akan kabarkan melewati email ini jika barang anda sudah sampai.

	Terima Kasih,  
	GoodKarma Team
	`,
		req.GetName(),
		time.Now().Format("January 2, 2006, 3:04 PM"),
		req.GetStatus(),
		req.GetAlamat(),
		req.GetAmmount())

	// Send the email
	utils.SendEmailNotification(to, subject, content)

	webResponse := pb.WebResponse{
		Message: fmt.Sprintf("Success sending the send notification to email: %v!", req.GetEmail()),
	}

	return &webResponse, nil
}

func SendGoodsArrivalNotification(req *pb.SendGoodsData) (*pb.WebResponse, error) {
	// Send email notification
	to := req.GetEmail()
	subject := fmt.Sprintf("Pengiriman Barang GoodKarma Telah Tiba", req.GetName())

	// Format the content to include detailed payment and booking information
	content := fmt.Sprintf(`
	Halo %s,

	Terima kasih atas donasi anda, kami ingin menyampaikan bahwa barang anda dengan detail:
	
	- Tanggal Pengiriman: %v

	Barang:
	- Status Pengiriman: %s
	- Alamat: %s
	- Ammount: Rp. %v

	Sudah tiba pada alamat yang tertera dan sedang diproses oleh tim kami.

	Terima Kasih,  
	GoodKarma Team
	`,
		req.GetName(),
		time.Now().Format("January 2, 2006, 3:04 PM"),
		req.GetStatus(),
		req.GetAlamat(),
		req.GetAmmount())

	// Send the email
	utils.SendEmailNotification(to, subject, content)

	webResponse := pb.WebResponse{
		Message: fmt.Sprintf("Success sending the arrival notification to email: %v!", req.GetEmail()),
	}

	return &webResponse, nil
}
