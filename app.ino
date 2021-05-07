#include "DHT.h"

#define DHTPIN 2

#define DHTTYPE DHT11

DHT dht(DHTPIN, DHTTYPE);

void setup() {
  Serial.begin(9600);
  dht.begin();
}

void loop() {
  delay(1000);

  Serial.print("humidity=");
  Serial.print(dht.readHumidity());
  Serial.print("\t temperature=");
  Serial.println(dht.readTemperature());
}
