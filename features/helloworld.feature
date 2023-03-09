Feature: Helloworld
  Scenario: Posting and checking a response
    When I GET "/scrubber/search?q=dentists"
    Then I should receive a scrubber search empty response
    
    When I GET "/scrubber/search?q=W00009754"
    Then I should receive a scrubber search response with OAC codes populated
    
    When I GET "/scrubber/search?q=26513"
    Then I should receive a scrubber search response with Industry codes populated
    
    When I GET "/scrubber/search?q=26513%20W00009754"
    Then I should receive a scrubber search response full response   