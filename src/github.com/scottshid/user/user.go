package user

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/bitly/go-simplejson"
	"fmt"
	"log"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
    "errors"
)

var _session *mgo.Session

type User struct {
	Username *string
	Password *string
	ID bson.ObjectId  `bson:"_id,omitempty"`
}

func init() {
	session, error := mgo.Dial("localhost:27017")
	if error != nil {
		panic(error)
	}

	_session = session
	fmt.Println("finished init")
}

func GetUserFromRequest(req *http.Request) (error, string) {
    username := req.Context().Value("username").(string)
    if len(username) > 0 {
        return errors.New("No user in context"), ""
    }
    return nil, username
}

func getUser(writer http.ResponseWriter, req *http.Request) (bool, *User) {
	user := new(User)
	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()
	error := decoder.Decode(&user); if error != nil {
		http.Error(writer, error.Error(), http.StatusInternalServerError)
		return false, nil
	}
	if user.Username == nil || user.Password == nil {
		json := simplejson.New()
		json.Set("Error", "Username and password are required")
		payload, err := json.MarshalJSON(); if err != nil {
			fmt.Println("Error marshallings")
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return false, nil
		}
		http.Error(writer, string(payload), http.StatusBadRequest)
		return false, nil
	}
	return true, user
}



func HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	success, user := getUser(w, r); if !success {
		return
	}

	c := _session.DB("dev").C("users")
	count, err := c.Find(bson.M{"username": bson.M{"$exists": true, "$eq": user.Username } } ).Count();	if err != nil {
		fmt.Println("Error querying for user")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count > 0 {
		json := simplejson.New()
		json.Set("Error", "Username already exists")
		payload, err := json.MarshalJSON(); if err != nil {
			fmt.Println("Error marshalling json reponse")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(payload)
	} else {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*user.Password), 10); if err != nil {
			fmt.Println("Could not hash password")
			http.Error(w,
				"Could not hash your password. We don't store unhashed passwords since that is bad!",
				http.StatusInternalServerError)
			return
		}
		p := string(hashedPassword)
		newUser := User{
			Username: user.Username,
			Password: &p,
			ID: bson.NewObjectId(),
		}
		insertErr := c.Insert(newUser); if insertErr != nil {
			fmt.Println("Could not insert error")
			http.Error(w, "Could not insert new user", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func CreateToken(username string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
	})
	signedToken, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
		panic(error)
	}
	return signedToken
}

func HandleAuthenticate(w http.ResponseWriter, req *http.Request) {
	success, user := getUser(w, req); if !success {
		return
	}
	dbUser := new(User)
	c := _session.DB("dev").C("users")
	err := c.Find(bson.M{"username": user.Username} ).One(dbUser);	if err != nil {
		fmt.Println("Error querying for user")
		json := simplejson.New()
		json.Set("Error", "Username not found")
		payload, err := json.MarshalJSON(); if err != nil {
			fmt.Println("Error marshalling json reponse")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(payload)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	invalidPassword := bcrypt.CompareHashAndPassword([]byte(*dbUser.Password), []byte(*user.Password))
	if invalidPassword != nil {
		fmt.Println("Password invalid for user " + *user.Username)
		http.Error(w, "Invalid Password", http.StatusUnauthorized)
		return
	}

	token := CreateToken(*user.Username)

	json := simplejson.New()
	json.Set("token", token)
	payload, err := json.MarshalJSON(); if err != nil {
		fmt.Println("Error marshalling json reponse")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/user", HandleCreateUser).Methods("POST")
	router.HandleFunc("/authorize", HandleAuthenticate).Methods("POST")
	log.Fatal(http.ListenAndServe(":12345", router))
}