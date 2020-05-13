package controller

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"qok.com/identity/db"
	"qok.com/identity/logwrapper"
	"qok.com/identity/model"

	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, req *http.Request) {
	logger := logwrapper.Load()
	w.Header().Set("Content-Type", "application/json")

	body, ioerr := ioutil.ReadAll(req.Body)
	if ioerr != nil {
		logger.Fatal("could not read from io")
		return
	}
	var response model.ResponseResult
	var user model.User

	err := json.Unmarshal(body, &user)
	if err != nil {
		logger.Printf("error in umarshalling")
		response.Error = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	collection, err := db.GetDBCollection()

	if err != nil {
		logger.Printf("error in connecting to DB")
		response.Error = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	var queryResult model.User

	err = collection.FindOne(context.TODO(), bson.D{{"username", user.Username}}).Decode(&queryResult)
	if err == nil {
		response.Result = "Username already Exists!!"
		json.NewEncoder(w).Encode(response)
		return
	}

	if err.Error() == "mongo: no documents in result" {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)
		if err != nil {
			response.Error = "Error While Hashing Password, Try Again"
			json.NewEncoder(w).Encode(response)
			return
		}

		user.Password = string(hash)

		_, err = collection.InsertOne(context.TODO(), user)

		if err != nil {
			response.Error = "Error while inserting to DB, try again"
			json.NewEncoder(w).Encode(response)
			return
		}

		//LAST STATE WHICH NOTHING WENT WRONG
		response.Result = "Successfuly Registered"

		json.NewEncoder(w).Encode(response)

	}
	response.Error = err.Error()
	json.NewEncoder(w).Encode(response)

}

func LoginHandler(w http.ResponseWriter, req *http.Request) {
	logger := logwrapper.Load()
	w.Header().Set("Content-Type", "application/json")

	body, readErr := ioutil.ReadAll(req.Body)

	if readErr != nil {
		logger.Fatalf("error in reading io : %v", readErr)
	}

	var user model.User

	err := json.Unmarshal(body, &user)

	if err != nil {
		logger.Fatal(err)
	}

	collection, db_connection_error := db.GetDBCollection()

	if db_connection_error != nil {
		logger.Fatal(db_connection_error)
	}

	var userObjectForResponse model.User
	var res model.ResponseResult

	db_find_error := collection.FindOne(context.TODO(), bson.D{{"username", user.Username}}).Decode(&userObjectForResponse)

	if db_find_error != nil {
		logger.Printf("Invalid username : %q", user.Username)
		res.Error = "Invalid username"
		json.NewEncoder(w).Encode(res)
		return
	}

	password_error := bcrypt.CompareHashAndPassword([]byte(userObjectForResponse.Password), []byte(user.Password))

	if password_error != nil {
		logger.Printf("Password not match for username: %q", user.Username)
		res.Error = "Password does not match"
		json.NewEncoder(w).Encode(res)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": userObjectForResponse.Username,
	})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		logger.Printf("error on generating the token")
		res.Error = "Error while generating token,Try again"
		json.NewEncoder(w).Encode(res)
		return
	}

	userObjectForResponse.Token = tokenString
	userObjectForResponse.Password = ""

	json.NewEncoder(w).Encode(userObjectForResponse)

}

func UserInfoHandler(w http.ResponseWriter, req *http.Request) {
	logger := logwrapper.Load()
	w.Header().Set("Content-Type", "application/json")
	tokenString := req.Header.Get("Authorization")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Fatal("Unexpected signing method")
			return nil, fmt.Errorf("Unexpected signing method %v", ok)
		}
		return []byte("secret"), nil
	})

	var userResponseObject model.User
	var res model.ResponseResult

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userResponseObject.Username = claims["username"].(string)
		json.NewEncoder(w).Encode(userResponseObject)
		return
	} else {
		res.Error = err.Error()
		json.NewEncoder(w).Encode(res)
		return
	}
}
