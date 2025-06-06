#!/bin/bash

# Redis connection details
REDIS_HOST="localhost"
REDIS_PORT="6379"

# Sample data - GitHub usernames and scores
declare -A leaderboard=(
  ["IAmRiteshKoushik"]=100
  ["aditya-menon-r"]=75 
  ["FLASH2332"]=150
  ["KiranRajeev-KV"]=200
  ["vijay-sb"]=125
)

# Add entries to Redis sorted set
for username in "${!leaderboard[@]}"; do
  # Randomly decide whether to add decimal points (1 in 3 chance)
  if [ $((RANDOM % 3)) -eq 0 ]; then
    # Generate random 3 digit decimal between 0-999
    decimal=$((RANDOM % 1000))
    # Pad with leading zeros if needed
    decimal=$(printf "%03d" $decimal)
    score=${leaderboard[$username]}.$decimal
  else
    score=${leaderboard[$username]}.000
  fi
  
  # Add to leaderboard sorted set
  redis-cli -h $REDIS_HOST -p $REDIS_PORT ZADD "leaderboard-sset" $score "$username"
  
  if [ $? -eq 0 ]; then
    echo "Added $username with score $score"
  else
    echo "Failed to add $username"
  fi
done

echo "Leaderboard population complete!"
