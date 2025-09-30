'use client';

import { useState } from 'react';
import { useQuery, useMutation } from '@tanstack/react-query';
import { api } from '@/lib/api-client';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Database, Play, Download, Search, Table as TableIcon, Plus, Trash2 } from 'lucide-react';
import { toast } from 'sonner';
import { CreateTableDialog } from './create-table';

export default function DatabasePage() {
    const [selectedTable, setSelectedTable] = useState<string>('');
    const [sqlQuery, setSqlQuery] = useState('SELECT * FROM users LIMIT 10;');
    const [searchTerm, setSearchTerm] = useState('');
    const [openTabs, setOpenTabs] = useState<Array<{ id: string; name: string; type: 'table' | 'query' }>>([]);
    const [activeTab, setActiveTab] = useState<string>('');
    const [createTableOpen, setCreateTableOpen] = useState(false);

    // Tab management functions
    const openTableTab = (tableName: string) => {
        const tabId = `table-${tableName}`;
        const existingTab = openTabs.find(tab => tab.id === tabId);

        if (!existingTab) {
            const newTab = { id: tabId, name: tableName, type: 'table' as const };
            setOpenTabs(prev => [...prev, newTab]);
        }
        setActiveTab(tabId);
        setSelectedTable(tableName);
    };

    const closeTab = (tabId: string) => {
        setOpenTabs(prev => {
            const newTabs = prev.filter(tab => tab.id !== tabId);
            if (activeTab === tabId) {
                setActiveTab(newTabs.length > 0 ? newTabs[newTabs.length - 1].id : '');
            }
            return newTabs;
        });
    };

    const openQueryTab = () => {
        const tabId = 'query-tab';
        const existingTab = openTabs.find(tab => tab.id === tabId);

        if (!existingTab) {
            const newTab = { id: tabId, name: 'SQL Query', type: 'query' as const };
            setOpenTabs(prev => [...prev, newTab]);
        }
        setActiveTab(tabId);
    };

    // Use the centralized API client
    const apiClient = {
        getTables: () => api.getTables(),
        getTableData: (tableName: string, page = 1, limit = 20) =>
            api.getTableData(tableName, limit, (page - 1) * limit),
        getTableSchema: (tableName: string) =>
            api.getTableSchema(tableName),
        executeQuery: (query: string) => api.executeQuery(query),
    };

    const { data: tables, isLoading: tablesLoading } = useQuery({
        queryKey: ['database-tables'],
        queryFn: () => apiClient.getTables(),
    });

    const { data: tableData, isLoading: tableDataLoading } = useQuery({
        queryKey: ['table-data', selectedTable],
        queryFn: () => apiClient.getTableData(selectedTable),
        enabled: !!selectedTable,
    });

    const { data: schema } = useQuery({
        queryKey: ['table-schema', selectedTable],
        queryFn: () => apiClient.getTableSchema(selectedTable),
        enabled: !!selectedTable,
    });

    const executeQueryMutation = useMutation({
        mutationFn: (query: string) => apiClient.executeQuery(query),
        onSuccess: () => {
            toast.success('Query executed successfully');
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
        executeQueryMutation.mutate(sqlQuery);
    };

    const filteredTables = tables?.data?.data?.filter((table: any) =>
        table.name.toLowerCase().includes(searchTerm.toLowerCase())
    );

    return (
        <div className="space-y-6">
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">Database Explorer</h1>
                    <p className="text-gray-500 mt-1">
                        Browse tables, execute queries, and view schema
                    </p>
                </div>
                <div className="flex gap-2">
                    <Button onClick={openQueryTab} variant="outline">
                        <Database className="h-4 w-4 mr-2" />
                        SQL Query
                    </Button>
                    <Button onClick={() => setCreateTableOpen(true)}>
                        <Plus className="h-4 w-4 mr-2" />
                        Create Table
                    </Button>
                </div>
            </div>


            <CreateTableDialog open={createTableOpen} onOpenChange={setCreateTableOpen} />

            {/* Main Content Area */}
            <div className="space-y-4">
                <div className="grid md:grid-cols-4 gap-4">
                    {/* Tables List */}
                    <Card className="md:col-span-1">
                        <CardHeader>
                            <CardTitle className="text-sm">Tables</CardTitle>
                            <div className="relative mt-2">
                                <Search className="absolute left-2 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400" />
                                <Input
                                    placeholder="Search tables..."
                                    value={searchTerm}
                                    onChange={(e) => setSearchTerm(e.target.value)}
                                    className="pl-8 h-9"
                                />
                            </div>
                        </CardHeader>
                        <CardContent>
                            <div className="space-y-1">
                                {tablesLoading ? (
                                    <p className="text-sm text-gray-500">Loading...</p>
                                ) : filteredTables?.length === 0 ? (
                                    <p className="text-sm text-gray-500">No tables found</p>
                                ) : (
                                    filteredTables?.map((table: any) => (
                                        <button
                                            key={table.name}
                                            onClick={() => openTableTab(table.name)}
                                            className={`w-full text-left px-3 py-2 rounded-md text-sm transition-colors ${openTabs.some(tab => tab.name === table.name)
                                                ? 'bg-blue-50 text-blue-600 font-medium'
                                                : 'hover:bg-gray-100'
                                                }`}
                                        >
                                            <div className="flex items-center justify-between">
                                                <span>{table.name}</span>
                                                <Badge variant="secondary" className="text-xs">
                                                    {table.row_count || 0}
                                                </Badge>
                                            </div>
                                        </button>
                                    ))
                                )}
                            </div>
                        </CardContent>
                    </Card>

                    {/* Main Content Area with Tabs */}
                    <div className="md:col-span-3">
                        <Card>
                            <CardHeader>
                                <div className="flex items-center justify-between">
                                    <div className="flex items-center space-x-4 relative ">

                                        {/* Tabs */}
                                        <div className='w-full overflow-x-auto
                                        max-w-[850px]'>
                                            {openTabs.length > 0 && (
                                                <div className="flex space-x-1 overflow-x-auto scrollbar-thin scrollbar-thumb-gray-300 scrollbar-track-gray-100">
                                                    {openTabs.map((tab) => (
                                                        <button
                                                            key={tab.id}
                                                            onClick={() => {
                                                                setActiveTab(tab.id);
                                                                if (tab.type === 'table') {
                                                                    setSelectedTable(tab.name);
                                                                }
                                                            }}
                                                            className={`px-3 py-1 text-sm rounded-md flex items-center space-x-2 flex-shrink-0 ${activeTab === tab.id
                                                                ? 'bg-blue-100 text-blue-700'
                                                                : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                                                                }`}
                                                        >
                                                            <TableIcon className="h-4 w-4" />
                                                            <span>{tab.name}</span>
                                                            <span
                                                                onClick={(e) => {
                                                                    e.stopPropagation();
                                                                    closeTab(tab.id);
                                                                }}
                                                                className="ml-1 hover:text-red-600 cursor-pointer"
                                                                role="button"
                                                                aria-label={`Close ${tab.name}`}
                                                            >
                                                                <Trash2 className="h-3 w-3" />
                                                            </span>
                                                        </button>
                                                    ))}
                                                </div>
                                            )}

                                        </div>
                                      
                                    </div>
                                    {activeTab && (
                                        <Button size="sm" variant="outline">
                                            <Download className="h-4 w-4 mr-2" />
                                            Export
                                        </Button>
                                    )}
                                </div>
                            </CardHeader>
                            <CardContent>
                                {!activeTab ? (
                                    <div className="text-center py-12 text-gray-500">
                                        <Database className="h-12 w-12 mx-auto mb-4 opacity-50" />
                                        <p>Select a table from the list to view data</p>
                                    </div>
                                ) : (
                                    <div className="space-y-6">
                                        {/* Table Data */}
                                        {openTabs.find(tab => tab.id === activeTab)?.type === 'table' && (
                                            <div>
                                                <h3 className="font-semibold text-lg mb-4">Table Data</h3>
                                                {tableDataLoading ? (
                                                    <p className="text-sm text-gray-500">Loading...</p>
                                                ) : tableData?.data?.data?.data?.length === 0 ? (
                                                    <p className="text-sm text-gray-500">No data found</p>
                                                ) : (
                                                    <div className="overflow-x-auto">
                                                        <Table>
                                                            <TableHeader>
                                                                <TableRow>
                                                                    {tableData?.data?.data?.columns?.map((col: string) => (
                                                                        <TableHead key={col}>{col}</TableHead>
                                                                    ))}
                                                                </TableRow>
                                                            </TableHeader>
                                                            <TableBody>
                                                                {tableData?.data?.data?.data?.map((row: any, idx: number) => (
                                                                    <TableRow key={idx}>
                                                                        {tableData?.data?.data?.columns?.map((col: string) => (
                                                                            <TableCell key={col}>
                                                                                {row[col] !== null ? String(row[col]) : 'NULL'}
                                                                            </TableCell>
                                                                        ))}
                                                                    </TableRow>
                                                                ))}
                                                            </TableBody>
                                                        </Table>
                                                    </div>
                                                )}
                                            </div>
                                        )}

                                        {/* SQL Query Editor */}
                                        {openTabs.find(tab => tab.id === activeTab)?.type === 'query' && (
                                            <div>
                                                <h3 className="font-semibold text-lg mb-4">SQL Query Editor</h3>
                                                <div className="space-y-4">
                                                    <div>
                                                        <label className="block text-sm font-medium mb-2">SQL Query</label>
                                                        <textarea
                                                            value={sqlQuery}
                                                            onChange={(e) => setSqlQuery(e.target.value)}
                                                            className="w-full h-32 p-3 border rounded-md font-mono text-sm"
                                                            placeholder="Enter your SQL query here..."
                                                        />
                                                    </div>
                                                    <div className="flex gap-2">
                                                        <Button onClick={handleExecuteQuery} disabled={!sqlQuery.trim()}>
                                                            <Play className="h-4 w-4 mr-2" />
                                                            Execute Query
                                                        </Button>
                                                        <Button variant="outline" onClick={() => setSqlQuery('')}>
                                                            Clear
                                                        </Button>
                                                    </div>
                                                    {executeQueryMutation.isPending && (
                                                        <p className="text-sm text-gray-500">Executing query...</p>
                                                    )}
                                                    {executeQueryMutation.data && (
                                                        <div className="mt-4">
                                                            <h4 className="font-medium mb-2">Query Results</h4>
                                                            <div className="overflow-x-auto border rounded-md">
                                                                <Table>
                                                                    <TableHeader>
                                                                        <TableRow>
                                                                            {executeQueryMutation.data.data?.data?.columns?.map((col: string) => (
                                                                                <TableHead key={col}>{col}</TableHead>
                                                                            ))}
                                                                        </TableRow>
                                                                    </TableHeader>
                                                                    <TableBody>
                                                                        {executeQueryMutation.data.data?.data?.data?.map((row: any, idx: number) => (
                                                                            <TableRow key={idx}>
                                                                                {executeQueryMutation.data.data?.data?.columns?.map((col: string) => (
                                                                                    <TableCell key={col}>
                                                                                        {row[col] !== null ? String(row[col]) : 'NULL'}
                                                                                    </TableCell>
                                                                                ))}
                                                                            </TableRow>
                                                                        ))}
                                                                    </TableBody>
                                                                </Table>
                                                            </div>
                                                        </div>
                                                    )}
                                                </div>
                                            </div>
                                        )}

                                        {/* Schema */}
                                        {/* {openTabs.find(tab => tab.id === activeTab)?.type === 'table' && (
                                            <div>
                                                <h3 className="font-semibold text-lg mb-4">Schema</h3>
                                                {schema?.data?.data?.columns ? (
                                                    <div className="border rounded-lg p-4">
                                                        <h4 className="font-semibold text-lg mb-3">{schema.data.data.table}</h4>
                                                        <Table>
                                                            <TableHeader>
                                                                <TableRow>
                                                                    <TableHead>Column</TableHead>
                                                                    <TableHead>Type</TableHead>
                                                                    <TableHead>Nullable</TableHead>
                                                                    <TableHead>Default</TableHead>
                                                                </TableRow>
                                                            </TableHeader>
                                                            <TableBody>
                                                                {schema.data.data.columns.map((col: any) => (
                                                                    <TableRow key={col.name}>
                                                                        <TableCell className="font-medium">{col.name}</TableCell>
                                                                        <TableCell>
                                                                            <code className="text-xs bg-gray-100 px-2 py-1 rounded">
                                                                                {col.type}
                                                                            </code>
                                                                        </TableCell>
                                                                        <TableCell>
                                                                            <Badge variant={col.nullable ? 'secondary' : 'outline'}>
                                                                                {col.nullable ? 'Yes' : 'No'}
                                                                            </Badge>
                                                                        </TableCell>
                                                                        <TableCell className="text-xs">
                                                                            {col.default_value || 'NULL'}
                                                                        </TableCell>
                                                                    </TableRow>
                                                                ))}
                                                            </TableBody>
                                                        </Table>
                                                    </div>
                                                ) : (
                                                    <p className="text-sm text-gray-500">No schema data available</p>
                                                )}
                                            </div>
                                        )} */}
                                    </div>
                                )}
                            </CardContent>
                        </Card>
                    </div>
                </div>
            </div>
        </div>
    );
}
