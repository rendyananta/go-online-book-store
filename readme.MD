# Online Book Store Project
This is a golang project of online book store backend API. 
Database used in this project is SQLite due to its portability and simplicity and zero dependencies/installation. 
This project layout structure designed to have the best experience of running go service which embrace defaults. 
The only configuration we need to set is the application security key that will be used to encrypt and decrypt the data.

The project layout are highly inspired by Clean Architecture that divides the structure into three main components such as repository, use case, and presenter. 
These components live under `internal/` directory as `repo`, `usecase`, and `presenter` respectively.
- Repo
  - Repo layer is responsible for the interactions of data fetching related process. 
  - Relation and polymorphic data shape need to be handled and fetched in this layer. 
  - API Call and any other infrastructure calls should be live in the repository layer. 
- Use Case
  - Use case layer is responsible for the main business logic of the application. If the business process needs the data, calling repository would be enough.
- Presenter
  - Presenter layer is responsible to mutate any user request, whether it is HTTP or RPC based protocol
    to the use case format. Any user request validations should live here, as the clean data will be processed in the use case layer.

Besides of the app main structure, this project also separates the usable tools into the `pkg/` directory. Packages that lives in the `pkg` directory
responsible for creating a nice little abstraction to reduce vendor lock in. In this project we provide database connection manager to manage multiple database connections and drivers.

Another example in this project have auth manager and cache manager to manage auth and user session, as well as validate their token. 
This small amount of abstraction is less likely to change in the future due to business process change, which makes the best use case to place the auth manager into this package.

## Dependencies
- [sqlite3](https://github.com/mattn/go-sqlite3/) 
- [go-validate](https://github.com/go-playground/validator)
- [sqlx](https://github.com/jmoiron/sqlx)
- [uuid](https://github.com/google/uuid)
- [gomock](https://github.com/golang/mock)

API Features:
- User registration
- User authentication
- Get All Books using cursor
- Search books by using text
- Place an order of the book
- Review user orders

Technical Features (Future?):
- Well-defined data structure that can be developed further with minimum amount of existing code changes.  
  - User:
    - To add admin capability, we can add a new staff / seller table
  - Order related:
    - Payments can be handled by adding new payment-related status and order_transactions table
      - Payment process can contain the payment details (gateway management)
    - Promotion
      - By utilizing order lines polymorphic data definition, we can add discount as a new order line.
      - To calculate the discount, it needs to introduce the discount rule engine that responsible to check and store
        discount rule, in which attributes that can apply to the voucher.
    - Platform Fee
      - By utilizing order lines polymorphic data definition, we can platform fee, insurances, additional handling or any features fee as a new order line.
    - Deliveries can be handled by adding new delivery-related status
      - By utilizing order lines polymorphic data definition, we can add shipping fee as a new order line.
      - By introducing delivery capability, it implies that the app should support the addresses both of the users and the warehouses (inventory)  
      - Delivery process can contain the logistic partner integration (shipping fee, booking and tracking gateway)

Things can be Improved:
- Full Text search using sqlite [fts5 extension](https://www.sqlite.org/fts5.html), or using any database support full text search capability. Currently, the app using `LIKE` where clause which
  degrades search performances. 
- Metric and traces can be implemented. 
- Caching layer in the book repository, currently the database is using SQLite on NVME machine which make it less important to use separate cache layer.
  - SQLite in NVME-based machine is fast enough, so in memory cache solution such as redis can be implemented later if bottleneck happened.
- Better database capability, currently this project only focuses on SQLite due to its simplicity and portability. 
  As the user increases, the needs of battle tested and scalable database solution such as postgresql or mysql is mandatory. 
  Need to retest the query, since the connection abstracts has been built under database connection manager package.  
- Migration CLI tool, currently any SQL query changes requires to rebuild and re-run the database migrations (drop all tables and recreate them).
  - This tool currently doesn't support incremental up and down commands.
  - Migrate to golang-migrate? SQLite seems not supported due to CGO binding that need to be enabled to compile the library [issue](https://github.com/mattes/migrate/issues/165).

## How to use
### Makefile!!
Makefile is essential tool for building binary.  

`make db-refresh` command to refresh the database migration.
`make db-seed` command to seed the data using given csv files in the repository. You don't need to because the app already ship with sqlite database included.
`make run-http` command to run the app.
`make test` command to run unit tests inside the entire app.

You only need to run make run-http to execute the app. The server will be served on localhost:8080.

## Sample requests

### User registration
```shell
curl --request POST \
  --url http://localhost:8080/auth/token \
  --header 'Content-Type: application/json' \
  --data '{
	"email": "rendy@email.com",
	"password": "password"
}'
```

### User login / request for token
```shell
curl --request POST \
  --url http://localhost:8080/auth/token \
  --header 'Content-Type: application/json' \
  --data '{
	"email": "rendy@email.com",
	"password": "password"
}'
```

### Get all books
```shell
curl --request GET \
  --url http://localhost:8080/books
```

### Get all books
```shell
curl --request GET \
  --url http://localhost:8080/books
```

### Get all books for the next page
```shell
curl --request GET \
  --url 'http://localhost:8080/books?last_id=01926c92-1827-7817-975f-5f5d6db677af'
```

### Get all books for the next page
```shell
curl --request GET \
  --url 'http://localhost:8080/books?last_id=01926c92-1827-7817-975f-5f5d6db677af'
```

### Search books
```shell
curl --request GET \
  --url 'http://localhost:8080/books/search?q=hunger'
```

> For logged-in user requests, you may need to install jq as it is a dependencies to querying a JSON response in a shell.

> installing on mac os: `brew install jq`

> installing on arch linux: `sudo pacman -S jq`

> installing on debian-based linux distro: `sudo apt install jq`

### Place order
```shell
curl --request POST \
  --url http://localhost:8080/orders/place \
  --header "Authorization: Bearer $(curl --request POST --url http://localhost:8080/auth/token \
                                              --header 'Content-Type: application/json' \
                                              --data '{"email": "rendy@email.com","password": "password"}' | jq  ".data.token" | tr -d '"')" \
  --header 'Content-Type: application/json' \
  --data '{
	"lines": [
		{
			"line_reference_id": "01926c92-1843-7b9e-a7a9-1b4accc9bdcd",
			"line_reference_type": "book",
			"quantity": 1
		},
		{
			"line_reference_id": "01926c92-189a-79e5-b7b5-d6f46e30dd0a",
			"line_reference_type": "book",
			"quantity": 3
		}
	]
}'
```

### User orders
```shell
curl -v --request GET \
  --url http://localhost:8080/orders \
  --header "Authorization: Bearer $(curl --request POST --url http://localhost:8080/auth/token \
                                              --header 'Content-Type: application/json' \
                                              --data '{"email": "rendy@email.com","password": "password"}' | jq  ".data.token" | tr -d '"')"
```

### Order details
```shell
curl -v --request GET \
  --url http://localhost:8080/orders/01926cb0-bdd5-7cad-aeaa-cb2764c010a6 \
  --header "Authorization: Bearer $(curl --request POST --url http://localhost:8080/auth/token \
                                              --header 'Content-Type: application/json' \
                                              --data '{"email": "rendy@email.com","password": "password"}' | jq  ".data.token" | tr -d '"')"
```

Dataset used:
https://www.kaggle.com/datasets/thedevastator/comprehensive-overview-of-52478-goodreads-best-b
