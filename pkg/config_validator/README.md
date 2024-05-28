# How to generate JSON schema from go struct playbook

In order to generate the JSON schema for the content configuration, the following steps need to be followed:

- generate the JSON schema from the go struct without `omitempty` tags

- run the `Test_createJson` unit test to fetch in REPL mode the JSON schema
