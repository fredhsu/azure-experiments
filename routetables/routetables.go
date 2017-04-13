package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/arm/examples/helpers"
	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

func withInspection() autorest.PrepareDecorator {
	return func(p autorest.Preparer) autorest.Preparer {
		return autorest.PreparerFunc(func(r *http.Request) (*http.Request, error) {
			fmt.Printf("Inspecting Request: %s %s\n", r.Method, r.URL)
			return p.Prepare(r)
		})
	}
}

func byInspecting() autorest.RespondDecorator {
	return func(r autorest.Responder) autorest.Responder {
		return autorest.ResponderFunc(func(resp *http.Response) error {
			fmt.Printf("Inspecting Response: %s for %s %s\n", resp.Status, resp.Request.Method, resp.Request.URL)
			return r.Respond(resp)
		})
	}
}

func main() {

	c := map[string]string{
		"AZURE_CLIENT_ID":       os.Getenv("AZURE_CLIENT_ID"),
		"AZURE_CLIENT_SECRET":   os.Getenv("AZURE_CLIENT_SECRET"),
		"AZURE_SUBSCRIPTION_ID": os.Getenv("AZURE_SUBSCRIPTION_ID"),
		"AZURE_TENANT_ID":       os.Getenv("AZURE_TENANT_ID")}
	if err := checkEnvVar(&c); err != nil {
		log.Fatalf("Error: %v", err)
		return
	}
	fmt.Println("check")
	spt, err := helpers.NewServicePrincipalTokenFromCredentials(c, azure.PublicCloud.ResourceManagerEndpoint)
	if err != nil {
		log.Fatalf("Error: %v", err)
		return
	}
	fmt.Println("token")

	rt := network.NewRouteTablesClient(c["AZURE_SUBSCRIPTION_ID"])
	rt.Authorizer = spt
	rtlist, err := rt.ListAll()
	if err != nil {
		log.Fatalf("Error: %v", err)
		return
	}
	fmt.Printf("%+v\n", (*rtlist.Value)[0])
	rtentry := (*rtlist.Value)[0]
	fmt.Printf("route table name: %s\n", *rtentry.Name)

	rgn := "test"
	rc := network.NewRoutesClient(c["AZURE_SUBSCRIPTION_ID"])
	rc.Authorizer = spt
	routeslist, err := rc.List(rgn, *rtentry.Name)
	if err != nil {
		log.Fatalf("Error: %v", err)
		return
	}
	route := (*routeslist.Value)[0]
	fmt.Println(route)
	fmt.Printf("route name: %s\n", *route.Name)
	fmt.Printf("rpf: %+v\n", route.RoutePropertiesFormat)
}

func checkEnvVar(envVars *map[string]string) error {
	var missingVars []string
	for varName, value := range *envVars {
		if value == "" {
			missingVars = append(missingVars, varName)
		}
	}
	if len(missingVars) > 0 {
		return fmt.Errorf("Missing environment variables %v", missingVars)
	}
	return nil
}
