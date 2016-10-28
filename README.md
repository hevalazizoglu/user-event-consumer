# User Event Consumer

This application contains a RESTful API that will allow us to process user generated data and an endpoint to observe response times of this API.

### Request Handling
Golang http package has been used for server side implementation. However, since URL routing capabilities of this package is limited, implementation has been enhanced with "gorilla/mux". To get the package
```
go get -u github.com/gorilla/mux
```

### Data Storage
MongoDB has been used as main datastore for both user events and response times. Reasoning:
  - It is schemaless which enables us to store events that can change in terms of content in the same collection
  - It is scalable and high performant
  - New versions has lookup feature will let us to join collections that means although we're grouping events based on clients, we can still do combined analysis on multiple clients based on this data.
To get Golang client:
```
go get gopkg.in/mgo.v2
```

However, Go client is still very low level.

### Histogram - Response Time Visualization
I used a third party library, go chart. This library provides various charting options with easy use.
To install chart:
```
go get -u github.com/wcharczuk/go-chart
```

## How to Run
Setup GOPATH: put this line in .bash_profile
```
export GOPATH=PATH_TO_THIS_APPLICATION/user-event-consumer-application
```
Run inside project directory:
```
go run *.go
```

## How to Make Requests and Expected Responses
This project provides two endpoints:
- /user-event : To consume user events with POST requests
- /user-event/\_stats : To get response time stats Histogram

### /user-event
Sample curl request when running this application locally:
```
curl -XPOST 'localhost:8080/user-event' -d '{"APIKey": "A", "Timestamp": 1477602071, "UserID": 1}' -i
```
Expected response:
```
HTTP/1.1 200 OK
Content-Type: text/plain; charset=utf-8
Date: Fri, 28 Oct 2016 10:46:21 GMT
Content-Length: 0
```

Error Cases:
- Invalid APIKey
```
curl -XPOST 'localhost:8080/user-event' -d '{"APIKey": "D", "Timestamp": 1477602071, "UserID": 7}' -i
```
Expected response:
```
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 28 Oct 2016 10:47:44 GMT
Content-Length: 15

Invalid APIKey
```

- Missing Required Fields
```
curl -XPOST 'localhost:8080/user-event' -d '{"APIKey": "A", "Timestamp": 1477602071}' -i
```
Expected response:
```
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Fri, 28 Oct 2016 10:48:37 GMT
Content-Length: 24

Missing Required Fields
```

### Known Clients
Currently known clients information is maintained in **clients.go** file. Since authorization implementation is beyond the scope of this project, I maintained a simple list to do verifications. Requests made with APIKeys not listed will not be processed.

### /user-event/\_stats
This endpoint consumes the response time data that has been created during the calls of /user-event endpoint.
To see histogram when running application, open browser then access:
```
http://localhost:8080/user-event/_stats
```

Example output:
!["Response Times"](http://i65.tinypic.com/wpbd.png)


## Assumptions and Details During Implementation
- I assumed Timestamp field in user events will not be used during this process and stored it as Unix time directly. It is easy to do conversion on application side after retrieval from database though.
- Client APIKey have been map to collection name directly. In some cases, this may be erroneous as invalid types or characters can be provided.
- Comments are being avoided unless extremely necessary as I do believe they reduce readability when not needed.
- I've learnt In Golang world, conventional spacing is 8 but I found 4 more readable for this project.
