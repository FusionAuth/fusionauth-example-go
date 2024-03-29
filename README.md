# Archived
A new version of this repo is available here: [Go Resource Server with FusionAuth](https://github.com/FusionAuth/fusionauth-quickstart-golang-api).

# The Example For using FusionAuth and Golang

Simple demo for using [FusionAuth](http://fusionauth.io/) with Golang.

You can view the corresponding blog post: https://fusionauth.io/blog/2020/10/22/securing-a-golang-app-with-oauth/

This application will use an OAuth Authorization Code workflow and the PKCE extension to log users in. PKCE stands for Proof Key for Code Exchange, and is often pronounced “pixie”.

## Installation

clone with git first

```bash
git clone https://github.com/fusionauth/fusionauth-example-go
```

## Usage

Assuming you've configured FusionAuth with a new application, update `main.go` with the client ID and client secret.

```shell
go get github.com/thanhpk/randstr
go get golang.org/x/oauth2
go get github.com/nirasan/go-oauth-pkce-code-verifier
go run main.go
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## Who made this?

[Krissanawat Kaewsanmaung](https://github.com/krissnawat) - Creator


## License
[APACHE 2.0](https://www.apache.org/licenses/LICENSE-2.0)

