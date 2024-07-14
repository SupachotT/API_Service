# DATABEASES CONNECTION
This project is for create database from Image Docker "postgres" and connect postgres to Go Language
Reference form **[This](https://youtu.be/Y7a0sNKdoQk?si=j9jtPJX2drl7oB90)**.

## Create Database server
1. Pull Image Postgres from Docker hub.

    ```
    docker pull postgres
    ```
2. Run docker with command
   
    ```
    docker run --name pg-container -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres
    ```
3. Create Database server and connect to server

    ```
    docker exec -ti pg-container createdb -U postgres gopgtest    
    ```
    After that connect to database server with

    ```
    docker exec -ti pg-container psql -U postgres
    ```

4. Others

    for check that it connected to database with sql command  `\c` and `\q` for exit from databases