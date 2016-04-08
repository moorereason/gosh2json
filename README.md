# gosh2json

Google Spreadsheet to JSON converter written in Go.

![Project Status](https://img.shields.io/badge/status-beta-yellow.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/moorereason/gosh2json)](https://goreportcard.com/report/github.com/moorereason/gosh2json)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`gosh2json` is a simple tool that converts the content of a Google Spreadsheet into a JSON object.

## Usage

```
Usage of gosh2json:
  -c string
        path to JSON config file
  -i string
        spreadsheet ID
  -p    pretty print JSON
```

## Example

```
$ gosh2json -i 34ZfqfhHiynWxs2TT6ocpX4tm7D3T1nkI2rRL3sgZkPM -c client.json -p
{
  "Sheet1": [
    [
      "Date",
      "Name",
      "Title",
      "Website",
    ],
    [
      "Jan 19",
      "Tom Johnson",
      "Life Groups Pastor",
      "example.org"
    ],
    [
      "Jan 26",
      "Russ Whitehead",
      "Pastor",
      "example.org"
    ],
    [
      "Feb 2",
      "Juan Juarez",
      "Missions Director",
      "www.example.edu"
    ]
  ]
}
```

## License

`gosh2jon` is licensed under the MIT license.  See [LICENSE](LICENSE) file for
details.
