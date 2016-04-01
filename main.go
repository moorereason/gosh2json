package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/Iwark/spreadsheet"
	"github.com/Jeffail/gabs"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
)

func main() {
	sheetID := flag.String("i", "", "spreadsheet ID")
	config := flag.String("c", "", "path to JSON config file")
	pretty := flag.Bool("p", false, "pretty print JSON")
	flag.Parse()

	if *sheetID == "" {
		flag.PrintDefaults()
		return
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	client, err := newClient(*config)
	if err != nil {
		log.Fatal(err)
	}

	err = generate(client, *sheetID, *pretty, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}

// newClient returns a new HTTP client.  If no config file is passed, return a
// default client; otherwise, return a JWT-enabled client using the given JSON
// config file.
func newClient(config string) (*http.Client, error) {
	if config != "" {
		return newJWTClient(config)
	}
	return newDefaultClient()
}

// newJWTClient returns a new JWT-enabled HTTP client.
func newJWTClient(config string) (*http.Client, error) {
	clientData, err := ioutil.ReadFile(config)
	if err != nil {
		return nil, err
	}

	conf, err := google.JWTConfigFromJSON(clientData, spreadsheet.Scope)
	if err != nil {
		return nil, err
	}

	return conf.Client(context.TODO()), nil
}

// newDefaultClient returns a default HTTP client.
func newDefaultClient() (*http.Client, error) {
	fmt.Printf("Default\n")
	return google.DefaultClient(context.TODO(), spreadsheet.Scope)
}

// getSheets retrieves a spreadsheet by ID and returns it as a
// spreadsheet.Worksheets object.
func getSheets(client *http.Client, ID string) (*spreadsheet.Spreadsheet, error) {
	s := &spreadsheet.Service{Client: client}
	return s.Get(ID)
}

// generate is the main entry point in generating the resulting resources JSON
// content.  The results are written to the given io.Writer interface.
func generate(client *http.Client, ID string, pretty bool, w io.Writer) error {
	sheets, err := getSheets(client, ID)
	if err != nil {
		return err
	}

	jsonObj, err := getRecords(sheets)
	if err != nil {
		return err
	}

	var b []byte
	if pretty {
		b = jsonObj.BytesIndent("", "  ")
	} else {
		b = jsonObj.Bytes()
	}
	_, err = w.Write(b)
	if err != nil {
		return err
	}

	return nil
}

// getResourcesData parses a given set of Worksheets and returns a slice of
// Resource objects.
func getRecords(sheets *spreadsheet.Spreadsheet) (*gabs.Container, error) {
	obj := gabs.New()

	for i := 0; i < len(sheets.Worksheets); i++ {

		s, err := sheets.Get(i)
		if err != nil {
			return nil, err
		}

		_, err = obj.ArrayOfSize(s.MaxRowNum, s.Title)
		if err != nil {
			return nil, err
		}

		for j, r := range s.Rows {
			_, err = obj.S(s.Title).ArrayOfSizeI(s.MaxColNum, j)
			if err != nil {
				return nil, err
			}

		CellLoop:
			for k, c := range r {
				if c == nil {
					continue CellLoop
				}

				_, err = obj.S(s.Title).Index(j).SetIndex(c.Content, k)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return obj, nil
}
