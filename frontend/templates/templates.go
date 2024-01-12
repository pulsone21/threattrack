package templates

import (
	"github.com/pulsone21/threattrack/lib/entities"

	"github.com/a-h/templ"
)

func IncPlaningPage(incident entities.Incident, backlog []entities.Task, doing []entities.Task, done []entities.Task) templ.Component {
	return IncidentPlaningPage(incident, backlog, doing, done)
}

func IncSummaryPage(incident entities.Incident, wl []entities.Worklog) templ.Component {
	return IncidentSummaryPage(incident, wl)
}
