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

# Todo
- [ ] Documentation
    - [x] Swagger
    - [ ] Postman
    - [x] Database Schema
- [ ] Architecture diagram flow CRUD and Login
- [ ] CRUD
    - [ ] Get all of list monsters
        - [ ] Search by name
        - [ ] Filter by type
        - [ ] Filter by type catched or uncatched (options)
        - [ ] Sort by name, id, ascending or descending
    - [ ] Get profile detail monster
    - [ ] Add (admin only)
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