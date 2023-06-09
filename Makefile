up:
	docker compose up -d

down:
	docker compose down

fresh:
	docker compose down --remove-orphans
	docker compose build --no-cache
	docker compose up -d --build -V

logs:
	docker compose logs -f

build:
	docker build --platform linux/amd64 -t $(name) -f Dockerfile .
