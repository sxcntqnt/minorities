#!/env python

import datetime
import paho.mqtt.client as mqtt
import json

#getting date and time
current_time = str(datetime.datetime.now())
current_time_minute = current_time[:]
current_time_minute = str(time.time())
#read iot data from sensors
location, rainstatus = os.open('/dev/ttyACM0')
#payload data to be sent
payload ="location:" str(longitude,latitude,\n"rainStatus:"status,\n"timestamp:"current_time_minute)
#ip addr of the broker
broker_address="192.168.0.29"
#create new client instance
client = mqtt.Client("pi_client")
#connect client to broker
client.conect(broker_address)
#publish iot data to subscribed toptic
client.publish("IOTDATA-Topic", json.dumps(d))
time.sleep(1000
