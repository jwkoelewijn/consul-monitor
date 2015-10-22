package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/consul/api"
)

const ownServiceName = "internal-query-service"

func main() {
	config := api.DefaultConfig()
	client, _ := api.NewClient(config)
	agent := client.Agent()

	members, err := agent.Members(false)
	if err != nil {
		panic(err)
	}

	if len(members) > 0 {
		log.Println("Members:")
		log.Println("==============================================")
	}

	for _, member := range members {
		log.Println(fmt.Sprintf("%+v", member))
	}

	services, err := agent.Services()
	if err != nil {
		panic(err)
	}

	if len(services) > 0 {
		log.Println("Local Services:")
		log.Println("==============================================")
	}

	for name, serv := range services {
		log.Println(fmt.Sprintf("%s:\t\t%+v", name, serv))
	}

	catalog := client.Catalog()
	catalogServiceNames, _, err := catalog.Services(nil)
	if err != nil {
		log.Printf("Could not query services from Catalog: %+v", err)
		return
	}
	if len(catalogServiceNames) > 0 {
		log.Println("Catalog Services:")
		log.Println("==============================================")
	}
	for name, tags := range catalogServiceNames {
		log.Printf("%s:\t\t%s", name, strings.Join(tags, ", "))
		catalogServices, _, err := catalog.Service(name, "", nil)
		if err != nil {
			panic(err)
		}
		for _, catalogService := range catalogServices {
			log.Printf("\t+-->\t%+v", *catalogService)
		}
	}
}
