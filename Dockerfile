# golang image where workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container’s workspace.
ADD .            /go/src/github.com/jacsmith21/lukabox
ADD ./api        /go/src/github.com/jacsmith21/lukabox/api
ADD ./domain     /go/src/github.com/jacsmith21/lukabox/domain
ADD ./ext/db     /go/src/github.com/jacsmith21/lukabox/ext/db
ADD ./ext/log    /go/src/github.com/jacsmith21/lukabox/ext/log
ADD ./ext/render /go/src/github.com/jacsmith21/lukabox/ext/render
ADD ./mock       /go/src/github.com/jacsmith21/lukabox/mock
ADD ./stc        /go/src/github.com/jacsmith21/lukabox/stc

RUN go get github.com/go-chi/jwtauth
RUN go get github.com/go-chi/render
RUN go get github.com/go-chi/chi
RUN go get github.com/Sirupsen/logrus
RUN go get github.com/go-errors/errors
RUN go get github.com/go-playground/validator

# Build the lukabox command inside the container.
RUN go install github.com/jacsmith21/lukabox

# Run the lukabox command when the container starts.
ENTRYPOINT /go/bin/lukabox

# http server listens on port 3001
EXPOSE 3001
