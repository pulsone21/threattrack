package gui

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pulsone21/threattrack/lib/entities"
	"github.com/pulsone21/threattrack/lib/utils"
)

type IndicatorHandler struct {
	backendAdress string
}

type indTableViewData struct {
	Indicators []Indicator
}

type Indicator struct {
	Id      string        `json:"id"`
	Value   string        `json:"value"`
	Type    IndicatorType `json:"iocType"`
	Verdict Verdict       `json:"verdict"`
}

type IndicatorType struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type Verdict string

const (
	Benigne   Verdict = "Benigne"
	Malicious Verdict = "Malicious"
	Neutral   Verdict = "Neutral"
)

func (iH *IndicatorHandler) createHandles(s *Server) {
	s.Router.HandleFunc("/indicatorTable/", utils.CreateApiHandleFunc(iH.serveIncidentTable)).Methods("GET")
	s.Router.HandleFunc("/indicator/", utils.CreateApiHandleFunc(iH.serveIndicatorPage)).Methods("GET")
}

func CreateIndicatorHandler(ser *Server, backendBase string) *IndicatorHandler {
	iH := &IndicatorHandler{
		backendAdress: fmt.Sprintf("%s/ioc", backendBase),
	}
	iH.createHandles(ser)
	return iH
}

func (iH *IndicatorHandler) serveIncidentTable(ctx context.Context, w http.ResponseWriter, r *http.Request) (*entities.ApiResponse, *entities.ApiError) {
	uri := ctx.Value("uri").(string)
	fmt.Printf("\nrequesting backend with %s \n", iH.backendAdress)
	res, err := http.Get(iH.backendAdress)
	if err != nil {
		log.Fatalln(err.Error())
		return nil, entities.InternalServerError(err, uri)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		// Handle the Error
		return nil, entities.InternalServerError(http.ErrAbortHandler, uri)
	}

	resbody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err.Error())
		return nil, entities.InternalServerError(err, uri)
	}
	tmpl, err := template.ParseFiles("./templates/indicatorTable.html")
	if err != nil {
		log.Fatalln(err.Error())
		return nil, entities.InternalServerError(err, uri)
	}
	var inds []Indicator

	if err = json.Unmarshal(resbody, &inds); err != nil {
		return nil, entities.InternalServerError(err, uri)
	}
	fmt.Println(inds)

	if err = tmpl.Execute(w, indTableViewData{
		Indicators: inds,
	}); err != nil {
		return nil, entities.InternalServerError(err, uri)
	}
	return entities.NewApiResponse(200, uri, ""), nil
}

func (iH *IndicatorHandler) serveIndicatorPage(ctx context.Context, w http.ResponseWriter, r *http.Request) (*entities.ApiResponse, *entities.ApiError) {
	uri := ctx.Value("uri").(string)
	incId := mux.Vars(r)["id"]
	url := fmt.Sprintf("%s/%s", iH.backendAdress, incId)
	fmt.Printf("\nrequesting backend with %s \n", url)

	res, err := http.Get(url)
	if err != nil {

		log.Fatalln(err.Error())
		return nil, entities.InternalServerError(err, uri)
	}

	var ind Indicator
	if err = json.NewDecoder(res.Body).Decode(&ind); err != nil {
		log.Fatalln(err.Error())
		return nil, entities.InternalServerError(err, uri)
	}
	fmt.Println(ind)
	tmpl, err := template.ParseFiles("./templates/indicator.html")
	if err != nil {
		log.Fatalln(err.Error())
		return nil, entities.InternalServerError(err, uri)
	}
	if err = tmpl.Execute(w, ind); err != nil {
		return nil, entities.InternalServerError(err, uri)
	}
	return entities.NewApiResponse(200, uri, ""), nil
}
