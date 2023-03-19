# e2e test coverage 

There are several great test reporting tools that can provide an overview of automated test results. 

However, from my perspective, they often lack a comprehensive view of which product areas are covered by automated tests, along with the corresponding test status. 

I have yet to come across a tool that meets this requirement.

**e2e test coverage** offers a way to define a product overview that includes different areas and their features, as well as to upload test results of automatic tests for these areas.

![e2e test coverage](e2e-test-coverage.png "e2e test coverage")

For automated tests to be properly mapped to their corresponding areas and features, it is essential that the same identifiers are used. The results of these automated tests can be uploaded through a REST endpoint, and at present, only Mocha reports are supported.

**e2e test coverage** takes into account the test results from the last 28 days when displaying the coverage information.

It's important to note that this representation is solely a quantitative view and does not provide any insights into the quality of the tests. However, even so, having such an overview can still be helpful in my opinion.

In addition to displaying test coverage information, **e2e test coverage** also enables the collection of feedback from exploratory testing.

# Installation
* Download the Docker image and upload it via ```docker load -i e2ecoverage_<version>.tar.gz```

* **e2e test coverage** needs a MySQL database. The configuration can be set via Docker environment variables. The databases *e2ecoverage* and *user* are created when it is started for the first time.

* To start it, please enter ```docker run --env-file docker_env_vars --add-host host.docker.internal:host-gateway -d  -p 127.0.0.1:8080:8080 e2ecoverage```

  Example ```docker_env_vars```:
  ```
  DB_USER=root
  DB_PASSWORD=your db password
  DB_HOST=host.docker.internal:3306
  JWT_KEY=please enter a random value
  ```
  
# Guide 

## **e2e test coverage** user interface
* To access **e2e test coverage**, open the URL in your browser. Upon your first connection to the database, the user ```admin``` with the password ```e2ecoverage``` will be automatically created. It is highly recommended to change the password on the *My Account* page.

* As a user with the Admin role, you have the ability to create an API key on the *My Account* page. This API key is necessary to upload test results through HTTP requests, such as from a CI/CD pipeline. The admin can also manage other users and assign them roles:

  * Admin: has the capability to create new users and API keys
  * Maintainer: responsible for maintaining products, areas, and features
  * Tester: can view test coverage and add exploratory test results

* On the *Product* page, you can select a product and enter its name. Please note that the product name will not be visible later on. You can also add areas and features to the selected product.

* The *Coverage* page displays the number of test cases and their status over the past 28 days, including the impact of exploratory testing on each product area and feature. 

* The *Tests* page provides a complete overview of all tests conducted on the product, including their status over the last 28 days. This includes tests that have not been assigned to a specific area or feature. Only the most recent test for each suite name will be displayed."

## CI/CD integration
* Please adapt your test files to include the following format for the title: ```{area name}|{feature name}|{suite name}```, e.g. for Cypress Tests ```describe('{area name}|{feature name}|{suite name}', () => {```
* Upload the mocha report using the REST API endpoint (directly from the CI/CD pipeline):

  Example:
  ```curl -d @mocha-report1_1.json -H "apiKey: <your api key>" -H "testReportUrl: <Url where the generated Mocha report can be found>" http://localhost:8080/api/v1/coverage/1/upload-mocha-summary-report```

# Development
Please bear with me, this is my first Golang & Vue 3 project. I used

* Golang 1.20
* Vue 3
* Bootstrap 5

## Folder Structure
* ```api/cmd/coverage``` Main Golang Application
* ```api/config``` Read config
* ```api/coverage``` Endpoints to maintain areas and features, upload tests and coverage implementation
* ```api/db``` Common DB function to connect to a database
* ```api/docs``` API Doc generated by Swag
* ```api/router``` Route to the different endpoints
* ```api/ui``` Helper to include the Vue application in the Golang binary
* ```api/user``` Endpoints to maintain user and to log-in
* ```ui/``` Vue 3 application

## Development 
By running the command ```make help```, you can obtain a comprehensive overview of the various targets available, such as building the app, starting it locally in development mode, and others."

# TO-DO
- Add forgot password feature

# Further Ideas
- Allow to re-order areas and features
- Create area + feature automatically when uploading reports
- Allow to add bugs in production (Ticket Number, Short Desc, severity). Or better: get this automatically
- Allow to use more than one product / maintain products
- Include SLA data, e.g. from Datadog