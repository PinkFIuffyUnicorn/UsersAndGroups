package main

import (
	"log"
	"os"

	"github.com/alexflint/go-arg"

	"github.com/UsersAndGroups/swagger/restapi"
	"github.com/UsersAndGroups/swagger/restapi/operations"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
)

type cliArgs struct {
	Port int `arg:"-p,help:port to listen to"`
}

var (
	args = &cliArgs{
		Port: 8080,
	}
)

func main() {
	arg.MustParse(args)

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewHelloAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	server.Port = args.Port

	api.GetHostnameHandler = operations.GetHostnameHandlerFunc(
		func(params operations.GetHostnameParams) middleware.Responder {
			response, err := os.Hostname()
			if err != nil {
				log.Fatalln(err)
			}

			return operations.NewGetHostnameOK().WithPayload(response)
		})

	// Start listening using having the handlers and port
	// already set up.
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
