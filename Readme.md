make migrate-create name=add_comment_table
this is the way you create a new migration
make migrate-up
for migrations up
and for down 

air -d  // for debug mode

air  // for normal mode


docker exec -it postgres-db psql -U admin -d socialnetwork
docker exec -it <container-name> psql -U <db-user> -d <db-name>
environment:
  POSTGRES_USER: admin
  POSTGRES_PASSWORD: secret
  POSTGRES_DB: socialnetwork

And container name:

container_name: postgres-db

docker exec -it  → enter container terminal
container-name   → which container
psql             → postgres CLI
-U               → user
-d               → database