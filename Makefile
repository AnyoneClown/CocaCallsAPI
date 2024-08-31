BINARY_NAME = CocaCalls

build:
	cd backend && go build -o /tmp/bin/${BINARY_NAME} .

clean:
	cd backend && go clean
	rm /tmp/bin/${BINARY_NAME}

run: build
	cd backend && /tmp/bin/${BINARY_NAME} & \
	cd frontend && npm run dev

tidy:
	cd backend && go mod tidy -v && gofmt -w .
