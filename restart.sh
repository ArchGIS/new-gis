#!/bin/bash
# Stop server
killall new-gis

# Start server
go install
new-gis &