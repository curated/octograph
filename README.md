# octograph

Golang worker responsible for fetching issues from the GitHub GraphQL endpoint and indexing them into Elasticsearch.

Issues are filtered by a minimum amount of user reactions and indexed with the minimum footprint to allow for scaling with open source plans by generous hosting providers such as Bonsai and Heroku.
