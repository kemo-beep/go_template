# Admin Dashboard - Developer Console

A comprehensive admin dashboard for managing your Go backend, similar to Firebase Console or Supabase Studio.

## 🎯 Features

### 1. **User Management** (`/dashboard/users`)
- View all users with search and filtering
- Create new users
- Delete users
- Reset user passwords
- View user status (active/inactive)
- Manage admin roles

### 2. **Database Explorer** (`/dashboard/database`)
- List all database tables
- Browse table data with pagination
- Execute custom SQL queries
- View database schema and relationships
- Export data

### 3. **Storage Explorer (R2)** (`/dashboard/storage`)
- Browse uploaded files
- Upload new files
- Delete files
- Generate signed download URLs
- View file metadata (size, type, owner)

### 4. **Logs & Monitoring** (`/dashboard/logs`)
- View API request logs
- Filter logs by level, time range
- Real-time log streaming
- Search through logs

### 5. **Metrics** (`/dashboard/metrics`)
- Requests per minute
- Error rates
- Average latency
- Database performance
- Custom metric charts

### 6. **Developer Tools** (`/dashboard/dev-tools`)
- Run database migrations (up/down)
- View migration status
- Execute background jobs
- Manage feature flags
- Toggle system settings

## 🛠️ Tech Stack

- **Next.js 14** - React framework
- **TypeScript** - Type safety
- **Shadcn UI** - Beautiful components
- **Tailwind CSS** - Styling
- **React Query** - Data fetching & caching
- **Zustand** - State management
- **Axios** - HTTP client
- **Better Auth** - Authentication

## 📁 Project Structure

```
frontend/
├── app/
│   ├── dashboard/
│   │   ├── users/            # User management
│   │   ├── database/         # Database explorer
│   │   ├── storage/          # Storage (R2) management
│   │   ├── logs/             # API logs viewer
│   │   ├── metrics/          # Metrics & analytics
│   │   ├── dev-tools/        # Developer tools
│   │   └── settings/         # Settings
│   ├── login/                # Login page
│   ├── layout.tsx            # Root layout
│   └── providers.tsx         # React Query & global providers
├── components/
│   ├── ui/                   # Shadcn UI components
│   └── dashboard-layout.tsx  # Main dashboard layout
├── lib/
│   ├── api-client.ts         # API client with Axios
│   ├── store.ts              # Zustand stores
│   └── utils.ts              # Utility functions
└── public/                   # Static assets
```

## 🚀 Getting Started

### 1. Install Dependencies

```bash
cd frontend
npm install
```

### 2. Environment Variables

Create `.env.local`:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

### 3. Run Development Server

```bash
npm run dev
```

Visit `http://localhost:3000`

### 4. Build for Production

```bash
npm run build
npm start
```

## 🔐 Authentication

The dashboard uses JWT tokens for authentication:

1. Login page at `/login`
2. Token stored in localStorage
3. Axios interceptor adds token to all requests
4. Automatic redirect on 401 errors

## 📊 State Management

### Auth Store (Zustand)
```typescript
const { user, token, login, logout } = useAuthStore();
```

### UI Store (Zustand)
```typescript
const { sidebarOpen, theme, toggleSidebar, setTheme } = useUIStore();
```

## 🔌 API Integration

### Using React Query Hooks

```typescript
// Fetch users
const { data, isLoading, error } = useQuery({
  queryKey: ['users'],
  queryFn: () => api.getUsers(),
});

// Create user
const mutation = useMutation({
  mutationFn: (data) => api.createUser(data),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['users'] });
  },
});
```

### Available API Functions

See `lib/api-client.ts` for all available API functions:
- `api.getUsers()`, `api.createUser()`, `api.updateUser()`, etc.
- `api.getFiles()`, `api.uploadFile()`, `api.deleteFile()`, etc.
- `api.getTables()`, `api.executeQuery()`, etc.
- `api.getLogs()`, `api.getMetrics()`, etc.

## 🎨 Styling

Uses Tailwind CSS with Shadcn UI components:

- **Colors**: Neutral theme by default
- **Typography**: Geist font family
- **Dark Mode**: Ready (toggle in UI store)
- **Responsive**: Mobile-first design

## 📦 Components

### Pre-built UI Components (Shadcn)
- Button, Card, Input, Table
- Dialog, Alert Dialog, Dropdown Menu
- Badge, Avatar, Tooltip
- Tabs, Select, Sonner (toast)

### Custom Components
- `dashboard-layout.tsx` - Main layout with sidebar
- More components in respective feature pages

## 🔧 Customization

### Adding New Pages

1. Create page in `app/dashboard/[feature]/page.tsx`
2. Add route to navigation in `components/dashboard-layout.tsx`
3. Create API functions in `lib/api-client.ts`
4. Use React Query hooks for data fetching

### Adding New API Endpoints

```typescript
// lib/api-client.ts
export const api = {
  // Add new endpoint
  getAnalytics: () => apiClient.get('/api/v1/admin/analytics'),
};
```

### Customizing Theme

Edit `app/globals.css` and `components.json`

## 🚦 Backend Requirements

Your Go backend should provide these admin endpoints:

```
GET    /api/v1/admin/users
POST   /api/v1/admin/users
GET    /api/v1/admin/users/:id
PUT    /api/v1/admin/users/:id
DELETE /api/v1/admin/users/:id
POST   /api/v1/admin/users/:id/reset-password

GET    /api/v1/admin/database/tables
GET    /api/v1/admin/database/tables/:table
POST   /api/v1/admin/database/query
GET    /api/v1/admin/database/schema

GET    /api/v1/admin/logs
GET    /api/v1/admin/metrics

POST   /api/v1/admin/migrations/run
GET    /api/v1/admin/migrations
GET    /api/v1/admin/feature-flags
PUT    /api/v1/admin/feature-flags/:name
POST   /api/v1/admin/jobs/run
```

## 🔒 Security Considerations

1. **Admin-only routes**: Protect all `/admin/*` endpoints with admin middleware
2. **RBAC**: Implement role-based access control
3. **SQL injection**: Sanitize all SQL queries
4. **Rate limiting**: Protect sensitive endpoints
5. **Audit logs**: Log all admin actions

## 📝 TODO: Features to Implement

- [ ] Database Explorer page
- [ ] Storage Explorer page
- [ ] Logs Viewer page
- [ ] Metrics Dashboard page
- [ ] Developer Tools page
- [ ] Settings page
- [ ] Login page
- [ ] Dark mode toggle
- [ ] Export functionality
- [ ] Bulk operations
- [ ] Advanced filtering
- [ ] Real-time updates (WebSockets)

## 🐛 Debugging

### React Query Devtools
Enabled in development mode - click the React Query icon in bottom corner

### Error Handling
- API errors: Shown via Sonner toasts
- Auth errors: Automatic redirect to login
- Network errors: Retry logic built-in

## 📚 Resources

- [Next.js Docs](https://nextjs.org/docs)
- [Shadcn UI](https://ui.shadcn.com/)
- [React Query](https://tanstack.com/query/latest)
- [Zustand](https://zustand-demo.pmnd.rs/)
- [Tailwind CSS](https://tailwindcss.com/)

## 🤝 Contributing

1. Create feature branch
2. Implement feature
3. Add tests
4. Submit PR

## 📄 License

MIT
