# StormBox REST API

StormBox is a chat application developed in Golang using the Echo framework. This project includes functionality for user registration, authentication, user management, message managing.

## Installation and Running

### Requirements

- Go 1.22+
- Docker

### Installation and run

1. Clone the repository:
    ```sh
    git clone https://github.com/LzAxel/stormbox-backend
    cd stormbox-backend
    ```

2. Install dependencies:
    ```sh
    go mod tidy
    ```

3. Start Postgresql in docker:
    ```sh
    docker-compose up -d
    ```

4. Run the server:
    ```sh
    go run cmd/main.go
    ```

## Config
All configs are located in ```configs/``` directory. 
`dev.yaml` is a basic config that you can use to preview project.
`config.example.yaml` is a preset for your own configs.

## Dev-Features
- Auto db migrations on start app

## Features
- User registration
- Authentication and token refresh
- Retrieve user list and self-information
- Add users as friends and manage friends list
- Send and receive messages
- Subscribe to receive new messages via WebSocket
