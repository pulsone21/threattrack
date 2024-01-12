SELECT
	incidents.id AS incident_id,
	incidents.name as name,
	incidents.severity as severity,
	incidents.status as status,
	UNIX_TIMESTAMP(incidents.creationdate) as createdAt,
	incidenttypes.id AS type_id,
	incidenttypes.name AS type_name,
	users.id as usr_id,
	users.firstName as usr_firstname,
	users.lastName as usr_lastname,
	users.email as usr_email,
	users.fullname as usr_fullname,
	users.created_at as usr_createdAt
FROM
	incidents
LEFT JOIN
	incidenttypes ON incidents.type = incidenttypes.id
LEFT JOIN
	users ON incidents.owner_id = users.id
ORDER BY incidents.creationdate DESC
LIMIT ? OFFSET ?;
	