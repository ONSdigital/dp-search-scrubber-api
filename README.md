# dp-nlp-search-scrubber
## Description

This API allows users to identify Output Areas (OA) and Industry Classification (SIC) associated with a given location. OAs are small geographical areas in the UK used for statistical purposes, while SIC codes are a system of numerical codes used to identify and categorize industries.

The API takes a single, multiple or partial OA/SIC codes as input and returns a list of associated OAs and SIC information. Additionally, users can retrieve detailed information about the areas associated with each OA code.

### Available scripts

- `make help` - Displays a help menu with available `make` scripts
- `make update` - Go gets all of the dependencies and downloads them
- `make build` - Builds ./Dockerfile image name: test-project
- `make run` - First builds ./Dockerfile with image name: test-project and then runs a container, with name: test_api, on port 5000

### Configuration

| Environment variable         | Default   | Description
| ---------------------------- | --------- | -----------
| BIND_ADDR                    | :3002     | The host and port to bind to
| GRACEFUL_SHUTDOWN_TIMEOUT    | 5s        | The graceful shutdown timeout in seconds (`time.Duration` format)
| HEALTHCHECK_INTERVAL         | 30s       | Time between self-healthchecks (`time.Duration` format)
| HEALTHCHECK_CRITICAL_TIMEOUT | 90s       | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format)
|	AREA_DATA_FILE               | `envconfig:"AREA_DATA_FILE"` | The data files with the areas
|	INDUSTRY_DATA_FILE           | `envconfig:"INDUSTRY_DATA_FILE"` |The data files with the industries

## Quick setup

### Docker

```shell
make run
```

### Locally

```shell
make update
go run .
```

## Dependencies

- `github.com/ONSdigital/log.go/v2 v2.3.0`
- `github.com/alediaferia/prefixmap v1.0.1`
- `github.com/gocarina/gocsv v0.0.0-20230123225133-763e25b40669`
- `github.com/gorilla/mux v1.8.0`
- `github.com/invopop/jsonschema v0.7.0`
- `github.com/joho/godotenv v1.5.1`
- `github.com/kelseyhightower/envconfig v1.4.0`
- `go version go1.19.5 linux/amd64 `

## Usage

Running the project either locally or in docker will expose port 3002.

```shell
curl 'http://localhost:3002/health' 
```
This will return results of the form:

```shell
OK
```

```shell
curl 'http://localhost:3002/scrubber/search?q=dentists%20in%20london'
```
This will return results of the form:

```json
{
    "time": "4µs",
    "query": "dentists",
    "results": {
        "areas": null,
        "industries": null
    }
}
```

If you search for an area output code like: E00000014 and an industry code like: 01140
```shell
curl 'http://localhost:3002/scrubber/search?q=dentists%20in%20E00000014%2001140'
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

```shell
curl 'http://localhost:3002/json-schema'
```
This will return results of the form:

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://gitlab.com/flaxandteal/onyx/dp-nlp-search-scrubber/payloads/scrubber-resp",
  "$ref": "#/$defs/ScrubberResp",
  "$defs": {
    "AreaResp": {
      "properties": {
        "name": {
          "type": "string"
        },
        "region": {
          "type": "string"
        },
        "region_code": {
          "type": "string"
        },
        "codes": {
          "patternProperties": {
            ".*": {
              "type": "string"
            }
          },
          "type": "object"
        }
      },
      "additionalProperties": false,
      "type": "object",
      "required": [
        "name",
        "region",
        "region_code",
        "codes"
      ]
    },
    "IndustryResp": {
      "properties": {
        "code": {
          "type": "string"
        },
        "name": {
          "type": "string"
        }
      },
      "additionalProperties": false,
      "type": "object",
      "required": [
        "code",
        "name"
      ]
    },
    "Results": {
      "properties": {
        "areas": {
          "items": {
            "$ref": "#/$defs/AreaResp"
          },
          "type": "array"
        },
        "industries": {
          "items": {
            "$ref": "#/$defs/IndustryResp"
          },
          "type": "array"
        }
      },
      "additionalProperties": false,
      "type": "object",
      "required": [
        "areas",
        "industries"
      ]
    },
    "ScrubberResp": {
      "properties": {
        "time": {
          "type": "string"
        },
        "query": {
          "type": "string"
        },
        "results": {
          "$ref": "#/$defs/Results"
        }
      },
      "additionalProperties": false,
      "type": "object",
      "required": [
        "time",
        "query",
        "results"
      ]
    }
  }
}
```

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2023, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.

