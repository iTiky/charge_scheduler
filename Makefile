all: install

lint:
	golangci-lint run --exclude 'unused'

test:
	go test -v ./... --count=1

pack_migrations:
	#brew install go-bindata
	go-bindata -o ./storage/sqlite_base/resources/migrations.go -prefix "./storage/sqlite_base/migrations/" -pkg resources ./storage/sqlite_base/migrations/

install:
	go build -o ./build/charge-scheduler ./cmd
