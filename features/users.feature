Feature: users endpoint
	In order to use the API
	As a user of timrourke.com
	I need to be able to access the users data

	Scenario: should get empty users
		When I send "GET" request to "/api/users"
		Then the response code should be 200
		And the response should match json:
			"""
			{
				"data": [],
				"meta": {
					"version": "0"
				}
			}
			"""
