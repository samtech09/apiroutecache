package apiroutecache

import (
	"fmt"
	"mahendras/common/apiroutecache/models"
	"testing"
)

var ses *MongoSession

func init() {
	// initialize config
	config := MongoConfig{}
	config.DbName = "testdb"
	config.Host = "192.168.60.206"
	config.Port = 27017
	config.User = ""
	config.Pwd = ""

	ses = InitSession(config)
}

func Test_SaveRoutes(t *testing.T) {
	fmt.Println("\n\nTestSaveRoutes ***")

	var routes []models.RouteInfo
	routes = append(routes,
		NewRouteInfo("ADMIN", "Test", "Getall", "GET", 1),
		NewRouteInfo("USER", "Test", "Getall", "GET", 2),
		NewRouteInfo("ADMIN", "Test", "GetList", "GET", 1))

	err := ses.SaveRoutes(&routes)
	if err != nil {
		t.Fail()
		t.Errorf("SaveReoutes failed. %v", err)
	}
}

func Test_GetAll(t *testing.T) {
	fmt.Println("\n\nTestGetAll ***")

	routes, err := ses.GetAllRoutes()
	if err != nil {
		t.Fail()
	}

	if len(*routes) != 3 {
		t.Fail()
	}
}

func Test_FindRoutes(t *testing.T) {
	fmt.Println("\n\nTestFindRoutes ***")

	routes, err := ses.FindRoutes("Test", "Getall", "GET")
	if err != nil {
		t.Fail()
	}

	if len(*routes) != 2 {
		t.Errorf("Records found: %d\n", len(*routes))
		t.Fail()
	}
}

func Test_GetScopesFromRoute(t *testing.T) {
	fmt.Println("\n\nTestGetScopesFromRoute ***")

	scopes, err := ses.GetScopesFromRoute("Test", "Getall", "GET")
	if err != nil {
		t.Fail()
	}

	if len(*scopes) != 2 {
		t.Errorf("Records found: %d\n", len(*scopes))
		t.Fail()
	}
}

func Test_Cleanup(t *testing.T) {
	ses.Cleanup()
}
