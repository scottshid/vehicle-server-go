package vehicle

import (
    "github.com/scottshid/db"
    "gopkg.in/mgo.v2"
    "net/http"
    "github.com/scottshid/app"
    "gopkg.in/mgo.v2/bson"
    "github.com/gorilla/mux"
)

var vehiclesCollection *mgo.Collection

func init() {
    vehiclesCollection = db.GetDB().C("vehicles").With(db.GetSession())
}

func HandleGetVehicleMakes(writer http.ResponseWriter, req *http.Request) *app.RestResponse {
    var result []string
    vehiclesCollection.Find(bson.M{}).Distinct("make", &result)
    payload := make(map[string]interface{})
    payload["results"] = result
    return &app.RestResponse{Code:http.StatusOK, Error: nil, Payload: payload }
}

func HandleGetVehicleModels(writer http.ResponseWriter, req *http.Request) *app.RestResponse {

    var results []struct{
        Make string
        Model string
    }
    vars := mux.Vars(req)
    vehiclesCollection.Find(bson.M{"make": vars["make"]}).All(&results)
    payload := make(map[string]interface{})
    payload["results"] = results
    return &app.RestResponse{Code: http.StatusOK, Error: nil, Payload: payload}
}
