// Package server is used for initialization, configuration, and execution of HTTPs server for the application.
//
// This package could be used just by running RunAndConf function, which initializes server,
// attaches paths and handlers to it, and runs a server on the user-defined or default port.
// Handlers are used to manage sources, and automatically update file storing them, or
// in order to retrieve latest news.
//
// The server is using the following paths:
// - /admin/sources - to manage sources
// - /news - to retrieve latest news
//
// DEPRECATED:
// FetchNewsJob was used to fetch and parse articles feeds, and then writes the parsed data to a
// JSON file named with the current date.
// This job is deprecated now, and it's functionality is moved to the cron_job folder.
// You can run this job in separate container, as a Kubernetes CronJob object.
//
// It was updated in the way that FetchNewsJob is running as separate container, and uses PersistentVolume for storage.
// That means that data is stored in the volume, and it is not lost after the container is deleted.
package server
