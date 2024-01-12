USE threattrack;
INSERT INTO ioctypes (name) VALUES ("URL"), ("DOMAIN");
INSERT INTO incidenttypes (name) VALUES ("CSIRTaaS"), ("RapidResponse");
INSERT INTO iocs (id, value, iocType) VALUES ("49c9793c-8492-468b-8ae0-64e37eb01fa0", "youtube.com", 2), ("b7a8ae7e-55ae-4983-bd36-ba26c5320487", "google.com", 2);
INSERT INTO incidents (id, name, severity, type) VALUES ("b595bee6-d8c7-4b3b-b071-1e45c2103002", "RapidResponse Case 1", "Low", 1), ("b595bee6-d8c6-4b3b-b072-1e45c2103002", "RapidResponse Case 2", "Critical", 1);