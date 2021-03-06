package user

import (
	"encoding/json"
	"fmt"
	"go-api-ws/addresses"
	"go-api-ws/auth"
	"go-api-ws/cart"
	"go-api-ws/config"
	"go-api-ws/core"
	"go-api-ws/helpers"
	"net/http"
	"strconv"
)

var userModule core.ApiModule

func init() {
	userModule = core.ApiModule{
		Name:        "User module",
		Description: "User module. Supports username and email authentication. Categories are stored as a flat list.",
		Version:     "0.1",
		Author:      "Matas Cereskevicius @ JivaLabs",
	}

}

const adminRole = "admin"
const userRole = "user"

// Get Order History
// Path: /api/user/order-history
func getOrderHistory(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get order history")
	urlToken, err := helpers.GetTokenFromUrl(r)
	helpers.PanicErr(err)

	token := auth.ParseToken(urlToken)

	claims, err := auth.GetTokenClaims(token)
	helpers.CheckErr(err)

	if err != nil {
		helpers.WriteResultWithStatusCode(w, "Invalid token", http.StatusForbidden)
	} else {
		if auth.CheckIfTokenIsNotExpired(claims) {

			orderHistory := getUserOrderHistoryFromMongo(claims["sub"].(string))
			response := helpers.Response{
				Code:   http.StatusOK,
				Result: orderHistory}
			response.SendResponse(w)
		} else {
			helpers.WriteResultWithStatusCode(w, "Token is expired", http.StatusForbidden)
		}
	}
}

// Me endpoint
// Path /api/user/me
func meEndpoint(w http.ResponseWriter, r *http.Request) {
	urlToken, err := helpers.GetTokenFromUrl(r)
	helpers.PanicErr(err)
	token := auth.ParseToken(urlToken)

	claims, err := auth.GetTokenClaims(token)
	helpers.CheckErr(err)

	if err != nil {
		helpers.WriteResultWithStatusCode(w, "Invalid token", http.StatusBadRequest)
	} else {
		if auth.CheckIfTokenIsNotExpired(claims) {
			userId, err := strconv.Atoi(claims["sub"].(string))
			helpers.PanicErr(err)
			userId64 := int64(userId)
			userInfo := GetUserFromMySQLById(userId64)
			userInfo.Addresses = addresses.GetAddressesFromMySQL(userId64)
			response := helpers.Response{
				Code:   http.StatusOK,
				Result: userInfo}
			response.SendResponse(w)
		} else {
			helpers.WriteResultWithStatusCode(w, "Token expired", http.StatusForbidden)
		}
	}
}

// RegisterUser function
// Path: api/user/create
func registerUser(w http.ResponseWriter, r *http.Request) {
	var user User
	_ = json.NewDecoder(r.Body).Decode(&user)
	validationResult := helpers.CheckJSONSchemaWithGoStruct("file://user/jsonSchemaModels/userRegister.schema.json", user)
	if validationResult.Valid() {
		id := insertUserIntoMySQL(user)
		customer := GetUserFromMySQLById(id)
		cart.CreateCartInMongoDB(user.ID)
		response := helpers.Response{
			Code:   http.StatusOK,
			Result: customer}
		response.SendResponse(w)
	} else {
		helpers.WriteResultWithStatusCode(w, validationResult.Errors(), http.StatusBadRequest)
	}
}

// Path: /api/user/update
//Method: post
func updateUser(w http.ResponseWriter, r *http.Request) {
	var user UpdatedCustomer
	err := json.NewDecoder(r.Body).Decode(&user)
	helpers.PanicErr(err)

	for i := range user.UpdateUser.Addresses {
		user.UpdateUser.Addresses[i].InsertOrUpdateAddressIntoMySQL(user.UpdateUser.ID)
	}
	for _, address := range user.UpdateUser.Addresses {
		if address.DefaultShipping == true {
			user.UpdateUser.DefaultShipping = strconv.Itoa(int(address.ID))
		}
	}
	user.UpdateUser.UpdateUserByIdMySQL()
	response := helpers.Response{
		Result: user.UpdateUser,
		Code:   http.StatusOK}
	response.SendResponse(w)
}

// Path: /api/user/refresh
func refreshToken(w http.ResponseWriter, req *http.Request) {
	var jsonBody map[string]string
	_ = json.NewDecoder(req.Body).Decode(&jsonBody)
	token := auth.ParseToken(jsonBody["refreshToken"])

	claims, err := auth.GetTokenClaims(token)
	helpers.CheckErr(err)

	if err != nil {
		helpers.WriteResultWithStatusCode(w, "Invalid token", http.StatusBadRequest)
	} else {
		if auth.CheckIfTokenIsNotExpired(claims) {

			groupId := GetGroupIdFromMySQLById(claims["sub"].(int))

			authToken := auth.GetNewAuthToken(claims["sub"].(string), groupId)

			authTokenString, err := authToken.SignedString([]byte(config.MySecret))
			helpers.PanicErr(err)

			refreshToken := auth.GetNewRefreshToken(claims["sub"].(string))
			refreshTokenString, err := refreshToken.SignedString([]byte(config.MySecret))
			helpers.PanicErr(err)

			response := helpers.Response{
				Code:   http.StatusOK,
				Result: authTokenString,
				Meta: map[string]string{
					"refreshToken": refreshTokenString}}
			response.SendResponse(w)

		} else {
			helpers.WriteResultWithStatusCode(w, "Token expired", http.StatusForbidden)
		}
	}
}

// Path: /api/user/login
func loginEndpoint(w http.ResponseWriter, req *http.Request) {
	var userLogin LoginForm

	_ = json.NewDecoder(req.Body).Decode(&userLogin)
	validationResult := helpers.CheckJSONSchemaWithGoStruct("file://user/jsonSchemaModels/userLogin.schema.json", userLogin)

	pswd := userLogin.Password
	userLogin.Password = ""

	if validationResult.Valid() {
		userFromDb := getUserFromMySQLByEmail(userLogin.Username)
		if checkPasswordHash(pswd, userFromDb.Password) {

			authToken := auth.GetNewAuthToken(userFromDb.ID, userFromDb.GroupId)
			authTokenString, err := authToken.SignedString([]byte(config.MySecret))
			helpers.PanicErr(err)

			refreshToken := auth.GetNewRefreshToken(userFromDb.ID)
			refreshTokenString, err := refreshToken.SignedString([]byte(config.MySecret))
			helpers.PanicErr(err)

			response := helpers.Response{
				Code:   http.StatusOK,
				Result: authTokenString,
				Meta: map[string]string{
					"refreshToken": refreshTokenString,
				}}
			response.SendResponse(w)
		} else {
			helpers.WriteResultWithStatusCode(w, "Password is invalid", http.StatusUnauthorized)
		}
	} else {
		helpers.WriteResultWithStatusCode(w, validationResult.Errors(), http.StatusBadRequest)

	}
}
