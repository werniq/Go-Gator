## Go Gator
### News Aggregator Server, implemented in golang
<hr>

## Implemented:
1. Dynamic reading from different feeds
2. Storing that data into sorted files
3. Reading from that files whenever user asked for the news
4. Filtering that news by parameters
5. Logging

## Documentation
Here I would like to explain the purpose of each directory.
<br />
### 1. Certs 
Certificates used to run server in secure mode (Enabling HTTPs)

### 2. Cmd
Main logic of the application is there. Cmd has few sub folders:
1. Parsers  - Parsing various data types
2. Parsing Instructions - Filtering news by various parameters
3. Types - Types folder is used to avoid cycling imports. Object that are used in multilpe packages (e.g. Parsing Parameters) lives there
4. Utils - Helper functions, in order to separate them with main logic
5. Server - server initialization and configuration
6. Server/handlers - handlers attached to the server

### 3. Docs 
Documentation, specfile and C4 model.

Available parameters: <br/>
> `ts-from 2024-05-12` News will be retrieved starting from that timestamp <br/>
> `ts-to 2024-05-18` No news will be retrieved, where publication date is bigger than provided parameter <br/>
> `sources bbc,washingtontimes` News will be retrieved ONLY from mentioned sources (separated by ',') <br/> 
> `keywords Ukraine,Chine` News will be filtered by existence of keywords in title or description <br/>

