# Sanntid

Go application for retrieving realtime arrival data for public transportation in Norway.
Data is provided by the [Ruter API](https://ruter.no/labs/).

![Screenshot](https://raw.githubusercontent.com/michaelenger/sanntid/master/screenshot.png)

## Installation

Get the package.

```shell
go get github.com/michaelenger/sanntid
```

## Usage

Run the `sanntid` application with the location name as the first parameter to get the arrivals for that location.

```shell
sanntid "alexander kiellands plass"
```
