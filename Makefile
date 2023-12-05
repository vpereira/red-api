all: client server


# if you want to build to run evil calculators or any other UI payload, you need windowsgui
client:
	GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui"  -o ./injectus/injectus.exe ./injectus/main.go

client-console:
	GOOS=windows GOARCH=amd64 go build -o ./injectus/injectus.exe ./injectus/main.go

# we use garble to scramble the literals, to make harder static analysis
client-release:
	GOOS=windows GOARCH=amd64 garble -literals -seed=random build -ldflags -H=windowsgui -ldflags "-s -w" -o ./injectus/injectus.exe ./injectus/main.go



server:
	CGO_ENABLED=0 GOOS=linux go build -o ./webapi/webapi ./webapi/main.go
tidy:
	go mod tidy

clean:
	go clean

