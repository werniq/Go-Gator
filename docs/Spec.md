- Project Name: Go-Gator
- Author: Oleksandr Matviienko

## Purpose

This project is created in order to simplify the process of gathering news from different sources.

## Goals 

Goals for this project is to simplify life for journalists, by making a centralized 
solution for gathering different news.
In addition, while implementing this service, I am improving my Cloud & Golang skills.

## API Design

This API will work in following way:
1. It retrieves and stores data from files
2. Then program is parsing that data into an array of news
3. Afterward, in case where user have provided arguments for filtering, news are filtered by that arguments
4. Displays all news that passed validation

News model looks like this:

```
type News struct {
   Title            string
   Description      string
   PublishedDate    time.Time  
}
```
Fields here represent respective properties of article: headline, description, and publication date.

#### Filtering news
Currently, go-gator supports filtering by keywords, publication date range, and sources.

1. **Keyword Filtering**: This filter checks if the title or description of an article contains the specified keywords.
2. **Date Range Filtering**:
    - **Start Date**: Displays news articles published from this timestamp onwards.
    - **End Date**: Displays news articles published up to this timestamp.
3. **Source Filtering**: Displays news articles only from selected public service broadcasters.
