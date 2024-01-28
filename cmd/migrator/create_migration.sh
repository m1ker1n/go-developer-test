#!/bin/bash
#Usage: create_migration <migration name>

migrationName=$1
goose -dir ../../migrations create "$migrationName" sql