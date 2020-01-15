package models

import "strconv"

//RouteInfo is used to return api route information
type RouteInfo struct {
	ID         string `bson:"_id" json:"omitempty"`
	Scope      string `bson:"scope"`
	Controller string `bson:"controller"`
	Endpoint   string `bson:"endpoint"`
	Method     string `bson:"method"`
	Precedence int    `bson:"precedence"`
}

//SetID is used to create unique ID for route
func (r *RouteInfo) SetID() {
	r.ID = r.Scope + "." + r.Controller + "." + r.Endpoint + "." + r.Method + "." + strconv.Itoa(r.Precedence)
}
