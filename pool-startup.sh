#!/bin/bash
cd "$(dirname "$0")"

# For swim season when the pool is heated:
#HEATER_TEMP=87

# For winter when the heater is not on:
HEATER_TEMP=40

LOGFILE=/home/nb/go/src/pool/logfile

./pool --heater_temp=$HEATER_TEMP > $LOGFILE 2>&1
