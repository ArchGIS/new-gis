MATCH (n:Monument) SET n:Site;
MATCH (n:MonumentType) SET n:SiteType;
MATCH (n:SpatialReference) SET n:Spatial;
MATCH (n:SpatialReferenceType) SET n:SpatialType;

MATCH (n:ResearchType) SET n.ru_name = n.name;
MATCH (n:Epoch) SET n.ru_name = n.name;
MATCH (n:MonumentType) SET n.ru_name = n.name;
MATCH (n:SpatialReferenceType) SET n.ru_name = n.name;
MATCH (n:Culture) SET n.ru_name = n.name;

MATCH (n:Knowledge) SET n.name = n.monument_name;

MATCH (e:Excavation)<--(s:Site)<--(k:Knowledge)<--(r:Research)-->(e) CREATE (e)<-[:EXCAVATED]-(k);