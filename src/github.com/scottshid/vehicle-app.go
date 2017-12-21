package main

import (
    "github.com/gorilla/mux"
    "log"
    "net/http"
    "github.com/scottshid/user"
    "github.com/scottshid/auth"
    "github.com/scottshid/app"
    "github.com/scottshid/vehicle"
)

func HandleTest(w http.ResponseWriter, r *http.Request) *app.RestResponse {
    m := make(map[string]interface{})
    m["Key"] = "Success"
    return &app.RestResponse{Error: nil, Code: http.StatusOK, Payload: m}
}

func main() {
    router := mux.NewRouter()
    router.HandleFunc("/user", user.HandleCreateUser).Methods("POST")
    router.HandleFunc("/login", user.HandleAuthenticate).Methods("POST")
    router.HandleFunc("/auth", auth.ValidateMiddleware(HandleTest))
    router.HandleFunc("/vehicle/make", auth.ValidateMiddleware(vehicle.HandleGetVehicleMakes)).Methods("GET")
    router.HandleFunc("/vehicle/make/{make}/model", auth.ValidateMiddleware(vehicle.HandleGetVehicleModels)).Methods("GET")
    log.Fatal(http.ListenAndServe(":12345", router))
}
