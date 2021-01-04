#include <TinyGPS++.h>
#include <SoftwareSerial.h>


// Choose two Arduino pins to use for software serial
int RXPin = 19;
int TXPin = 18;

// Sensor pins
#define sensorPower 7
#define sensorPin 8


int GPSBaud = 9600;

// Create a TinyGPS++ object
TinyGPSPlus gps;

// Create a software serial port called "Serial1"
SoftwareSerial gpsSerial(RXPin, TXPin);


void setup()
{
  // Start the Arduino hardware serial port at 9600 baud
  Serial.begin(GPSBaud);
  // Start the software serial port at the GPS's default baud
  Serial1.begin(GPSBaud);

  pinMode(sensorPower, OUTPUT);
  // Initially keep the sensor OFF
  digitalWrite(sensorPower, LOW);
  Serial.begin(GPSBaud);
}

// This function returns the sensor output
int readSensor() {
  digitalWrite(sensorPower, HIGH);  // Turn the sensor ON
  delay(10);              // Allow power to settle
  int val = digitalRead(sensorPin); // Read the sensor output
  digitalWrite(sensorPower, LOW);   // Turn the sensor OFF
  return val;             // Return the value
}

void displayInfo() {
  if (gps.location.isValid())
  {
    Serial.print("Latitude: ");
    Serial.println(gps.location.lat(), 6);
    Serial.print("Longitude: ");
    Serial.println(gps.location.lng(), 6);
    Serial.print("Altitude: ");
    Serial.println(gps.altitude.meters());
  }
  else
  {
    Serial.println("Location: Not Available");
  }

  Serial.print("Date: ");
  if (gps.date.isValid())
  {
    Serial.print(gps.date.month());
    Serial.print("/");
    Serial.print(gps.date.day());
    Serial.print("/");
    Serial.println(gps.date.year());
  }
  else
  {
    Serial.println("Not Available");
  }

  Serial.print("Time: ");
  if (gps.time.isValid())
  {
    if (gps.time.hour() < 10) Serial.print(F("0"));
    Serial.print(gps.time.hour());
    Serial.print(":");
    if (gps.time.minute() < 10) Serial.print(F("0"));
    Serial.print(gps.time.minute());
    Serial.print(":");
    if (gps.time.second() < 10) Serial.print(F("0"));
    Serial.print(gps.time.second());
    Serial.print(".");
    if (gps.time.centisecond() < 10) Serial.print(F("0"));
    Serial.println(gps.time.centisecond());
  }
  else
  {
    Serial.println("Not Available");
  }

  Serial.println();
  Serial.println();
  delay(1000);
}


void loop() {
  // This sketch displays information every time a new sentence is correctly encoded.
  while (Serial1.available() > 0)
    if (gps.encode(Serial1.read()))
      displayInfo();


// If 5000 milliseconds pass and there are no characters coming in
  // over the software serial port, show a "No GPS detected" error
  if (millis() > 5000 && gps.charsProcessed() < 10)
  {
    Serial.println("No GPS detected");
    while (true);
  }

  //get the reading from the function below and print it
  int val = readSensor();
  Serial.print("Digital Output: ");
  Serial.println(val);

  // Determine status of rain
  if (val) {
    Serial.println("Status: Clear");
  } else {
    Serial.println("Status: It's raining");
  }

  delay(1000);  // Take a reading every minute
  Serial.println();




}
