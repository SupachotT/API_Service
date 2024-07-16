# Start Services
service:
	docker pull postgres
	docker run --name pg-container -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres

openDB:
	docker exec -ti pg-container createdb -U postgres gopgtest
	docker exec -ti pg-container psql -U postgres

clean:
	docker stop pg-container
	docker rm pg-container
	docker rmi postgres:latest



# Seperate Commands
pullDB:
	docker pull postgres

startDB:
	docker run --name pg-container -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres

createDB:
	docker exec -ti pg-container createdb -U postgres gopgtest

joinDB:
	docker exec -ti pg-container psql -U postgres

go:
	go run ./main.go