#!/bin/sh
set -e


# validate
if [[ ! "$VITE_API_URL" =~ ^https?://[a-zA-Z0-9./?&=#_-]+$ ]]; then
  echo "Invalid VITE_API_URL"
  exit 1
fi


# Script to replace environment variables in files in runtime
# Generate the env.js file from the template and the environment variables
if [ -f /usr/share/nginx/html/env.js.template ]; then
  echo "Generating /usr/share/nginx/html/env.js..."
  envsubst < /usr/share/nginx/html/env.js.template > /usr/share/nginx/html/env.js
else
  echo "env.js.template not found"
fi

# Start Nginx in the foreground
exec nginx -g 'daemon off;'
