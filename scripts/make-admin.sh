#!/bin/bash

# Make a user admin by email
EMAIL=${1:-"admin@example.com"}

echo "Making user $EMAIL an admin..."

# Connect to PostgreSQL and update the user
PGPASSWORD=secret psql -h localhost -p 5433 -U app -d myapp -c "UPDATE users SET is_admin = true WHERE email = '$EMAIL';"

if [ $? -eq 0 ]; then
    echo "✅ User $EMAIL is now an admin!"
else
    echo "❌ Failed to make user admin. Make sure the user exists and PostgreSQL is running."
fi
