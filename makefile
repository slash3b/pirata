build:
	docker build -f ./services/scraper/Dockerfile -t registry.digitalocean.com/pirata/scraper:latest .
	docker build -f ./services/imdb/Dockerfile -t registry.digitalocean.com/pirata/imdb:latest .
	docker build -f ./services/api/Dockerfile -t registry.digitalocean.com/pirata/api:latest .

push-registry:
	docker push registry.digitalocean.com/pirata/scraper:latest
	docker push registry.digitalocean.com/pirata/imdb:latest
	docker push registry.digitalocean.com/pirata/api:latest
