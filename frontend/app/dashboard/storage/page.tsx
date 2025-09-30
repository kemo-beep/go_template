'use client';

import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from '@/components/ui/dialog';
import {
    Upload,
    Download,
    Trash2,
    MoreVertical,
    File,
    Folder,
    Link as LinkIcon,
    Search,
    Image as ImageIcon,
} from 'lucide-react';
import { api } from '@/lib/api-client';
import { toast } from 'sonner';

export default function StoragePage() {
    const queryClient = useQueryClient();
    const [searchTerm, setSearchTerm] = useState('');
    const [uploadDialogOpen, setUploadDialogOpen] = useState(false);
    const [selectedFile, setSelectedFile] = useState<File | null>(null);
    const [uploadPath, setUploadPath] = useState('/');

    const { data: files, isLoading } = useQuery({
        queryKey: ['storage-files'],
        queryFn: () => api.getFiles(),
    });

    const uploadMutation = useMutation({
        mutationFn: (data: { file: File; path: string }) => api.uploadFile(data),
        onSuccess: () => {
            toast.success('File uploaded successfully');
            queryClient.invalidateQueries({ queryKey: ['storage-files'] });
            setUploadDialogOpen(false);
            setSelectedFile(null);
        },
        onError: (error: any) => {
            toast.error(error.response?.data?.error || 'Upload failed');
        },
    });

    const deleteMutation = useMutation({
        mutationFn: (fileId: string) => api.deleteFile(fileId),
        onSuccess: () => {
            toast.success('File deleted successfully');
            queryClient.invalidateQueries({ queryKey: ['storage-files'] });
        },
        onError: (error: any) => {
            toast.error(error.response?.data?.error || 'Delete failed');
        },
    });

    const generateSignedUrlMutation = useMutation({
        mutationFn: (fileId: string) => api.generateSignedUrl(fileId),
        onSuccess: (data) => {
            navigator.clipboard.writeText(data.data.data.url);
            toast.success('Signed URL copied to clipboard');
        },
        onError: (error: any) => {
            toast.error(error.response?.data?.error || 'Failed to generate URL');
        },
    });

    const handleUpload = () => {
        if (!selectedFile) {
            toast.error('Please select a file');
            return;
        }
        uploadMutation.mutate({ file: selectedFile, path: uploadPath });
    };

    const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.files && e.target.files[0]) {
            setSelectedFile(e.target.files[0]);
        }
    };

    // The API response is wrapped: { data: { success, message, data: [...] } }
    // So we need to access files?.data?.data to get the actual array
    const filesList = Array.isArray(files?.data?.data) ? files.data.data :
        Array.isArray(files?.data) ? files.data : [];

    const filteredFiles = filesList.filter((file: any) =>
        file.file_name.toLowerCase().includes(searchTerm.toLowerCase())
    );

    const formatFileSize = (bytes: number) => {
        if (bytes === 0) return '0 B';
        const k = 1024;
        const sizes = ['B', 'KB', 'MB', 'GB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
    };

    const getFileIcon = (type: string) => {
        if (type.startsWith('image/')) return <ImageIcon className="h-4 w-4" />;
        if (type === 'folder') return <Folder className="h-4 w-4" />;
        return <File className="h-4 w-4" />;
    };

    return (
        <div className="space-y-6">
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Storage Explorer</h1>
                    <p className="text-gray-500 mt-1">
                        Browse, upload, and manage files in Cloudflare R2
                    </p>
                </div>

                <Dialog open={uploadDialogOpen} onOpenChange={setUploadDialogOpen}>
                    <DialogTrigger asChild>
                        <Button>
                            <Upload className="h-4 w-4 mr-2" />
                            Upload File
                        </Button>
                    </DialogTrigger>
                    <DialogContent>
                        <DialogHeader>
                            <DialogTitle>Upload File</DialogTitle>
                        </DialogHeader>
                        <div className="space-y-4 py-4">
                            <div>
                                <label className="text-sm font-medium mb-2 block">
                                    Upload Path
                                </label>
                                <Input
                                    value={uploadPath}
                                    onChange={(e) => setUploadPath(e.target.value)}
                                    placeholder="/path/to/folder/"
                                />
                            </div>
                            <div>
                                <label className="text-sm font-medium mb-2 block">
                                    Select File
                                </label>
                                <Input
                                    type="file"
                                    onChange={handleFileSelect}
                                />
                                {selectedFile && (
                                    <p className="text-sm text-gray-500 mt-2">
                                        Selected: {selectedFile.name} ({formatFileSize(selectedFile.size)})
                                    </p>
                                )}
                            </div>
                            <div className="flex gap-2">
                                <Button
                                    onClick={handleUpload}
                                    disabled={uploadMutation.isPending || !selectedFile}
                                    className="flex-1"
                                >
                                    {uploadMutation.isPending ? 'Uploading...' : 'Upload'}
                                </Button>
                                <Button
                                    variant="outline"
                                    onClick={() => setUploadDialogOpen(false)}
                                    className="flex-1"
                                >
                                    Cancel
                                </Button>
                            </div>
                        </div>
                    </DialogContent>
                </Dialog>
            </div>

            {/* Stats */}
            <div className="grid md:grid-cols-4 gap-4">
                <Card>
                    <CardHeader className="pb-3">
                        <CardTitle className="text-sm font-medium text-gray-500">
                            Total Files
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">{filesList.length || 0}</div>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="pb-3">
                        <CardTitle className="text-sm font-medium text-gray-500">
                            Total Size
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">
                            {formatFileSize(
                                filesList.reduce((acc: number, file: any) => acc + (file.file_size || 0), 0)
                            )}
                        </div>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="pb-3">
                        <CardTitle className="text-sm font-medium text-gray-500">
                            Images
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">
                            {filesList.filter((f: any) => f.file_type?.startsWith('image/')).length || 0}
                        </div>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader className="pb-3">
                        <CardTitle className="text-sm font-medium text-gray-500">
                            Storage Used
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">42%</div>
                        <p className="text-xs text-gray-500 mt-1">of 10 GB</p>
                    </CardContent>
                </Card>
            </div>

            {/* Files Table */}
            <Card>
                <CardHeader>
                    <div className="flex items-center justify-between">
                        <CardTitle>Files</CardTitle>
                        <div className="relative w-64">
                            <Search className="absolute left-2 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
                            <Input
                                placeholder="Search files..."
                                value={searchTerm}
                                onChange={(e) => setSearchTerm(e.target.value)}
                                className="pl-8"
                            />
                        </div>
                    </div>
                </CardHeader>
                <CardContent>
                    {isLoading ? (
                        <p className="text-center py-8">Loading files...</p>
                    ) : filteredFiles?.length === 0 ? (
                        <div className="text-center py-12 text-gray-500">
                            <Upload className="h-12 w-12 mx-auto mb-4 opacity-50" />
                            <p>No files found</p>
                            <p className="text-sm mt-1">Upload your first file to get started</p>
                        </div>
                    ) : (
                        <Table>
                            <TableHeader>
                                <TableRow>
                                    <TableHead>Name</TableHead>
                                    <TableHead>Type</TableHead>
                                    <TableHead>Size</TableHead>
                                    <TableHead>Uploaded</TableHead>
                                    <TableHead>Public</TableHead>
                                    <TableHead className="w-[100px]">Actions</TableHead>
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                {filteredFiles?.map((file: any) => (
                                    <TableRow key={file.id}>
                                        <TableCell>
                                            <div className="flex items-center gap-2">
                                                {getFileIcon(file.file_type)}
                                                <span className="font-medium">{file.file_name}</span>
                                            </div>
                                        </TableCell>
                                        <TableCell>
                                            <code className="text-xs bg-gray-100 px-2 py-1 rounded">
                                                {file.file_type}
                                            </code>
                                        </TableCell>
                                        <TableCell>{formatFileSize(file.file_size)}</TableCell>
                                        <TableCell className="text-sm text-gray-500">
                                            {new Date(file.created_at).toLocaleDateString()}
                                        </TableCell>
                                        <TableCell>
                                            <Badge variant={file.is_public ? 'default' : 'secondary'}>
                                                {file.is_public ? 'Public' : 'Private'}
                                            </Badge>
                                        </TableCell>
                                        <TableCell>
                                            <DropdownMenu>
                                                <DropdownMenuTrigger asChild>
                                                    <Button variant="ghost" size="sm">
                                                        <MoreVertical className="h-4 w-4" />
                                                    </Button>
                                                </DropdownMenuTrigger>
                                                <DropdownMenuContent align="end">
                                                    <DropdownMenuItem
                                                        onClick={() => window.open(file.r2_url, '_blank')}
                                                    >
                                                        <Download className="h-4 w-4 mr-2" />
                                                        Download
                                                    </DropdownMenuItem>
                                                    <DropdownMenuItem
                                                        onClick={() => generateSignedUrlMutation.mutate(file.id)}
                                                    >
                                                        <LinkIcon className="h-4 w-4 mr-2" />
                                                        Copy Signed URL
                                                    </DropdownMenuItem>
                                                    <DropdownMenuItem
                                                        onClick={() => {
                                                            if (confirm('Are you sure you want to delete this file?')) {
                                                                deleteMutation.mutate(file.id);
                                                            }
                                                        }}
                                                        className="text-red-600"
                                                    >
                                                        <Trash2 className="h-4 w-4 mr-2" />
                                                        Delete
                                                    </DropdownMenuItem>
                                                </DropdownMenuContent>
                                            </DropdownMenu>
                                        </TableCell>
                                    </TableRow>
                                ))}
                            </TableBody>
                        </Table>
                    )}
                </CardContent>
            </Card>
        </div>
    );
}
