Feature: Helloworld
    Feature: Scrubber
    Scenario: When Searching for Without OAC or SIC codes I get empty resp
        When I GET "/v1/scrubber/search?q=dentists"
        And the response body is the same as the json in "./features/testdata/expecteddata/emptyResponse.json"

    Scenario: When Searching for With only OAC I get resp as in json
        When I GET "/v1/scrubber/search?q=E00000001"
        And the response body is the same as the json in "./features/testdata/expecteddata/onlyOACResponse.json"

    Scenario: When Searching for With only SIC I get resp as in json
        When I GET "/v1/scrubber/search?q=01230"
        And the response body is the same as the json in "./features/testdata/expecteddata/onlySICResponse.json"

    Scenario: When Searching for With both OAC and SIC I get resp as in json
        When I GET "/v1/scrubber/search?q=01230%20E00000001"
        And the response body is the same as the json in "./features/testdata/expecteddata/fullResponse.json"

    Scenario: When Searching for With both OAC and SIC that have special characters between then I get resp as in json
        When I GET "/v1/scrubber/search?q=01230,E00000001"
        And the response body is the same as the json in "./features/testdata/expecteddata/fullResponse.json"

    Scenario: When Searching for With both OAC and SIC that have special characters and multiple OAC then I get resp as in json
        When I GET "/v1/scrubber/search?q=01230,E00000001,E00000014,E00000017,E00000016"
        And the response body is the same as the json in "./features/testdata/expecteddata/fullResponseMultipleOAC.json"

    Scenario: When Searching for With both OAC and SIC that have special characters and multiple SIC then I get resp as in json
        When I GET "/v1/scrubber/search?q=01230,01240,01250,E00000001"
        And the response body is the same as the json in "./features/testdata/expecteddata/fullResponseMultipleSIC.json"
        
    Scenario: When Searching for With both OAC and SIC that have special characters and multiple SIC/OAC then I get resp as in json
        When I GET "/v1/scrubber/search?q=01230,01240,01250,E00000001,E00000014,E00000016"
        And the response body is the same as the json in "./features/testdata/expecteddata/fullResponseMultipleOACandSIC.json"
        
    Scenario: When Searching for With both OAC and SIC where OAC is has a typo but is correct length then I get resp as in json
        When I GET "/v1/scrubber/search?q=01230,01240,01250,E00000015"
        And the response body is the same as the json in "./features/testdata/expecteddata/fullResponseIfOACHasATypo.json"
        
    Scenario: When Searching for With both OAC and SIC where SIC is has a typo but is correct length then I get resp as in json
        When I GET "/v1/scrubber/search?q=01230,01240,01251,E00000014"
        And the response body is the same as the json in "./features/testdata/expecteddata/fullResponseIfSICHasATypo.json"