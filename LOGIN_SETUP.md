# 🔐 Login Setup Guide

## Problem: Cannot Login to Admin Dashboard

The login page requires a **user account in your database**. Here's how to set it up:

---

## ✅ **Quick Setup (3 Steps)**

### **Step 1: Start Your Backend**
Make sure your Go backend is running:
```bash
cd /Users/wonder/Documents/experimental/go_template
make dev
```
Backend should be running on: **http://localhost:8080**

### **Step 2: Create Admin User**
Run this script to create an admin account:
```bash
bash scripts/create-admin.sh
```

You should see:
```
✅ Admin user created!
Email: admin@example.com
Password: Admin123!
```

### **Step 3: Login**
Go to: **http://localhost:3001/login**

Use these credentials:
- **Email**: `admin@example.com`
- **Password**: `Admin123!`

---

## 🔧 **Manual User Creation (Alternative)**

If the script doesn't work, create a user manually via curl:

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "Admin123!",
    "name": "Admin User"
  }'
```

---

## 🐛 **Troubleshooting**

### Error: "Connection refused"
**Cause**: Backend is not running  
**Fix**: Start the backend with `make dev`

### Error: "Invalid credentials"
**Cause**: Wrong password or user doesn't exist  
**Fix**: 
1. Check if user exists in database
2. Re-run the create-admin script
3. Try registering via Swagger UI: http://localhost:8080/docs/index.html

### Error: "User already exists"
**Cause**: User was already created  
**Fix**: Just login with `admin@example.com` / `Admin123!`

### Database Not Available
**Cause**: PostgreSQL is not running  
**Fix**: Start PostgreSQL with `docker-compose up -d postgres`

---

## 📊 **Verify User in Database**

Connect to PostgreSQL:
```bash
docker exec -it go_template_postgres psql -U postgres -d go_mobile_backend
```

Check if user exists:
```sql
SELECT id, email, name, is_active, is_admin FROM users;
```

Exit:
```
\q
```

---

## 🔒 **Security Notes**

### **Production Setup**
For production, you should:
1. Use a strong admin password (not `Admin123!`)
2. Set `is_admin = true` in the database
3. Enable 2FA
4. Use environment variables for credentials

### **Set Admin Role**
If your user is not an admin, update in database:
```sql
UPDATE users SET is_admin = true WHERE email = 'admin@example.com';
```

---

## 🎯 **What Happens After Login**

Once logged in successfully:
1. JWT access token is stored in localStorage
2. User data is stored in Zustand state
3. You're redirected to `/dashboard`
4. All API calls will include the JWT token in headers
5. Protected routes will work automatically

---

## 📝 **API Response Structure**

Your backend returns this on successful login:
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": 1,
      "email": "admin@example.com",
      "name": "Admin User",
      "is_active": true,
      "is_admin": true
    },
    "expires_in": 1440
  }
}
```

The frontend extracts `access_token` and `user` from this response.

---

## ✅ **Next Steps**

After successful login, you can:
- View the dashboard overview
- Manage users
- Explore database tables
- View API logs
- Monitor metrics
- And much more!

---

**Need help?** Check the browser console for error messages.

