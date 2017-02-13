echo "Downloading dependencies..."
go get -u github.com/llgcode/draw2d

mkdir -p ./bin

echo "Building Linux binaries..."
GOOS=linux GOARCH=amd64 go build -o ./bin/town-gen_linux_64bit
GOOS=linux GOARCH=386 go build -o ./bin/town-gen_linux_32bit
echo "Building macOS binaries..."
GOOS=darwin GOARCH=amd64 go build -o ./bin/town-gen_macos_64bit
GOOS=darwin GOARCH=386 go build -o ./bin/town-gen_macos_32bit
echo "Building Windows binaries..."
GOOS=windows GOARCH=amd64 go build -o ./bin/town-gen_windows_64bit
GOOS=windows GOARCH=386 go build -o ./bin/town-gen_windows_32bit
