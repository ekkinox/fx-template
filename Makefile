up:
	docker compose up -d

down:
	docker compose down

fresh:
	docker compose down
	docker compose build --no-cache
	docker compose up -d --build -V

logs:
	docker compose logs -f

build:
	docker build --platform linux/amd64 -t $(name) -f Dockerfile .

delve:
	docker compose exec -it app dlv debug