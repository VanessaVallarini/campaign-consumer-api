# Campaign Consumer Api
Service worker responsible for:
- consumption of campaign events and
- campaign management

## Impacts of service interruption
The information displayed on the Partner Portal may be out of date. As well as the campaigns displayed on the APP.

## Doc
[Documentation Campaign Consumer Api](https://github.com/VanessaVallarini/campaign-consumer-api)

## Technologies and Dependencies
* Golang 1.22

## Running local
### Requirements:
- docker

### Secrets Env:
Make a copy of the secret-env.yml.template file by running the command below and replace the values ​​with your credentials.
The environments must be the same names as the k8s files (values-production and values-sandbox).
```shell
cp ./local/secret-env.yml.template ./local/secret-local.yml
```

### Override Config:
Run the command below according to the environment you want (local, sandbox or production):
| Command                | Environment|
|------------------------|------------|
| make config-local      | local      |
| make config-sandbox    | sandbox    |
| make config-production | production |

### Starting docker-compose (using kafka local):
```shell
make compose-infra-up
```

### Launch:
```shell
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Package",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "cwd": "${workspaceFolder}",
            "program": "${cwd}/cmd/campaign-consumer-api/main.go"
        }
    ]
}
```

### Starting App:
Now just run the `make air` command (with air, in addition to building the application, it also allows live reloading. If it has not been started, just run `make air-init`) to run the application.

### Stop docker-compose:
```shell
make compose-infra-down
```

## Architecture
![Architecture Diagram](docs/diagrams/src/architecture/campaign_consumer_api.png)