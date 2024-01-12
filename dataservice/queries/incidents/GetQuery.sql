SELECT
	incidents.id AS incident_id,
	incidents.name as name,
	incidents.severity as severity,
	incidents.status as status,
	UNIX_TIMESTAMP(incidents.creationdate) as createdAt,
	incidenttypes.id AS type_id,
	incidenttypes.name AS type_name
FROM
	incidents
LEFT JOIN
	incidenttypes ON incidents.type = incidenttypes.id
%s
ORDER BY incidents.creationdate DESC
LIMIT ? OFFSET ?;
	