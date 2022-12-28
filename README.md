# Xeneta Task
This branch contains the solution for Xeneta's HTTP based Rate API.

# Project structure
This projcet includes docker configs, code and corresponding unit test.
```
.
├── application                     # Business logic
│   ├── application.go              # the gin app
│   └── db.go                       # operations on db
│   └── param.go                    # param extract and check
│   └── methods.go                  # http server methods, /rates, /
│   └── toPort.go                   # util funcs, distinguish ports from stubs
│   └── rspFormat.go                # response format to json
│   └── toPort_test.go              # unit test for toPort.go
│   └── param_test.go               # unit test for param.go  
|
├── dbDocker                        # postgreSQL docker conf
│   ├── Dockerfile                  
│   └── rates.sql                   
|
├── conf                             # server conf      
│   └── config.toml                  # sql conf
├── docker-compose.yml               # docker-compose.yml
├── Dockerfile                       # server docker file
├── go.mod                           # go.mod
├── Makefile                         # Makefile  
└── README.md                        # README.md 
```

# Quick Deployment
```
git clone git@github.com:zzuRingo/ratetask.git
cd ratetask
docker-compose up --build
```

# API Usage
```
curl "http://127.0.0.1:8080/rates?date_from=2016-01-01&date_to=2016-01-04&origin=CNSGH&destination=north_europe_main"
```

### URL
```
/rates
```
### Method

```
GET
```

### Rsp
```
[
    {"price":1112,"date":"2016-01-01"},
    {"price":1112,"date":"2016-01-02"},
    {"price":null,"date":"2016-01-03"},
    {"price":null,"date":"2016-01-04"}
]
```

### Error Code
```
200 -- OK
400 -- Error Param 
500 -- Server Internal Error
```

# Some other Info
## Thoughts on the project
![image](https://user-images.githubusercontent.com/21214705/209851125-2deead92-cf4d-45ea-8b4e-a9c39149a1b2.png)

## Time
I spent 8 hours on this task.
### Time allocate
- 2 hours on learning docker & debug
- 0.5 hour on familiar with PostgreSQL 
- 4 hours on coding 
- 1 hour on writting unit tests & self test
- 0.5 hour on writting README