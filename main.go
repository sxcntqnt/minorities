package main

import ("fmt"
    "github.com/hyperledger/fabric/core/chaincode/shim"
    pb "github.com/hyperledger/fabric/protos/peer"
    "math"
)
/**Calculates distance btn two points
  Given latitude/longitude of those points
  Uses GeodataSource products
  Definitions:
    South latitudes are negative,east longitudes are positive
  Passed to function:
  lat1, lon1 = Latitude and Longitude of point 1 (in decimal degrees)
    lat2, lon2 = Latitude and Longitude of point 2 (in decimal degrees)
    unit = the unit you desire for results
            where: 'K' is kilometers
  **/
  func distance (lat1 float64, lng1 float64, lat2 float64 lng2 float64 , unit ...string) float64{
    const PI float64 = 3.141592653589793

    radlat1 := float64(PI * lat1 /180)
    radlat2 := float64(PI * lat2 /180)

    theta := float64(lng1-lng2)
    radtheta := float64(PI * theta /180)
    dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

    if dist > 1 {
      dist =1
    }

    dist = math.Acos(dist)
    dist = dist *180/PI
    dist = dist * 60 * 1.1515

    if len(unit) > 0 {
      if unit[0] =="K"{
        dist = dist * 1.609344
      }
      return dist
}
func speed (s float64, dist float64) float64 {
	spd := dist / s

	return spd
}

func main() {
  err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
	fmt.Printf("%f Kilometers\n", distance(32.9697, -96.80322, 29.46786, -98.53506, "K"))
	fmt.Printf("%f TTGH" in minutes, speed(80.89,422.738931))

}
