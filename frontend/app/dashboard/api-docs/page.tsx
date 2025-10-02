'use client';

import { useState, useEffect, useCallback } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion';
import {
    ExternalLink,
    Copy,
    Download,
    Filter,
    Search,
    Globe,
    Shield,
    BookOpen,
    FileText,
    Zap,
    Code,
    Play,
    Loader2,
    XCircle,
    Database,
    Type,
    FileCode
} from 'lucide-react';
import { useToast } from '../../../hooks/use-toast';

interface APIEndpoint {
    path: string;
    method: string;
    summary: string;
    description: string;
    tags: string[];
    parameters?: unknown[];
    responses?: Record<string, unknown>;
    requestBody?: {
        content?: Record<string, {
            schema?: unknown;
        }>;
    };
}

interface TestResult {
    status: number;
    statusText: string;
    data: unknown;
    headers: Record<string, string>;
    duration: number;
    error?: string;
}

interface TestState {
    isLoading: boolean;
    result?: TestResult;
    params: Record<string, string>;
    body: string;
    headers: Record<string, string>;
}

interface APIDocumentation {
    info: {
        title: string;
        version: string;
        description: string;
    };
    paths: Record<string, Record<string, APIEndpoint>>;
    tags: Array<{
        name: string;
        description: string;
    }>;
}

export default function APIDocsPage() {
    const router = useRouter();
    const searchParams = useSearchParams();
    const [apiDoc, setApiDoc] = useState<APIDocumentation | null>(null);
    const [loading, setLoading] = useState(true);
    const [searchTerm, setSearchTerm] = useState('');
    const [selectedTag, setSelectedTag] = useState('all');
    const [selectedMethod, setSelectedMethod] = useState('all');
    const [showCodeExamples, setShowCodeExamples] = useState(false);
    const [openAccordions, setOpenAccordions] = useState<string[]>([]);
    const [baseUrl] = useState('http://localhost:8080/api/v1');
    const [testStates, setTestStates] = useState<Record<string, TestState>>({});
    const { toast } = useToast();

    // URL parameter management
    const updateURLParams = useCallback((updates: Record<string, string | string[] | null>) => {
        const params = new URLSearchParams(searchParams.toString());

        Object.entries(updates).forEach(([key, value]) => {
            if (value === null || value === '' || (Array.isArray(value) && value.length === 0)) {
                params.delete(key);
            } else if (Array.isArray(value)) {
                params.set(key, value.join(','));
            } else {
                params.set(key, value);
            }
        });

        const newURL = `${window.location.pathname}?${params.toString()}`;
        router.replace(newURL, { scroll: false });
    }, [searchParams, router]);

    const getURLParam = useCallback((key: string, defaultValue: string = '') => {
        return searchParams.get(key) || defaultValue;
    }, [searchParams]);

    const getURLArrayParam = useCallback((key: string, defaultValue: string[] = []) => {
        const value = searchParams.get(key);
        return value ? value.split(',').filter(Boolean) : defaultValue;
    }, [searchParams]);

    // Initialize state from URL parameters
    useEffect(() => {
        const urlSearchTerm = getURLParam('search', '');
        const urlTag = getURLParam('tag', 'all');
        const urlMethod = getURLParam('method', 'all');
        const urlShowCode = getURLParam('showCode', 'false') === 'true';
        const urlAccordions = getURLArrayParam('accordions', []);
        const urlScrollY = getURLParam('scrollY', '');

        setSearchTerm(urlSearchTerm);
        setSelectedTag(urlTag);
        setSelectedMethod(urlMethod);
        setShowCodeExamples(urlShowCode);
        setOpenAccordions(urlAccordions);

        // Restore scroll position
        if (urlScrollY) {
            setTimeout(() => {
                window.scrollTo(0, parseInt(urlScrollY, 10));
            }, 100);
        }
    }, [getURLParam, getURLArrayParam]);

    // Save scroll position to URL
    useEffect(() => {
        const handleScroll = () => {
            const scrollY = window.scrollY;
            if (scrollY > 0) {
                updateURLParams({ scrollY: scrollY.toString() });
            } else {
                updateURLParams({ scrollY: null });
            }
        };

        let timeoutId: NodeJS.Timeout;
        const throttledScroll = () => {
            clearTimeout(timeoutId);
            timeoutId = setTimeout(handleScroll, 100);
        };

        window.addEventListener('scroll', throttledScroll);
        return () => {
            window.removeEventListener('scroll', throttledScroll);
            clearTimeout(timeoutId);
        };
    }, [updateURLParams]);

    const loadAPIDocumentation = useCallback(async () => {
        try {
            const response = await fetch('http://localhost:8080/docs/doc.json');
            const data = await response.json();
            setApiDoc(data);
        } catch (error) {
            console.error('Failed to load API documentation:', error);
            toast({
                title: "Error",
                description: "Failed to load API documentation",
                variant: "destructive",
            });
        } finally {
            setLoading(false);
        }
    }, [toast]);

    const getTestState = (path: string, method: string): TestState => {
        const key = `${path}-${method}`;
        return testStates[key] || {
            isLoading: false,
            params: {},
            body: '',
            headers: { 'Content-Type': 'application/json' }
        };
    };

    const updateTestState = (path: string, method: string, updates: Partial<TestState>) => {
        const key = `${path}-${method}`;
        setTestStates(prev => ({
            ...prev,
            [key]: { ...getTestState(path, method), ...updates }
        }));
    };

    const testEndpoint = async (path: string, method: string) => {
        const testState = getTestState(path, method);

        updateTestState(path, method, { isLoading: true });

        try {
            const startTime = Date.now();
            const fullPath = path.startsWith('/') ? path : `/${path}`;

            // Build URL with query parameters
            const url = new URL(`${baseUrl}${fullPath}`);
            Object.entries(testState.params).forEach(([key, value]) => {
                if (value) url.searchParams.append(key, value);
            });

            // Prepare headers
            const headers: Record<string, string> = {
                'Content-Type': 'application/json',
                ...testState.headers
            };

            // Add auth token if available
            const token = localStorage.getItem('auth_token');
            if (token) {
                headers.Authorization = `Bearer ${token}`;
            }

            // Make the request
            const response = await fetch(url.toString(), {
                method: method.toUpperCase(),
                headers,
                body: testState.body || undefined
            });

            const duration = Date.now() - startTime;
            const responseData = await response.text();

            let parsedData: unknown;
            try {
                parsedData = JSON.parse(responseData);
            } catch {
                parsedData = responseData;
            }

            const result: TestResult = {
                status: response.status,
                statusText: response.statusText,
                data: parsedData,
                headers: Object.fromEntries(response.headers.entries()),
                duration
            };

            updateTestState(path, method, {
                isLoading: false,
                result
            });

            toast({
                title: "Request completed",
                description: `Status: ${response.status} ${response.statusText}`,
                variant: response.ok ? "default" : "destructive",
            });

        } catch (error) {
            const duration = Date.now() - Date.now();
            const result: TestResult = {
                status: 0,
                statusText: 'Network Error',
                data: null,
                headers: {},
                duration,
                error: error instanceof Error ? error.message : 'Unknown error'
            };

            updateTestState(path, method, {
                isLoading: false,
                result
            });

            toast({
                title: "Request failed",
                description: error instanceof Error ? error.message : 'Unknown error',
                variant: "destructive",
            });
        }
    };

    useEffect(() => {
        loadAPIDocumentation();
    }, [loadAPIDocumentation]);

    // Update URL when filters change
    useEffect(() => {
        updateURLParams({
            search: searchTerm,
            tag: selectedTag,
            method: selectedMethod,
            showCode: showCodeExamples ? 'true' : null,
            accordions: openAccordions.length > 0 ? openAccordions : null
        });
    }, [searchTerm, selectedTag, selectedMethod, showCodeExamples, openAccordions, updateURLParams]);

    // Handle filter changes
    const handleSearchChange = useCallback((value: string) => {
        setSearchTerm(value);
    }, []);

    const handleTagChange = useCallback((value: string) => {
        setSelectedTag(value);
    }, []);

    const handleMethodChange = useCallback((value: string) => {
        setSelectedMethod(value);
    }, []);

    const handleShowCodeChange = useCallback((value: boolean) => {
        setShowCodeExamples(value);
    }, []);

    const handleAccordionChange = useCallback((value: string[]) => {
        setOpenAccordions(value);
    }, []);

    const clearFilters = useCallback(() => {
        setSearchTerm('');
        setSelectedTag('all');
        setSelectedMethod('all');
        setShowCodeExamples(false);
        setOpenAccordions([]);
    }, []);

    const copyCurrentURL = useCallback(() => {
        const currentURL = window.location.href;
        navigator.clipboard.writeText(currentURL);
        toast({
            title: "URL Copied!",
            description: "Current view URL copied to clipboard",
        });
    }, [toast]);

    const copyToClipboard = (text: string, label: string) => {
        navigator.clipboard.writeText(text);
        toast({
            title: "Copied!",
            description: `${label} copied to clipboard`,
        });
    };

    const getMethodColor = (method: string) => {
        const colors = {
            GET: 'bg-green-100 text-green-800 border-green-200',
            POST: 'bg-blue-100 text-blue-800 border-blue-200',
            PUT: 'bg-yellow-100 text-yellow-800 border-yellow-200',
            DELETE: 'bg-red-100 text-red-800 border-red-200',
            PATCH: 'bg-purple-100 text-purple-800 border-purple-200',
        };
        return colors[method as keyof typeof colors] || 'bg-gray-100 text-gray-800 border-gray-200';
    };


    const formatJSON = (data: unknown) => {
        try {
            return JSON.stringify(data, null, 2);
        } catch {
            return String(data);
        }
    };

    // Schema generation functions
    const generateGoStruct = (schema: Record<string, unknown>, structName: string = 'Model'): string => {
        if (!schema || !schema.properties) {
            return `type ${structName} struct {\n    // No properties defined\n}`;
        }

        let struct = `type ${structName} struct {\n`;

        Object.entries(schema.properties as Record<string, unknown>).forEach(([key, prop]) => {
            const propObj = prop as Record<string, unknown>;
            const goType = mapToGoType(propObj);
            const jsonTag = `json:"${key}"`;
            const gormTag = generateGormTag(propObj, key);
            const omitempty = propObj.required ? '' : ',omitempty';

            struct += `    ${toPascalCase(key)} ${goType} \`${jsonTag}${omitempty}\` ${gormTag}\n`;
        });

        struct += '}';
        return struct;
    };

    const generateTypeScriptInterface = (schema: Record<string, unknown>, interfaceName: string = 'Model'): string => {
        if (!schema || !schema.properties) {
            return `interface ${interfaceName} {\n    // No properties defined\n}`;
        }

        let interfaceStr = `interface ${interfaceName} {\n`;

        Object.entries(schema.properties as Record<string, unknown>).forEach(([key, prop]) => {
            const propObj = prop as Record<string, unknown>;
            const tsType = mapToTypeScriptType(propObj);
            const optional = propObj.required ? '' : '?';
            interfaceStr += `    ${key}${optional}: ${tsType};\n`;
        });

        interfaceStr += '}';
        return interfaceStr;
    };

    const mapToGoType = (prop: Record<string, unknown>): string => {
        const type = prop.type || 'string';
        const format = prop.format;

        switch (type) {
            case 'integer':
                return 'int';
            case 'number':
                return 'float64';
            case 'boolean':
                return 'bool';
            case 'string':
                if (format === 'date-time') return 'time.Time';
                if (format === 'date') return 'time.Time';
                if (format === 'email') return 'string';
                if (format === 'uuid') return 'string';
                return 'string';
            case 'array':
                const itemType = prop.items ? mapToGoType(prop.items as Record<string, unknown>) : 'interface{}';
                return `[]${itemType}`;
            case 'object':
                return 'map[string]interface{}';
            default:
                return 'interface{}';
        }
    };

    const mapToTypeScriptType = (prop: Record<string, unknown>): string => {
        const type = prop.type || 'string';
        const format = prop.format;

        switch (type) {
            case 'integer':
                return 'number';
            case 'number':
                return 'number';
            case 'boolean':
                return 'boolean';
            case 'string':
                if (format === 'date-time') return 'string'; // ISO date string
                if (format === 'date') return 'string';
                if (format === 'email') return 'string';
                if (format === 'uuid') return 'string';
                return 'string';
            case 'array':
                const itemType = prop.items ? mapToTypeScriptType(prop.items as Record<string, unknown>) : 'any';
                return `${itemType}[]`;
            case 'object':
                return 'Record<string, any>';
            default:
                return 'any';
        }
    };

    const generateGormTag = (prop: Record<string, unknown>, fieldName: string): string => {
        const tags = [];

        // Primary key
        if (fieldName === 'id' || fieldName.endsWith('_id')) {
            tags.push('primaryKey');
        }

        // Auto increment
        if (prop.type === 'integer' && fieldName === 'id') {
            tags.push('autoIncrement');
        }

        // Not null
        if (prop.required) {
            tags.push('not null');
        }

        // Unique
        if (prop.unique) {
            tags.push('unique');
        }

        // Size
        if (prop.maxLength) {
            tags.push(`size:${prop.maxLength}`);
        }

        // Default value
        if (prop.default !== undefined) {
            tags.push(`default:${prop.default}`);
        }

        return tags.length > 0 ? `gorm:"${tags.join(';')}"` : '';
    };

    const toPascalCase = (str: string): string => {
        return str.replace(/(?:^|_)([a-z])/g, (_, letter) => letter.toUpperCase());
    };

    const getSchemaFromEndpoint = (endpoint: APIEndpoint, type: 'request' | 'response') => {
        if (type === 'request') {
            // Look for request body schema
            return endpoint.requestBody?.content?.['application/json']?.schema;
        } else {
            // Look for 200 response schema
            if (!endpoint.responses) return null;
            const successResponse = endpoint.responses['200'] || endpoint.responses['201'];
            return (successResponse as { content?: Record<string, { schema?: unknown }> })?.content?.['application/json']?.schema;
        }
    };

    const renderTestInterface = (path: string, method: string, endpoint: APIEndpoint) => {
        const testState = getTestState(path, method);
        const hasAuth = !!localStorage.getItem('auth_token');

        return (
            <div className="space-y-4 p-4 bg-muted/50 rounded-lg">
                <div className="flex items-center justify-between">
                    <h4 className="text-sm font-medium flex items-center gap-2">
                        <Play className="h-4 w-4" />
                        Test Endpoint
                        {hasAuth && (
                            <Badge variant="outline" className="text-xs">
                                <Shield className="h-3 w-3 mr-1" />
                                Authenticated
                            </Badge>
                        )}
                    </h4>
                    <Button
                        size="sm"
                        onClick={() => testEndpoint(path, method)}
                        disabled={testState.isLoading}
                    >
                        {testState.isLoading ? (
                            <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                        ) : (
                            <Play className="h-4 w-4 mr-2" />
                        )}
                        {testState.isLoading ? 'Testing...' : 'Test'}
                    </Button>
                </div>

                {/* Parameters */}
                {endpoint.parameters && endpoint.parameters.length > 0 && (
                    <div className="space-y-2">
                        <label className="text-sm font-medium">Parameters</label>
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
                            {endpoint.parameters.map((param: unknown, index) => {
                                const paramObj = param as Record<string, unknown>;
                                return (
                                    <div key={index} className="space-y-1">
                                        <label className="text-xs text-muted-foreground">
                                            {String(paramObj.name)} {Boolean(paramObj.required) && <span className="text-red-500">*</span>}
                                        </label>
                                        <Input
                                            placeholder={String(paramObj.description || paramObj.name)}
                                            value={testState.params[String(paramObj.name)] || ''}
                                            onChange={(e) => updateTestState(path, method, {
                                                params: { ...testState.params, [String(paramObj.name)]: e.target.value }
                                            })}
                                        />
                                    </div>
                                );
                            })}
                        </div>
                    </div>
                )}

                {/* Request Body */}
                {(method.toUpperCase() === 'POST' || method.toUpperCase() === 'PUT' || method.toUpperCase() === 'PATCH') && (
                    <div className="space-y-2">
                        <label className="text-sm font-medium">Request Body</label>
                        <Textarea
                            placeholder="Enter JSON request body..."
                            value={testState.body}
                            onChange={(e) => updateTestState(path, method, { body: e.target.value })}
                            rows={4}
                            className="font-mono text-sm"
                        />
                    </div>
                )}

                {/* Custom Headers */}
                <div className="space-y-2">
                    <label className="text-sm font-medium">Custom Headers</label>
                    <div className="space-y-2">
                        {Object.entries(testState.headers).map(([key, value], index) => (
                            <div key={index} className="flex gap-2">
                                <Input
                                    placeholder="Header name"
                                    value={key}
                                    onChange={(e) => {
                                        const newHeaders = { ...testState.headers };
                                        delete newHeaders[key];
                                        newHeaders[e.target.value] = value;
                                        updateTestState(path, method, { headers: newHeaders });
                                    }}
                                    className="flex-1"
                                />
                                <Input
                                    placeholder="Header value"
                                    value={value}
                                    onChange={(e) => updateTestState(path, method, {
                                        headers: { ...testState.headers, [key]: e.target.value }
                                    })}
                                    className="flex-1"
                                />
                                <Button
                                    size="sm"
                                    variant="outline"
                                    onClick={() => {
                                        const newHeaders = { ...testState.headers };
                                        delete newHeaders[key];
                                        updateTestState(path, method, { headers: newHeaders });
                                    }}
                                >
                                    <XCircle className="h-4 w-4" />
                                </Button>
                            </div>
                        ))}
                        <Button
                            size="sm"
                            variant="outline"
                            onClick={() => updateTestState(path, method, {
                                headers: { ...testState.headers, '': '' }
                            })}
                        >
                            Add Header
                        </Button>
                    </div>
                </div>

                {/* Test Result */}
                {testState.result && (
                    <div className="space-y-2">
                        <div className="flex items-center gap-2">
                            <h4 className="text-sm font-medium">Response</h4>
                            <Badge
                                variant={testState.result.status >= 200 && testState.result.status < 300 ? "default" : "destructive"}
                                className="text-xs"
                            >
                                {testState.result.status} {testState.result.statusText}
                            </Badge>
                            <span className="text-xs text-muted-foreground">
                                {testState.result.duration}ms
                            </span>
                        </div>

                        <Tabs defaultValue="response" className="w-full">
                            <TabsList className="grid w-full grid-cols-3">
                                <TabsTrigger value="response">Response</TabsTrigger>
                                <TabsTrigger value="headers">Headers</TabsTrigger>
                                <TabsTrigger value="raw">Raw</TabsTrigger>
                            </TabsList>

                            <TabsContent value="response" className="mt-2">
                                <div className="bg-background border rounded-lg p-3">
                                    <pre className="text-sm overflow-auto max-h-64">
                                        {formatJSON(testState.result.data)}
                                    </pre>
                                </div>
                            </TabsContent>

                            <TabsContent value="headers" className="mt-2">
                                <div className="bg-background border rounded-lg p-3">
                                    <pre className="text-sm overflow-auto max-h-64">
                                        {formatJSON(testState.result.headers)}
                                    </pre>
                                </div>
                            </TabsContent>

                            <TabsContent value="raw" className="mt-2">
                                <div className="bg-background border rounded-lg p-3">
                                    <pre className="text-sm overflow-auto max-h-64">
                                        {JSON.stringify(testState.result, null, 2)}
                                    </pre>
                                </div>
                            </TabsContent>
                        </Tabs>
                    </div>
                )}
            </div>
        );
    };

    const renderSchemaDisplay = (endpoint: APIEndpoint, path: string) => {
        const requestSchema = getSchemaFromEndpoint(endpoint, 'request');
        const responseSchema = getSchemaFromEndpoint(endpoint, 'response');

        if (!requestSchema && !responseSchema) {
            return null;
        }

        const generateStructName = (type: string) => {
            const pathParts = path.split('/').filter(Boolean);
            const resource = pathParts[pathParts.length - 1] || 'Model';
            const capitalizedResource = resource.charAt(0).toUpperCase() + resource.slice(1);
            return `${capitalizedResource}${type}`;
        };

        return (
            <div className="space-y-4 p-4 bg-muted/30 rounded-lg">
                <div className="flex items-center gap-2">
                    <Database className="h-4 w-4" />
                    <h4 className="text-sm font-medium">Data Schemas</h4>
                </div>

                <Tabs defaultValue="response" className="w-full">
                    <TabsList className="grid w-full grid-cols-2">
                        <TabsTrigger value="request" disabled={!requestSchema}>
                            Request Schema
                        </TabsTrigger>
                        <TabsTrigger value="response" disabled={!responseSchema}>
                            Response Schema
                        </TabsTrigger>
                    </TabsList>

                    {requestSchema ? (
                        <TabsContent value="request" className="mt-4">
                            <div className="space-y-4">
                                <div className="flex items-center justify-between">
                                    <h5 className="text-sm font-medium">Request Body Schema</h5>
                                    <div className="flex gap-2">
                                        <Button
                                            size="sm"
                                            variant="outline"
                                            onClick={() => copyToClipboard(
                                                generateGoStruct(requestSchema as Record<string, unknown>, generateStructName('Request')),
                                                'Go Struct'
                                            )}
                                        >
                                            <FileCode className="h-3 w-3 mr-1" />
                                            Go
                                        </Button>
                                        <Button
                                            size="sm"
                                            variant="outline"
                                            onClick={() => copyToClipboard(
                                                generateTypeScriptInterface(requestSchema as Record<string, unknown>, generateStructName('Request')),
                                                'TypeScript Interface'
                                            )}
                                        >
                                            <Type className="h-3 w-3 mr-1" />
                                            TS
                                        </Button>
                                        <Button
                                            size="sm"
                                            variant="outline"
                                            onClick={() => copyToClipboard(
                                                formatJSON(requestSchema),
                                                'JSON Schema'
                                            )}
                                        >
                                            <FileText className="h-3 w-3 mr-1" />
                                            JSON
                                        </Button>
                                    </div>
                                </div>

                                <Tabs defaultValue="go" className="w-full">
                                    <TabsList className="grid w-full grid-cols-3">
                                        <TabsTrigger value="go">Go Struct</TabsTrigger>
                                        <TabsTrigger value="ts">TypeScript</TabsTrigger>
                                        <TabsTrigger value="json">JSON Schema</TabsTrigger>
                                    </TabsList>

                                    <TabsContent value="go" className="mt-2">
                                        <div className="bg-background border rounded-lg p-3">
                                            <pre className="text-sm overflow-auto max-h-64">
                                                {generateGoStruct(requestSchema as Record<string, unknown>, generateStructName('Request'))}
                                            </pre>
                                        </div>
                                    </TabsContent>

                                    <TabsContent value="ts" className="mt-2">
                                        <div className="bg-background border rounded-lg p-3">
                                            <pre className="text-sm overflow-auto max-h-64">
                                                {generateTypeScriptInterface(requestSchema as Record<string, unknown>, generateStructName('Request'))}
                                            </pre>
                                        </div>
                                    </TabsContent>

                                    <TabsContent value="json" className="mt-2">
                                        <div className="bg-background border rounded-lg p-3">
                                            <pre className="text-sm overflow-auto max-h-64">
                                                {formatJSON(requestSchema)}
                                            </pre>
                                        </div>
                                    </TabsContent>
                                </Tabs>
                            </div>
                        </TabsContent>
                    ) : null}

                    {responseSchema ? (
                        <TabsContent value="response" className="mt-4">
                            <div className="space-y-4">
                                <div className="flex items-center justify-between">
                                    <h5 className="text-sm font-medium">Response Schema</h5>
                                    <div className="flex gap-2">
                                        <Button
                                            size="sm"
                                            variant="outline"
                                            onClick={() => copyToClipboard(
                                                generateGoStruct(responseSchema as Record<string, unknown>, generateStructName('Response')),
                                                'Go Struct'
                                            )}
                                        >
                                            <FileCode className="h-3 w-3 mr-1" />
                                            Go
                                        </Button>
                                        <Button
                                            size="sm"
                                            variant="outline"
                                            onClick={() => copyToClipboard(
                                                generateTypeScriptInterface(responseSchema as Record<string, unknown>, generateStructName('Response')),
                                                'TypeScript Interface'
                                            )}
                                        >
                                            <Type className="h-3 w-3 mr-1" />
                                            TS
                                        </Button>
                                        <Button
                                            size="sm"
                                            variant="outline"
                                            onClick={() => copyToClipboard(
                                                formatJSON(responseSchema),
                                                'JSON Schema'
                                            )}
                                        >
                                            <FileText className="h-3 w-3 mr-1" />
                                            JSON
                                        </Button>
                                    </div>
                                </div>

                                <Tabs defaultValue="go" className="w-full">
                                    <TabsList className="grid w-full grid-cols-3">
                                        <TabsTrigger value="go">Go Struct</TabsTrigger>
                                        <TabsTrigger value="ts">TypeScript</TabsTrigger>
                                        <TabsTrigger value="json">JSON Schema</TabsTrigger>
                                    </TabsList>

                                    <TabsContent value="go" className="mt-2">
                                        <div className="bg-background border rounded-lg p-3">
                                            <pre className="text-sm overflow-auto max-h-64">
                                                {generateGoStruct(responseSchema as Record<string, unknown>, generateStructName('Response'))}
                                            </pre>
                                        </div>
                                    </TabsContent>

                                    <TabsContent value="ts" className="mt-2">
                                        <div className="bg-background border rounded-lg p-3">
                                            <pre className="text-sm overflow-auto max-h-64">
                                                {generateTypeScriptInterface(responseSchema as Record<string, unknown>, generateStructName('Response'))}
                                            </pre>
                                        </div>
                                    </TabsContent>

                                    <TabsContent value="json" className="mt-2">
                                        <div className="bg-background border rounded-lg p-3">
                                            <pre className="text-sm overflow-auto max-h-64">
                                                {formatJSON(responseSchema)}
                                            </pre>
                                        </div>
                                    </TabsContent>
                                </Tabs>
                            </div>
                        </TabsContent>
                    ) : null}
                </Tabs>
            </div>
        );
    };

    const filteredEndpoints = () => {
        if (!apiDoc) return {};

        const groupedEndpoints: Record<string, Array<{ path: string; method: string; endpoint: APIEndpoint }>> = {};

        Object.entries(apiDoc.paths).forEach(([path, methods]) => {
            Object.entries(methods).forEach(([method, endpoint]) => {
                const matchesSearch = endpoint.summary?.toLowerCase().includes(searchTerm.toLowerCase()) ||
                    endpoint.description?.toLowerCase().includes(searchTerm.toLowerCase()) ||
                    path.toLowerCase().includes(searchTerm.toLowerCase());

                const matchesTag = selectedTag === 'all' || endpoint.tags?.includes(selectedTag);
                const matchesMethod = selectedMethod === 'all' || method.toUpperCase() === selectedMethod.toUpperCase();

                if (matchesSearch && matchesTag && matchesMethod) {
                    const tag = endpoint.tags?.[0] || 'other';
                    if (!groupedEndpoints[tag]) {
                        groupedEndpoints[tag] = [];
                    }
                    groupedEndpoints[tag].push({ path, method: method.toUpperCase(), endpoint });
                }
            });
        });

        return groupedEndpoints;
    };

    const generateCodeSnippet = (endpoint: APIEndpoint, method: string, path: string) => {
        const fullUrl = `${baseUrl}${path}`;

        const curlExample = `curl -X ${method} "${fullUrl}" \\
  -H "Content-Type: application/json" \\
  -H "Authorization: Bearer YOUR_TOKEN"`;

        const fetchExample = `fetch('${fullUrl}', {
  method: '${method}',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer YOUR_TOKEN'
  }
})
.then(response => response.json())
.then(data => console.log(data));`;

        return { curlExample, fetchExample };
    };

    if (loading) {
        return (
            <div className="flex items-center justify-center h-64">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
            </div>
        );
    }

    if (!apiDoc) {
        return (
            <div className="text-center py-12">
                <FileText className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
                <h3 className="text-lg font-semibold mb-2">API Documentation Not Available</h3>
                <p className="text-muted-foreground mb-4">
                    Unable to load API documentation. Please check if the server is running.
                </p>
                <Button onClick={loadAPIDocumentation}>
                    <Zap className="h-4 w-4 mr-2" />
                    Retry
                </Button>
            </div>
        );
    }

    const groupedEndpoints = filteredEndpoints();
    const uniqueTags = Object.keys(groupedEndpoints);

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex items-center justify-between">
                <div>
                    <h1 className="text-3xl font-bold tracking-tight">API Documentation</h1>
                    <p className="text-muted-foreground">
                        Interactive API reference and testing tools
                    </p>
                </div>
                <div className="flex gap-2">
                    <Button variant="outline" onClick={() => window.open('http://localhost:8080/swagger-ui', '_blank')}>
                        <ExternalLink className="h-4 w-4 mr-2" />
                        Open Swagger UI
                    </Button>
                    <Button variant="outline" onClick={copyCurrentURL}>
                        <Copy className="h-4 w-4 mr-2" />
                        Share View
                    </Button>
                    <Button onClick={() => copyToClipboard(JSON.stringify(apiDoc, null, 2), 'API Documentation')}>
                        <Download className="h-4 w-4 mr-2" />
                        Export JSON
                    </Button>
                </div>
            </div>

            {/* API Info Card */}
            <Card>
                <CardHeader>
                    <div className="flex items-center gap-3">
                        <div className="p-2 bg-primary/10 rounded-lg">
                            <BookOpen className="h-6 w-6 text-primary" />
                        </div>
                        <div>
                            <CardTitle className="text-xl">{apiDoc.info.title}</CardTitle>
                            <CardDescription>Version {apiDoc.info.version}</CardDescription>
                        </div>
                    </div>
                </CardHeader>
                <CardContent>
                    <p className="text-muted-foreground">{apiDoc.info.description}</p>
                    <div className="mt-4 flex items-center gap-4 text-sm text-muted-foreground">
                        <div className="flex items-center gap-1">
                            <Globe className="h-4 w-4" />
                            <span>Base URL: {baseUrl}</span>
                        </div>
                        <div className="flex items-center gap-1">
                            <Shield className="h-4 w-4" />
                            <span>Bearer Token Authentication</span>
                        </div>
                    </div>
                </CardContent>
            </Card>

            {/* Filters */}
            <Card>
                <CardHeader>
                    <CardTitle className="text-lg">Filters & Search</CardTitle>
                </CardHeader>
                <CardContent>
                    <div className="grid grid-cols-1 md:grid-cols-5 gap-4">
                        <div className="relative">
                            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
                            <Input
                                placeholder="Search endpoints..."
                                value={searchTerm}
                                onChange={(e) => handleSearchChange(e.target.value)}
                                className="pl-10"
                            />
                        </div>
                        <Select value={selectedTag} onValueChange={handleTagChange}>
                            <SelectTrigger>
                                <SelectValue placeholder="Filter by tag" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="all">All Tags</SelectItem>
                                {uniqueTags.map(tag => (
                                    <SelectItem key={tag} value={tag}>{tag}</SelectItem>
                                ))}
                            </SelectContent>
                        </Select>
                        <Select value={selectedMethod} onValueChange={handleMethodChange}>
                            <SelectTrigger>
                                <SelectValue placeholder="Filter by method" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectItem value="all">All Methods</SelectItem>
                                <SelectItem value="GET">GET</SelectItem>
                                <SelectItem value="POST">POST</SelectItem>
                                <SelectItem value="PUT">PUT</SelectItem>
                                <SelectItem value="DELETE">DELETE</SelectItem>
                                <SelectItem value="PATCH">PATCH</SelectItem>
                            </SelectContent>
                        </Select>
                        <Button variant="outline" onClick={clearFilters}>
                            <Filter className="h-4 w-4 mr-2" />
                            Clear Filters
                        </Button>
                        <Button
                            variant={showCodeExamples ? "default" : "outline"}
                            onClick={() => handleShowCodeChange(!showCodeExamples)}
                        >
                            <Code className="h-4 w-4 mr-2" />
                            {showCodeExamples ? 'Hide' : 'Show'} Code Examples
                        </Button>
                    </div>
                </CardContent>
            </Card>

            {/* Endpoints List - Grouped by Tags */}
            <div className="space-y-4">
                {Object.entries(groupedEndpoints).map(([tag, endpoints]) => (
                    <Card key={tag}>
                        <CardHeader>
                            <CardTitle className="text-xl capitalize flex items-center gap-2">
                                <Badge variant="outline" className="text-sm">
                                    {endpoints.length} endpoint{endpoints.length !== 1 ? 's' : ''}
                                </Badge>
                                {tag}
                            </CardTitle>
                        </CardHeader>
                        <CardContent>
                            <Accordion
                                type="multiple"
                                className="space-y-2"
                                value={openAccordions}
                                onValueChange={handleAccordionChange}
                            >
                                {endpoints.map(({ path, method, endpoint }, index) => {
                                    const codeSnippets = generateCodeSnippet(endpoint, method, path);

                                    return (
                                        <AccordionItem key={`${path}-${method}-${index}`} value={`${path}-${method}-${index}`}>
                                            <AccordionTrigger className="hover:no-underline">
                                                <div className="flex items-center gap-3 w-full">
                                                    <Badge className={`${getMethodColor(method)} border`}>
                                                        {method}
                                                    </Badge>
                                                    <code className="text-sm font-mono bg-muted px-2 py-1 rounded">
                                                        {path}
                                                    </code>
                                                    <span className="text-sm font-medium text-muted-foreground">
                                                        {endpoint.summary}
                                                    </span>
                                                </div>
                                            </AccordionTrigger>
                                            <AccordionContent>
                                                <div className="space-y-4">
                                                    <div className="flex items-center justify-between">
                                                        <div>
                                                            <h4 className="text-lg font-semibold">{endpoint.summary}</h4>
                                                            <p className="text-sm text-muted-foreground mt-1">{endpoint.description}</p>
                                                        </div>
                                                        <div className="flex gap-2">
                                                            <Button
                                                                variant="outline"
                                                                size="sm"
                                                                onClick={() => copyToClipboard(`${baseUrl}${path}`, 'Endpoint URL')}
                                                            >
                                                                <Copy className="h-4 w-4" />
                                                            </Button>
                                                            <Button
                                                                variant="outline"
                                                                size="sm"
                                                                onClick={() => window.open(`http://localhost:8080/swagger-ui#${path.replace(/\//g, '')}`, '_blank')}
                                                            >
                                                                <ExternalLink className="h-4 w-4" />
                                                            </Button>
                                                        </div>
                                                    </div>

                                                    {/* Test Interface */}
                                                    {renderTestInterface(path, method, endpoint)}

                                                    {/* Schema Display */}
                                                    {renderSchemaDisplay(endpoint, path)}

                                                    {showCodeExamples && (
                                                        <Tabs defaultValue="curl" className="w-full">
                                                            <TabsList className="grid w-full grid-cols-2">
                                                                <TabsTrigger value="curl">cURL</TabsTrigger>
                                                                <TabsTrigger value="fetch">JavaScript</TabsTrigger>
                                                            </TabsList>
                                                            <TabsContent value="curl" className="mt-4">
                                                                <div className="space-y-2">
                                                                    <div className="flex items-center justify-between">
                                                                        <label className="text-sm font-medium">cURL Example</label>
                                                                        <Button
                                                                            variant="outline"
                                                                            size="sm"
                                                                            onClick={() => copyToClipboard(codeSnippets.curlExample, 'cURL command')}
                                                                        >
                                                                            <Copy className="h-4 w-4" />
                                                                        </Button>
                                                                    </div>
                                                                    <Textarea
                                                                        value={codeSnippets.curlExample}
                                                                        readOnly
                                                                        className="font-mono text-sm"
                                                                        rows={4}
                                                                    />
                                                                </div>
                                                            </TabsContent>
                                                            <TabsContent value="fetch" className="mt-4">
                                                                <div className="space-y-2">
                                                                    <div className="flex items-center justify-between">
                                                                        <label className="text-sm font-medium">JavaScript Example</label>
                                                                        <Button
                                                                            variant="outline"
                                                                            size="sm"
                                                                            onClick={() => copyToClipboard(codeSnippets.fetchExample, 'JavaScript code')}
                                                                        >
                                                                            <Copy className="h-4 w-4" />
                                                                        </Button>
                                                                    </div>
                                                                    <Textarea
                                                                        value={codeSnippets.fetchExample}
                                                                        readOnly
                                                                        className="font-mono text-sm"
                                                                        rows={6}
                                                                    />
                                                                </div>
                                                            </TabsContent>
                                                        </Tabs>
                                                    )}
                                                </div>
                                            </AccordionContent>
                                        </AccordionItem>
                                    );
                                })}
                            </Accordion>
                        </CardContent>
                    </Card>
                ))}
            </div>

            {Object.keys(groupedEndpoints).length === 0 && (
                <Card>
                    <CardContent className="text-center py-12">
                        <Search className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
                        <h3 className="text-lg font-semibold mb-2">No endpoints found</h3>
                        <p className="text-muted-foreground">
                            Try adjusting your search criteria or filters.
                        </p>
                    </CardContent>
                </Card>
            )}
        </div>
    );
}
