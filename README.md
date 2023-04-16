# ReadyWorker
A GetNinjas-like website made to conclude my technical high school :)
 
# Running:

With docker-compose:
```sh
$ make compose
```

With the go compiler, ensure you have mysql/mariadb running:
```sh
$ echo "127.0.0.1 mariadb.local mariadb" | sudo tee >> /etc/hosts 
$ MARIADB_PASSWORD=<your-password> JWT_SECRETKEY=<your-jwt-secret> go run .
```

for the front-end check [here](https://github.com/trqt/readyworker)

## WARNING!
DO NOT USE THIS IS PRODUCTION, IT HAS SEVERE SECURITY ISSUES! YOU HAVE BEEN WARNED!

