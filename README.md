# instavote
Instavote app for Linux Foundation training LFS261

As part of the Linux Foundation trainings, the LFS261 course re-uses one of Docker's official sample applications, called [example-voting-app](https://github.com/dockersamples/example-voting-app).

Instead of simply re-using the provided application, I built a similar one in Go and created the necessary Dockerfile and docker-compose.yml files.

The application provides:
- a web server that allows users to choose between 2 options and vote, and pushes the votes to a Redis queue,
- a worker that gets votes from that Redis queue and updates the MongoDB database,
- a web server that builds the results from the MongoDB database and displays them in the browser, refreshing the results automatically.
