package userrepository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"qok.com/identity/db"
	"qok.com/identity/logwrapper"
	"qok.com/identity/model"
)

//FindOne return a User struct or an error
func FindOne(username string) (model.User, error) {
	var returningUser model.User
	collection, err := db.GetDBCollection()
	if err != nil {
		logwrapper.Load().Fatalf("error in connecting to DB : %q", err.Error())
		return model.User{}, err
	}
	err = collection.FindOne(context.TODO(), bson.D{{"username", username}}).Decode(&returningUser)
	return returningUser, err
}

//Store store user struct into database
func Store(user model.User) error {
	collection, err := db.GetDBCollection()
	if err != nil {
		logwrapper.Load().Fatalf("error in connecting to DB : %q", err.Error())
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)
	if err != nil {
		logwrapper.Load().Fatalf("error hashing the password : %q", err.Error())
		return err
	}

	user.Password = string(hash)

	_, err = collection.InsertOne(context.TODO(), user)

	if err != nil {
		logwrapper.Load().Fatalf("error inserting DB : %q", err.Error())
		return err
	}

	return nil
}
