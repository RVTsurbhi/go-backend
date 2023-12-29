steps for initiating the project

step 1: go mod init <module_name>
this is similar like npm init
step 2: install the package like npm install - go get github.com/gin-gonic/gin
you will see some files got added automatically
2.1 go.sum - its similar like package.lock.json
2.2 go.mod - similar like package.json

step3: install mongodb driver to store data
go get go.mongodb.org/mongo-driver/mongo

step4 : Install godotenv package to access env variables from .env file
go get github.com/joho/godotenv

step 5: run server cmd
go run main.go

step 6: To install all packages like npm install
go mod tidy
