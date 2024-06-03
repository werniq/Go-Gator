- Project Name: Go-Gator
- Author: Oleksandr Matviienko

## Summary

This project is created in order to simplify the process of gathering news from different sources.

## Motivation

Goals for this project is to simplify life for journalists, by making a centralized 
solution for gathering news from different sources.
In addition, while implementing this service, I will improve my Cloud & Golang skills.

## API Design

Our API will be used to send our clients filtered array of news containing specific keywords, 
starting or ending with specific date, or retrieved from particular resource.

This API will work in following way:
1. It will communicate with APIs of public service broadcasters, retrieving the newest news
2. Then program is parsing that data into an array of news
3. Afterward, if user have provided arguments for filtering - news are filtered by that parameters
4. Returns filtered array of news

### API Endpoints

1. Get Filtered News

- Endpoint: /api/news
- Request Method: GET
- Parameters: 
  - keywords (optional): Comma-separated keywords to filter articles
  - startDate (optional): Start date to filter articles (format: YYYY-MM-DD)
  - endDate (optional): End date to filter articles (format: YYYY-MM-DD)
  - sources (optional): Comma-separated list of news sources

- Response: 
```json
{
    "articles": [
        {
            "title": "Article Title",
            "description": "Article Description",
            "publishedDate": "2024-05-30T14:00:00Z",
            "source": "Source"
        }
    ]
}
```


## Step-By-Step Workflow

### Step 1: 
Program handles user input, and creates **Filters** model. It will be later used for filtering.

### Step 2:
Program opens files with prepared data, and parses it into array of **Article**.

### Step 3:
Pass an array of articles to filtering service. 
It will validate user input, and include only **Articles**
that passed validation.

### Step 4:
Display resulting array in the terminal.

### External APIs
We will use external APIs in order to retrieve latest news.

<hr />

## Models

We will have two main models in this program: News and Filters.
<br />

**Article entity** looks like this:

```go
type Article struct {
   Title            string
   Description      string
   PublishedDate    time.Time  
}
```

Fields here represent respective properties of article: headline, description, and publication date.
It will be used to actually display and/or return the resulting array of news

**Model filters** will be used in order to validate articles by certain conditions, 
and decide whether to remove it or leave.
If some article has not passed the validation - it is removed from array.

```go    
type FilteringParams struct {
   Keywords          string   
   StartingTimestamp string   
   EndingTimestamp   string   
   Sources           []string 
}
```

Currently, go-gator supports filtering by keywords, publication date range, and sources.

1. **Keyword Filtering**: This filter checks if the title or description of an article contains the specified keywords.
2. **Date Range Filtering**:
    - **Start Date**: Displays news articles published from this timestamp onwards.
    - **End Date**: Displays news articles published up to this timestamp.
3. **Source Filtering**: Displays news articles only from selected public service broadcasters.
