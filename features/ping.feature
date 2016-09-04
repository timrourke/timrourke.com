Feature: ping server
	In order to make sure the site is running
	As a user of timrourke.com
	I need to be able to ping the site

	Scenario: Ping the site
		When I send "GET" request to "/ping"
		Then the response code should be 200
		And the response should match text "pong"
