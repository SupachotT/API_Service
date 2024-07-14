startDB:
	docker run --name pg-container -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres

createDB:
	docker exec -ti pg-container createdb -U postgres gopgtest

joinDB:
	docker exec -ti pg-container psql -U postgres

stopDB:
	docker stop pg-container

rmDB:
	docker rm pg-container

go:
	go run ./main.go