// Package server is used for initialization, configuration, and execution of HTTPs server for the application.
//
// This package could be used just by running RunAndConf function, which initializes server,
// attaches paths and handlers to it, and runs a server together with concurrent FetchNewsJob.Run()
// Handlers are used to manage sources, and automatically update file storing them, or
// in order to retrieve latest news.
//
// FetchNewsJob makes request to news feeds with latest articles, then parses it to the array of types.News
// and writes that data to a filename with current data into user-defined or default data folder.
package server
