FROM alpine:3.9

COPY sender/sender /

CMD ["/sender"]