# User microservice

[![Build](https://travis-ci.com/Microkubes/microservice-user.svg?token=UB5yzsLHNSbtjSYrGbWf&branch=master)](https://travis-ci.com/Microkubes/microservice-user)
[![Test Coverage](https://api.codeclimate.com/v1/badges/2cf4d5d4a0ade7c5c358/test_coverage)](https://codeclimate.com/repos/5971a9ed730e750274000347/test_coverage)
[![Maintainability](https://api.codeclimate.com/v1/badges/2cf4d5d4a0ade7c5c358/maintainability)](https://codeclimate.com/repos/5971a9ed730e750274000347/maintainability)

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
git clone https://github.com/Microkubes/microservice-user.git
```
Then compile and run:
```
cd microservice-user
go build -o user
./user
```

## Change the desing
If you change the desing then you should regenerate the files. Run:
```
cd $GOPATH/src/microservice-user
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
cd $GOPATH/src/microservice-user
go build -o user
./user
```

## Tests
For testing we use controller_test.go files which call the generated test helpers which package that data into HTTP requests and calls the actual controller functions. The test helpers retrieve the written responses, deserialize them, validate the generated data structures (against the validations written in the design) and make them available to the tests. Run:
```
go test -v
```

## Set up MongoDB
Create users database with default username and password.
See: [Set up MongoDB](https://github.com/Microkubes/jormungandr-infrastructure#mongodb--v346-)
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

Find a key that you'll use to acceess Microkubes organization on github. Then copy the
private key to the directory you created above. The build would use this key to
access ```Microkubes/microservice-tools``` repository.

```bash
cp ~/.ssh/id_rsa keys/
```

**WARNING!** Make sure you don't commit or push this key to the repository!

To build the docker image of the microservice, run the following command:
```bash
docker build -t microservice-user .
```

# Running the microservice

To run the microservice-user you'll need to set up some ENV variables:

 * **SERVICE_CONFIG_FILE** - Location of the configuration JSON file
 * **API_GATEWAY_URL** - Kong API url (default: http://localhost:8001)
 * **MONGO_URL** - Host IP(example: 192.168.1.10:27017)
 * **MS_USERNAME** - Mongo username (default: restapi)
 * **MS_PASSWORD** - Mongo password (default: restapi)
 * **MS_DBNAME** - Mongo database name (default: users)

Run the docker image:
```bash
docker run microservice-user
```

## Check if the service is self-registering on Kong Gateway

First make sure you have started Kong. See [Jormungandr Infrastructure](https://github.com/Microkubes/jormungandr-infrastructure)
on how to set up Kong locally.

If you have Kong admin endpoint running on http://localhost:8001 , you're good to go.
Build and run the service:
```bash
go build -o user
./user
```

To access the user service, then instead of calling the service on :8080 port,
make the call to Kong:

```bash
curl -v --header "Host: user.services.jormugandr.org" http://localhost:8000/user/1
```

You should see a log on the terminal running the service that it received and handled the request.

## Running with the docker image

Assuming that you have Kong and it is availabel od your host (ports: 8001 - admin, and 8000 - proxy) and
you have build the service docker image (microservice-user), then you need to pass
the Kong URL as an ENV variable to the docker run. This is needed because by default
the service will try http://localhost:8001 inside the container and won't be able to connect to kong.

Find your host IP using ```ifconfig``` or ```ip addr```.
Assuming your host IP is 192.168.1.10, then run:

```bash
docker run -ti -e API_GATEWAY_URL=http://192.168.1.10:8001 -e MONGO_URL=192.168.1.10:27017 microservice-user
```

If there are no errors, on a different terminal try calling Kong on port :8000

```bash
curl -v --header "Host: user.services.jormugandr.org" http://localhost:8000/user/1
```

You should see output (log) in the container running the service.



# Service configuration

The service loads the gateway configuration from a JSON file /run/secrets/microservice_user_config.json. To change the path set the
**SERVICE_CONFIG_FILE** env var.
Here's an example of a JSON configuration file:

```json
{
  "name": "microservice-user",
  "port": 8080,
  "virtual_host": "user.services.jormugandr.org",
  "hosts": ["localhost", "user.services.jormugandr.org"],
  "weight": 10,
  "slots": 100
}
```

Configuration properties:
 * **name** - ```"microservice-user"``` - the name of the service, do not change this.
 * **port** - ```8080``` - port on which the microservice is running.
 * **virtual_host** - ```"user.services.jormugandr.org"``` domain name of the service group/cluster. Don't change if not sure.
 * **hosts** - list of valid hosts. Used for proxying and load balancing of the incoming request. You need to have at least the **virtual_host** in the list.
 * **weight** - instance weight - user for load balancing.
 * **slots** - maximal number of service instances under ```"user.services.jormugandr.org"```.
