# Currency

Currency provides API services for currency. 

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development.

### Prerequisites

* [go 1.9.3](https://golang.org/dl/)
* [dep](https://github.com/golang/dep)
* [docker](https://docs.docker.com/engine/installation/)
* [docker-compose](https://docs.docker.com/compose/install/)

### Installing

#### 1. Setting GOPATH

The GOPATH environment variable specifies the location of your [workspace](https://golang.org/doc/code.html#Workspaces). 
If no GOPATH is set, it is assumed to be $HOME/go on Unix systems and %USERPROFILE%\go on Windows. 
If you want to use a custom location as your workspace, you can set the GOPATH environment variable.

Please refer [here](https://github.com/golang/go/wiki/SettingGOPATH) for setting up GOPATH 

#### 2. Installing dependency management tool [dep](https://github.com/golang/dep)

```
go get -u github.com/golang/dep/cmd/dep
```

#### 3. cd to golang workspace/src

```
cd $GOPATH/src
```

#### 4. clone repo

```
git clone https://github.com/TinkLabs/Currency.git
```

#### 5. cd to project root dir

```
cd Currency
```

#### 6. Run application

```
docker-compose up
```

#### 7. Open another terminal and ping server

```
curl localhost:8080
```

## Running the tests

Please import postman collection and play with the APIs

## Deployment

Only support docker-compose for localhost development

## Contribution guidelines

[Coding Style Guideline](https://github.com/golang/go/wiki/CodeReviewComments)

[RESTful API Guideline](https://github.com/Microsoft/api-guidelines/blob/vNext/Guidelines.md)