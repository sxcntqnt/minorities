// Sensor pins
#define sensorPower 7
#define sensorPin 8

void setup() {
	pinMode(sensorPower, OUTPUT);

	// Initially keep the sensor OFF
	digitalWrite(sensorPower, LOW);

	Serial.begin(9600);
}

void loop() {
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

	delay(1000);	// Take a reading every second
	Serial.println();
}

//  This function returns the sensor output
int readSensor() {
	digitalWrite(sensorPower, HIGH);	// Turn the sensor ON
	delay(10);							// Allow power to settle
	int val = digitalRead(sensorPin);	// Read the sensor output
	digitalWrite(sensorPower, LOW);		// Turn the sensor OFF
	return val;							// Return the value
}
