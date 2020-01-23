package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	Expires     int    `json:"expires_in"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type CodePayload struct {
	Code string `json:"code"`
}

type SecureResponse struct {
	Authorization string `json:"authorization"`
}

var authCodes = make(map[string]bool)
var accessTokens = make(map[string]bool)

func Code(w http.ResponseWriter, r *http.Request) {
	code := RandomString()
	authCodes[code] = true
	log.Printf("Generated code = %s", code)

	queryValues := r.URL.Query()
	url := queryValues.Get("redirect_uri") + "?code=" + code
	log.Printf("Redirect URI = %s", url)

	http.Redirect(w, r, url, 301)
}

func Token(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var data CodePayload
	err := decoder.Decode(&data)
	if err != nil {
		panic(err)
	}

	code := data.Code

	log.Printf("Code = %s", code)
	_, ok := authCodes[code]

	if ok {
		delete(authCodes, code)
		accessToken := RandomString()
		accessTokens[accessToken] = true
		log.Printf("Generated token = %s", accessToken)
		tokenResp := TokenResponse{AccessToken: accessToken, Expires: 60 * 60 * 24}

		sendResponse(w, tokenResp, http.StatusOK)
	} else {
		sendResponse(w, ErrorResponse{Message: "Invalid auth token"}, 400)
	}
}

func Secure(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("authorization")
	log.Printf("Authorization = %s", authorization)

	_, ok := accessTokens[authorization]

	if ok {
		sendResponse(w, SecureResponse{Authorization: authorization}, http.StatusOK)
	} else {
		sendResponse(w, ErrorResponse{Message: "Unauthorized"}, 403)
	}
}

func sendResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}

func RandomString() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
