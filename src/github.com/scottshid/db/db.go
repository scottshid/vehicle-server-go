package db

import "gopkg.in/mgo.v2"

var Database *mgo.Database
var Session *mgo.Session

func init() {
    session, error := mgo.Dial("localhost:27017")
    if error != nil {
        panic(error)
    }
    Session = session
    Database = session.DB("dev")
}

func GetDB() *mgo.Database {
    return Database
}

func GetSession() *mgo.Session {
    return Session
}
