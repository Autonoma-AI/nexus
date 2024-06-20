# Nexus

A reverse proxy and load balancing with dynamic backend configuration.

## Description

This project implements a basic HTTP load balancer that distributes incoming requests to a set of backend services configured via environment variables. It features:

- **Load Balancing:** Distributes requests across multiple backend services.
- **Dynamic Configuration:** Backend services can be dynamically configured using environment variables.

## Installation

To install and run this project, follow these steps:

1. Clone the repository:
    ```bash  
    git clone https://github.com/Autonoma-AI/nexus  
    cd nexus 
    ```  

2. Build the project:
    ```bash  
    go build -o nexus
    ```  

3. Set the environment variables for backend services. For example:
    ```bash  
    export NEXUS_SERVICE1_LOCATION1_URL=http://backend1:8080  
    export NEXUS_SERVICE1_LOCATION1_HEADERS="Authorization=Bearer token"  
    export NEXUS_SERVICE2_LOCATION1_URL=http://backend2:8080  
    export NEXUS_SERVICE2_LOCATION1_HEADERS="Authorization=Bearer token;Logging=true"  
    export PORT=8080  
    ```  

4. Run the executable:
    ```bash  
    ./nexus
    ```  

## Usage

Once Nexus is running, you can use the following endpoints:

- **Health Check:** `GET /health`
    ```bash  
    curl http://localhost:8080/health  
    ```  
  This should return:
    ```  
    Healthy!  
    ```  

- **Load Balanced Requests:** Any other requests will be proxied to one of the backend services.

## Environment Variables

The backend services can be configured via environment variables with the following format:

- **URL Configuration**: `NEXUS_<SERVICE>_<LOCATION>_URL`
    - `SERVICE`: Name of the service.
    - `LOCATION`: A unique identifier for the backend location.
    - `URL`: The URL of the backend service.

- **Headers Configuration**: `NEXUS_<SERVICE>_<LOCATION>_HEADERS`
    - `HEADERS`: Semicolon-separated key=value pairs of headers to add to requests.

## Contributing

If you would like to contribute to this project, please follow these steps:
- Fork the repository.
- Create a new branch for your feature or bug fix.
- Write your code.
- Submit a pull request.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.
