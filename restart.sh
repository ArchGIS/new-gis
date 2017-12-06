#!/bin/bash
# Stop server
killall new-gis

# Start server
go install
export GIN_MODE=release new-gis &