package table

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var showing = []int{5, 10, 25, 50}
var replacer = strings.NewReplacer("%2C", ",", "%5B", "[", "%5D", "]")

// Info holds the pagination fields.
type Info struct {
	Page         int
	TotalPages   int
	PerPage      int
	Pages        []int
	Offset       int
	TotalResults int
	ShowResults  []int
	Condition    []Filter
	HTTPReq      *http.Request
	SQLString    string
}

//Filter condition for table
type Filter struct {
	Key       string
	Operation string
	Value     string
	Query     string
}

// New returns a pagination struct.
func New(r *http.Request) *Info {
	var err error
	info := &Info{
		ShowResults: showing,
		HTTPReq:     r,
	}

	show := r.URL.Query().Get("show")

	if len(show) > 0 {
		e := strings.Split(show, ",")
		if len(e) == 1 {
			info.Page, err = strconv.Atoi(e[0])
			if err != nil || info.Page < 1 {
				info.Page = 1
			}
			info.PerPage = 10
		}
		if len(e) > 1 {
			info.Page, err = strconv.Atoi(e[0])
			if err != nil || info.Page < 1 {
				info.Page = 1
			}
			info.PerPage, err = strconv.Atoi(e[1])
			if err != nil || info.PerPage < 1 {
				info.PerPage = 10
			}
		}
	} else {
		info.Page = 1
		info.PerPage = 10
	}

	if info.Page > 1 {
		info.Offset = (info.Page - 1) * info.PerPage
	}

	info.SQLString = fmt.Sprintf(` LIMIT %v OFFSET %v `, info.PerPage, info.Offset)

	return info
}

// CalculatePages calculates the number of pages by passing in the item total.
func (i *Info) CalculatePages(itemTotal int) {
	i.TotalPages = itemTotal / i.PerPage
	for x, val := range showing {
		if itemTotal == val && x < len(showing) {
			i.ShowResults = showing[:x]
		}
		if itemTotal > val && x+2 < len(showing) {
			i.ShowResults = showing[:x+1]
		}
	}
	if itemTotal%i.PerPage != 0 {
		i.TotalPages++
	}
	maxshow := 5
	if maxshow > i.TotalPages {
		maxshow = i.TotalPages
	}
	arr := make([]int, maxshow+2)
	if i.Page > 1 {
		arr[0] = i.Page - 1
	}
	if i.Page < i.TotalPages {
		arr[maxshow+1] = i.Page + 1
	}
	for pi := 0; pi < maxshow; pi++ {
		if i.Page <= i.TotalPages-maxshow+1 {
			arr[pi+1] = pi + i.Page
		} else {
			arr[pi+1] = i.TotalPages - maxshow + pi + 1
		}
	}

	i.Pages = arr

}

// Map returns a template.FuncMap PAGINATION which makes it easy to navigate
// between pages of results.
func Map() template.FuncMap {
	f := make(template.FuncMap)

	f["TABLE"] = func(option string, info Info, m map[string]interface{}) template.HTML {

		currentURI, ok := m["CurrentURI"]
		if !ok {
			log.Println("Issue")
			return template.HTML("Pagination could not load because CurrentURI is missing.")
		}
		var output string
		req := info.HTTPReq
		values := req.URL.Query()
		values.Del("show")
		var qs string
		if len(values) > 0 {
			qs = "&" + replacer.Replace(values.Encode())
		}
		show := info.PerPage

		if option == "PAGINATION" {

			pageleft := `<div class="btn-group pv5 pl30 pull-left">`
			pagemiddle := ""
			pageright := `</div>`

			for i, page := range info.Pages {
				var state string
				if page == 0 {
					state = "disabled"
				}
				if i == 0 {
					pagemiddle += fmt.Sprintf(`<a href="%v?show=%v,%v%v" class="btn btn-default dark  btn-sm %v"><i class="fa fa-chevron-left"></i></a>`, currentURI, page, show, qs, state)
				}
				if i > 0 && i < len(info.Pages)-1 {
					if page == info.Page {
						pagemiddle += fmt.Sprintf(`<a href="%v?show=%v,%v%v" class="btn btn-info light  active btn-sm">%v</a>`, currentURI, page, show, qs, page)
					} else {
						pagemiddle += fmt.Sprintf(`<a href="%v?show=%v,%v%v" class="btn btn-default  btn-sm">%v</a>`, currentURI, page, show, qs, page)
					}
				}
				if i == len(info.Pages)-1 {
					pagemiddle += fmt.Sprintf(`<a href="%v?show=%v,%v%v" class="btn btn-default dark light btn-sm %v"><i class="fa fa-chevron-right"></i></a>`, currentURI, page, show, qs, state)
				}
			}

			pagediv := pageleft + pagemiddle + pageright

			showleft := `<div class="btn-group pv5 pr20 pull-right">`
			showmiddle := ""
			showright := `</div>`

			for _, res := range info.ShowResults {
				if res == info.PerPage {
					showmiddle += fmt.Sprintf(`<a href="%v?show=%v,%v%v" class="btn btn-info active dark  btn-sm">%v</a>`, currentURI, info.Page, res, qs, res)
				} else {
					showmiddle += fmt.Sprintf(`<a href="%v?show=%v,%v%v" class="btn btn-default dark btn-sm">%v</a>`, currentURI, info.Page, res, qs, res)
				}
			}

			showdiv := showleft + showmiddle + showright
			output = pagediv + showdiv
		}
		//if info.TotalPages == 1 {
		//	pagediv = ""
		//}
		return template.HTML(output)
	}

	return f
}
