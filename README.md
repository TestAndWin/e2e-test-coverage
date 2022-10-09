# Overview
...
**e2e test coverage**
Currently only Mocha reports are supported.

# Installation
* Download the binary matching the operating system.
* **e2e test coverage** needs a MySQL database. The DB connection is configured in the ```db.env``` file. This file is stored in the same directory as the binary.

  Example:
  ```
  DB_USER = "root"
  DB_PASSWORD = "root"
  DB_HOST = "127.0.0.1:3306"
  DB_NAME   = "e2ecoverage"
  ```
* Run ```e2e-coverage```

  This starts **e2e test coverage** on port 8080. 

# Guide 
* Open **e2e test coverage** URL in the browser
* Select _Product_ and enter a product name
* Add areas and features to the product
* Adapt the test files of your test automation to have a title with this build-up: ```{area name}|{feature name}|{suite name}```, e.g. for Cypress Tests ```describe('{area name}|{feature name}|{suite name}', () => {```
* Upload the mocha report using the REST API endpoint (directly from the CI/CD pipeline)

  Example:
  ```curl -d @mocha-report1_1.json http://localhost:8080/api/v1/coverage/1/upload-mocha-report```
* Explore the _Coverage_ section
* Start logging exploratory tests by clicking on _Log new_


# Further Ideas
- [ ] Allow to re-order areas and features
- [ ] Create area + feature automatically when uploading reports
- [ ] Allow to add bugs in production (Ticket Number, Short Desc, severity). Or better: get this automatically
- [ ] Allow to upload a Mocha summary file
- [ ] Allow to use more than one product / maintain products
- [ ] Security
- [ ] Login
- [ ] Include SLA data, e.g. from Datadog