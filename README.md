# HollywoodStarsCRUD
### Before using you need to install:
1. docker-compose
2. make
3. migrate
4. golang
5. swagger
### Installing:
```
git clone https://github.com/AngelicaNice/HollywoodStarsCRUD.git
```
```
make build && make run

or

docker-compose up --build app
```

### If this is the first launch then:
```
make migrate
```

### This Rest API contains the following methods:
[post]   /auth/sign-up   - to create new user.<br />
[post]   /auth/sign-in   - user authentication.<br />
[get]    /actors         - get all actors.<br />
[post]   /actors         - create new actor.<br />
[get]    /actors/id/{id} - get actor by id.<br />
[put]    /actors/id/{id} - update actor by id.<br />
[delete] /actors/id/{id} - delete actor by id.<br />

#### Or after launching the application visit the page localhost:8080/swagger/index.html where all available methods are described.
