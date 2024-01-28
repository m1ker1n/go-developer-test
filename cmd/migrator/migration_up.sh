#!/bin/bash
#Usage: migration_up <postgres db string>

postgresDbString=$1
goose -dir ../../migrations postgres "$postgresDbString" up