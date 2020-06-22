FROM golang:rc-alpine3.12@sha256:98187f74d7837b7ad75acc09390cd31d87904d43ca590a7ff4f6f56bc3f31710 as builder
#Use go modules
ENV GO111MODULE=on
#user and usergroup for application- non root user
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group
RUN apk update && apk add --no-cache git ca-certificates tzdata
RUN mkdir /build
COPY . /build/
WORKDIR /build

COPY ./.netrc /root/.netrc
RUN chmod 600 /root/.netrc

COPY go.mod .
COPY go.sum .

RUN go mod download

#delete the netrc file after use
RUN rm -f /root/.netrc
#import code to build - not very optimistic at this stage
COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o authclient /cmd/client
FROM scratch AS final
LABEL author="David Solly"

#timezone info
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

#import user and group files to use
COPY --from=builder /user/group /user/passwd /etc/

#import certification
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

#import the built go binary of our program
COPY --from=builder /build/authclient /

#change work directory to root location
WORKDIR /

#Run the program as an unprivileged user
#Set current user
USER nobody:nobody
#Call the program
ENTRYPOINT [ "/authclient" ]
CMD [ "-consul.addr","localhost" ,"-consul.port","8500" ]

#Expose port through container to our application
EXPOSE 8080