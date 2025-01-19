package helpers

import (
	"encoding/json"
	"goodkarma-notification-service/entity"
	"log"
)

func AssertJsonToUserStruct(body []byte) entity.UserRegistData {
	var link entity.UserRegistData
	if err := json.Unmarshal(body, &link); err != nil {
		log.Fatal(err)
	}
	return link
}

func AssertJsonToInvoiceStruct(body []byte) entity.InvoiceData {
	var invoice entity.InvoiceData
	if err := json.Unmarshal(body, &invoice); err != nil {
		log.Fatal(err)
	}
	return invoice
}

func AssertJsonToGoodsStruct(body []byte) entity.GoodsData {
	var goods entity.GoodsData
	if err := json.Unmarshal(body, &goods); err != nil {
		log.Fatal(err)
	}
	return goods
}
