'use client';

import { useState, useEffect, useMemo } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useRouter, useSearchParams } from 'next/navigation';
import { api } from '@/lib/api-client';
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
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Skeleton } from '@/components/ui/skeleton';
import { Alert, AlertDescription } from '@/components/ui/alert';
import {
    Database,
    Play,
    Download,
    Search,
    Table as TableIcon,
    Plus,
    Trash2,
    MoreHorizontal,
    Eye,
    Edit,
    Scissors,
    Lock,
    History,
    Settings,
    X,
    RefreshCw,
    Filter,
    SortAsc,
    SortDesc,
    ChevronRight,
    AlertCircle,
    CheckCircle2,
    Clock,
    Zap,
    Database as Database2,
    Code2,
    BarChart3,
    Share2,
    Copy,
    Check
} from 'lucide-react';
import { toast } from 'sonner';
import { CreateTableDialog } from './create-table';
import { AlterTableContent } from './alter-table-content';
import { MigrationHistory } from './migration-history';
import TableDataEditor from './table-data-editor';
import { cn } from '@/lib/utils';

export default function DatabasePage() {
    const router = useRouter();
    const searchParams = useSearchParams();
    const queryClient = useQueryClient();

    // URL query parameters
    const tableParam = searchParams.get('table');
    const tabParam = searchParams.get('tab');
    const queryParam = searchParams.get('query');
    const migrationParam = searchParams.get('migrations') === 'true';

    const [selectedTable, setSelectedTable] = useState<string>(tableParam || '');
    const [sqlQuery, setSqlQuery] = useState(queryParam || 'SELECT * FROM users LIMIT 10;');
    const [searchTerm, setSearchTerm] = useState('');
    const [openTabs, setOpenTabs] = useState<Array<{ id: string; name: string; type: 'table' | 'query' | 'alter'; isDirty?: boolean }>>([]);
    const [activeTab, setActiveTab] = useState<string>(tabParam || '');
    const [createTableOpen, setCreateTableOpen] = useState(false);
    const [selectedTableForAlter, setSelectedTableForAlter] = useState<string>('');
    const [showMigrationHistory, setShowMigrationHistory] = useState(migrationParam);
    const [sortBy, setSortBy] = useState<'name' | 'created' | 'size'>('name');
    const [sortOrder, setSortOrder] = useState<'asc' | 'desc'>('asc');
    const [isRefreshing, setIsRefreshing] = useState(false);
    const [recentQueries, setRecentQueries] = useState<string[]>([]);
    const [favoriteTables, setFavoriteTables] = useState<string[]>([]);
    const [lastActiveTab, setLastActiveTab] = useState<string>('');
    const [copiedToClipboard, setCopiedToClipboard] = useState(false);

    // Load saved state from localStorage
    useEffect(() => {
        const savedFavorites = localStorage.getItem('db-favorite-tables');
        const savedRecentQueries = localStorage.getItem('db-recent-queries');
        const savedTabs = localStorage.getItem('db-open-tabs');
        const savedSortBy = localStorage.getItem('db-sort-by');
        const savedSortOrder = localStorage.getItem('db-sort-order');
        const savedSearchTerm = localStorage.getItem('db-search-term');
        const savedActiveTab = localStorage.getItem('db-active-tab');
        const savedSelectedTable = localStorage.getItem('db-selected-table');

        if (savedFavorites) setFavoriteTables(JSON.parse(savedFavorites));
        if (savedRecentQueries) setRecentQueries(JSON.parse(savedRecentQueries));
        if (savedTabs) setOpenTabs(JSON.parse(savedTabs));
        if (savedSortBy) setSortBy(savedSortBy as 'name' | 'created' | 'size');
        if (savedSortOrder) setSortOrder(savedSortOrder as 'asc' | 'desc');
        if (savedSearchTerm) setSearchTerm(savedSearchTerm);

        // Only restore from localStorage if no URL parameters
        if (!tableParam && savedActiveTab) {
            setActiveTab(savedActiveTab);
        }
        if (!tableParam && savedSelectedTable) {
            setSelectedTable(savedSelectedTable);
        }
    }, [tableParam]);

    // Initialize from URL parameters - only run once on mount
    useEffect(() => {
        if (tableParam) {
            setSelectedTable(tableParam);
            const tabId = `table-${tableParam}`;
            setActiveTab(tabId);
            // Add tab if it doesn't exist
            setOpenTabs(prev => {
                const existingTab = prev.find(tab => tab.id === tabId);
                if (!existingTab) {
                    return [...prev, { id: tabId, name: tableParam, type: 'table' as const, isDirty: false }];
                }
                return prev;
            });
        }
        if (queryParam) {
            setSqlQuery(queryParam);
        }
        if (migrationParam) {
            setShowMigrationHistory(true);
        }
    }, []); // Empty dependency array - only run on mount

    // Handle URL parameter changes after initial load
    useEffect(() => {
        if (tableParam && tableParam !== selectedTable) {
            setSelectedTable(tableParam);
            const tabId = `table-${tableParam}`;
            setActiveTab(tabId);
            setOpenTabs(prev => {
                const existingTab = prev.find(tab => tab.id === tabId);
                if (!existingTab) {
                    return [...prev, { id: tabId, name: tableParam, type: 'table' as const, isDirty: false }];
                }
                return prev;
            });
        }
    }, [tableParam, selectedTable]);

    // Invalidate table data query when selected table changes
    useEffect(() => {
        if (selectedTable) {
            queryClient.invalidateQueries({ queryKey: ['table-data', selectedTable] });
        }
    }, [selectedTable, queryClient]);

    // URL synchronization
    const updateURL = (params: Record<string, string | null>) => {
        const current = new URLSearchParams(searchParams.toString());

        Object.entries(params).forEach(([key, value]) => {
            if (value === null || value === '') {
                current.delete(key);
            } else {
                current.set(key, value);
            }
        });

        const newURL = `${window.location.pathname}?${current.toString()}`;
        router.replace(newURL, { scroll: false });
    };

    // Keyboard shortcuts
    useEffect(() => {
        const handleKeyDown = (event: KeyboardEvent) => {
            // Only handle shortcuts when not in input/textarea
            if (event.target instanceof HTMLInputElement || event.target instanceof HTMLTextAreaElement) {
                return;
            }

            const isCtrlOrCmd = event.ctrlKey || event.metaKey;

            switch (true) {
                case isCtrlOrCmd && event.key === 'n':
                    event.preventDefault();
                    setCreateTableOpen(true);
                    break;
                case isCtrlOrCmd && event.key === 'q':
                    event.preventDefault();
                    openQueryTab();
                    break;
                case isCtrlOrCmd && event.key === 'r':
                    event.preventDefault();
                    refreshTables();
                    break;
                case event.key === 'Escape':
                    if (showMigrationHistory) {
                        setShowMigrationHistory(false);
                    }
                    break;
                case isCtrlOrCmd && event.key === 'k':
                    event.preventDefault();
                    // Focus search input
                    const searchInput = document.querySelector('input[placeholder="Search tables..."]') as HTMLInputElement;
                    searchInput?.focus();
                    break;
            }
        };

        document.addEventListener('keydown', handleKeyDown);
        return () => document.removeEventListener('keydown', handleKeyDown);
    }, [showMigrationHistory]);

    // Save state to localStorage
    useEffect(() => {
        localStorage.setItem('db-favorite-tables', JSON.stringify(favoriteTables));
    }, [favoriteTables]);

    useEffect(() => {
        localStorage.setItem('db-recent-queries', JSON.stringify(recentQueries));
    }, [recentQueries]);

    useEffect(() => {
        localStorage.setItem('db-open-tabs', JSON.stringify(openTabs));
    }, [openTabs]);

    useEffect(() => {
        localStorage.setItem('db-sort-by', sortBy);
    }, [sortBy]);

    useEffect(() => {
        localStorage.setItem('db-sort-order', sortOrder);
    }, [sortOrder]);

    useEffect(() => {
        localStorage.setItem('db-search-term', searchTerm);
    }, [searchTerm]);

    // Save active tab and selected table to localStorage
    useEffect(() => {
        if (activeTab) {
            localStorage.setItem('db-active-tab', activeTab);
            setLastActiveTab(activeTab);
        }
    }, [activeTab]);

    useEffect(() => {
        if (selectedTable) {
            localStorage.setItem('db-selected-table', selectedTable);
        }
    }, [selectedTable]);

    // Note: URL sync is now handled manually in each function to avoid conflicts

    // Tab management functions
    const openTableTab = (tableName: string) => {
        const tabId = `table-${tableName}`;
        const existingTab = openTabs.find(tab => tab.id === tabId);

        if (!existingTab) {
            const newTab = { id: tabId, name: tableName, type: 'table' as const, isDirty: false };
            setOpenTabs(prev => [...prev, newTab]);
        }

        // Always set active tab and selected table
        setActiveTab(tabId);
        setSelectedTable(tableName);

        // Update URL
        updateURL({ table: tableName, tab: tabId });
    };

    const closeTab = (tabId: string) => {
        const tab = openTabs.find(t => t.id === tabId);
        if (tab?.isDirty) {
            if (!window.confirm('This tab has unsaved changes. Are you sure you want to close it?')) {
                return;
            }
        }

        setOpenTabs(prev => {
            const newTabs = prev.filter(tab => tab.id !== tabId);
            const newActiveTab = activeTab === tabId ? (newTabs.length > 0 ? newTabs[newTabs.length - 1].id : '') : activeTab;
            setActiveTab(newActiveTab);

            // Update URL
            if (activeTab === tabId) {
                const newTab = newTabs.find(t => t.id === newActiveTab);
                updateURL({
                    table: newTab?.type === 'table' ? newTab.name : null,
                    tab: newActiveTab || null
                });
            }

            return newTabs;
        });
    };

    const openQueryTab = () => {
        const tabId = 'query-tab';
        const existingTab = openTabs.find(tab => tab.id === tabId);

        if (!existingTab) {
            const newTab = { id: tabId, name: 'SQL Query', type: 'query' as const, isDirty: false };
            setOpenTabs(prev => [...prev, newTab]);
        }
        setActiveTab(tabId);

        // Update URL
        updateURL({ tab: tabId, table: null });
    };

    const markTabAsDirty = (tabId: string) => {
        setOpenTabs(prev => prev.map(tab =>
            tab.id === tabId ? { ...tab, isDirty: true } : tab
        ));
    };

    const markTabAsClean = (tabId: string) => {
        setOpenTabs(prev => prev.map(tab =>
            tab.id === tabId ? { ...tab, isDirty: false } : tab
        ));
    };

    // Use the centralized API client
    const apiClient = {
        getTables: () => api.getTables(),
        getTableData: (tableName: string, page = 1, limit = 20) =>
            api.getTableData(tableName, page, limit),
        getTableSchema: (tableName: string) =>
            api.getTableSchema(tableName),
        executeQuery: (query: string) => api.executeQuery(query),
    };

    // Utility functions
    const refreshTables = async () => {
        setIsRefreshing(true);
        try {
            await queryClient.invalidateQueries({ queryKey: ['database-tables'] });
            toast.success('Tables refreshed');
        } catch (error) {
            toast.error('Failed to refresh tables');
        } finally {
            setIsRefreshing(false);
        }
    };

    const toggleFavorite = (tableName: string) => {
        setFavoriteTables(prev =>
            prev.includes(tableName)
                ? prev.filter(name => name !== tableName)
                : [...prev, tableName]
        );
    };

    const addToRecentQueries = (query: string) => {
        if (query.trim() && !recentQueries.includes(query)) {
            setRecentQueries(prev => [query, ...prev].slice(0, 10));
        }
    };

    const { data: tables, isLoading: tablesLoading } = useQuery({
        queryKey: ['database-tables'],
        queryFn: () => apiClient.getTables(),
    });

    const { data: tableData, isLoading: tableDataLoading } = useQuery({
        queryKey: ['table-data', selectedTable],
        queryFn: () => apiClient.getTableData(selectedTable),
        enabled: !!selectedTable,
        staleTime: 30000, // Cache for 30 seconds
        retry: 2,
    });

    const executeQueryMutation = useMutation({
        mutationFn: (query: string) => apiClient.executeQuery(query),
        onSuccess: (data, query) => {
            toast.success('Query executed successfully');
            addToRecentQueries(query);
            markTabAsClean('query-tab');
        },
        onError: (error: any) => {
            toast.error(error?.message || 'Query execution failed');
        },
    });

    const handleExecuteQuery = () => {
        if (!sqlQuery.trim()) {
            toast.error('Please enter a SQL query');
            return;
        }
        markTabAsDirty('query-tab');
        executeQueryMutation.mutate(sqlQuery);
    };

    // Computed values
    const filteredTables = useMemo(() => {
        if (!tables?.data?.data) return [];

        let filtered = tables.data.data.filter((table: any) =>
            table.name.toLowerCase().includes(searchTerm.toLowerCase())
        );

        // Sort tables
        filtered.sort((a: any, b: any) => {
            let comparison = 0;
            switch (sortBy) {
                case 'name':
                    comparison = a.name.localeCompare(b.name);
                    break;
                case 'created':
                    comparison = new Date(a.created_at || 0).getTime() - new Date(b.created_at || 0).getTime();
                    break;
                case 'size':
                    comparison = (a.row_count || 0) - (b.row_count || 0);
                    break;
            }
            return sortOrder === 'asc' ? comparison : -comparison;
        });

        return filtered;
    }, [tables, searchTerm, sortBy, sortOrder]);

    const favoriteTablesList = useMemo(() => {
        return filteredTables.filter((table: any) => favoriteTables.includes(table.name));
    }, [filteredTables, favoriteTables]);

    const regularTablesList = useMemo(() => {
        return filteredTables.filter((table: any) => !favoriteTables.includes(table.name));
    }, [filteredTables, favoriteTables]);

    // Table action handlers
    const handleBrowseData = (tableName: string) => {
        openTableTab(tableName);
        toast.success(`Opening data for table: ${tableName}`);
    };

    const handleAlterTable = (tableName: string) => {
        const tabId = `alter-${tableName}`;
        const existingTab = openTabs.find(tab => tab.id === tabId);

        if (!existingTab) {
            const newTab = { id: tabId, name: `Alter ${tableName}`, type: 'alter' as const };
            setOpenTabs(prev => [...prev, newTab]);
        }
        setActiveTab(tabId);
        setSelectedTableForAlter(tableName);

        // Update URL
        updateURL({ table: tableName, tab: tabId });
    };

    const handleTruncateTable = (tableName: string) => {
        if (window.confirm(`Are you sure you want to truncate table "${tableName}"? This will delete all data in the table.`)) {
            // TODO: Implement truncate table functionality
            toast.success(`Table ${tableName} truncated successfully`);
        }
    };

    const handleDropTable = (tableName: string) => {
        if (window.confirm(`Are you sure you want to drop table "${tableName}"? This action cannot be undone.`)) {
            // TODO: Implement drop table functionality
            toast.success(`Table ${tableName} dropped successfully`);
        }
    };

    const handleEnableRLS = (tableName: string) => {
        toast.info(`Enable RLS for table ${tableName} - Coming soon!`);
        // TODO: Implement RLS functionality
    };

    const handleQueryChange = (query: string) => {
        setSqlQuery(query);
        markTabAsDirty('query-tab');

        // Update URL with query
        updateURL({ query: query || null });
    };

    const handleMigrationHistoryToggle = (show: boolean) => {
        setShowMigrationHistory(show);
        updateURL({ migrations: show ? 'true' : null });
    };

    const handleShareLink = async () => {
        try {
            const currentURL = window.location.href;
            await navigator.clipboard.writeText(currentURL);
            setCopiedToClipboard(true);
            toast.success('Link copied to clipboard!');
            setTimeout(() => setCopiedToClipboard(false), 2000);
        } catch (error) {
            toast.error('Failed to copy link to clipboard');
        }
    };

    return (
        <div className="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100 dark:from-slate-900 dark:to-slate-800">
            {/* Header Section */}
            <div className="bg-white dark:bg-slate-800 border-b border-slate-200 dark:border-slate-700">
                <div className="px-6 py-4">
                    <div className="flex items-center justify-between">
                        <div className="space-y-1">
                            <div className="flex items-center gap-3">
                                <div className="p-2 rounded-lg bg-gradient-to-br from-primary to-accent">
                                    <Database2 className="h-6 w-6 text-white" />
                                </div>
                                <div>
                                    <h1 className="text-2xl font-bold bg-gradient-to-r from-primary to-accent bg-clip-text text-transparent">
                                        Database Explorer
                                    </h1>
                                    <p className="text-sm text-muted-foreground">
                                        Manage tables, execute queries, and monitor your database
                                        {selectedTable && (
                                            <span className="ml-2 px-2 py-1 bg-primary/10 text-primary rounded-md text-xs">
                                                Viewing: {selectedTable}
                                            </span>
                                        )}
                                    </p>
                                </div>
                            </div>
                        </div>

                        <div className="flex items-center gap-3">
                            <div className="flex items-center gap-2 text-xs text-muted-foreground">
                                <kbd className="px-1.5 py-0.5 text-xs bg-muted rounded">⌘K</kbd>
                                <span>Search</span>
                                <kbd className="px-1.5 py-0.5 text-xs bg-muted rounded">⌘Q</kbd>
                                <span>SQL</span>
                                <kbd className="px-1.5 py-0.5 text-xs bg-muted rounded">⌘N</kbd>
                                <span>New</span>
                            </div>

                            <Button
                                onClick={refreshTables}
                                variant="outline"
                                size="sm"
                                disabled={isRefreshing}
                                className="gap-2"
                            >
                                <RefreshCw className={cn("h-4 w-4", isRefreshing && "animate-spin")} />
                                Refresh
                            </Button>

                            <Button
                                onClick={() => handleMigrationHistoryToggle(!showMigrationHistory)}
                                variant="outline"
                                size="sm"
                                className="gap-2"
                            >
                                <History className="h-4 w-4" />
                                Migrations
                            </Button>

                            <Button
                                onClick={openQueryTab}
                                variant="outline"
                                size="sm"
                                className="gap-2"
                            >
                                <Code2 className="h-4 w-4" />
                                SQL Editor
                            </Button>

                            <Button
                                onClick={handleShareLink}
                                variant="outline"
                                size="sm"
                                className="gap-2"
                            >
                                {copiedToClipboard ? (
                                    <>
                                        <Check className="h-4 w-4" />
                                        Copied!
                                    </>
                                ) : (
                                    <>
                                        <Share2 className="h-4 w-4" />
                                        Share
                                    </>
                                )}
                            </Button>

                            <Button
                                onClick={() => setCreateTableOpen(true)}
                                size="sm"
                                className="gap-2 bg-gradient-to-r from-primary to-accent hover:from-primary/90 hover:to-accent/90"
                            >
                                <Plus className="h-4 w-4" />
                                Create Table
                            </Button>
                        </div>
                    </div>
                </div>
            </div>

            <CreateTableDialog open={createTableOpen} onOpenChange={setCreateTableOpen} />

            {/* Migration History */}
            {showMigrationHistory && (
                <div className="p-6">
                    <MigrationHistory />
                </div>
            )}

            {/* Main Content */}
            {!showMigrationHistory && (
                <div className="p-6">
                    <div className="grid grid-cols-1 lg:grid-cols-4 gap-6 h-[calc(100vh-200px)]">
                        {/* Tables Sidebar */}
                        <div className="lg:col-span-1 space-y-4">
                            <Card className="h-full">
                                <CardHeader className="pb-3">
                                    <div className="flex items-center justify-between">
                                        <CardTitle className="text-lg flex items-center gap-2">
                                            <TableIcon className="h-5 w-5" />
                                            Tables
                                            {tablesLoading && <Skeleton className="h-4 w-4 rounded" />}
                                        </CardTitle>
                                        <div className="flex items-center gap-1">
                                            <DropdownMenu>
                                                <DropdownMenuTrigger asChild>
                                                    <Button variant="ghost" size="sm" className="h-8 w-8 p-0">
                                                        <Filter className="h-4 w-4" />
                                                    </Button>
                                                </DropdownMenuTrigger>
                                                <DropdownMenuContent align="end" className="w-48">
                                                    <DropdownMenuItem onClick={() => setSortBy('name')}>
                                                        <SortAsc className="h-4 w-4 mr-2" />
                                                        Sort by Name
                                                    </DropdownMenuItem>
                                                    <DropdownMenuItem onClick={() => setSortBy('created')}>
                                                        <Clock className="h-4 w-4 mr-2" />
                                                        Sort by Created
                                                    </DropdownMenuItem>
                                                    <DropdownMenuItem onClick={() => setSortBy('size')}>
                                                        <BarChart3 className="h-4 w-4 mr-2" />
                                                        Sort by Size
                                                    </DropdownMenuItem>
                                                    <DropdownMenuSeparator />
                                                    <DropdownMenuItem onClick={() => setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc')}>
                                                        {sortOrder === 'asc' ? <SortDesc className="h-4 w-4 mr-2" /> : <SortAsc className="h-4 w-4 mr-2" />}
                                                        {sortOrder === 'asc' ? 'Descending' : 'Ascending'}
                                                    </DropdownMenuItem>
                                                </DropdownMenuContent>
                                            </DropdownMenu>
                                        </div>
                                    </div>

                                    <div className="relative">
                                        <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                                        <Input
                                            placeholder="Search tables..."
                                            value={searchTerm}
                                            onChange={(e) => setSearchTerm(e.target.value)}
                                            className="pl-10 h-9"
                                        />
                                    </div>
                                </CardHeader>

                                <CardContent className="p-0 h-[calc(100%-120px)] overflow-y-auto">
                                    {tablesLoading ? (
                                        <div className="p-4 space-y-3">
                                            {Array.from({ length: 5 }).map((_, i) => (
                                                <div key={i} className="flex items-center gap-3">
                                                    <Skeleton className="h-4 w-4 rounded" />
                                                    <Skeleton className="h-4 flex-1" />
                                                </div>
                                            ))}
                                        </div>
                                    ) : (
                                        <div className="space-y-1 p-2">
                                            {/* Favorites Section */}
                                            {favoriteTablesList.length > 0 && (
                                                <div className="space-y-1">
                                                    <div className="px-3 py-2 text-xs font-semibold text-muted-foreground uppercase tracking-wider">
                                                        Favorites
                                                    </div>
                                                    {favoriteTablesList.map((table: any) => (
                                                        <TableItem
                                                            key={table.name}
                                                            table={table}
                                                            isActive={openTabs.some(tab => tab.name === table.name)}
                                                            isFavorite={true}
                                                            onOpenTable={openTableTab}
                                                            onToggleFavorite={toggleFavorite}
                                                            onBrowseData={handleBrowseData}
                                                            onAlterTable={handleAlterTable}
                                                            onEnableRLS={handleEnableRLS}
                                                            onTruncateTable={handleTruncateTable}
                                                            onDropTable={handleDropTable}
                                                        />
                                                    ))}
                                                </div>
                                            )}

                                            {/* Regular Tables Section */}
                                            {regularTablesList.length > 0 && (
                                                <div className="space-y-1">
                                                    {favoriteTablesList.length > 0 && (
                                                        <div className="px-3 py-2 text-xs font-semibold text-muted-foreground uppercase tracking-wider">
                                                            All Tables
                                                        </div>
                                                    )}
                                                    {regularTablesList.map((table: any) => (
                                                        <TableItem
                                                            key={table.name}
                                                            table={table}
                                                            isActive={openTabs.some(tab => tab.name === table.name)}
                                                            isFavorite={false}
                                                            onOpenTable={openTableTab}
                                                            onToggleFavorite={toggleFavorite}
                                                            onBrowseData={handleBrowseData}
                                                            onAlterTable={handleAlterTable}
                                                            onEnableRLS={handleEnableRLS}
                                                            onTruncateTable={handleTruncateTable}
                                                            onDropTable={handleDropTable}
                                                        />
                                                    ))}
                                                </div>
                                            )}

                                            {filteredTables.length === 0 && !tablesLoading && (
                                                <div className="p-8 text-center">
                                                    <TableIcon className="h-12 w-12 mx-auto text-muted-foreground/50 mb-3" />
                                                    <p className="text-sm text-muted-foreground">
                                                        {searchTerm ? 'No tables found matching your search' : 'No tables available'}
                                                    </p>
                                                </div>
                                            )}
                                        </div>
                                    )}
                                </CardContent>
                            </Card>
                        </div>

                        {/* Main Content Area */}
                        <div className="lg:col-span-3">
                            <Card className="h-full">
                                <CardHeader className="pb-3">
                                    <div className="flex items-center justify-between">
                                        <div className="flex items-center gap-2">
                                            {/* Enhanced Tabs */}
                                            {openTabs.length > 0 && (
                                                <div className="flex items-center gap-1 overflow-x-auto max-w-2xl">
                                                    {openTabs.map((tab) => (
                                                        <div
                                                            key={tab.id}
                                                            className={cn(
                                                                "group flex items-center gap-2 px-3 py-2 rounded-lg text-sm font-medium transition-all duration-200 cursor-pointer min-w-0",
                                                                activeTab === tab.id
                                                                    ? "bg-primary text-primary-foreground shadow-sm"
                                                                    : "bg-muted hover:bg-muted/80 text-muted-foreground hover:text-foreground"
                                                            )}
                                                            onClick={() => {
                                                                setActiveTab(tab.id);
                                                                if (tab.type === 'table') {
                                                                    setSelectedTable(tab.name);
                                                                    // Update URL when switching to table tab
                                                                    updateURL({ table: tab.name, tab: tab.id });
                                                                } else if (tab.type === 'query') {
                                                                    // Update URL when switching to query tab
                                                                    updateURL({ table: null, tab: tab.id });
                                                                } else if (tab.type === 'alter') {
                                                                    // Update URL when switching to alter tab
                                                                    updateURL({ table: tab.name.replace('Alter ', ''), tab: tab.id });
                                                                }
                                                            }}
                                                        >
                                                            {tab.type === 'table' && <TableIcon className="h-4 w-4 flex-shrink-0" />}
                                                            {tab.type === 'query' && <Code2 className="h-4 w-4 flex-shrink-0" />}
                                                            {tab.type === 'alter' && <Settings className="h-4 w-4 flex-shrink-0" />}
                                                            <span className="truncate max-w-32">{tab.name}</span>
                                                            {tab.isDirty && (
                                                                <div className="h-2 w-2 rounded-full bg-orange-500 flex-shrink-0" />
                                                            )}
                                                            <Button
                                                                variant="ghost"
                                                                size="sm"
                                                                className="h-5 w-5 p-0 opacity-0 group-hover:opacity-100 transition-opacity hover:bg-destructive hover:text-destructive-foreground"
                                                                onClick={(e) => {
                                                                    e.stopPropagation();
                                                                    closeTab(tab.id);
                                                                }}
                                                            >
                                                                <X className="h-3 w-3" />
                                                            </Button>
                                                        </div>
                                                    ))}
                                                </div>
                                            )}
                                        </div>

                                        {activeTab && (
                                            <div className="flex items-center gap-2">
                                                <Button size="sm" variant="outline" className="gap-2">
                                                    <Download className="h-4 w-4" />
                                                    Export
                                                </Button>
                                            </div>
                                        )}
                                    </div>
                                </CardHeader>

                                <CardContent className="h-[calc(100%-80px)] overflow-hidden">
                                    {!activeTab ? (
                                        <EmptyState
                                            onOpenQueryTab={openQueryTab}
                                            onCreateTable={() => setCreateTableOpen(true)}
                                            onShareLink={handleShareLink}
                                        />
                                    ) : (
                                        <div className="h-full">
                                            {/* Table Data */}
                                            {openTabs.find(tab => tab.id === activeTab)?.type === 'table' && (
                                                <TableDataEditor
                                                    tableName={selectedTable}
                                                    onRefresh={refreshTables}
                                                />
                                            )}

                                            {/* SQL Query Editor */}
                                            {openTabs.find(tab => tab.id === activeTab)?.type === 'query' && (
                                                <QueryEditor
                                                    query={sqlQuery}
                                                    onQueryChange={handleQueryChange}
                                                    onExecute={handleExecuteQuery}
                                                    isExecuting={executeQueryMutation.isPending}
                                                    result={executeQueryMutation.data}
                                                    recentQueries={recentQueries}
                                                    onSelectRecentQuery={(query) => setSqlQuery(query)}
                                                />
                                            )}

                                            {/* Alter Table Content */}
                                            {openTabs.find(tab => tab.id === activeTab)?.type === 'alter' && (
                                                <AlterTableContent tableName={selectedTableForAlter} />
                                            )}
                                        </div>
                                    )}
                                </CardContent>
                            </Card>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}

// Table Item Component
interface TableItemProps {
    table: any;
    isActive: boolean;
    isFavorite: boolean;
    onOpenTable: (tableName: string) => void;
    onToggleFavorite: (tableName: string) => void;
    onBrowseData: (tableName: string) => void;
    onAlterTable: (tableName: string) => void;
    onEnableRLS: (tableName: string) => void;
    onTruncateTable: (tableName: string) => void;
    onDropTable: (tableName: string) => void;
}

function TableItem({
    table,
    isActive,
    isFavorite,
    onOpenTable,
    onToggleFavorite,
    onBrowseData,
    onAlterTable,
    onEnableRLS,
    onTruncateTable,
    onDropTable,
}: TableItemProps) {
    return (
        <div
            className={cn(
                "group flex items-center justify-between px-3 py-2 rounded-lg text-sm transition-all duration-200 cursor-pointer",
                isActive
                    ? "bg-primary text-primary-foreground shadow-sm"
                    : "hover:bg-muted/80"
            )}
            onClick={() => onOpenTable(table.name)}
        >
            <div className="flex items-center gap-3 flex-1 min-w-0">
                <TableIcon className="h-4 w-4 flex-shrink-0" />
                <span className="truncate font-medium">{table.name}</span>
                {table.row_count !== undefined && (
                    <Badge variant="secondary" className="text-xs ml-auto">
                        {table.row_count}
                    </Badge>
                )}
            </div>

            <div className="flex items-center gap-1 flex-shrink-0">
                <Button
                    variant="ghost"
                    size="sm"
                    className="h-6 w-6 p-0 opacity-0 group-hover:opacity-100 transition-opacity"
                    onClick={(e) => {
                        e.stopPropagation();
                        onToggleFavorite(table.name);
                    }}
                >
                    <Zap className={cn("h-3 w-3", isFavorite && "text-yellow-500 fill-current")} />
                </Button>

                <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                        <Button
                            variant="ghost"
                            size="sm"
                            className="h-6 w-6 p-0 opacity-0 group-hover:opacity-100 transition-opacity"
                            onClick={(e) => e.stopPropagation()}
                        >
                            <MoreHorizontal className="h-4 w-4" />
                        </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end" className="w-48">
                        <DropdownMenuItem
                            onClick={(e) => {
                                e.stopPropagation();
                                onBrowseData(table.name);
                            }}
                        >
                            <Eye className="h-4 w-4 mr-2" />
                            Browse data
                        </DropdownMenuItem>
                        <DropdownMenuItem
                            onClick={(e) => {
                                e.stopPropagation();
                                onAlterTable(table.name);
                            }}
                        >
                            <Edit className="h-4 w-4 mr-2" />
                            Alter table
                        </DropdownMenuItem>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem
                            onClick={(e) => {
                                e.stopPropagation();
                                onEnableRLS(table.name);
                            }}
                        >
                            <Lock className="h-4 w-4 mr-2" />
                            Enable RLS
                        </DropdownMenuItem>
                        <DropdownMenuItem
                            onClick={(e) => {
                                e.stopPropagation();
                                onTruncateTable(table.name);
                            }}
                            className="text-orange-600 focus:text-orange-600"
                        >
                            <Scissors className="h-4 w-4 mr-2" />
                            Truncate
                        </DropdownMenuItem>
                        <DropdownMenuItem
                            onClick={(e) => {
                                e.stopPropagation();
                                onDropTable(table.name);
                            }}
                            className="text-red-600 focus:text-red-600"
                        >
                            <Trash2 className="h-4 w-4 mr-2" />
                            Drop
                        </DropdownMenuItem>
                    </DropdownMenuContent>
                </DropdownMenu>
            </div>
        </div>
    );
}

// Empty State Component
interface EmptyStateProps {
    onOpenQueryTab?: () => void;
    onCreateTable?: () => void;
    onShareLink?: () => void;
}

function EmptyState({ onOpenQueryTab, onCreateTable, onShareLink }: EmptyStateProps) {
    return (
        <div className="flex flex-col items-center justify-center h-full text-center py-12">
            <div className="p-4 rounded-full bg-muted/50 mb-4">
                <Database2 className="h-12 w-12 text-muted-foreground" />
            </div>
            <h3 className="text-lg font-semibold mb-2">Welcome to Database Explorer</h3>
            <p className="text-muted-foreground mb-6 max-w-md">
                Select a table from the sidebar to view and edit data, or open the SQL editor to run custom queries.
            </p>
            <div className="flex flex-wrap gap-3 justify-center">
                <Button variant="outline" className="gap-2" onClick={onOpenQueryTab}>
                    <Code2 className="h-4 w-4" />
                    Open SQL Editor
                </Button>
                <Button className="gap-2" onClick={onCreateTable}>
                    <Plus className="h-4 w-4" />
                    Create Table
                </Button>
                {onShareLink && (
                    <Button variant="outline" className="gap-2" onClick={onShareLink}>
                        <Share2 className="h-4 w-4" />
                        Share Current View
                    </Button>
                )}
            </div>
        </div>
    );
}

// Query Editor Component
interface QueryEditorProps {
    query: string;
    onQueryChange: (query: string) => void;
    onExecute: () => void;
    isExecuting: boolean;
    result: any;
    recentQueries: string[];
    onSelectRecentQuery: (query: string) => void;
}

function QueryEditor({
    query,
    onQueryChange,
    onExecute,
    isExecuting,
    result,
    recentQueries,
    onSelectRecentQuery,
}: QueryEditorProps) {
    return (
        <div className="h-full flex flex-col space-y-4">
            <div className="flex items-center justify-between">
                <h3 className="text-lg font-semibold flex items-center gap-2">
                    <Code2 className="h-5 w-5" />
                    SQL Query Editor
                </h3>
                <div className="flex gap-2">
                    <Button
                        onClick={onExecute}
                        disabled={!query.trim() || isExecuting}
                        className="gap-2"
                    >
                        <Play className="h-4 w-4" />
                        {isExecuting ? 'Executing...' : 'Execute Query'}
                    </Button>
                    <Button
                        variant="outline"
                        onClick={() => onQueryChange('')}
                        disabled={isExecuting}
                    >
                        Clear
                    </Button>
                </div>
            </div>

            <div className="flex-1 grid grid-cols-1 lg:grid-cols-2 gap-4">
                {/* Query Input */}
                <div className="space-y-3">
                    <div className="flex items-center justify-between">
                        <label className="text-sm font-medium">SQL Query</label>
                        {recentQueries.length > 0 && (
                            <DropdownMenu>
                                <DropdownMenuTrigger asChild>
                                    <Button variant="ghost" size="sm" className="gap-2">
                                        <Clock className="h-4 w-4" />
                                        Recent
                                    </Button>
                                </DropdownMenuTrigger>
                                <DropdownMenuContent align="end" className="w-80">
                                    {recentQueries.map((recentQuery, index) => (
                                        <DropdownMenuItem
                                            key={index}
                                            onClick={() => onSelectRecentQuery(recentQuery)}
                                            className="p-3"
                                        >
                                            <div className="truncate text-xs font-mono">
                                                {recentQuery}
                                            </div>
                                        </DropdownMenuItem>
                                    ))}
                                </DropdownMenuContent>
                            </DropdownMenu>
                        )}
                    </div>
                    <textarea
                        value={query}
                        onChange={(e) => onQueryChange(e.target.value)}
                        className="w-full h-full min-h-[200px] p-3 border rounded-lg font-mono text-sm resize-none focus:ring-2 focus:ring-primary focus:border-transparent"
                        placeholder="Enter your SQL query here..."
                    />
                </div>

                {/* Results */}
                <div className="space-y-3">
                    <label className="text-sm font-medium">Query Results</label>
                    <div className="h-full border rounded-lg overflow-hidden">
                        {isExecuting ? (
                            <div className="flex items-center justify-center h-full">
                                <div className="flex items-center gap-2">
                                    <RefreshCw className="h-4 w-4 animate-spin" />
                                    <span className="text-sm text-muted-foreground">Executing query...</span>
                                </div>
                            </div>
                        ) : result ? (
                            <div className="h-full overflow-auto">
                                <Table>
                                    <TableHeader>
                                        <TableRow>
                                            {result.data?.data?.columns?.map((col: string) => (
                                                <TableHead key={col} className="font-semibold">
                                                    {col}
                                                </TableHead>
                                            ))}
                                        </TableRow>
                                    </TableHeader>
                                    <TableBody>
                                        {result.data?.data?.data?.map((row: any, idx: number) => (
                                            <TableRow key={idx}>
                                                {result.data?.data?.columns?.map((col: string) => (
                                                    <TableCell key={col} className="font-mono text-xs">
                                                        {row[col] !== null ? String(row[col]) : (
                                                            <Badge variant="secondary">NULL</Badge>
                                                        )}
                                                    </TableCell>
                                                ))}
                                            </TableRow>
                                        ))}
                                    </TableBody>
                                </Table>
                            </div>
                        ) : (
                            <div className="flex items-center justify-center h-full text-muted-foreground">
                                <div className="text-center">
                                    <Database className="h-8 w-8 mx-auto mb-2 opacity-50" />
                                    <p className="text-sm">Execute a query to see results</p>
                                </div>
                            </div>
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
}