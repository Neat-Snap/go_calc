# Calculation server on Golang

_This project was created as a final task for the first sprint in Yandex Lyceum Golang educational program._
# Usage

### Specify the port (optional)
Specify the port you want to run the server on by setting the "PORT" environment variable to the needed value. If not specified, the server will be running on 8080.
### Run the server
To run the application install dependencies (Golang (1.2 or higher recommended)) and type:
```bash
    go run cmd/main.go
```

You will see something like this:
```
2024/11/22 20:40:07 Port not specified, defaulting to 8080
2024/11/22 20:40:07 Starting server on port 8080
```

## Calculations
To calculate a problem you can use GET method with url query parameter or POST method.
### Get
To calculate using GET method and url query parameters you can use the following command:
```
127.0.0.1:8080/expression/?expression=your_expression
```

_Notice that in this method instead of "+" sign you have to use "%2B" because of how urls are interpreted_
### Post
To calculate using POST method you will have to send a request to the following url:
```
127.0.0.1:8080/expression
```
and provide the requests with the payload. For example like this:
```json
{"expression":"your_expression"}
```
In comparison with previous method, here you can use "+" sign as usual.