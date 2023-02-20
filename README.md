# Donut Backend
#### Built in ```Golang``` 🚀

## Running the project
* Clone the project
  ```git clone https://github.com/orgdonut/backend.git```
* Setup the  ```config.yml``` from given ```config-example.yml```
* install dependencies using the command: ```go get .```
* build the project: ```go build .```
* run the server: ```go run main.go```
## Conventions 🤌
* Branch Naming - braches should be named in snake_case with following tags separated by '/' at the start:
  * **feat** - indicating a new feature has been implemented in this branch
  * **fix** - indicating a bug fix implemented or a feature updated with fix
  * **ref** - indicating a refactor of directories, files, names etc.
  * **hot** - indicating hotfixes on branches branched from prod. These branches should be merged in prod and main.
* File Naming - all files should be named in snake_case
* Dir Naming - all dirs should be named in snake_case
* Avoid using long names for branch names, file names, class names or public variables
