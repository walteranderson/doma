repl: build
	rlwrap ./doma
build:
	go build -C cmd -o ../doma
