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


