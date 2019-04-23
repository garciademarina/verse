# verse test

Write a REST API in the language/framework of your choice for a money-transferring application like Verse. It should have the following features:


* An endpoint to get the current balance of a given user.

* An endpoint to make a transfer of money between two users.


## Notes
The API can be used as admin or as regular user. 

- Admin endpoints are available under /admin/* and need an url parameter ?key=foobar. 
- User endpoints require a valid JWT as a url parameter ?jwt={jwttoken}

At startup the application will load some test-data and show some valid jwt for debug.

* Repository pattern, inmemory implementation for now.
* Balance and Amount fields are in cents. (int64)
* JWT middelware to validate and authenticate the user. 

## Api response
- Successfull request will return http status code 200 along with a json with additional information.
- Failed request will return http status code 40X along with a json with additional information of the error.
```
Type (string) Posible values api_error/authentication_error
Code (string) Optional code
Message (string) A string representation for the error
```

## Admin endpoints

* required url parameter ?key=foobar

#### GET /admin/balance/{userID}
Get the current balance of a given user. {userID}
**Arguments**
userID (required) A string ID of a user to get the balance
**Response**
num (string) A string ID of a user to withdraw the money 
DestinationUser (string) A string ID of a user to send the money
Amount (int64)

```json
{
    num: "D5DP9",
    balance: 10000
}
```

#### POST /admin/transfers 
Make a transfer of money between two users.
**Arguments**
Amount (required) A positive integer representing how much money to transfer
UserOrigin (required) A string ID of a user to withdraw the money 
UserDestination (required) A string ID of a user to send the money
**Response**
OriginUser (string) A string ID of a user to withdraw the money 
DestinationUser (string) A string ID of a user to send the money
Amount (int64) The amount of money transfered

#### GET /admin/accounts
Get all accounts.
**Response**
A list of all accounts.
```json
[
    {
        num: "62AC2",
        userID: "01D3XZ89NFJZ9QT2DHVD462AC2",
        name: "Rainbow account",
        openAt: "0001-01-01T00:00:00Z",
        balance: 10000
    },
    {
        num: "D5DP9",
        userID: "01D3XZ7CN92AKS9HAPSZ4D5DP9",
        name: "Billy account",
        openAt: "0001-01-01T00:00:00Z",
        balance: 10000
    },
    ...
]
```



## User endpoints

#### GET /balance 
Get the current balance of the user.
**Response**
num (string) A string ID of a user to withdraw the money 
DestinationUser (string) A string ID of a user to send the money
Amount (int64)

```json
{
    num: "D5DP9",
    balance: 10000
}
```





#### POST /transfers 
Make a transfer of money between two users.
**Arguments**
Amount (required) A positive integer representing how much money to transfer
UserOrigin (required) A string ID of a user to withdraw the money 
UserDestination (required) A string ID of a user to send the money
**Response**
OriginUser (string) A string ID of a user to withdraw the money 
DestinationUser (string) A string ID of a user to send the money
Amount (int64) The amount of money transfered






#### GET /user
Get user information.
**Response**
The user information.

```json
{
    id: "01D3XZ7CN92AKS9HAPSZ4D5DP9",
    name: "Billy",
    email: "Billy@example.com"
}
```









Example: 
```
curl -v -X POST "http://localhost:8080/transfers?jwt={jwtoken}" -H "accept: application/json;" -d "{ \"UserDestination\": \"01D3XZ7CN92AKS9HAPSZ4D5DP9\", \"Amount\": 400}" -H "Content-Type: application/json;"
```




Carlos García de Marina Vilar - garciademarina(at)gmail.com
