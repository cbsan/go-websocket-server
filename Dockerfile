FROM alpine

WORKDIR /cbio

COPY ./bin/server .

COPY ./home.html .

EXPOSE 8080

CMD ["./server"]
