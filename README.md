# Hayward AquaConnect pool data collector

Pool operators with a Hayward Pro Logic or Aqua Plus system **with the AQ-CO-HOMENET AquaConnect networking unit** can retrieve pool equipment data and store in a variety of time series databases.

## Getting Started and Why

There exists an app and cloud storage ecosystem for these pool controllers.  This tool is for someone who wishes to export and use the pool data independently and wants to use the home networking hardware, e.g. where pool location makes wireless communication desirable.

### Supported Hardware

 * Hayward Pro Logic or Aqua Plus main control unit
 * Hayward Goldline AQL2-BASE-RF AquaConnect Wireless Antenna
 * Hayward AQ-CO-HOMENET AquaConnect Home Network, Internet and Wi-fi Remote Control

### Software Prerequisites

 * A supported time-series DB (influxdb is recommended, but also kairosdb and thingsboard API)
 * Grafana `apt install grafana`
 * golang to build the project `apt install golang`

### Installing

```
git clone git@github.com:ma-nathan/hayward-pool.git
cd hayward-pool
edit settings.ini
make
./pool
```

### How it works

This is an odd system that seems to span several generations of technology.  The AQL2-BASE-RF at the controller speaks some proprietary 900 MHz RF protocol with the AQ-CO-HOMENET indoors, which presents a Web UI (and an endpoint for the cloud integration) on your LAN.

I was able to "scrape" the Web UI and decypher the status string (up to a point) which resembles the serial output decoded by [draythomp](http://www.desert-home.com/p/swimming-pool.html).

There may be variations in how different devices are configured to the controller.  The code assumes your pool, like mine, controls the "Cleaner" via AUX2, but if it does not, you will have to edit scrape.go, where you will find my attempts to decode the status string.

Outside of this web-UI-scraping method, you can also connect a RS-485 to Ethernet adapter directly to the controller and use the [aqualogic python library](https://github.com/swilson/aqualogic) which also allows features control and is overall a more sophisticated approach.

### Next steps

Rework this tool into a golang library or HTTP API.  

Separate out the stats delivery such that the stats can be simple inputs to telegraf instead of delivered directly to time-series DB backend.

Support querying pump power usage.

### Results

Grafana dashboard available for download at [https://grafana.com/grafana/dashboards/11354](https://grafana.com/grafana/dashboards/11354)

![Example grafana dashboard](http://www.fumanchu.com/pool-dashboard-example.png)

### Contact

You may reach me at my personal address nb@fumanchu.com and i may be able to help.

### License

This project is licensed under the MIT License

