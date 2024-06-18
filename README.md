# go-app

This repository hosts Golang microservices for an API gateway, investment accounts service, and customer data service.

## Usage

### Build Docker Images

To build the Docker images for each service:

```bash
docker build -t customers-service -f path/to/Dockerfile.customers .
docker build -t gateway-service -f path/to/Dockerfile.gateway .
docker build -t invest-accounts-service -f path/to/Dockerfile.invest-accounts .


###To run each service container:

docker run -d -p 8080:8080 customers-service
docker run -d -p 8081:8081 gateway-service
docker run -d -p 8082:8082 invest-accounts-service

# Testing the API:

Use the following curl commands to interact with the API:

- Retrieve all customers:
curl -X GET http://localhost:8081/customer

- Retrieve a specific customer (e.g., customer with ID 4):
curl -X GET http://localhost:8081/customer/4

- Create a new customer:
curl -X POST \
  http://localhost:8081/customer \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John",
    "surname": "Doe",
    "age": 30,
    "phone_number": "1234567890",
    "debit_card": "1234-5678-9101-1121",
    "credit_card": "5432-1098-7654-3210",
    "date_of_birth": "1994-05-20T00:00:00Z",
    "date_of_issue": "2023-01-15T00:00:00Z",
    "issuing_authority": "Authority XYZ",
    "has_foreign_country_tax_liability": false
  }'

- Update an existing customer (e.g., customer with ID 7):
curl -X PUT \
  http://localhost:8081/customer/7 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Name",
    "surname": "Updated Surname",
    "age": 35,
    "phone_number": "9876543210",
    "debit_card": "5678-9101-1121-3141",
    "credit_card": "8765-4321-0987-6543",
    "date_of_birth": "1990-08-15T00:00:00Z",
    "date_of_issue": "2020-03-10T00:00:00Z",
    "issuing_authority": "Authority ABC",
    "has_foreign_country_tax_liability": true
  }'

- Delete a customer (e.g., customer with ID 7):
curl -X DELETE http://localhost:8081/customer/7
