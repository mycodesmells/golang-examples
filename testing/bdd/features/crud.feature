# features/crud.feature
Feature: CRUD

    Scenario: create and get quote
        Given I create quote "xxx" by "aaa"
        When I ask for last created quote
        Then I should get "xxx" as quote
        And I should get "bbb" as author

