To run some tests you need to copy the 'cfg.json' file to their directory

# Order cacher

This project is a Go-based application that integrates with Kafka and PostgreSQL to process and manage orders. It includes a producer, consumer, and a simple web UI for interacting with the data.

## Features

- **Kafka Integration**: Produces and consumes messages using Kafka.
- **PostgreSQL Integration**: Stores and retrieves order data in a PostgreSQL database.
- **Web UI**: A simple web interface to view and query orders.
- **Docker Support**: Includes a `docker-compose.yml` file for setting up Kafka, Zookeeper, and PostgreSQL.
