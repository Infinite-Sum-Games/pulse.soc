#!/bin/bash

# Array of sorted set names (excluding Leaderboard)
sets=(
  "cpp-ranking-sset"
  "java-ranking-sset" 
  "python-ranking-sset"
  "javascript-ranking-sset"
  "go-ranking-sset"
  "rust-ranking-sset"
  "zig-ranking-sset"
  "flutter-ranking-sset"
  "kotlin-ranking-sset"
  "haskell-ranking-sset"
)

# Sample usernames
usernames=(
  "IAmRiteshKoushik"
  "Ashrockzzz2003"
  "vijay-sb"
  "Leela0o5"
  "KiranRajeev-KV"
)

# For each sorted set
for set in "${sets[@]}"; do
  # Add 2-3 random users with random scores
  num_users=$((RANDOM % 2 + 2)) # 2 or 3 users
  
  for ((i=1; i<=num_users; i++)); do
    # Get random username
    username=${usernames[$RANDOM % ${#usernames[@]}]}
    
    # Generate random score between 1-10
    score=$((RANDOM % 10 + 1))
    
    # Add to sorted set
    redis-cli ZADD "$set" $score "$username"
  done
done
