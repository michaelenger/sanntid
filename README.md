# Sanntid

Program written in Go for retrieving realtime arrival data from the [Ruter API](http://labs.trafikanten.no/how-to-use-the-api.aspx).

## Installation

Get the package.

```shell
go get github.com/michaelenger/sanntid
```

## Usage

Run the `sanntid` application with the location ID as the first parameter to get the arrivals for that location. See [here](http://193.69.180.119:8080/tabledump/stops2.csv) for a list of available arrival IDs.

```shell
sanntid 3010536
```
