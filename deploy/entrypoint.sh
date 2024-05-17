#!/bin/sh

# Exit immediately if a command exits with a non-zero status
set -e

# Function to handle the termination signals
term_handler() {
  echo "Terminating Nginx..."
  nginx -s quit
  exit 0
}

# Trap the termination signals
trap 'term_handler' INT TERM

# Start Nginx in the background
nginx -g 'daemon off;' &

# Wait for Nginx to exit
wait $!
