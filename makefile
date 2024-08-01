name = monkey

all: build run

build: 
	go build -o ${name} .

run: 
	./${name}
