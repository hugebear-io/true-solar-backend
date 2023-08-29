# |=> API
app: 
	go run ./cmd/api/main.go

app_build:
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=1 go build -ldflags "-linkmode external" -o app ./cmd/api/main.go

# |=> Huawei
huawei:
	go run ./cmd/huawei/main.go

huawei_build:
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=1 go build -ldflags "-linkmode external" -o huawei ./cmd/huawei/main.go

# |=> INVT
invt:
	go run ./cmd/solarman/main.go

invt_build:
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=1 go build -ldflags "-linkmode external" -o solarman ./cmd/solarman/main.go