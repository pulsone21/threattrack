package templates

import (
	"github.com/pulsone21/threattrack/frontend/templates/incidentviews"
	"github.com/pulsone21/threattrack/lib/entities"

	"github.com/a-h/templ"
)

func IncPlaningPage(incident entities.Incident, backlog []entities.Task, doing []entities.Task, done []entities.Task) templ.Component {
	return incidentviews.IncidentPlaningPage(incident, backlog, doing, done)
}

func IncSummaryPage(incident entities.Incident, wl []entities.Worklog) templ.Component {
	return incidentviews.IncidentSummaryPage(incident, wl)
}

func IncTablePage(incs []entities.Incident) templ.Component {
	return incidentviews.IncTable(incs)
}

func IncWorklogPage(incident entities.Incident, wl []entities.Worklog) templ.Component {
	return incidentviews.IncWorklogPage(incident, wl)
}

func IncIndexPage(incident entities.Incident, wl []entities.Worklog) templ.Component {
	return incidentviews.IncidentIndex(incident, wl)
}
