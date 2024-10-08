/*
Package news_fetcher provides functionality to fetch and parse news articles from various sources.

This package is created to run a job that fetches news articles at a regular interval,
parses them, and stores them in JSON files. The main job is performed by the NewsFetchingJob type.

The primary function to execute the news fetching job is RunJob. This function
initializes and runs a NewsFetchingJob instance, which fetches and parses news articles,
and then writes the parsed data to a JSON file named with the current date.
It will help us to better retrieve that news later, since filenames are identifying news
from that particular date.

Types:

	NewsFetchingJob - Represents a job for fetching news articles.
	Contains a Date field to specify the job's timestamp.

Functions:

	RunJob - Initializes and runs a NewsFetchingJob, which parses data from feeds into respective files.

	(j *NewsFetchingJob) Execute - Fetches news, parses it, and writes the parsed data to a JSON file named

with the current date in the format YYYY-MM-DD.
*/
package main
