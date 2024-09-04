package main

import "github.com/wilfridterry/contact-list/internal/app"

//	@title			Swagger Contacts API
//	@version		1.0
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	app.Run()
}
