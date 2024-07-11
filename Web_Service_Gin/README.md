# Web Service
This project is for create Restful API
Reference form **[This](https://go.dev/doc/tutorial/web-service-gin)**.

## Format to GET or POST in PowerShell using crul.exe

### GET

```
curl.exe http://localhost:8080/albums `
    --header "Content-Type: application/json" `
    --request "GET"
```
### POST

``` 
curl.exe http://localhost:8080/albums `
    --include `
    --header "Content-Type: application/json" `
    --request "POST" `
    --data '{\"id\": \"4\", \"title\": \"The Modern Sound of Betty Carter\", \"artist\": \"Betty Carter\", \"price\": 49.99}'
```
