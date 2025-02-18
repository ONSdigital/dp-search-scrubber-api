# dp-search-scrubber-api

## Description

This API allows users to identify Output Areas (OA) and Industry Classification (SIC) associated with a given location. OAs are small geographical areas in the UK used for statistical purposes, while SIC codes are a system of numerical codes used to identify and categorize industries.

The API takes a single, multiple or partial OA/SIC codes as input and returns a list of associated OAs and SIC information. Additionally, users can retrieve detailed information about the areas associated with each OA code.

### Available scripts

- `make help` - Displays a help menu with available `make` scripts
- `make all` - Runs audit test and build commands
- `make audit` - Audits and finds vulnerable dependencies
- `make build` - Builds ./Dockerfile image name: nlp_hub
- `make build-bin` - Build bin file in folder build
- `make clean` - Removes /bin folder
- `make convey` - Runs only convey tests
- `make debug` - Runs application locally with debug mode on
- `make fmt` - Formats the code using go fmt and go vet
- `make lint` - Automated checking of your swagger spec and source code for programmatic and stylistic errors
- `make run` - Runs container name: hub from image name: nlp_hub
- `make run-locally` - Runs the app locally
- `make test` - Runs all tests with -cover -race flags
- `make test-component` - Test components
- `make update` - Go gets all of the dependencies and downloads them

NB. For running the swagger linter locally in your CLI (using 'make lint') you will need to have first installed redocly-cli, as follows:

```shell
npm install @redocly/cli -g
```

### Configuration

| Environment variable         | Default                                       | Description
| ---------------------------- | ---------                                     | -----------
| AREA_DATA_FILE               | `data/2011 OAC Clusters and Names csv v2.csv` | The data files with the areas
| BIND_ADDR                    | :28700                                        | The host and port to bind to
| GRACEFUL_SHUTDOWN_TIMEOUT    | 5s                                            | The graceful shutdown timeout in seconds (`time.Duration` format)
| HEALTHCHECK_INTERVAL         | 30s                                           | Time between self-healthchecks (`time.Duration` format)
| HEALTHCHECK_CRITICAL_TIMEOUT | 90s                                           | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format)
| INDUSTRY_DATA_FILE           | `data/SIC07_CH_condensed_list_en.csv`         |The data files with the industries

## Quick setup

### Running the API in Docker

```shell
make run
```

### Running the API Locally

```shell
make debug
```

## Usage

Running the project either locally or in docker will expose port 28700.

```shell
curl 'http://localhost:28700/health' 
```
This will return results of the form:

```json
{
    "status": "OK",
    "version": {
        "build_time": "2020-09-26T14:30:18+03:00",
        "git_commit": "6584b786caac36b6214ffe04bf62f058d4021538",
        "language": "go",
        "language_version": "go1.19.5",
        "version": "v0.1.0"
    },
    "uptime": 7771,
    "start_time": "2023-03-09T07:46:43.587143363Z",
    "checks": []
}
```

```shell
curl 'http://localhost:28700/scrubber?q=dentists%20in%20london'
```
This will return results of the form:

```json
{
    "time": "4µs",
    "query": "dentists",
    "results": {}
}
```

If you search for an area output code like: E00000014 and an industry code like: 01140
```shell
curl 'http://localhost:28700/scrubber?q=dentists%20in%20E00000014%2001140'
```
This will return results of the form:

```json
{
    "time": "55µs",
    "query": "dentists in E00000014 01140",
    "results": {
        "areas": [
            {
                "name": "City of London",
                "region": "London",
                "region_code": "E12000007",
                "codes": {
                    "E00000014": "E00000014"
                }
            }
        ],
        "industries": [
            {
                "code": "01140",
                "name": "Growing of sugar cane"
            }
        ]
    }
}
```

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2023, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.

