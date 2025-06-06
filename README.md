# dipMoistureMonitorGo
This simple application acts as a bridge between my ESP-based REST server and Prometheus.

Specifically, it runs a lightweight REST server that exposes a /metrics endpoint for Prometheus to scrape.
Each time Prometheus queries /metrics, the application performs a real-time GET request to the ESP REST server to retrieve the latest sensor data.
