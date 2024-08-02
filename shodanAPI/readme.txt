You need a shodan account to use this, go to https://www.shodan.io and create one.
Login and go to https://account.shodan.io/ , then you will see the "API Key" field, click in "show" to get it.
Create a .env file in /cmd/shodan/ and add the line:
SHODAN_API_KEY="yourApiKey"

Run: 
go run cmd/shodan/main.go [search_term]

For free, without configurations, you can use InternetDB API: 
go run cmd/idb/main.go [search_term]

