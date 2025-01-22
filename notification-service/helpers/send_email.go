package helpers

import (
	"fmt"
	"goodkarma-notification-service/entity"
	"goodkarma-notification-service/utils"
	"time"
)

func SendRegistrasiEmailNotification(data entity.UserRegistData) {
	to := data.Email
	confirmationLink := data.Link

	subject := "Register GoodKarma Successfully"
	content := fmt.Sprintf("Your register in our website is success! Please use this link to activate your account! <a href='%s'>click here</a>", confirmationLink)
	utils.SendEmailNotification(to, subject, content)
}

func SendInvoiceEmailNotification(data entity.InvoiceData) {
	to := data.Email
	subject := fmt.Sprintf("Pembayaran Donasi GoodKarma", data.Name)

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
		data.Name,
		time.Now().Format("2 Januari 2006, 15:04 WIB"),
		data.Status,
		data.InvoiceID,
		data.Description,
		data.Link,
		data.Ammount)

	// Kirim email
	utils.SendEmailNotification(to, subject, content)
}

func SendGoodslNotification(data entity.GoodsData) {
	to := data.Email
	subject := fmt.Sprintf("Pengiriman Barang GoodKarma", data.Name)

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
		data.Name,
		time.Now().Format("January 2, 2006, 3:04 PM"),
		data.Status,
		data.Alamat,
		data.Ammount)

	// Send the email
	utils.SendEmailNotification(to, subject, content)
}

func SendGoodsArrivalNotification(data entity.GoodsData) {
	// Send email notification
	to := data.Email
	subject := fmt.Sprintf("Pengiriman Barang GoodKarma Telah Tiba", data.Name)

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
		data.Name,
		time.Now().Format("January 2, 2006, 3:04 PM"),
		data.Status,
		data.Alamat,
		data.Ammount)

	// Send the email
	utils.SendEmailNotification(to, subject, content)
}
