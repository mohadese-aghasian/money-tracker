CompileDaemon -directory="./cmd/api" -build="go build -o main.exe" -command="./cmd/api/main.exe"
CompileDaemon -command="swag init && go run main.go"

"user_name": "porsio.admin",
"password": "12345"

swag init -g cmd/api/main.go -o cmd/api/docs

go run migrate_runner/main.go
go run migrate_runner/main.go -rollback

////---------------env variables
PORT=8088
GIN_MODE=debug
JWT_SECRET=
SWAGGER_HOST=localhost:8088

DB_HOST=
DB_PORT=
DB_USER=
DB_PASS=
DB_NAME=

IMAGE_MAX_SIZE=5
Video_MAX_SIZE=15
MAX_REQUEST_SIZE_MB=32
