module gateway-service

go 1.23.5

require (
	github.com/dharmasatrya/goodkarma/donation-service v0.0.0-20250120163630-0be5ef3fd993
	github.com/dharmasatrya/goodkarma/event-service v0.0.0-00010101000000-000000000000
	github.com/dharmasatrya/goodkarma/karma-service v0.0.0-00010101000000-000000000000
	github.com/dharmasatrya/goodkarma/payment-service v0.0.0-20250120132112-a701142f86ee
	github.com/dharmasatrya/goodkarma/user-service v0.0.0-20250120085545-5b1cef43f774
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/joho/godotenv v1.5.1
	github.com/labstack/echo/v4 v4.13.3
	github.com/swaggo/echo-swagger v1.4.1
	google.golang.org/grpc v1.69.4
	google.golang.org/protobuf v1.36.2
)

require (
	github.com/KyleBanks/depth v1.2.1 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/spec v0.20.4 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/swaggo/files/v2 v2.0.0 // indirect
	github.com/swaggo/swag v1.8.12 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	go.mongodb.org/mongo-driver v1.17.2 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	golang.org/x/time v0.8.0 // indirect
	golang.org/x/tools v0.21.1-0.20240508182429-e35e4ccd0d2d // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241015192408-796eee8c2d53 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/dharmasatrya/goodkarma/donation-service => ../donation-service

replace github.com/dharmasatrya/goodkarma/payment-service => ../payment-service

replace github.com/dharmasatrya/goodkarma/event-service => ../event-service

replace github.com/dharmasatrya/goodkarma/user-service => ../user-service

replace github.com/dharmasatrya/goodkarma/karma-service => ../karma-service
