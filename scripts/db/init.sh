#!/bin/bash
echo "*** CREATING TABLE AND SEED DATA STARTS ***"
psql -U $POSTGRES_USER -d $POSTGRES_DB -a -f /app/scripts/db/dump.sql
echo "*** CREATING TABLE AND SEED DATA FINISH ***"