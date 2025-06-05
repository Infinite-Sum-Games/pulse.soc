#!/bin/bash

# Populate Redis stream with sample live events
for i in {1..5}; do
  # Generate random event type
  event_types=("Bounty" "Issue" "Top-3")
  event_type=${event_types[$RANDOM % ${#event_types[@]}]}
  
  # Generate sample usernames
  usernames=("IAmRiteshKoushik" "Ashrockzzz2003" "vijay-sb" "Leela0o5" "KiranRajeev-KV")
  username=${usernames[$RANDOM % ${#usernames[@]}]}
  
  # Generate appropriate message based on event type
  case $event_type in
    "Bounty")
      message="$username claimed a bounty worth 500 points!"
      ;;
    "Issue") 
      message="$username opened a new issue"
      ;;
    "Top-3")
      message="$username made it to top 3 in Python category!"
      ;;
  esac

  # Current timestamp in milliseconds
  timestamp=$(date +%s%3N)

  # Create JSON payload
  json_data=$(cat <<EOF
{
  "github_username": "$username",
  "message": "$message", 
  "event_type": "$event_type",
  "time": $timestamp
}
EOF
)

  # Add to Redis stream
  redis-cli XADD live-update-stream \* data "$json_data"
  
  # Small delay between events
  sleep 1
done
