package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"qok.com/identity/logwrapper"
	"qok.com/identity/model"
	"qok.com/identity/userrepository"

	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

//RegisterHandler is a http Handler for /register  route
func RegisterHandler(w http.ResponseWriter, req *http.Request) {
	logger := logwrapper.Load()
	w.Header().Set("Content-Type", "application/json")
	var responseObj model.ResponseResult
	var user model.User

	body, ioerr := ioutil.ReadAll(req.Body)
	if ioerr != nil {
		logger.Println(ioerr.Error())
		responseObj.Error = ioerr.Error()
		json.NewEncoder(w).Encode(responseObj)
		return
	}

	err := json.Unmarshal(body, &user)
	if err != nil {
		logger.Printf("error in umarshalling : %v", err.Error())
		responseObj.Error = err.Error()
		json.NewEncoder(w).Encode(responseObj)
		return
	}

	_, findErr := userrepository.FindOne(user.Username)

	if findErr == nil {
		responseObj.Result = fmt.Sprintf("Username : %q  already Exists!!", user.Username)
		json.NewEncoder(w).Encode(responseObj)
		return
	}

	if findErr.Error() == "mongo: no documents in result" {
		err := userrepository.Store(user)
		if err != nil {
			responseObj.Error = err.Error()
			json.NewEncoder(w).Encode(responseObj)
			return
		}
		//LAST STATE WHICH NOTHING WENT WRONG
		responseObj.Result = "Successfuly Registered"
		json.NewEncoder(w).Encode(responseObj)
	}

}

// LoginHandler is a http handler for /login  route
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

	var userObjectForResponse model.User
	var res model.ResponseResult

	userObjectForResponse, dbFindError := userrepository.FindOne(user.Username)
	if dbFindError != nil {
		logger.Printf("Invalid username : %q", user.Username)
		res.Error = fmt.Sprintf("Invalid username : %v", user.Username)
		json.NewEncoder(w).Encode(res)
		return
	}

	passwordError := bcrypt.CompareHashAndPassword([]byte(userObjectForResponse.Password), []byte(user.Password))

	if passwordError != nil {
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

// UserInfoHandler is a http handler for /user_info route
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

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userResponseObject := model.User{
			Username: claims["username"].(string),
			Token:    tokenString,
		}
		json.NewEncoder(w).Encode(userResponseObject)
		return
	}

	res := model.ResponseResult{
		Error: err.Error(),
	}
	json.NewEncoder(w).Encode(res)
	return

}
