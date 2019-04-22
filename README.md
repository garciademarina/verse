# verse test

Write a REST API in the language/framework of your choice for a money-transferring application like Verse. It should have the following features:


* An endpoint to get the current balance of a given user.

* An endpoint to make a transfer of money between two users.


## TODO

* Transaction log (RequestID)
* Make idempotent all POST request to the API
* repository pattern, in the future implement sql db store

 
curl -i -X POST "http://localhost:8080/balance?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMDFEM1haM1pIQ1AzS0c5VlQ0RkdBRDhLRFIifQ.7IOWOubBcQ5e6LXV1kqvSX9hdGwEPl6hSfPhYuj4QSA" -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMDFEM1haM1pIQ1AzS0c5VlQ0RkdBRDhLRFIifQ.7IOWOubBcQ5e6LXV1kqvSX9hdGwEPl6hSfPhYuj4QSA" -H "accept: application/json;odata.metadata=minimal;odata.streaming=true"

curl -v -X POST "http://localhost:8080/transfers?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMDFEM1haM1pIQ1AzS0c5VlQ0RkdBRDhLRFIifQ.7IOWOubBcQ5e6LXV1kqvSX9hdGwEPl6hSfPhYuj4QSA" -H "authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMDFEM1haM1pIQ1AzS0c5VlQ0RkdBRDhLRFIifQ.7IOWOubBcQ5e6LXV1kqvSX9hdGwEPl6hSfPhYuj4QSA" -H "accept: application/json;odata.metadata=minimal;odata.streaming=true" -H "Content-Type: application/json;odata.metadata=minimal;odata.streaming=true" -d "{ \"UserDestination\": \"01D3XZ7CN92AKS9HAPSZ4D5DP9\", \"Amount\": 0.9}" 