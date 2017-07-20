#!/bin/bash
# Stop server
kill $(ps aux | grep '[n]ew-gis' | awk '{print $2}')

# Start server
new-gis &