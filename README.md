# Assignment subito.it

This is the implementation of the assignment requested from subito.it.

It's implemented in golang, with the only external dependency being `sqlite`.

Go is not my primary language, I've used it mostly in personal projects and some CI in my previous job, but I thought that it would be a nice way of showing you how I'd approach a new challenge in a new programming language.

## Project Structure

The project is structured as follows:

```
.
├── cmd/server/main.go          # Entrypoint — wires repos, services, handlers
├── internal/
│   ├── entity/                 # Domain structs and validation
│   │   ├── product.go
│   │   └── order.go
│   ├── repository/             # SQLite data access layer
│   │   ├── product_repository.go
│   │   └── order_repository.go
│   ├── service/                # Business logic
│   │   ├── product_service.go
│   │   └── order_service.go
│   └── handler/                # HTTP handlers
│       ├── product_handler.go
│       ├── order_handler.go
│       ├── health_handler.go
│       └── helpers.go
├── scripts/
│   ├── run.sh                  # Build and run via Docker
│   └── tests.sh                # Run tests via Docker
├── Dockerfile
└── Makefile
```

Each relevant file has its own test, following Go best practices.

## Usage

### Run

```sh
./scripts/run.sh
```

Builds the Docker image and starts the server on `http://localhost:8080`.

### Tests

```sh
./scripts/tests.sh
```

Builds the Docker image and runs all tests inside the container.

### Without Docker

| Command              | Description                                        |
| -------------------- | -------------------------------------------------- |
| `make build`         | Build the binary to `bin/server`                   |
| `make serve`         | Start the server                                   |
| `make serve-watch`   | Start the server with live reload (requires `air`) |
| `make test`          | Run all tests                                      |
| `make test-no-cache` | Run all tests bypassing the cache                  |
| `make test-watch`    | Run tests in watch mode (requires `gotestsum`)     |
| `make lint`          | Run the linter (requires `golangci-lint`)          |

## APIs

The server exposes the following APIs:

```
GET /products - get all the available products
GET /products/{id} - get a product by id
POST /products - create a new product
PUT /products/{id} - update a product
DELETE /products/{id} - delete a product

GET /orders - get all the available orders
GET /orders/{id} - get an order by id
POST /orders - create a new order
PUT /orders/{id} - update an order
DELETE /orders/{id} - delete an order
```

The APIs are designed to give the consumer less flexibility on how to manage the order itself. The consumer, while creating the order, can only specify the product to include and the quantity of items, the rest is calculated by the backend when creating an order.

### Examples

**List products**

```sh
curl http://localhost:8080/products
```

**Get a product by id**

```sh
curl http://localhost:8080/products/1
```

**Create a product**

```sh
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Keyboard","description":"Mechanical keyboard","price":15000,"vat_rate":0.22}'
```

**List orders**

```sh
curl http://localhost:8080/orders
```

**Get an order by id**

```sh
curl http://localhost:8080/orders/1
```

**Create an order**

```sh
curl -X POST http://localhost:8080/orders \
  -H "Content-Type: application/json" \
  -d '{"items":[{"product_id":1,"quantity":1},{"product_id":3,"quantity":2}]}'
```

### Considerations

Another possible structure would have exposed the CRUD apis for `/orders` and `/orders/{id}/items`, this would have required more time but it's more in line with the RESTful approach, that divides the entity on different APIs to give the consumer more flexibility on the usage.

## Price/VAT Calculations

I've decided to store the price as cents instead floating numbers to avoid rounding errors. Integer overflow should not be a problem here.

I've also decided to validate each entity and make the price calculations inside the `entity` itself, it makes possible to test every edge case regarding the calculation itself.

## Tests

I followed the general Go heuristic approach of `don't mock what you own and can run cheaply`.

I am used to mocking the internal dependencies of a handler/service/repository, but that ends up testing the implementation more than the logic itself. Go makes it easy to swap the underlying configuration, so using an in-memory db instead of a real database is straightforward.

I ended up changing the implementation multiple times without changing the tests (e.g. moving validation from the service to the entities), which ensures that previously implemented requirements stay in place after refactoring.

### Considerations

I did not implement any integration tests, since I did not mock anything in the handlers layer, so the tests cover the logic and the integration with the db as well.

What is not covered is verifying that the server starts correctly. I chose not to implement this for the demo, but in a proper setup a smoke test suite would run post-deploy to confirm the server is up and correctly configured.

## Storage

I decided to go for an in-memory `sqlite` database for simplicity. The database is reset on every restart and requires no particular configuration. Note that SQLite does not enforce foreign keys by default — they require explicitly enabling via `PRAGMA foreign_keys = ON`.

The server populates the db with some initial data on startup.

### `products`

| Column        | Type     | Constraints               |
| ------------- | -------- | ------------------------- |
| `id`          | INTEGER  | PRIMARY KEY AUTOINCREMENT |
| `name`        | TEXT     | NOT NULL                  |
| `description` | TEXT     | NOT NULL DEFAULT `''`     |
| `price`       | INTEGER  | NOT NULL                  |
| `vat_rate`    | REAL     | NOT NULL DEFAULT `0.22`   |
| `created_at`  | DATETIME | NOT NULL                  |
| `updated_at`  | DATETIME | NOT NULL                  |

### `orders`

| Column        | Type     | Constraints               |
| ------------- | -------- | ------------------------- |
| `id`          | INTEGER  | PRIMARY KEY AUTOINCREMENT |
| `total_price` | INTEGER  | NOT NULL                  |
| `total_vat`   | INTEGER  | NOT NULL                  |
| `created_at`  | DATETIME | NOT NULL                  |

### `order_items`

| Column         | Type    | Constraints                      |
| -------------- | ------- | -------------------------------- |
| `id`           | INTEGER | PRIMARY KEY AUTOINCREMENT        |
| `order_id`     | INTEGER | NOT NULL REFERENCES `orders(id)` |
| `product_id`   | INTEGER | NOT NULL                         |
| `product_name` | TEXT    | NOT NULL                         |
| `quantity`     | INTEGER | NOT NULL                         |
| `unit_price`   | INTEGER | NOT NULL                         |
| `vat_rate`     | REAL    | NOT NULL                         |
| `price`        | INTEGER | NOT NULL                         |
| `vat`          | INTEGER | NOT NULL                         |

I decided to duplicate the data from the `products` into a `order_items` table, so it would be possible to get a snapshot of the order when it was created. Referencing only the id of the product without duplicating the data would have resulted in a side effect when editing the product itself.

### Considerations

In the future, of course, changing the db into a PostgreSQL would be the minimum upgrade possible.

Even though, for projects like this one where you know the access pattern in advance, I like to use a DynamoDB (if I remember correctly you are currently on AWS), since it's extremely cheap and it scales really well with zero costs. I'll put a summary of the structure of the tables using dynamo db, just to be complete:

```json
{
  "Product": {
    "PK": "PRODUCT#<uuid>",
    "SK": "PRODUCT#<uuid>",
    "entity_type": "PRODUCT",
    "id": "<uuid>",
    "name": "string",
    "description": "string",
    "price": "number (cents)",
    "vat_rate": "number",
    "created_at": "<isodate>",
    "updated_at": "<isodate>"
  },
  "Order": {
    "PK": "ORDER#<uuid>",
    "SK": "ORDER#<uuid>",
    "entity_type": "ORDER",
    "id": "<uuid>",
    "total_price": "number (cents)",
    "total_vat": "number (cents)",
    "created_at": "<isodate>"
  },
  "OrderItem": {
    "PK": "ORDER#<order_uuid>",
    "SK": "ITEM#<item_uuid>",
    "entity_type": "ORDER_ITEM",
    "id": "<item_uuid>",
    "order_id": "<order_uuid>",
    "product_id": "<product_uuid>",
    "product_name": "string",
    "quantity": "number",
    "unit_price": "number (cents)",
    "vat_rate": "number",
    "price": "number (cents)",
    "vat": "number (cents)"
  }
}
```

By querying for `PK = ORDER#<order_uuid>` you'll get both the order and the order items.
