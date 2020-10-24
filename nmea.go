package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/jacobsa/go-serial/serial"
	"log"
	"math"
	"strconv"
	"strings"
	"time"
	
)

type NMEA struct {
	fixTimestamp      string
	latitude          string
	latitudeDirection string
	longitude	string
	longitudeDirection string
	fixQuality         string
	satellites         string
	horizontalDilution string
	antennaHeight      string
	updateAge          string
}

func ParseNmeaLine(line string) (NMEA, error) {
	tokens := strings.Split(line, ",")
	if tokens[0] == "$GPGGA" {
		return NMEA{
			fixTimestamp:       tokens[1],
			latitude:           tokens[2],
			latitudeDirection:  tokens[3],
			longitude:          tokens[4],
			longitudeDirection: tokens[5],
			fixQuality:         tokens[6],
			satellites:         tokens[7],
		}, nil
	}
	return NMEA{}, errors.New("Unsupported nmea string")
}
func ParseDegrees(value string, direction string) (string, error) {
	if value == "" || direction == "" {
		return "", errors.New("the location and /or direction value does not exits")
	}
	lat, _ := strconv.ParseFloat(value, 64)
	degrees := math.Floor(lat / 100)
	minutes := ((lat / 100) - math.Floor(lat/100)) * 100 / 60
	decimal := degrees + minutes
	if direction == "W" || direction == "S" {
		decimal *= -1
	}
	return fmt.Sprintf("%.6f", decimal), nil
}
func (nmea NMEA) GetLatitude() (string, error) {
	return ParseDegrees(nmea.latitude, nmea.latitudeDirection)
}
func (nmea NMEA) GetLongitude() (string, error) {
	return ParseDegrees(nmea.longitude, nmea.longitudeDirection)
}
func main() {
	fmt.Println("Starting the application...")

	options := serial.OpenOptions{
		PortName:        "/dev/ttyACM0",
		BaudRate:            9600,
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 4,
	}

	serialPort, err := serial.Open(options)
	if err != nil {
		log.Fatal("Serial.Open: %v", err)
	}
	defer serialPort.Close()
	reader := bufio.NewReader(serialPort)
	scanner := bufio.NewScanner(reader)
	geocoder := Geocoder{Bearer: "", ApiKey: "Se644nTA3g06gtu4G87Pqe92ODlEP6lSLyubbIyRTaQ"}
	for scanner.Scan() {
		gps, err := ParseNMEALine(scanner.Text())
		if err == nil {
			if gps.fixQuality == "1" || gps.fixQuality == "2" {
				latitude, _ := gps.GetLatitude()
				longitude, _ := gps.GetLongitude()
				fmt.Println(latitude + "," + longitude)
				result, _ := geocoder.reverse(Position{Latitude: latitude, Longitude: longitude})
				if len(result.Response.View) > 0 && len(result.Response.View[0].Result) > 0 {
					data, _ := json.Marshal(result.Response.View[0].Result[0])
					fmt.Println(string(data))
				} else {
					fmt.Println("no address estimates found for the position")
				}
			} else {
				fmt.Println("no gps fix available")
			}
			time.Sleep(2 * time.Secound)
		}
	}
}
