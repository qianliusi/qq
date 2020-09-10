package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func ParseRequestTest(r *http.Request, T interface{}) interface{} {
	CheckError(r.ParseForm())
	err := json.NewDecoder(r.Body).Decode(T)
	CheckError(err)
	return T
}

func SendResponseJson(w http.ResponseWriter, T interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(T)
	CheckError(err)
}

func ParseRequest(r *http.Request, T interface{}) interface{} {
	return parseForm(r, T)
}

func parseForm(r *http.Request, T interface{}) interface{} {
	s := reflect.ValueOf(T).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		f.SetString(r.FormValue(typeOfT.Field(i).Name))
	}
	return T
}

func parseFormPost(r *http.Request, T interface{}) interface{} {
	fmt.Println(json.Marshal(r.Form))
	err := json.NewDecoder(r.Body).Decode(T)
	CheckError(err)
	return T
}

func parseFormGet(r *http.Request, T interface{}) interface{} {
	s := reflect.ValueOf(T).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		f.Set(reflect.ValueOf(r.FormValue(typeOfT.Field(i).Name)))
	}
	return T
}
