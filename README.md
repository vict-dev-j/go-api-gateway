# go-app

This repository hosts Golang microservices for an API gateway, investment accounts service, and customer data service. The repository is currently in development.

## Usage

### Build Docker Images

To build the Docker images for each service:

```bash
docker build -t customers-service -f path/to/Dockerfile.customers .
docker build -t gateway-service -f path/to/Dockerfile.gateway .
docker build -t invest-accounts-service -f path/to/Dockerfile.invest-accounts .
```

### To run each service container:

```bash
docker run -d -p 8080:8080 customers-service
docker run -d -p 8081:8081 gateway-service
docker run -d -p 8082:8082 invest-accounts-service
```

### Testing the API:

# Authorization

```bash
curl -X POST http://localhost:8081/login `
  -H "Content-Type: application/json" `
  -d "{\"username\": \"admin\", \"password\": \"password\"}"
```
Replace <token> with the actual JWT token obtained from the previous step

```bash
# Example GET request to fetch customer data
curl -X GET http://localhost:8081/customer `
  -H "Authorization: Bearer $TOKEN"

# Example POST request to create a customer
curl -X POST http://localhost:8081/customer `
  -H "Authorization: Bearer $TOKEN" `
  -H "Content-Type: application/json" `
  -d "{\"name\": \"V N\", \"email\": \"v.n@example.com\"}"

# Example PUT request to update a customer
curl -X PUT http://localhost:8081/customer/1 `
  -H "Authorization: Bearer $TOKEN" `
  -H "Content-Type: application/json" `
  -d "{\"name\": \"Updated Name\"}"

# Example DELETE request to delete a customer
curl -X DELETE http://localhost:8081/customer/1 `
  -H "Authorization: Bearer $TOKEN"
```

Use the following curl commands to interact with the API:

- Retrieve all customers:
```bash
curl -X GET http://localhost:8081/customer
```

- Retrieve a specific customer (e.g., customer with ID 4):
```bash
curl -X GET http://localhost:8081/customer/4
```

- Create a new customer:
```bash
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
```

- Update an existing customer (e.g., customer with ID 7):
```bash
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
```

- Delete a customer (e.g., customer with ID 7):
```bash
curl -X DELETE http://localhost:8081/customer/7
```

- Retrieve all investment accounts:
```bash
curl -i http://localhost:8081/invest-account
```

- Retrieve a specific account (e.g., account with ID 1):
```bash
curl -i http://localhost:8081/invest-account/1
```

- Create a new account:
```bash
curl -i -X POST http://localhost:8081/invest-account -H "Content-Type: application/json" -d "{\"owner_id\": \"1\", \"client_survey_number\": \"12345678\", \"share\": 100, \"invested_amount_of_money\": 5000.0, \"free_amount_of_money\": 2000.0}"
```

- Update an existing account (e.g., account with ID 7):
```bash
curl -i -X PUT http://localhost:8081/invest-account/7 -H "Content-Type: application/json" -d "{\"owner_id\": \"1\", \"client_survey_number\": \"87654321\", \"share\": 150, \"invested_amount_of_money\": 6000.0, \"free_amount_of_money\": 2500.0}"
```

- Delete an account (e.g., account with ID 1):
```bash
curl -i -X DELETE http://localhost:8081/invest-account/1
```
