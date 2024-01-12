package gui

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/pulsone21/threattrack/frontend/templates"
	"github.com/pulsone21/threattrack/lib/entities"
	"github.com/pulsone21/threattrack/lib/utils"
)

type IncidentHandler struct {
	backendAdress string
	templatePath  string
}
type incTableViewData struct {
	Incidents []entities.Incident
}

type backendData struct {
	StatusCode int
	RequestUrl string
	Data       []entities.Incident
}

type incPlaningViewData struct {
	Incident entities.Incident
	Backlog  []entities.Task
	Doing    []entities.Task
	Done     []entities.Task
	// TODO Playbooks
}

func (iH *IncidentHandler) createHandles(s *Server) {
	s.Router.HandleFunc("/incidentTable/", utils.CreateHtmlHandleFunc(iH.serveIncidentTable)).Methods("GET")
	s.Router.HandleFunc("/incident/summary", utils.CreateHtmlHandleFunc(iH.serveIncidentPage)).Methods("GET")
	s.Router.HandleFunc("/incident/worklog", utils.CreateHtmlHandleFunc(iH.serveIncidentWorklog)).Methods("GET")
	s.Router.HandleFunc("/incident/planing", utils.CreateHtmlHandleFunc(iH.serveIncidentPlaning)).Methods("GET")
	s.Router.HandleFunc("/incident/iocView", utils.CreateHtmlHandleFunc(iH.serveIncidentiocView)).Methods("GET")
}

func CreateIncidentHandler(ser *Server, backendBase string) *IncidentHandler {
	wd, _ := os.Getwd()
	iH := &IncidentHandler{
		backendAdress: fmt.Sprintf("%s/incidents", backendBase),
		templatePath:  "frontend/templates/incident",
	}
	fmt.Printf("%s/%s\n", wd, iH.templatePath)
	iH.createHandles(ser)
	return iH
}

func (iH *IncidentHandler) serveIncidentTable(ctx context.Context, w http.ResponseWriter, r *http.Request) *entities.ApiError {
	uri := ctx.Value("uri").(string)
	fmt.Printf("\nrequesting backend with %s \n", iH.backendAdress)
	res, err := http.Get(iH.backendAdress)
	if err != nil {
		fmt.Println(err.Error())
		return entities.InternalServerError(err, uri)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		// Handle the Error
		return entities.InternalServerError(http.ErrAbortHandler, uri)
	}
	resbody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return entities.InternalServerError(err, uri)
	}
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/incidentTable.html", iH.templatePath))
	if err != nil {
		fmt.Println(err.Error())
		return entities.InternalServerError(err, uri)
	}
	var data backendData
	fmt.Println("defining the struct")
	if err = json.Unmarshal(resbody, &data); err != nil {
		fmt.Println(err.Error())
		return entities.InternalServerError(err, uri)
	}
	fmt.Println("struct unmarshaled")
	fmt.Println(data.Data)
	if err = tmpl.Execute(w, incTableViewData{
		Incidents: data.Data,
	}); err != nil {
		return entities.InternalServerError(err, uri)
	}
	return nil
}

func (iH *IncidentHandler) serveIncidentPage(ctx context.Context, w http.ResponseWriter, r *http.Request) *entities.ApiError {
	uri := ctx.Value("uri").(string)
	incId := r.URL.Query().Get("id")
	url := fmt.Sprintf("%s/%s", iH.backendAdress, incId)
	fmt.Printf("\nrequesting backend with %s \n", url)

	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Error in backend request")
		return entities.InternalServerError(err, uri)
	}
	var data backendData
	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		fmt.Println("Error in decoding of response content")
		return entities.InternalServerError(err, uri)
	}
	//TODO Implement the worklog stuff
	page := templates.IncSummaryPage(data.Data[0], nil)

	if err = page.Render(ctx, w); err != nil {
		return entities.InternalServerError(err, uri)
	}
	return nil
}

func (iH *IncidentHandler) serveIncidentWorklog(ctx context.Context, w http.ResponseWriter, r *http.Request) *entities.ApiError {
	uri := ctx.Value("uri").(string)
	incId := r.URL.Query().Get("id")
	url := fmt.Sprintf("%s/%s", iH.backendAdress, incId)
	fmt.Printf("\nrequesting backend with %s \n", url)

	res, err := http.Get(url)
	if err != nil {
		return entities.InternalServerError(err, uri)
	}
	var inc entities.Incident
	if err = json.NewDecoder(res.Body).Decode(&inc); err != nil {
		return entities.InternalServerError(err, uri)
	}

	tmpl, err := iH.loadTemplate("incidentWorklogs.html")
	if err != nil {
		return entities.InternalServerError(err, uri)
	}
	if err = tmpl.Execute(w, inc); err != nil {
		return entities.InternalServerError(err, uri)
	}
	return nil
}

func (iH *IncidentHandler) serveIncidentPlaning(ctx context.Context, w http.ResponseWriter, r *http.Request) *entities.ApiError {
	uri := ctx.Value("uri").(string)
	incId := r.URL.Query().Get("id")
	url := fmt.Sprintf("%s/%s", iH.backendAdress, incId)
	fmt.Printf("\nrequesting backend with %s \n", url)

	res, err := http.Get(url)
	if err != nil {
		return entities.InternalServerError(err, uri)
	}
	var data backendData
	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		return entities.InternalServerError(err, uri)
	}

	vD := iH.createPlaningViewData(data.Data[0])
	fmt.Println(vD)
	page := templates.IncPlaningPage(vD.Incident, vD.Backlog, vD.Doing, vD.Done)

	if err = page.Render(ctx, w); err != nil {
		return entities.InternalServerError(err, uri)
	}
	return nil
}

func (iH *IncidentHandler) serveIncidentiocView(ctx context.Context, w http.ResponseWriter, r *http.Request) *entities.ApiError {
	uri := ctx.Value("uri").(string)
	incId := r.URL.Query().Get("id")
	url := fmt.Sprintf("%s/%s", iH.backendAdress, incId)
	fmt.Printf("\nrequesting backend with %s \n", url)

	res, err := http.Get(url)
	if err != nil {
		return entities.InternalServerError(err, uri)
	}
	var inc entities.Incident
	if err = json.NewDecoder(res.Body).Decode(&inc); err != nil {
		return entities.InternalServerError(err, uri)
	}

	tmpl, err := iH.loadTemplate("incidentIOCView.html")
	if err != nil {
		return entities.InternalServerError(err, uri)
	}
	if err = tmpl.Execute(w, inc); err != nil {
		return entities.InternalServerError(err, uri)
	}
	return nil
}

func (iH *IncidentHandler) createPlaningViewData(inc entities.Incident) *incPlaningViewData {
	var backlog, doing, done []entities.Task
	for _, v := range inc.Tasks {
		switch v.State {
		case entities.Backlog:
			backlog = append(backlog, v)
		case entities.Doing:
			doing = append(doing, v)
		case entities.Done:
			done = append(done, v)
		}
	}
	return &incPlaningViewData{
		Incident: inc,
		Backlog:  backlog,
		Doing:    doing,
		Done:     done,
	}
}

func (iH *IncidentHandler) loadTemplate(templ string) (*template.Template, error) {
	return template.ParseFiles(fmt.Sprintf("%s/%s", iH.templatePath, templ), fmt.Sprintf("%s/partials/com_incNavbar.html", iH.templatePath), fmt.Sprintf("%s/partials/com_incHeader.html", iH.templatePath))
}
