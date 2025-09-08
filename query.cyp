//-- The structure of railway and road connections between cities is given in the Neo4J database.
//-- Cities are nodes, connections are edges.
//-- The type of connection is distinguished by label.
//-- Each city has a unique name represented as an attribute.
//--
//--
//-- 1. Write a query that creates the following railway connections: Wrocław - Kraków - Warszawa - Gdańsk.
//-- 2. Write a query that adds road connections: Kraków - Gdańsk and Warszawa - Gdańsk to existing cities (Wrocław, Kraków, Warsaw, Gdańsk).
//-- 3. Write a query that returns the names of cities that have at least one railway and one road connection with another city.
//-- 4. Write a query that finds all road paths between 2 cities

UNWIND [
  {name: 'Karkow'},
  {name: 'Gdnask'},
  {name: 'Warszawa'},
  {name: 'Wroclaw'}
] AS cities
CREATE (c:City {name: cities.name});

MATCH (a:City {name: 'Wroclaw'}),(b:City {name: 'Karkow'}),(c:City {name: 'Warszawa'}),(d:City {name: 'Gdnask'})
CREATE (a)-[:RAILWAY]->(b)-[:RAILWAY]->(c)-[:RAILWAY]->(d);

MATCH (a:City {name: 'Wroclaw'}),(b:City {name: 'Karkow'}),(c:City {name: 'Warszawa'}),(d:City {name: 'Gdnask'})
CREATE (d)-[:ROAD]->(b), (d)-[:ROAD]->(c),(d)-[:ROAD]->(a)

MATCH (c:City) WHERE (
(c)-[:ROAD]-() AND (c)-[:RAILWAY]-()
) RETURN c.name;

MATCH p = ((start:City {name: "Krakow"})-[:ROAD*1..2]->(end:City {name: "Warsaw"}))
RETURN p;