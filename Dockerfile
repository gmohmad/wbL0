FROM scratch

WORKDIR /app

COPY . ./

CMD ["./builds/linux/main"]
