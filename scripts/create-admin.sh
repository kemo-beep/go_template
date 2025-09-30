#!/bin/bash

# Create admin user via API
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "Admin123!",
    "name": "Admin User"
  }'

echo ""
echo "✅ Admin user created!"
echo "Email: admin@example.com"
echo "Password: Admin123!"

