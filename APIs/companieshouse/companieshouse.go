package companieshouse

import (
	"fmt"
	"encoding/gob"
	"net/http"
	"github.com/dghubble/sling"
)




const baseURL = "https://api.companieshouse.gov.uk/"

// Issue is a simplified Github issue
// https://api.companieshouse.gov.uk/#response
type Company struct {
	Address1 	string	`json:"address_line_1"`
	Address2 	string	`json:"address_line_2"`
	Locality 	string	`json:"locality"`
	PostCode 	string	`json:"postal_code"`
}

// Session is an interface for typical sessions
type Session interface {
	Save(*http.Request, http.ResponseWriter) error
	Flashes(vars ...string) []interface{}
}

func init() {
	// Magic goes here to allow serializing maps in securecookie
	// http://golang.org/pkg/encoding/gob/#Register
	// Source: http://stackoverflow.com/questions/21934730/gob-type-not-registered-for-interface-mapstringinterface
	gob.Register(Company{})
}



type SearchError struct {
	Errors  []struct {
		ErrorValues []struct {
			Argument interface{} `json:"<argument>"`
		} `json:"error_values"`
		Location    	string `json:"location"`
		LocationType    string `json:"location_type"`
		Type    		string `json:"type"`
	} `json:"errors"`
}

func (e SearchError) Error() string {
	return fmt.Sprintf("%+v", e.Errors)
}


// CompanySearch provides methods for creating and reading issues.
type CompanySearch struct {
	sling *sling.Sling
}

// NewCompanySearch returns a new CompanySearch.
func NewCompanySearch(httpClient *http.Client) *CompanySearch {
	return &CompanySearch{
		sling: sling.New().Client(httpClient).Base(baseURL).SetBasicAuth("gYezs_qIXZGOCYByIXHPuxsAHJHJN0niqwvg54d4", ""),
	}
}


// ListByRepo returns a repository's issues.
func (s *CompanySearch) CompanyAddress(reg string) (Company, *http.Response, error) {
	company := new(Company)
	searchError := new(SearchError)
	path := fmt.Sprintf("company/%s/registered-office-address", reg)
	resp, err := s.sling.New().Get(path).Receive(company, searchError)
	if err == nil {
		err = searchError
	}
	return *company, resp, err
}


// Client is a tiny Github client
type Client struct {
	CompanySearch *CompanySearch
	// other service endpoints...
}

// NewClient returns a new Client
func NewClient(httpClient *http.Client) *Client {
	return &Client{
		CompanySearch: NewCompanySearch(httpClient),
	}
}

/*func main() {
	// Github Unauthenticated API
	client := NewClient(nil)
	company, resp, _ := client.CompanySearch.CompanyAddress("11483737")
	//fmt.Printf("HTTP Resp: %v\n", resp.Status[:3])
	//fmt.Printf("Error: %v\n", err)
	if resp.Status[:3] == "200" {
		fmt.Printf("Result: %v\n", company.POBox)
	} else {
		fmt.Println("Not Found")
	}
	


	
}*/
