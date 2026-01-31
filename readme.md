Initialize Project
go mod init go-cashier
go build

Run App (debug)
GODEBUG=netdns=go+v4 go run main.go

Test API
GET /api/products
POST /api/products
GET /api/products/:id
PUT /api/products/:id
DELETE /api/products/:id