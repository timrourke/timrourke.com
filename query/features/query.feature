Feature: build simple query
	In order to make a simple query
	As the developer of timrourke.com
	I need to be able to build SQL strings

	Scenario: Build a simple query
		When I create a new Query
		And I select "t.*" from "dinosaurs t"
		And I compile the Query
		Then the SQL should match "SELECT t.* FROM dinosaurs t  WHERE 1  ORDER BY id ASC LIMIT :limit,:offset" 

	Scenario: Build a query with multiple FROM clauses
		When I create a new Query
		And I select "a.*"
		And I select "a.leaves"
		And I select "b.bunches"
		And I add the FROM clause "apples a"
		And I add the FROM clause "bananas b"
		And I compile the Query
		Then the SQL should match "SELECT a.*, a.leaves, b.bunches FROM apples a, bananas b  WHERE 1  ORDER BY id ASC LIMIT :limit,:offset"

	Scenario: Build a query with WHERE clauses
		When I create a new Query
		And I select "d.cookies" from "d.desserts"
		And I add the WHERE clause "d.raisins = 0"
		And I add the WHERE clause "d.nuts = :nuts"
		And I compile the Query
		Then the SQL should match "SELECT d.cookies FROM d.desserts  WHERE 1 AND d.raisins = 0 AND d.nuts = :nuts ORDER BY id ASC LIMIT :limit,:offset"

	Scenario: Build a query with JOIN clauses
		When I create a new Query
		And I select "s.rivers" from "state s"
		And I select "c.county_name"
		And I add the join "LEFT JOIN county c" on "c.state_name = s.name"
		And I add the join "INNER JOIN state s2" on "s2.capitol = s2.largest_city"
		And I compile the Query
		Then the SQL should match "SELECT s.rivers, c.county_name FROM state s LEFT JOIN county c ON (c.state_name = s.name), INNER JOIN state s2 ON (s2.capitol = s2.largest_city) WHERE 1  ORDER BY id ASC LIMIT :limit,:offset"
