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

## Set up MongoDB
Create users database with default username and password.
See: [Set up MongoDB](https://github.com/JormungandrK/jormungandr-infrastructure#mongodb--v346-)
```
export MS_DBNAME=users
./mongo/run.sh
```
Then install mgo package:
```
cd $GOPATH
go get gopkg.in/mgo.v2
```

# Docker Builds

First, create a directory for the shh keys:
```bash
mkdir keys
```

Find a key that you'll use to acceess JormungandrK organization on github. Then copy the
private key to the directory you created above. The build would use this key to
access ```JormungandrK/microservice-tools``` repository.

```bash
cp ~/.ssh/id_rsa keys/
```

**WARNING!** Make sure you don't commit or push this key to the repository!

To build the docker image of the microservice, run the following command:
```bash
docker build -t user-microservice .
```

# Running the microservice

To run the user-microservice you'll need to set up some ENV variables:

 * **SERVICE_CONFIG_FILE** - Location of the configuration JSON file (default: config.json)
 * **API_GATEWAY_URL** - Kong API url (default: http://localhost:8001)

Run the docker image:
```bash
docker run user-microservice
```

# Service configuration

The service loads the gateway configuration from a JSON file (See **SERVICE_CONFIG_FILE**
ENV variable).
Here's an example of a JSON configuration file:

```json
{
  "name": "user-microservice",
  "port": 8080,
  "virtual_host": "user.services.jormugandr.org",
  "hosts": ["localhost", "user.services.jormugandr.org"],
  "weight": 10,
  "slots": 100
}
```

Configuration properties:
 * **name** - ```"user-microservice"``` - the name of the service, do not change this.
 * **port** - ```8080``` - port on which the microservice is running
 * **virtual_host** - ```"user.services.jormugandr.org"``` domain name of the service group/cluster. Don't change if not sure.
 * **hosts** - list of valid hosts. Used for proxying and load-balancing of the incoming request. You need to have at least the **virtual_host** in the list.
 * **weight** - instance weight - user for load balancing
 * **slots** - maximal number of service instances under ```"user.services.jormugandr.org"```
