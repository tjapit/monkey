name = monkey
code = ./src/code/
compiler = ./src/compiler/
vm = ./src/vm/

all: build run

test: testcode testcomp testvm

build: 
	go build -o ${name} .

run: 
	./${name}

testcode: 
	go test ${code}
testcomp: 
	go test ${compiler}
testvm: 
	go test ${vm}
