## Go Gator
### News Aggregator Server
<hr>

## Implemented:
1. Handlers for managing sources (Admin API)
2. Handler for retrieving news by given parameters:
3. Verifying user input, displaying appropriate messages/errors
4. Dynamic fetching news from external APIs:
    - Making requests to external APIs (ABC,BBC, etc.), and store it. Not from static files as were before
5. Covered code with unit and integration tests
6. Made server secure with HTTPS
7. Added logging
8. Containerized application with docker

## Documentation
Here I would like to explain the purpose of each directory.
<br />

### 1. Cmd
Main logic of the application is there. Cmd has few sub folders:
1. Parsers  - Parsing various data types
2. Filters - Filtering news by various parameters
3. Types - Types folder is used to avoid cycling imports. Object that are used in multiple packages (e.g. Parsing Parameters) lives there
4. Server - server initialization and configuration
5. Server/handlers - handlers attached to the server
6. Validator - Validating layer using chain of responsibility pattern

### 2. Docs 
Documentation, specfile and C4 model, and usage/response examples with images 

## Server Handlers
Go-Gator server contains few handlers. Few of them are used by admins to manage list of available
sources. <br />
News handler is used by clients to fetch news.
P.S. This example assumes that server is running on port :443, however you can change it any time.

1. GET: `/news` - Returns list of news, filtering them by parameters.
- Available parameters: <br/>
> `ts-from 2024-05-12` News will be retrieved starting from that timestamp <br/>
> `ts-to 2024-05-18` No news will be retrieved, where publication date is bigger than provided parameter <br/>
> `sources bbc,washingtontimes` News will be retrieved ONLY from mentioned sources (separated by ',') <br/>
> `keywords Ukraine,Chine` News will be filtered by existence of keywords in title or description <br/>

- Request example: 
![img.png](docs/images/get_news_request.png)

- Response example:
  ![img_2.png](docs/images/get_news_response.png)

2. GET: `/admin/sources` - Returns list of available sources

- Request example: 
![get_sources_request.png](docs/images/get_sources_request.png)

- Response example:
![img_1.png](docs/images/get_sources_response.png)

3. POST `/admin/sources` - Add new sources to the list <br />
If were provided already existing source - will return an error.

- Request example: 
![img_2.png](docs/images/register_source_request.png)

- Response example:
![img_3.png](docs/images/register_source_response.png)

4. PUT '/admin/sources' - Update already existing sources <br />
In source, you can update either format, and/or endpoint. 
If were provided not-existing source - will return an error 

- Request example:
![img_4.png](docs/images/put_source_request.png)

- Response example:
![img_5.png](docs/images/put_source_response.png)

5. DELETE '/admin/sources' - Update already existing sources <br />
If were provided not-existing source - will return an error

- Request example:
![img_6.png](docs/images/delete_source_request.png)

- Response example:
![img_7.png](docs/images/delete_source_response.png)

## Usage:
1. Using Golang: <br />
> `go build -o ./bin/go-gator` - Build golang binary <br />
> `./bin/go-gator`
You can change default parameters, such as updates frequency, server port and certificates:
1. -f - Changes news updates frequency
2. -p - Specify port on which server will be operating
3. -c and -k - Are used for SSL certificate and key

