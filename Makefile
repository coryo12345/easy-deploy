all: build

build: buildWeb
	@echo "Building..."
	go build -o main cmd/main.go

run: buildWeb
	go run cmd/main.go

buildWeb:
	npm run build:css
	templ generate

test:
	@echo "Testing..."
	go test ./tests -v

install:
	go install github.com/a-h/templ/cmd/templ@latest
	npm install
	cp ./node_modules/htmx.org/dist/htmx.min.js ./web/static/htmx.min.js
	cp ./node_modules/alpinejs/dist/cdn.min.js ./web/static/alpine.cdn.min.js

clean:
	@echo "Cleaning..."
	rm -f main

watch:
	@if command -v air > /dev/null; then \
		air; \
		echo "Watching...";\
	else \
		read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/cosmtrek/air@latest; \
			air; \
			echo "Watching...";\
		else \
			echo "You chose not to install air. Exiting..."; \
			exit 1; \
		fi; \
	fi

.PHONY: all build run test clean watch install
