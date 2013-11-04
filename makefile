all: factoring

factoring: src/factoring.go
	go build -o bin/factoring src/factoring.go

clean: 
	rm bin/*
