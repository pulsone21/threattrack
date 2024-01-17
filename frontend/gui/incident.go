package gui

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/pulsone21/threattrack/frontend/templates"
	"github.com/pulsone21/threattrack/lib/entities"
	"github.com/pulsone21/threattrack/lib/utils"
)

type IncidentHandler struct {
	backendBase   string
	backendAdress string
	templatePath  string
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
	s.Router.HandleFunc("/incident/{id}", utils.CreateHtmlHandleFunc(iH.serveIncidentIndex)).Methods("GET")
	s.Router.HandleFunc("/incident/{id}/summary", utils.CreateHtmlHandleFunc(iH.serveIncidentPage)).Methods("GET")
	s.Router.HandleFunc("/incident/{id}/worklog", utils.CreateHtmlHandleFunc(iH.serveIncidentWorklog)).Methods("GET")
	s.Router.HandleFunc("/incident/{id}/planing", utils.CreateHtmlHandleFunc(iH.serveIncidentPlaning)).Methods("GET")
	s.Router.HandleFunc("/incident/{id}/iocView", utils.CreateHtmlHandleFunc(iH.serveIncidentiocView)).Methods("GET")
}

func CreateIncidentHandler(ser *Server, backendBase string) *IncidentHandler {
	wd, _ := os.Getwd()
	iH := &IncidentHandler{
		backendBase:   backendBase,
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
	incData, err := requestData[entities.Incident](iH.backendAdress)
	if err != nil {
		return entities.InternalServerError(err, uri)
	}
	page := templates.IncTablePage(incData.Data)
	if err = page.Render(ctx, w); err != nil {
		return entities.InternalServerError(err, uri)
	}
	return nil
}

func (iH *IncidentHandler) serveIncidentIndex(ctx context.Context, w http.ResponseWriter, r *http.Request) *entities.ApiError {
	uri := ctx.Value("uri").(string)
	incId := mux.Vars(r)["id"]
	url := fmt.Sprintf("%s/%s", iH.backendAdress, incId)
	fmt.Printf("\nrequesting backend with %s \n", url)
	incData, err := requestData[entities.Incident](url)
	if err != nil {
		return entities.InternalServerError(err, uri)
	}
	wlData, err := requestData[entities.Worklog](fmt.Sprintf("%s/worklogs?incident=%s", iH.backendBase, incId))
	if err != nil {
		return entities.InternalServerError(err, uri)
	}

	page := templates.IncIndexPage(incData.Data[0], wlData.Data)
	if err = page.Render(ctx, w); err != nil {
		return entities.InternalServerError(err, uri)
	}
	return nil
}

func (iH *IncidentHandler) serveIncidentPage(ctx context.Context, w http.ResponseWriter, r *http.Request) *entities.ApiError {
	uri := ctx.Value("uri").(string)
	incId := mux.Vars(r)["id"]
	url := fmt.Sprintf("%s/%s", iH.backendAdress, incId)
	fmt.Printf("\nrequesting backend with %s \n", url)
	incData, err := requestData[entities.Incident](url)
	if err != nil {
		return entities.InternalServerError(err, uri)
	}

	wlData, err := requestData[entities.Worklog](fmt.Sprintf("%s/worklogs?incident=%s", iH.backendBase, incId))
	if err != nil {
		return entities.InternalServerError(err, uri)
	}
	page := templates.IncSummaryPage(incData.Data[0], wlData.Data)

	if err = page.Render(ctx, w); err != nil {
		return entities.InternalServerError(err, uri)
	}
	return nil
}

func (iH *IncidentHandler) serveIncidentWorklog(ctx context.Context, w http.ResponseWriter, r *http.Request) *entities.ApiError {
	uri := ctx.Value("uri").(string)
	incId := mux.Vars(r)["id"]
	url := fmt.Sprintf("%s/%s", iH.backendAdress, incId)
	fmt.Printf("\nrequesting backend with %s \n", url)
	incData, err := requestData[entities.Incident](url)
	if err != nil {
		return entities.InternalServerError(err, uri)
	}
	wlData, err := requestData[entities.Worklog](fmt.Sprintf("%s/worklogs?incident=%s", iH.backendBase, incId))
	if err != nil {
		return entities.InternalServerError(err, uri)
	}
	page := templates.IncWorklogPage(incData.Data[0], wlData.Data)
	if err = page.Render(ctx, w); err != nil {
		return entities.InternalServerError(err, uri)
	}
	return nil
}

func (iH *IncidentHandler) serveIncidentPlaning(ctx context.Context, w http.ResponseWriter, r *http.Request) *entities.ApiError {
	uri := ctx.Value("uri").(string)
	incId := mux.Vars(r)["id"]
	url := fmt.Sprintf("%s/%s", iH.backendAdress, incId)
	fmt.Printf("\nrequesting backend with %s \n", url)

	incData, err := requestData[entities.Incident](url)
	if err != nil {
		return entities.InternalServerError(err, uri)
	}

	vD := iH.createPlaningViewData(incData.Data[0])
	fmt.Println(vD)
	page := templates.IncPlaningPage(vD.Incident, vD.Backlog, vD.Doing, vD.Done)

	if err = page.Render(ctx, w); err != nil {
		return entities.InternalServerError(err, uri)
	}
	return nil
}

func (iH *IncidentHandler) serveIncidentiocView(ctx context.Context, w http.ResponseWriter, r *http.Request) *entities.ApiError {
	uri := ctx.Value("uri").(string)
	incId := mux.Vars(r)["id"]
	url := fmt.Sprintf("%s/%s", iH.backendAdress, incId)
	fmt.Printf("\nrequesting backend with %s \n", url)

	_, err := requestData[entities.Incident](url)
	if err != nil {
		return entities.InternalServerError(err, uri)
	}

	// tmpl, err := iH.loadTemplate("incidentIOCView.html")
	// if err != nil {
	// 	return entities.InternalServerError(err, uri)
	// }
	// if err = tmpl.Execute(w, incData.Data[0]); err != nil {
	// 	return entities.InternalServerError(err, uri)
	// }
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
