package apiroutecache

import (
	"fmt"
	"log"
	"mahendras/common/apiroutecache/models"
	"reflect"
	"strconv"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//NewRouteInfo create new RouteInfo struct
func NewRouteInfo(scope, controller, endpoint, method string, precedence int) models.RouteInfo {
	r := models.RouteInfo{
		Scope:      scope,
		Controller: controller,
		Endpoint:   endpoint,
		Method:     method,
		Precedence: precedence,
	}
	//r.ID = bson.NewObjectId()
	r.ID = scope + "." + controller + "." + endpoint + "." + method + "." + strconv.Itoa(precedence)
	//r.SetID()
	return r
}

// InsertRoutes only insert list of given routes to database
//  if route already exist, it will panic
func (m *MongoSession) InsertRoutes(routes *[]models.RouteInfo, setIDs bool) error {
	s := m.Copy()
	defer s.Close()

	if setIDs {
		for _, rout := range *routes {
			rout.SetID()
		}
	}

	err := routesColl(s, m.DBname).Insert(structToInterface(*routes)...)
	if err != nil {
		panic(err)
	}
	return err
}

//SaveRoutes insert or update one or more routes to database
func (m *MongoSession) SaveRoutes(routes *[]models.RouteInfo) error {
	if routes == nil {
		return fmt.Errorf("nothing to save")
	}

	s := m.Copy()
	defer s.Close()

	col := routesColl(s, m.DBname)
	var err error
	for _, rout := range *routes {
		rout.SetID()
		_, err = col.UpsertId(rout.ID, rout)
	}
	if err != nil {
		panic(err)
	}
	return err
}

//TruncateRoutes removes all routes from database
func (m *MongoSession) TruncateRoutes() error {
	s := m.Copy()
	defer s.Close()

	err := routesColl(s, m.DBname).DropCollection()
	// if err != nil {
	// 	g.Logger.Errorf("Routes.Truncate: Error %v", err)
	// }
	return err
}

//DeleteRoutesByScope remove all routes defined in given scope
func (m *MongoSession) DeleteRoutesByScope(scope string) error {
	s := m.Copy()
	defer s.Close()

	err := routesColl(s, m.DBname).Remove(bson.M{"scope": "\"" + scope + "\""})
	// if err != nil {
	// 	g.Logger.Errorf("Routes.DeleteByScope(%s): %v", scope, err)
	// }
	return err
	//mdb().C("routeinfo").Remove(bson.M{"name": "Foo Bar"})
}

//GetAllRoutes returns list of all routes
func (m *MongoSession) GetAllRoutes() (*[]models.RouteInfo, error) {
	s := m.Copy()
	defer s.Close()

	var routes []models.RouteInfo
	err := routesColl(s, m.DBname).Find(nil).All(&routes)
	if err != nil {
		//g.Logger.Errorf("Routes.GetAll: %v", err)
		return nil, err
	}
	return &routes, nil
}

//FindRoutes lookup for route for given controller/endpoint and method
func (m *MongoSession) FindRoutes(controller, endpoint, method string) (*[]models.RouteInfo, error) {
	s := m.Copy()
	defer s.Close()

	query := bson.M{
		"controller": controller,
		"endpoint":   endpoint,
		"method":     method,
	}

	//g.Logger.Debugm("FindRoutes", "Query: %v", query)

	var routes []models.RouteInfo
	err := routesColl(s, m.DBname).Find(query).Sort("precedence").All(&routes)
	if err != nil {
		//g.Logger.Errorf("Routes.FindRoutes([%s] %s/%s): %v", method, controller, endpoint, err)
		return nil, err
	}

	//g.Logger.Debugm("FindRoutes", "QueryResult: %v", &routes)

	return &routes, nil
	//change := bson.M{"$set": bson.M{"phone": "+86 99 8888 7777", "timestamp": time.Now()}}

	//err = c.Find(bson.M{"name": "Ale"}).Select(bson.M{"phone": 0}).One(&result)
	//err = c.Find(bson.M{"name": "Ale"}).Sort("-timestamp").All(&results)
}

//GetScopesFromRoute query and return []string of scoped for given route
func (m *MongoSession) GetScopesFromRoute(controller, endpoint, method string) (*[]string, error) {
	s := m.Copy()
	defer s.Close()

	query := bson.M{
		"controller": controller,
		"endpoint":   endpoint,
		"method":     method,
	}

	//g.Logger.Debugm("GetScopesFromRoute", "Query: %v", query)

	//var scopes []models.RouteInfo
	var scopes []struct {
		Scope string `bson:"scope"`
	}
	err := routesColl(s, m.DBname).Find(query).Select(bson.M{"scope": 1}).All(&scopes)
	//err := coll(s).Find(query).Distinct("scope", &scopes)
	if err != nil {
		//g.Logger.Errorf("Routes.GetScopesFromRoute([%s] %s/%s): %v", method, controller, endpoint, err)
		return nil, err
	}

	//g.Logger.Debugm("GetScopesFromRoute", "QueryResult: %q", &scopes)

	// convert to []string
	var slist []string
	for _, str := range scopes {
		slist = append(slist, str.Scope)
	}
	return &slist, nil

}

func routesColl(ses *mgo.Session, dbname string) *mgo.Collection {
	return ses.DB(dbname).C("routeinfo")
}

// structToInterface converts []struct{} to []interface{}
func structToInterface(array interface{}) []interface{} {

	v := reflect.ValueOf(array)
	t := v.Type()

	if t.Kind() != reflect.Slice {
		log.Panicf("`array` should be %s but got %s", reflect.Slice, t.Kind())
	}

	result := make([]interface{}, v.Len(), v.Len())

	for i := 0; i < v.Len(); i++ {
		result[i] = v.Index(i).Interface()
	}

	return result
}
