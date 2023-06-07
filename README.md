## XM Companies
A microservice to handle companies.

To launch the application, mount the root directory of the project and run docker:
```
docker compose up -d
```
and
```
make start
```

To migrate up
```
make migrate-up
```

To migrate down
```
make migrate-down
```

To run tests (coverage)
```
make coverage
```

To run tests with integration (coverage)
```
make coverage-integration
```

To run tests
```
make test
```

To run tests with integration
```
make test-integration
```

While the application is running, we can make requests to get, add, update and remove companies.

Get Company (GET) to `localhost:8000/companies/{id}`

Create, Update and Delete routes are protected with jwt-authorization. Therefor, the following HTTP Headers must be included to the requests:

```json
    "Authorization" : "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IlBldHJvcyBUcmFrYWRhcyIsImlhdCI6MTUxNjIzOTAyMn0.qSUO5wUFgkrvp5C96R_LSMy6tkTVGYQ74ELMrX4Zeyw"
```

Create Company (POST) to `localhost:8000/companies` with request body:

```json
{
    "name": "Petros Inc.",
    "description": "A short desc of my company",
    "number_of_employees": 50,
    "registered": false,
    "type": "Sole Proprietorship"
}
```

Update Company (PATCH) to `localhost:8000/companies/{id}` with request body:

```json
{
    "name": "Petros GmbH.",
    "number_of_employees": 5,
    "registered": true,
    "type": "NonProfit"
}
```

Delete Company (DELETE) to `localhost:8000/companies/{id}`

Create, Update and Delete handlers will also produce kafka events that can be monitored in kowl:
```
    http://localhost:8080/topics
```