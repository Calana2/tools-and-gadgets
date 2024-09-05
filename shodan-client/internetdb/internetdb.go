package internetdb

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const URL = "https://internetdb.shodan.io"

type IDBHost struct {
 Cpes []string `json:"cpes"`
 Hostnames []string `json:"hostnames"`
 Ip string `json:"ip"`
 Ports []int `json:"ports"`
 Tags []string `json:"tags"`
 Vulnerabilities []string `json:"vulns"`
}

type InternetDB struct { }

func New() *InternetDB {
 return &InternetDB{}
}

func (s* InternetDB) IpLookup(ip string) (*IDBHost,error) {
  res,err := http.Get(fmt.Sprintf("%s/%s",URL,ip))

  if err != nil {
   return nil,err
  } 

  defer res.Body.Close()

  if (res.StatusCode != 200) {
   var p map[string]interface{}
   if err := json.NewDecoder(res.Body).Decode(&p); err != nil{
    panic(err)
   }
   fmt.Println(p["detail"])
   return nil,nil
  }
  
  var ret IDBHost
  
  if err := json.NewDecoder(res.Body).Decode(&ret); err != nil {
   return nil,err
  }
  return &ret,nil
}


