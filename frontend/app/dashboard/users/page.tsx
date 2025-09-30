'use client';

import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from '@/components/ui/table';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from '@/components/ui/dialog';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { MoreVertical, Search, UserPlus, Shield, Users, UserCheck, UserX, Loader2 } from 'lucide-react';
import { toast } from 'sonner';
import { format } from 'date-fns';

// Types
interface User {
    id: number;
    email: string;
    name: string;
    is_active: boolean;
    is_admin: boolean;
    created_at: string;
    roles?: string[];
}

interface Role {
    id: number;
    name: string;
    description: string;
}

// API Client functions
const apiClient = {
    async getUsers(page = 1, limit = 20, search = '', role = '', isActive?: boolean) {
        const params = new URLSearchParams({
            page: page.toString(),
            limit: limit.toString(),
            ...(search && { search }),
            ...(role && { role }),
            ...(isActive !== undefined && { is_active: isActive.toString() }),
        });

        const response = await fetch(`http://localhost:8080/api/v1/admin/users?${params}`, {
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('auth_token')}`,
            },
        });

        if (!response.ok) {
            throw new Error('Failed to fetch users');
        }

        return response.json();
    },

    async updateUser(id: number, data: Partial<User>) {
        const response = await fetch(`http://localhost:8080/api/v1/admin/users/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('auth_token')}`,
            },
            body: JSON.stringify(data),
        });

        if (!response.ok) {
            throw new Error('Failed to update user');
        }

        return response.json();
    },

    async deleteUser(id: number) {
        const response = await fetch(`http://localhost:8080/api/v1/admin/users/${id}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('auth_token')}`,
            },
        });

        if (!response.ok) {
            throw new Error('Failed to delete user');
        }

        return response.json();
    },

    async getRoles() {
        const response = await fetch('http://localhost:8080/api/v1/admin/roles', {
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('auth_token')}`,
            },
        });

        if (!response.ok) {
            throw new Error('Failed to fetch roles');
        }

        return response.json();
    },

    async assignRole(userId: number, roleId: number) {
        const response = await fetch(`http://localhost:8080/api/v1/admin/users/${userId}/roles`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('auth_token')}`,
            },
            body: JSON.stringify({ role_id: roleId }),
        });

        if (!response.ok) {
            throw new Error('Failed to assign role');
        }

        return response.json();
    },

    async removeRole(userId: number, roleId: number) {
        const response = await fetch(`http://localhost:8080/api/v1/admin/users/${userId}/roles/${roleId}`, {
            method: 'DELETE',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('auth_token')}`,
            },
        });

        if (!response.ok) {
            throw new Error('Failed to remove role');
        }

        return response.json();
    },
};

export default function UsersPage() {
    const [search, setSearch] = useState('');
    const [page, setPage] = useState(1);
    const [roleFilter, setRoleFilter] = useState<string>('all');
    const [statusFilter, setStatusFilter] = useState<string>('all');
    const [selectedUser, setSelectedUser] = useState<User | null>(null);
    const [isEditOpen, setIsEditOpen] = useState(false);
    const [isRoleDialogOpen, setIsRoleDialogOpen] = useState(false);
    const [selectedRoleId, setSelectedRoleId] = useState<string>('');

    const queryClient = useQueryClient();

    // Queries
    const { data: usersData, isLoading } = useQuery({
        queryKey: ['admin-users', page, search, roleFilter, statusFilter],
        queryFn: () => apiClient.getUsers(
            page,
            20,
            search,
            roleFilter === 'all' ? undefined : roleFilter,
            statusFilter === 'all' ? undefined : statusFilter === 'active'
        ),
    });

    const { data: rolesData } = useQuery({
        queryKey: ['roles'],
        queryFn: () => apiClient.getRoles(),
    });

    // Mutations
    const updateUserMutation = useMutation({
        mutationFn: ({ id, data }: { id: number; data: Partial<User> }) =>
            apiClient.updateUser(id, data),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['admin-users'] });
            setIsEditOpen(false);
            toast.success('User updated successfully');
        },
        onError: () => {
            toast.error('Failed to update user');
        },
    });

    const deleteUserMutation = useMutation({
        mutationFn: (id: number) => apiClient.deleteUser(id),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['admin-users'] });
            toast.success('User deleted successfully');
        },
        onError: () => {
            toast.error('Failed to delete user');
        },
    });

    const assignRoleMutation = useMutation({
        mutationFn: ({ userId, roleId }: { userId: number; roleId: number }) =>
            apiClient.assignRole(userId, roleId),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['admin-users'] });
            setIsRoleDialogOpen(false);
            setSelectedRoleId('');
            toast.success('Role assigned successfully');
        },
        onError: () => {
            toast.error('Failed to assign role');
        },
    });

    const removeRoleMutation = useMutation({
        mutationFn: ({ userId, roleId }: { userId: number; roleId: number }) =>
            apiClient.removeRole(userId, roleId),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['admin-users'] });
            toast.success('Role removed successfully');
        },
        onError: () => {
            toast.error('Failed to remove role');
        },
    });

    const users = usersData?.data?.users || [];
    const total = usersData?.data?.total || 0;
    const roles: Role[] = rolesData?.data?.roles || [];

    const handleEditUser = (user: User) => {
        setSelectedUser(user);
        setIsEditOpen(true);
    };

    const handleManageRoles = (user: User) => {
        setSelectedUser(user);
        setIsRoleDialogOpen(true);
    };

    const handleSaveUser = () => {
        if (!selectedUser) return;

        updateUserMutation.mutate({
            id: selectedUser.id,
            data: {
                name: selectedUser.name,
                is_active: selectedUser.is_active,
                is_admin: selectedUser.is_admin,
            },
        });
    };

    const handleAssignRole = () => {
        if (!selectedUser || !selectedRoleId) return;

        assignRoleMutation.mutate({
            userId: selectedUser.id,
            roleId: parseInt(selectedRoleId),
        });
    };

    const handleRemoveRole = (roleId: number) => {
        if (!selectedUser) return;

        removeRoleMutation.mutate({
            userId: selectedUser.id,
            roleId,
        });
    };

    // Stats
    const activeUsers = users.filter((u: User) => u.is_active).length;
    const adminUsers = users.filter((u: User) => u.is_admin).length;

    return (
        <div className="space-y-6">
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">User Management</h1>
                    <p className="text-gray-500 mt-1">Manage user accounts, roles, and permissions</p>
                </div>
            </div>

            {/* Stats Cards */}
            <div className="grid gap-4 md:grid-cols-3">
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Total Users</CardTitle>
                        <Users className="h-4 w-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{total}</div>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Active Users</CardTitle>
                        <UserCheck className="h-4 w-4 text-green-600" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{activeUsers}</div>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Administrators</CardTitle>
                        <Shield className="h-4 w-4 text-purple-600" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{adminUsers}</div>
                    </CardContent>
                </Card>
            </div>

            {/* Filters */}
            <Card>
                <CardHeader>
                    <div className="flex items-center gap-4 flex-wrap">
                        <div className="relative flex-1 min-w-[200px]">
                            <Search className="absolute left-2 top-2.5 h-4 w-4 text-gray-400" />
                            <Input
                                placeholder="Search by email or name..."
                                value={search}
                                onChange={(e) => setSearch(e.target.value)}
                                className="pl-8"
                            />
                        </div>
                        <Select value={roleFilter} onValueChange={setRoleFilter}>
                            <SelectTrigger className="w-[180px]">
                                <SelectValue placeholder="Filter by role" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="all">All Roles</SelectItem>
                                {roles.map((role) => (
                                    <SelectItem key={role.id} value={role.name}>
                                        {role.name}
                                    </SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                        <Select value={statusFilter} onValueChange={setStatusFilter}>
                            <SelectTrigger className="w-[180px]">
                                <SelectValue placeholder="Filter by status" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="all">All Status</SelectItem>
                                <SelectItem value="active">Active</SelectItem>
                                <SelectItem value="inactive">Inactive</SelectItem>
                            </SelectContent>
                        </Select>
                    </div>
                </CardHeader>
                <CardContent>
                    {isLoading ? (
                        <div className="flex items-center justify-center py-8">
                            <Loader2 className="h-8 w-8 animate-spin text-gray-400" />
                        </div>
                    ) : (
                        <Table>
                            <TableHeader>
                                <TableRow>
                                    <TableHead>User</TableHead>
                                    <TableHead>Email</TableHead>
                                    <TableHead>Roles</TableHead>
                                    <TableHead>Status</TableHead>
                                    <TableHead>Created</TableHead>
                                    <TableHead className="text-right">Actions</TableHead>
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                {users.length === 0 ? (
                                    <TableRow>
                                        <TableCell colSpan={6} className="text-center text-gray-500 py-8">
                                            No users found
                                        </TableCell>
                                    </TableRow>
                                ) : (
                                    users.map((user: User) => (
                                        <TableRow key={user.id}>
                                            <TableCell>
                                                <div>
                                                    <div className="font-medium">{user.name}</div>
                                                    {user.is_admin && (
                                                        <Badge variant="outline" className="mt-1 bg-purple-50 text-purple-700 border-purple-200">
                                                            <Shield className="h-3 w-3 mr-1" />
                                                            Admin
                                                        </Badge>
                                                    )}
                                                </div>
                                            </TableCell>
                                            <TableCell className="text-sm text-gray-600">
                                                {user.email}
                                            </TableCell>
                                            <TableCell>
                                                <div className="flex gap-1 flex-wrap">
                                                    {user.roles && user.roles.length > 0 ? (
                                                        user.roles.map((role, idx) => (
                                                            <Badge key={idx} variant="secondary">
                                                                {role}
                                                            </Badge>
                                                        ))
                                                    ) : (
                                                        <span className="text-sm text-gray-400">No roles</span>
                                                    )}
                                                </div>
                                            </TableCell>
                                            <TableCell>
                                                {user.is_active ? (
                                                    <Badge className="bg-green-100 text-green-700 hover:bg-green-100">
                                                        <UserCheck className="h-3 w-3 mr-1" />
                                                        Active
                                                    </Badge>
                                                ) : (
                                                    <Badge variant="destructive">
                                                        <UserX className="h-3 w-3 mr-1" />
                                                        Inactive
                                                    </Badge>
                                                )}
                                            </TableCell>
                                            <TableCell className="text-sm text-gray-600">
                                                {format(new Date(user.created_at), 'MMM d, yyyy')}
                                            </TableCell>
                                            <TableCell className="text-right">
                                                <DropdownMenu>
                                                    <DropdownMenuTrigger asChild>
                                                        <Button variant="ghost" size="sm">
                                                            <MoreVertical className="h-4 w-4" />
                                                        </Button>
                                                    </DropdownMenuTrigger>
                                                    <DropdownMenuContent align="end">
                                                        <DropdownMenuLabel>Actions</DropdownMenuLabel>
                                                        <DropdownMenuItem onClick={() => handleEditUser(user)}>
                                                            Edit User
                                                        </DropdownMenuItem>
                                                        <DropdownMenuItem onClick={() => handleManageRoles(user)}>
                                                            Manage Roles
                                                        </DropdownMenuItem>
                                                        <DropdownMenuSeparator />
                                                        <DropdownMenuItem
                                                            onClick={() => deleteUserMutation.mutate(user.id)}
                                                            className="text-red-600"
                                                        >
                                                            Delete User
                                                        </DropdownMenuItem>
                                                    </DropdownMenuContent>
                                                </DropdownMenu>
                                            </TableCell>
                                        </TableRow>
                                    ))
                                )}
                            </TableBody>
                        </Table>
                    )}

                    {/* Pagination */}
                    {total > 20 && (
                        <div className="flex items-center justify-between mt-4">
                            <div className="text-sm text-gray-500">
                                Showing {(page - 1) * 20 + 1} to {Math.min(page * 20, total)} of {total} users
                            </div>
                            <div className="flex gap-2">
                                <Button
                                    variant="outline"
                                    size="sm"
                                    onClick={() => setPage(p => Math.max(1, p - 1))}
                                    disabled={page === 1}
                                >
                                    Previous
                                </Button>
                                <Button
                                    variant="outline"
                                    size="sm"
                                    onClick={() => setPage(p => p + 1)}
                                    disabled={page * 20 >= total}
                                >
                                    Next
                                </Button>
                            </div>
                        </div>
                    )}
                </CardContent>
            </Card>

            {/* Edit User Dialog */}
            <Dialog open={isEditOpen} onOpenChange={setIsEditOpen}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>Edit User</DialogTitle>
                        <DialogDescription>
                            Update user information and permissions
                        </DialogDescription>
                    </DialogHeader>
                    {selectedUser && (
                        <div className="space-y-4 py-4">
                            <div className="space-y-2">
                                <Label htmlFor="name">Name</Label>
                                <Input
                                    id="name"
                                    value={selectedUser.name}
                                    onChange={(e) =>
                                        setSelectedUser({ ...selectedUser, name: e.target.value })
                                    }
                                />
                            </div>
                            <div className="space-y-2">
                                <Label htmlFor="email">Email</Label>
                                <Input
                                    id="email"
                                    value={selectedUser.email}
                                    disabled
                                    className="bg-gray-50"
                                />
                            </div>
                            <div className="flex items-center justify-between space-x-2">
                                <Label htmlFor="is_active" className="flex flex-col space-y-1">
                                    <span>Active Status</span>
                                    <span className="text-sm font-normal text-gray-500">
                                        Enable or disable user account
                                    </span>
                                </Label>
                                <Switch
                                    id="is_active"
                                    checked={selectedUser.is_active}
                                    onCheckedChange={(checked) =>
                                        setSelectedUser({ ...selectedUser, is_active: checked })
                                    }
                                />
                            </div>
                            <div className="flex items-center justify-between space-x-2">
                                <Label htmlFor="is_admin" className="flex flex-col space-y-1">
                                    <span>Administrator</span>
                                    <span className="text-sm font-normal text-gray-500">
                                        Grant full administrative privileges
                                    </span>
                                </Label>
                                <Switch
                                    id="is_admin"
                                    checked={selectedUser.is_admin}
                                    onCheckedChange={(checked) =>
                                        setSelectedUser({ ...selectedUser, is_admin: checked })
                                    }
                                />
                            </div>
                        </div>
                    )}
                    <DialogFooter>
                        <Button variant="outline" onClick={() => setIsEditOpen(false)}>
                            Cancel
                        </Button>
                        <Button onClick={handleSaveUser} disabled={updateUserMutation.isPending}>
                            {updateUserMutation.isPending ? (
                                <>
                                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                    Saving...
                                </>
                            ) : (
                                'Save Changes'
                            )}
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>

            {/* Manage Roles Dialog */}
            <Dialog open={isRoleDialogOpen} onOpenChange={setIsRoleDialogOpen}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>Manage Roles</DialogTitle>
                        <DialogDescription>
                            Assign or remove roles for {selectedUser?.name}
                        </DialogDescription>
                    </DialogHeader>
                    {selectedUser && (
                        <div className="space-y-4 py-4">
                            <div className="space-y-2">
                                <Label>Current Roles</Label>
                                <div className="flex flex-wrap gap-2">
                                    {selectedUser.roles && selectedUser.roles.length > 0 ? (
                                        selectedUser.roles.map((roleName, idx) => {
                                            const role = roles.find(r => r.name === roleName);
                                            return (
                                                <Badge key={idx} variant="secondary" className="text-sm">
                                                    {roleName}
                                                    {role && (
                                                        <button
                                                            onClick={() => handleRemoveRole(role.id)}
                                                            className="ml-2 hover:text-red-600"
                                                        >
                                                            ×
                                                        </button>
                                                    )}
                                                </Badge>
                                            );
                                        })
                                    ) : (
                                        <span className="text-sm text-gray-400">No roles assigned</span>
                                    )}
                                </div>
                            </div>
                            <div className="space-y-2">
                                <Label htmlFor="assign-role">Assign New Role</Label>
                                <div className="flex gap-2">
                                    <Select value={selectedRoleId} onValueChange={setSelectedRoleId}>
                                        <SelectTrigger className="flex-1">
                                            <SelectValue placeholder="Select a role" />
                                        </SelectTrigger>
                                        <SelectContent>
                                            {roles
                                                .filter(r => !selectedUser.roles?.includes(r.name))
                                                .map((role) => (
                                                    <SelectItem key={role.id} value={role.id.toString()}>
                                                        {role.name}
                                                        {role.description && (
                                                            <span className="text-gray-500 text-xs ml-2">
                                                                - {role.description}
                                                            </span>
                                                        )}
                                                    </SelectItem>
                                                ))}
                                        </SelectContent>
                                    </Select>
                                    <Button
                                        onClick={handleAssignRole}
                                        disabled={!selectedRoleId || assignRoleMutation.isPending}
                                    >
                                        {assignRoleMutation.isPending ? (
                                            <Loader2 className="h-4 w-4 animate-spin" />
                                        ) : (
                                            'Assign'
                                        )}
                                    </Button>
                                </div>
                            </div>
                        </div>
                    )}
                    <DialogFooter>
                        <Button variant="outline" onClick={() => {
                            setIsRoleDialogOpen(false);
                            setSelectedRoleId('');
                        }}>
                            Close
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </div>
    );
}