# Loopline challenge

What I must point out here is that I took a lot of shortcuts (while focusing maybe too much on testing the go).
The app is definitely not ready for production, but it's a good start nonetheless, and I hope the main requirements were met.

### backend
This is the more "solid" part of the app. 
Though I skipped testing negative paths (like error handling), this could be done and test coverage could rise to even 100%.

Storage is in memory.

For production, this 
* should only expose secure comms (https, certs), only in the private container network.
* should have tighter CORS rules
* should use Swagger (OpenAPI)
* should store data in a db or redis

### frontend
Learning to use react with typescript was quite a pleasant experience, and I realised I
missed a lot of interesting web stuff in the last years while being focused on devops and lower level.

The solution doesn't have tests and runs in dev mode. 
For prod, it would need a different Dockerfile to `npm build` static files and maybe serve them with nginx).
From a product point of view, this barely works. It was a nice journey to get here nonetheless
and it's a good starting point for further development 
(e.g validating things, a modal, displaying (ajax) errors, only generating uuids in the backend, etc).

## Usage

First, make sure your host doesn't use ports 3000 and 8080 (or change them yourself), then
```bash
docker-compose up -d --build
```
Open http://localhost:3000 in the browser.

You can also use curl to add, list and delete notes: GET, POST and DELETE are exposed on localhost:8080 (route /)

Then, when you're done with it,
```bash
docker-compose stop
```

## Tests

```bash
cd back
go test -v -race ./...
```

`npm test` for front makes no sense at this point, since I didn't write any tests for the frontend app.
