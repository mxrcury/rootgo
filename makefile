
BINARY_FILE="main"

dev:
	go build -o ${BINARY_FILE}
	./${BINARY_FILE}
	rm ${BINARY_FILE}
	