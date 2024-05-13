# News Feed

The purpose of this API is to build a simple RESTful microservice that will read RSS feeds from a public source and store it in Mongo for future queries

### GET /load

Accepts a query param `feedUrl` that will be used to load the articles. If not provided it will load from all default sources.

#### Default Sources
- https://feeds.bbci.co.uk/news/uk/rss.xml
- https://feeds.bbci.co.uk/news/technology/rss.xml
- https://feeds.skynews.com/feeds/rss/uk.xml
- https://feeds.skynews.com/feeds/rss/technology.xml

Example:

Request:

    curl -X GET http://localhost:8080/load'
    curl -X GET http://localhost:8080/load?feedUrl=https://feeds.skynews.com/feeds/rss/technology.xml


Response:

    [

        {
            "id": "https://www.bbc.co.uk/news/business-61644033",
            "title": "Could flat tyres soon be a thing of the past?",
            "description": "Airless tyres that do not puncture are getting close to market but some remain sceptical about them.",
            "link": "https://www.bbc.co.uk/news/business-61644033?at_medium=RSS&at_campaign=KARANGA",
            "source": {
                "category": "uk",
                "feedUrl": "https://feeds.bbci.co.uk/news/uk/rss.xml",
                "provider": "bbc"
            },
            "publishedDateTime": "2022-06-13T23:18:58Z"
        },
        {
            "id": "https://www.bbc.co.uk/news/business-61483491",
            "title": "Could nuclear desalination plants beat water scarcity?",
            "description": "Engineers are developing mobile, floating nuclear desalination plants to help solve water shortages.",
            "link": "https://www.bbc.co.uk/news/business-61483491?at_medium=RSS&at_campaign=KARANGA",
            "source": {
                "category": "uk",
                "feedUrl": "https://feeds.bbci.co.uk/news/uk/rss.xml",
                "provider": "bbc"
            },
            "publishedDateTime": "2022-06-20T23:16:31Z"
        },
        ...

    ]



### GET /find

Accepts query params or arbitrary JSON document as payload. The endpoint should find all of the articles based on the filter provided.

| Parameter     | Type     | Description                                                                        |
| ------------- | -------- | -----------------------------------------------------------------------------------|
| category      | string   | Article's category                                                                 |
| provider      | string   | Article's provider                                                                 |
| limit         | int      | Max Limit is 1000 per page                                                         |
| page          | int      | Index Page. (First Page is 0)                                                      |
| sort          | string   | Sort column. e.g publishedDateTime (You can sort by any article's model property)  |
| order         | string   | Sort order. e.g. asc (It defaults to asc)                                          |


Example:

Request:

    curl http://localhost:8080/find'
    curl http://localhost:8080/find?sort=publishedDateTime&order=desc&category=uk&provider=bbc'

Response:
    

    {
        "criteria": {
            "category": "uk",
            "provider": "bbc",
            "sort": "publishedDateTime",
            "order": "desc"
        },
        "articles": [
            {
                "id": "https://www.bbc.co.uk/news/uk-62874346",
                "title": "King Charles III promises to follow Queen's selfless duty",
                "description": "The King hears condolences at Westminster before travelling to Edinburgh to mount a vigil for the Queen.",
                "link": "https://www.bbc.co.uk/news/uk-62874346?at_medium=RSS&at_campaign=KARANGA",
                "source": {
                    "category": "uk",
                    "feedUrl": "https://feeds.bbci.co.uk/news/uk/rss.xml",
                    "provider": "bbc"
                },
                "publishedDateTime": "2022-09-12T12:47:57Z"
            },
            {
                "id": "https://www.bbc.co.uk/news/uk-scotland-62869534",
                "title": "King arrives in Edinburgh ahead of Queen procession and tributes",
                "description": "The public will be able to view the Queen's coffin after a procession and a service of remembrance.",
                "link": "https://www.bbc.co.uk/news/uk-scotland-62869534?at_medium=RSS&at_campaign=KARANGA",
                "source": {
                    "category": "uk",
                    "feedUrl": "https://feeds.bbci.co.uk/news/uk/rss.xml",
                    "provider": "bbc"
                },
                "publishedDateTime": "2022-09-12T12:43:26Z"
            },
            ...
        ],
        "total": 79
    }


## Getting Set Up

Before running the application, you will need to ensure that you have a few requirements installed;
You will need Go.

### Go
### Docker (optional)

[Go](https://golang.org/) is an open source programming language that makes it easy to build simple, reliable, and efficient software.

[Docker](https://docker.com/) is used to build and sharing containerized applications.

## Project Structure

Following [Standard Go Project Layout](https://github.com/golang-standards/project-layout). (apologies in advance if it is too much for the purpose of the assignment)

### `/cmd`

Main application for this project.

### `/internal`

Internal application logic

### `/pkg`

Model that is okay to be shared with external applications.

### `/script`

Scripts to perform various build, install, analysis, etc operations. (for this example we have only tha start.sh script to start the application)

## Server

It uses the standard HTTP request multiplexer.

### Running the server
    go run cmd/news/main.go

### Running the server using start.sh script
    ./script/start.sh

### Running using docker-compose

The docker-compose file that resides under the root folder. To run it just run the following command: `docker compose up --build`. (update the port and add `PORT` env variable if you are running in another port rather than the default port 8080)

### Test

[Table-driven tests using subtests](https://blog.golang.org/subtests) were used as the approach to reduce the amount of repetitive code compared to repeating the same code for each test and makes it straightforward to add more test cases.

[Testify](https://github.com/stretchr/testify) the [assert](https://github.com/stretchr/testify#assert-package) package assert provides a set of comprehensive testing tools for use with the normal Go testing system.

## Running the tests

    go test ./... -v

[@maxalencar](https://github.com/maxalencar)