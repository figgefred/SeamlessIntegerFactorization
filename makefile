all: factoring

factoring: src/factoring.go src/task.go src/pollardrho.go
	go build -o bin/factoring src/*.go

clean: 
	rm bin/*
