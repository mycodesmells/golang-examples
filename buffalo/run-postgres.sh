#!/bin/bash
docker run --name some-postgres -e POSTGRES_PASSWORD=mysecretpassword -p 15432:5432 -d postgres
