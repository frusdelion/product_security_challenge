# Zendesk Product Security
### The Zendesk Product Security Challenge

- [Instructions](#instructions)
- [Features](#features)
- [Installation](#installation)

## Installation
### Requirements
* Docker Desktop
* Go 1.13

### macOS and Linux install
Open a bash shell and run the following:  
```bash
$ docker run --rm -it -p 2525:1025 -p 8080:1080 djfarrelly/maildev --incoming-user=test --incoming-pass=test
```

After that, run the binary as follows:
```bash
$ ./zendesk-product_security_challenge
```

The following list are the web endpoints you should visit to interact with the site:
* [https://localhost](https://localhost) - Website
* [http://127.0.0.1:8080](http://127.0.0.1:8080) - Mail

### Windows install
Open a terminal and run the following:
```
> docker run --rm -it -p 2525:1025 -p 8080:1080 djfarrelly/maildev --incoming-user=test --incoming-pass=test
```

After that, run the binary as follows:
```
> ./zendesk-product_security_challenge.exe
```

The following list are the web endpoints you should visit to interact with the site:
* [https://localhost](https://localhost) - Website
* [http://127.0.0.1:8080](http://127.0.0.1:8080) - Mail

### Build from source (macOS)
Clone the repository, and run the following commands. Commands to download the compromised password list
is optional, as the release `build.companion.zip` contains the `common-passwords.db`.
```bash
$ go mod download
$ make run

# Compromised password list
$ make import-passwords 
$ make setup-passwords

# Generate the releases
$ make pkg
```

## Features
The web application has the following features:
- [x] Brand new look (Bootstrap)
- [x] Input sanitization and validation
- [x] Password hashed
    - bcrypt
- [x] Prevention of timing attacks
    - Rate-limiting for all requests to the site
    - Maximum failed attempts on login
    - Temporary ban when maximum failed attempts threshold reached
- [x] Logging
    - Logging all DB (optional)
    - Logging all requests
- [x] CSRF prevention
- [x] Multi factor authentication
    - MFA through [emails](https://github.com/matcornic/hermes)
- [x] Password reset / forget password mechanism
    - Password reset through [emails](https://github.com/matcornic/hermes)
- [ ] Account lockout
    - Temporary ban by IP & Browser
    - User account is not locked when temporary ban is placed on offending client
- [x] Cookie
    - Secure, HttpOnly
- [x] HTTPS
    - Self-signed certificates
- [x] Known password check
    - [zxcvbn by Dropbox](https://github.com/dropbox/zxcvbn) for client-side password advice
    - [Valve fingerprintjs2](https://github.com/Valve/fingerprintjs2) for browser fingerprinting
    - [CrackStation's password cracking dictionary](https://crackstation.net/buy-crackstation-wordlist-password-cracking-dictionary.htm) for matching against well-known passwords.

The application starts with an empty database. Please sign up as a user before trying out the features.

### Future Steps
- [ ] Test coverage
- [ ] OAuth2

## Instructions

Hello friend,

We are super excited that you want to be part of the Product Security team at Zendesk.

**To get started, you need to fork this repository to your own Github profile and work off that copy.**

In this repository, there are the following files:
1. README.md - this file
2. project/ - the folder containing all the files that you require to get started
3. project/index.html - the main HTML file containing the login form
4. project/assets/ - the folder containing supporting assets such as images, JavaScript files, Cascading Style Sheets, etc. You shouldn’t need to make any changes to these but you are free to do so if you feel it might help your submission

As part of the challenge, you need to implement an authentication mechanism with as many of the following features as possible. It is a non exhaustive list, so feel free to add or remove any of it as deemed necessary.

1. Input sanitization and validation
2. Password hashed
3. Prevention of timing attacks
4. Logging
5. CSRF prevention
6. Multi factor authentication
7. Password reset / forget password mechanism
8. Account lockout
9. Cookie
10. HTTPS
11. Known password check

You will have to create a simple binary (platform of your choice) to provide any server side functionality you may require. Please document steps to run the application. Your submission should be a link to your Github repository which you've already forked earlier together with the source code and binaries.

Thank you!
