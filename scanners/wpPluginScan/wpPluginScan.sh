#!/bin/bash

if [[ $# != 1 ]]; then
 echo "Usage: program [http://site.com|https://site.com]"
 exit 0
fi

wordlist="wordpress-popular-plugins.txt"    # CHANGE THIS IF NEEDED
output_file="plugins_detected.txt"
base_url="$1/wp-content/plugins/"
> "$output_file"

check_plugin() {
  plugin_name="$1"
  url="$base_url$plugin_name/"
  echo $url
  response=$(curl -s -o /dev/null -w "%{http_code}" "$url")
  
  if [ "$response" == "200" ]; then
    echo "Plugin detected: $plugin_name" >> "$output_file"
    echo "Plugin detected: $plugin_name"
  else
    echo "Plugin not detected: $plugin_name"
  fi
}

# Leer la
while read -r plugin; do
  check_plugin "$plugin"
done < "$wordlist"

echo "Scan completed,results stored in $output_file."
