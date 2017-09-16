# Lukabox
This is the backend RESTful API for Lukabox!

## Setup
1. Install [vod](https://github.com/jacsmith21/vod) & make sure it's in your path
2. Create the folder stc `src/github.com/jacsmith21` in your GOPATH
3. Clone this repository in the `jacsmith21` folder
4. Run:
```
$ go get github.com/go-chi/jwtauth
$ go get github.com/go-chi/render
$ go get github.com/go-chi/chi
$ go get github.com/Sirupsen/logrus
$ go get github.com/go-errors/errors
$ go get github.com/go-playground/validator
```
5. In the `jacsmith21` folder, run:
```
$ vod run main.go
```

## TODO
* Check out OAuth for authentication & check out other authentication stuff
* Get code on google web services for testing with the box

## References
* https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1
* https://forum.golangbridge.org/t/comparing-the-structure-of-web-applications/1198/16
