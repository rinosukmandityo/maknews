# Maknews
Sample of hexagonal architecture to handle news creation and news retrieval  
Tech stack for this sample will be like this:

1. Kafka for message broker
2. Elasticsearch for searching
3. Redis for result caching
4. MySQL to save data
5. Go Chi for HTTP endpoint router

Prerequisites
---
You need docker compose installed in your local machine.  
And in this project already included docker compose file to setup following server:  

- Kafka  
- ElasticSearch  
- Redis  
- MySQL

To run all those servers on your local just go to this root directory project and run this following command:  
`docker-compose up -d` 

But in case this docker compose is not working, you can install all those servers as a standalone application.

How to run
---
### Set Environment Variable
This application support 2 kind of database MySQL and MongoDB to prove our ports is completely agnostic from the implementation.  
By default it will connect into our MySQL database with default host & port `127.0.0.1:3306` and database `news`.  
To connect into different database we need to set database information in environment variable like following example:

##### Persistence Database
	
1. MongoDB
```cli
set url=mongodb://localhost:27017/local   
set timeout=10   
set db=local   
set driver=mongo   
```
2. MySQL
```cli
set url=root:root@tcp(127.0.0.1:3306)/news  
set timeout=10  
set db=news  
set driver=mysql  
```

##### Cache Database
```cli
set redis_url=redis://:@localhost:6379/0  
set redis_timeout=10  
```

##### Elasticsearch
```cli
set elastic_url=http://localhost:9200  
set elastic_timeout=10  
set elastic_index=news  
```

##### Message Broker
```cli
set kafka_url=localhost:9092  
set kafka_timeout=10  
set kafka_topic=news  
```

After setting the database information we only need to run the main.go file  
`go run main.go`  

### API List & Payloads
Here is our API List and its payload:  

1. [GET] **/news?offset=0&limit=10**  
`/news?offset=0&limit=10`
2. [POST] **/news**  
```javascript
{
	ID: 	 15,
	Author:  "Rest",
	Body: 	 "Hello this is news from REST",
	Created: "2020-03-01T22:59:59.999Z"
}
```
3. [PUT] **/news/{_news\_id_}**  
`/news/15`
```javascript
{
	ID: 	 15,
	Author:  "Rest",
	Body: 	 "Hello this is news from REST",
	Created: "2020-03-01T22:59:59.999Z"
}
```
4. [DELETE] **/news/{_news\_id_}**  
`/news/15`

### The service that we are going to build  

We have our service which is a news and it will connect to serializer which will either serialize the data into json or message pack before serving it through REST API  
And then on the other side we have our repository which will either choose to use MySQL or MongoDB based on how we start the application from command line.  
So basically our API will be able to accept JSON or message pack format and also our repository is able to use both MySQL and MongoDB and it won't really affect our service  

#### Table Structure
Here is table structure for MySQL table:  
- id INT  
- author TEXT  
- body TEXT  
- created TIMESTAMP

#### Apps Flow
The apps flow would be like this:

1. Create news using  [POST] /news url and it will be sent to kafka producer
2. Kafka consumer will get the data from kafka producer and will store the complete data into mySQL database and for ID & created data will be stored in ElasticSearch (ES)
3. If we want to retrieve the data we can use [GET] /news:
	- it will fetch the data from redis and return the data to user
	- if data in redis already expired or empty it will fetch the data from elasticsearch
	- data get from elasticsearch will have offset and limit and it will be ordered descending by date creation (created field)
	- after get data from elasticsearch, it will fetch the data from database one by one using go routine worker
	- after get the data from database it will store the data into redis as a cache data

Project Structure
---
By implementing Hexagonal Architecture we also implement Dependency Inversion and Dependency Injection. Here is some explanations about project structure:

1. **api**  
contains handler for API
2. **models**  
contains data models
3. **repositories**  
contains **Port** interface for repository adapter
   - **mysql**  
contains mySQL **Adapter** that implement NewsRepository interface. This package will store mySQL client and connect to mySQL database to handle database query or command. Complete news data will be stored here.
	- **mongodb**  
contains mongoDB **Adapter** that implement NewsRepository interface. This package will store mongoDB client and connect to mongoDB database to handle database query or command. Complete news data will be stored here.
   - **redis**  
contains redis **Adapter** that implement CacheRepository interface. This package will store redis client and connect to redis server to handle database query or data manipulation
   - **elasticsearch**  
contains elasticsearch **Adapter** that implement ElasticRepository interface. This package will store elasticsearch client and connect to elasticsearch server to handle database query or command. ID and news date creation will be stored here.
   - **kafka**  
contains kafka **Adapter** that store kafka connection and has several methods to handle message write and message read from kafka server.
4. **serializer**  
contains **Port** interface for decode and encode serializer. It will be used in our API to decode and encode data.
   - **json**  
contains json **Adapter** that implement serializer interface to encode and decode data
   - **msgpack**  
contains message pack **Adapter** that implement serializer interface to encode and decode data
5. **services**  
contains **Port** interface for our domain service and logic 
6. **logic**  
contains service **Adapter** that implement service interface to handle service logic like constructing repository parameter and calling repository interface to do data manipulation or query