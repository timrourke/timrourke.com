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

	Scenario: should get single user
		Given there are users:
			| id | username  | email                 | created_at            | updated_at           | password_hash |
			| 1  | testuser1 | testuser1@example.com | 2016-02-07T03:27:16Z  | 2016-03-17T12:27:49Z | fakehash      |
		When I send "GET" request to "/api/users/1"
		Then the response code should be 200
		And the response should match json:
			"""
			{
				"data": {
					"type": "users",
					"id": "1",
					"attributes": {
						"created-at": "2016-02-07T03:27:16Z",
						"updated-at": "2016-03-17T12:27:49Z",
						"username": "testuser1",
						"email": "testuser1@example.com"
					}
				},
				"meta": {
					"version": "0"
				}
			}
			"""


	Scenario: should get list of users
		Given there are users:
			| id | username  | email                 | created_at            | updated_at           | password_hash |
			| 1  | testuser1 | testuser1@example.com | 2016-02-07T03:27:16Z  | 2016-03-17T12:27:49Z | fakehash      |
			| 2  | testuser2 | testuser2@example.com | 2016-02-07T04:27:16Z  | 2016-04-17T12:27:49Z | fakehash2     |
			| 3  | testuser3 | testuser3@example.com | 2016-02-07T05:27:16Z  | 2016-05-17T12:27:49Z | fakehash3     |
			| 4  | testuser4 | testuser4@example.com | 2016-02-07T06:27:16Z  | 2016-06-17T12:27:49Z | fakehash4     |
			| 5  | testuser5 | testuser5@example.com | 2016-02-07T07:27:16Z  | 2016-07-17T12:27:49Z | fakehash5     |
		When I send "GET" request to "/api/users"
		Then the response code should be 200
		And the response should match json:
			"""
			{
				"data": [
					{
						"type": "users",
						"id": "1",
						"attributes": {
							"created-at": "2016-02-07T03:27:16Z",
							"updated-at": "2016-03-17T12:27:49Z",
							"username": "testuser1",
							"email": "testuser1@example.com"
						}
					},
					{
						"type": "users",
						"id": "2",
						"attributes": {
							"created-at": "2016-02-07T04:27:16Z",
							"updated-at": "2016-04-17T12:27:49Z",
							"username": "testuser2",
							"email": "testuser2@example.com"
						}
					},
					{
						"type": "users",
						"id": "3",
						"attributes": {
							"created-at": "2016-02-07T05:27:16Z",
							"updated-at": "2016-05-17T12:27:49Z",
							"username": "testuser3",
							"email": "testuser3@example.com"
						}
					},
					{
						"type": "users",
						"id": "4",
						"attributes": {
							"created-at": "2016-02-07T06:27:16Z",
							"updated-at": "2016-06-17T12:27:49Z",
							"username": "testuser4",
							"email": "testuser4@example.com"
						}
					},
					{
						"type": "users",
						"id": "5",
						"attributes": {
							"created-at": "2016-02-07T07:27:16Z",
							"updated-at": "2016-07-17T12:27:49Z",
							"username": "testuser5",
							"email": "testuser5@example.com"
						}
					}
				],
				"meta": {
					"version": "0"
				}
			}
			"""
