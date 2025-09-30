# Admin Dashboard - Quick Start Guide

Get your admin dashboard up and running in minutes!

## 🚀 Quick Setup

### 1. Install Dependencies

```bash
cd frontend
npm install
```

### 2. Configure Environment

Create `.env.local`:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

### 3. Run Development Server

```bash
npm run dev
```

Visit: `http://localhost:3000`

## 📝 What's Included

✅ **Installed & Configured:**
- Next.js 14 with TypeScript
- Shadcn UI components
- React Query for data fetching
- Zustand for state management
- Axios for API calls
- Tailwind CSS for styling
- All required UI components

✅ **Pages Created:**
- `/` - Landing page
- `/dashboard` - Dashboard overview
- `/dashboard/users` - User management (COMPLETE)
- `/dashboard/database` - Database explorer (TODO)
- `/dashboard/storage` - Storage/R2 management (TODO)
- `/dashboard/logs` - API logs viewer (TODO)
- `/dashboard/metrics` - Metrics dashboard (TODO)
- `/dashboard/dev-tools` - Developer tools (TODO)
- `/dashboard/settings` - Settings (TODO)

✅ **Core Features:**
- Responsive sidebar navigation
- User authentication flow
- API client with automatic auth
- Toast notifications
- Loading states
- Error handling

## 🎯 Current Status

### ✅ Completed
1. **Project Setup**
   - Dependencies installed
   - Shadcn UI configured
   - Providers set up (React Query, Zustand)

2. **API Client**
   - Axios client configured
   - JWT auth interceptor
   - All API functions defined
   - TypeScript types

3. **State Management**
   - Auth store (Zustand)
   - UI store (Zustand)
   - Persistent storage

4. **Dashboard Layout**
   - Responsive sidebar
   - Top navigation bar
   - Route highlighting
   - User profile section

5. **User Management Page**
   - List all users
   - Search & filter
   - Create new user
   - Delete user
   - Reset password
   - View user details

### 🚧 To Be Implemented

1. **Database Explorer** (`/dashboard/database`)
   - List tables
   - Browse table data
   - Execute SQL queries
   - View schema

2. **Storage Explorer** (`/dashboard/storage`)
   - Browse files
   - Upload files
   - Delete files
   - Download files

3. **Logs Viewer** (`/dashboard/logs`)
   - View API logs
   - Filter by level/time
   - Search logs
   - Export logs

4. **Metrics Dashboard** (`/dashboard/metrics`)
   - Request rate charts
   - Error rate graphs
   - Latency metrics
   - Custom dashboards

5. **Developer Tools** (`/dashboard/dev-tools`)
   - Run migrations
   - Execute background jobs
   - Manage feature flags
   - System settings

6. **Authentication**
   - Login page
   - Logout functionality
   - Token refresh
   - Protected routes

## 📂 File Structure

```
frontend/
├── app/
│   ├── dashboard/
│   │   ├── users/page.tsx      ✅ Complete
│   │   ├── database/page.tsx   📝 TODO
│   │   ├── storage/page.tsx    📝 TODO
│   │   ├── logs/page.tsx       📝 TODO
│   │   ├── metrics/page.tsx    📝 TODO
│   │   ├── dev-tools/page.tsx  📝 TODO
│   │   ├── settings/page.tsx   📝 TODO
│   │   ├── layout.tsx          ✅ Complete
│   │   └── page.tsx            ✅ Complete
│   ├── layout.tsx              ✅ Complete
│   ├── page.tsx                ✅ Complete
│   └── providers.tsx           ✅ Complete
├── components/
│   ├── ui/                     ✅ 13 components
│   └── dashboard-layout.tsx    ✅ Complete
├── lib/
│   ├── api-client.ts           ✅ Complete
│   ├── store.ts                ✅ Complete
│   └── utils.ts                ✅ Complete
└── ...
```

## 🔧 Next Steps

### 1. Backend API Setup
Follow `ADMIN_API_SPEC.md` to implement backend endpoints:

```go
// Example: Create admin routes
adminGroup := router.Group("/api/v1/admin")
adminGroup.Use(middleware.Auth())
adminGroup.Use(middleware.AdminOnly())
adminGroup.Use(middleware.AuditLog())

adminGroup.GET("/users", admin.ListUsers)
adminGroup.POST("/users", admin.CreateUser)
// ... more routes
```

### 2. Implement Remaining Pages

**Database Explorer** (Priority: High)
```bash
# Create the page
touch app/dashboard/database/page.tsx
```

**Storage Explorer** (Priority: High)
```bash
# Create the page
touch app/dashboard/storage/page.tsx
```

**Logs Viewer** (Priority: Medium)
```bash
# Create the page
touch app/dashboard/logs/page.tsx
```

### 3. Add Authentication

```bash
# Install Better Auth (optional)
npm install better-auth

# Or implement custom JWT auth
# Create login page
touch app/login/page.tsx
```

### 4. Customize Styling

```css
/* app/globals.css */
/* Add your custom styles */
:root {
  --primary: your-color;
}
```

### 5. Add Tests

```bash
# Install testing libraries
npm install -D @testing-library/react @testing-library/jest-dom jest

# Create test files
touch app/dashboard/users/__tests__/page.test.tsx
```

## 🎨 Customization Guide

### Change Theme Colors

Edit `components.json`:
```json
{
  "style": "default",
  "tailwind": {
    "baseColor": "neutral"  // Change to: slate, gray, zinc, etc.
  }
}
```

### Add New Page

1. Create page file:
```tsx
// app/dashboard/my-feature/page.tsx
'use client';

export default function MyFeaturePage() {
  return <div>My Feature</div>;
}
```

2. Add to navigation:
```tsx
// components/dashboard-layout.tsx
const navigation = [
  // ... existing items
  { name: 'My Feature', href: '/dashboard/my-feature', icon: Star },
];
```

3. Add API functions:
```typescript
// lib/api-client.ts
export const api = {
  // ... existing functions
  getMyData: () => apiClient.get('/api/v1/my-feature'),
};
```

### Add New UI Component

```bash
# Using Shadcn CLI
npx shadcn@latest add [component-name]

# Example: Add form component
npx shadcn@latest add form
```

## 🐛 Troubleshooting

### API Connection Issues
```typescript
// Check API URL in .env.local
NEXT_PUBLIC_API_URL=http://localhost:8080

// Verify CORS in backend
// Enable CORS for frontend URL
```

### Auth Token Issues
```typescript
// Clear auth storage
localStorage.removeItem('auth_token');
localStorage.removeItem('auth-storage');
```

### Component Not Found
```bash
# Reinstall dependencies
rm -rf node_modules
npm install
```

## 📚 Learning Resources

- **Next.js**: https://nextjs.org/docs
- **Shadcn UI**: https://ui.shadcn.com/
- **React Query**: https://tanstack.com/query/latest
- **Zustand**: https://zustand-demo.pmnd.rs/
- **Tailwind CSS**: https://tailwindcss.com/

## 🤝 Getting Help

1. Check `ADMIN_DASHBOARD.md` for detailed documentation
2. Review `ADMIN_API_SPEC.md` for backend requirements
3. Check component examples at https://ui.shadcn.com/

## 🎉 You're Ready!

Your admin dashboard foundation is set up. Now:

1. ✅ Run `npm run dev`
2. ✅ Visit `http://localhost:3000`
3. ✅ Start building remaining pages
4. ✅ Implement backend APIs
5. ✅ Customize to your needs

Happy coding! 🚀
