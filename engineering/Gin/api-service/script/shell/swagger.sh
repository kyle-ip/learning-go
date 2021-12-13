
# swag
mkdir -p "$GOPATH"/src/github.com/swaggo
cd "$GOPATH"/src/github.com/swaggo
git clone https://github.com/swaggo/swag
cd swag/cmd/swag/
go install -v

# gin-swagger
cd "$GOPATH"/src/github.com/swaggo
git clone https://github.com/swaggo/gin-swagger

# swag init