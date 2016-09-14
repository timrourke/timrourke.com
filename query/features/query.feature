Feature: build simple query
	In order to make a simple query
	As the developer of timrourke.com
	I need to be able to build SQL strings

	Scenario: Build a simple query
		When I create a new Query
		And I select "t.*" from "dinosaurs t"
		And I compile the Query
		Then the SQL should match "SELECT t.* FROM dinosaurs t  WHERE 1  ORDER BY id ASC LIMIT :limit,:offset" 
