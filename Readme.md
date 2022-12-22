# How To Run
## Run development mode
Run development mode is use this code for run app in local machine.
- Use makefile
```go
make up_build
```

- If don't use makefile
```go
docker-compose up -d
```

## Run production mode
Running production mode uses an image from docker hub which is created and pushed by workflows github ci and also used database with AWS RDS.
- Use makefile
```go
make up_prod
```

- If don't use makefile
```go
docker-compose -f docker-compose.prod.yml up -d 
```

## Run test
**Note: To run test please run [Run development mode](##run-development-mode) first, for running database into container.**
- Use makefile
```go
make test 
```

- If don't use makefile
```go
go test -v ./...
```

# Documentation
[Database Schema](https://dbdiagram.io/d/63934a1abae3ed7c4545dab5)

[Postman Documentation](https://documenter.getpostman.com/view/12132212/2s8Z6scGJ9)

[Swagger/API Spesification](https://app.swaggerhub.com/apis/DARMAWANRIZKY43/POKEDEX/1.0.0#/Monsters/get_api_v1_monsters)

# Tech Stack
- **Golang**
- **PostgreSQL**
- **AWS RDS**
- **AWS S3**
- **Docker**

# Todo
- [x] Documentation
    - [x] Swagger
    - [x] Postman
    - [x] Database Schema
- [x] CRUD
    - [x] Get all of list monsters
        - [x] Search by name
        - [x] Filter by type
        - [x] Filter by type catched or uncatched (options)
        - [x] Sort by name, id
        - [x] Order by ascending or descending
    - [x] Get profile detail monster by id
    - [x] Add (admin only)
        - [x] Upload image to aws s3
    - [x] Update (admin only)
    - [x] Update as mark a moster as captured (user only)
    - [x] Delete (admin only)
    - [x] Get all of list categories (admin only)
    - [x] Get all of list types (admin only)
    - [x] Login
        - [x] Generate Token JWT
- [x] Containerization
- [x] Github Workflows
    - [x] Test
    - [x] Create image docker and push into docker hub
        - [x] Database use aws RDS