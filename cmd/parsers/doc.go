// Package parsers is used to parse data in multiple formats. By its idea package parsers is limited to
// parse incoming data (decoded files, or array of bytes) into array of articles.
//
// This package provides a factory to create parsers, which implement the Parser interface.
// Objects implemented this interface will have a method Parse which will parse the data into an array of articles.
// All these objects should have a source string as private field, in order to be able to map the source name
// to the desired file.
//
// Factory can create objects to parse RSS, HTML or JSON data. Because of factory approach it is
// easier to update code and add new parsers.
// They will be used in other parts of the program to decode data into an array of articles.
package parsers
