all:
	go build -o build/apps/bossd cmd/bossd/bossd.go

linux:
	GOOS=linux go build -o build/apps/bossd cmd/bossd/bossd.go
