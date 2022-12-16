# How To Run

# Documentation
[Database Schema](https://dbdiagram.io/d/63934a1abae3ed7c4545dab5)

[Postman Documentation](https://documenter.getpostman.com/view/12132212/2s8YevnUpD)

[Swagger/API Spesification](https://app.swaggerhub.com/apis/DARMAWANRIZKY43/POKEDEX/1.0.0)


# Architecture Diagram
## Local Development
![architecture diagram local development](/assets/use-deall-architecture-diagram-local-development.png)

# Tech Stack
- **Golang**
- **PostgreSQL**
- **AWS RDS**
- **AWS S3**
- **Docker**

# Todo
- [ ] Documentation
    - [x] Swagger
    - [ ] Postman
    - [x] Database Schema
- [ ] Architecture diagram flow CRUD and Login
- [ ] CRUD
    - [x] Get all of list monsters
        - [x] Search by name
        - [x] Filter by type
        - [x] Filter by type catched or uncatched (options)
        - [x] Sort by name, id
        - [x] Order by ascending or descending
    - [ ] Get profile detail monster by id
    - [x] Add (admin only)
        - [x] Upload image to aws s3
    - [ ] Update (admin only)
    - [ ] Update as mark a moster as captured (user only)
    - [ ] Delete (admin only)
    - [x] Get all of list categories (admin only)
    - [x] Get all of list types (admin only)
    - [x] Login
        - [x] Generate Token JWT
- [x] Containerization
- [ ] Github Workflows
    - [x] Test
    - [ ] Create image docker and push into docker hub