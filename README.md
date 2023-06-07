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

*   Get Company (GET) to `localhost:8000/companies/{id}`

*   Create Company (POST) to `localhost:8000/companies` with request body:

    ```
    {
        "name": "",
        "currency": "EUR"
    }
    ```

*   Update Company (PATCH) to `localhost:8000/companies/{id}` with request body:

    ```
    {
        "balance": 15000,
    }
    ```

*   Delete Company (DELETE) to `localhost:8000/companies/{id}`