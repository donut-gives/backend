# Donut Backend
#### Built in ```Golang``` 🚀
___
## Setup
1. Clone repository
2. Install ```mkcert``` - [How to Install](https://github.com/FiloSottile/mkcert#installation)
3. In your terminal, run `mkcert -install`
4. In your terminal, run `mkcert localhost` from **project folder**
5. Add the following to your **env variables**
   * GOOGLE_OAUTH_CLIENT_ID= _Google Oauth Client ID_
   * GOOGLE_OAUTH_CLIENT_SECRET= _Google Oauth Client Secret_
   * CERTIFICATE=localhost.pem
   * KEY=localhost-key.pem
6. Run `go run main.go`


## Guidelines
* Branch Naming - braches should be named in snake_case with following tags separated by '/' at the start:
  * feat - indicating a new feature has been implemented in this branch
  * fix - indicating a bug fix implemented or a feature updated with fix
  * ref - indicating a refactor of directories, files, names etc.
  * hot - indicating hotfixes on branches branched from prod. These branches should be merged in prod and main.
* File Naming - all files should be named in snake_case
* Class Naming - all class should be named in UpperCamelCase.
* Variable Naming - all variables should be names in lowerCamelCase.
* Avoid using long names for branch names, file names, class names or public variables
