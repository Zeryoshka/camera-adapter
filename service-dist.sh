go build -o upr-conf-api ./cmd/confapi/main.go
sudo cp -f upr-conf-api /usr/local/bin/
go build -o cam-upr ./cmd/main.go
sudo cp -f cam-upr /usr/local/bin/