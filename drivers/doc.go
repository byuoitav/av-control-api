/*
Package drivers provides the interfaces that driver libraries implement to simplify the library code as well as the server code for the driver server.

Each interface has an associated create function:
	type Create(Interface)Func func(context.Context, string) ((Interface), error)
as well as an associated function to create a HTTP server:
	func Create(Interface)Server(create Create(Interface)Func) (Server, error)

A driver library provides a struct that implements the appropriate interface. For example, a driver library for the QSC DSP should implement the DSP interface.

This package includes subdirectories with all of the driver servers that BYU maintains. A driver server puts together dependencies (a driver library, logging libraries, auth libraries, etc) and runs an HTTP server using standard endpoints for each device type. It will provide a Create... function and should have its own go mod file to manage its dependencies versions. The examples below are examples of the main function of a driver server.
*/
package drivers
