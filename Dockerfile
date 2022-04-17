FROM golang:1.17.9-alpine as build

RUN apk add --update --no-cache bash git
WORKDIR /k8s-env-manager

COPY . .

RUN go build -o k8s-env-manager

FROM golang:1.17.9-alpine

WORKDIR /k8s-env-manager

RUN adduser --disabled-password  --gecos "" shivam
USER shivam

COPY --chown=shivam --from=build /k8s-env-manager/k8s-env-manager .
COPY --chown=shivam --from=build /k8s-env-manager/config ./config

EXPOSE 8080

CMD ./k8s-env-manager
