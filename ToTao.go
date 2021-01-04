//get current location
package main

import ("bytes"
      "encoding/json"
      "fmt"
      "os/exec"
      "regexp"
      "runtime"
      "strconv"
      "strings"
  )
type En_Route{
	Field1	`json:"FIELD1"`
	Field2	`json:"FIELD2"`
}

func checkRouteNo(file,Address) string{
  /*grep this list from the location
    and determine matatu that plies that route
    from FIELD1 and FIELD2
  */
    match, _ := regexp.MatchString("^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$", columns[1])
	file,_ := os.Open("xray/YesBana.json")
	var Twende En_Route
	Twende := unmarshal.json(file)
	match, _ := regexp.MatchString("^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$",Twende)
	

}
func requestMat(location){
  /*Find the nearest matatu
  Do a blockchain lookup from location specified
  */
  //Query matatu's going to town only
}
