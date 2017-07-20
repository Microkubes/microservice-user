# User microservice
User microservice

## Prerequisite
Create a project directory. Set GOPATH enviroment variable to that project. Add $GOPATH/bin to the $PATH
```
export GOPATH=/path/to/project
export PATH=$GOPATH/bin:$PATH
```
Install goa and goagen:
```
cd $GOPATH
go get -u github.com/goadesign/goa/...
```

## Compile and run the service:
Clone the repo:
```
cd $GOPATH/src
git clone https://github.com/JormungandrK/user-microservice.git
``` 
Then compile and run:
```
cd user-microservice
go build -o user
./user
```

## Change the desing
If you change the desing then you should regenerate the files. Run:
```
cd $GOPATH/src/user-microservice
go generate
```
Also, recompile the service and start it again:
```
go build -o user
./user
```

## Other changes, not related to the design
For all other changes that are not related to the design just recompile the service and start it again:
```
cd $GOPATH/src/user-microservice
go build -o user
./user
```
