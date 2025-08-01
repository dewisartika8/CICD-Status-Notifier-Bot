// src/types/index.ts

export interface Project {
    id: string;
    name: string;
    repository_url: string;
    status: ProjectStatus;
    telegram_chat_id?: number;
    created_at: string;
    updated_at: string;
    
    // Extended properties for frontend display
    lastBuild?: BuildEvent;
    buildHistory: BuildEvent[];
    description?: string;
    branch?: string;
    isActive: boolean;
    lastBuildNumber?: number;
    lastUpdated?: string;
}

export type ProjectStatus = 'active' | 'inactive' | 'archived';

// OpenAPI ProjectResponse
export interface ProjectResponse {
    id: string;
    name: string;
    repository_url: string;
    status: ProjectStatus;
    telegram_chat_id?: number;
    created_at: string;
    updated_at: string;
}

export interface ListProjectResponse {
    projects: ProjectResponse[];
    total: number;
    limit: number;
    offset: number;
}

export interface BuildEvent {
    id: string;
    projectId: string;
    status: 'pending' | 'building' | 'success' | 'failed' | 'cancelled';
    duration?: number; // in seconds
    timestamp: string; // ISO date string
    commitHash?: string;
    commitMessage?: string;
    author?: string;
    branch?: string;
    buildNumber?: number;
    logUrl?: string;
    errorMessage?: string;
}

// API Request types - matching OpenAPI specification
export interface CreateProjectRequest {
    name: string;
    repository_url: string;
    webhook_secret: string;
    telegram_chat_id?: number;
}

export interface UpdateProjectRequest {
    name?: string;
    repository_url?: string;
    webhook_secret?: string;
    telegram_chat_id?: number;
}

export interface ProjectStatusUpdateRequest {
    status: 'active' | 'inactive' | 'archived';
}

export interface ProjectStatistics {
    project_id: string;
    project_name: string;
    total_builds: number;
    successful_builds: number;
    failed_builds: number;
    success_rate: number;
    average_duration: number;
    last_build_time: string;
    build_trends: BuildTrendData[];
}

export interface BuildTrendData {
    date: string;
    count: number;
    success_count: number;
}

export interface BuildAnalytics {
    time_range: string;
    total_builds: number;
    successful_builds: number;
    failed_builds: number;
    success_rate: number;
    average_duration: number;
    builds_by_day: BuildAnalyticsData[];
    duration_trends: DurationTrendData[];
}

export interface BuildAnalyticsData {
    date: string;
    total_builds: number;
    successful_builds: number;
    failed_builds: number;
}

export interface DurationTrendData {
    date: string;
    average_duration: number;
}

// Telegram API types - matching OpenAPI specification
export interface TelegramSubscription {
    id: string;
    project_id: string;
    chat_id: number;
    is_active: boolean;
    created_at: string;
    updated_at: string;
}

export interface CreateTelegramSubscriptionRequest {
    project_id: string;
    chat_id: number;
}

export interface TelegramSubscriptionResponse {
    id: string;
    project_id: string;
    chat_id: number;
    is_active: boolean;
    created_at: string;
    updated_at: string;
}

export interface TelegramSubscriptionListResponse {
    data: TelegramSubscriptionResponse[];
    total: number;
    page: number;
    limit: number;
    total_pages: number;
}

export interface UpdateTelegramSubscriptionRequest {
    chat_id?: number;
    is_active?: boolean;
}

// Webhook Event types - matching OpenAPI specification
export interface WebhookEvent {
    id: string;
    project_id: string;
    event_type: 'workflow_run' | 'push' | 'pull_request';
    delivery_id: string;
    processed_at: string | null;
    created_at: string;
}

export interface WebhookEventResponse {
    id: string;
    project_id: string;
    event_type: 'workflow_run' | 'push' | 'pull_request';
    delivery_id: string;
    processed_at: string | null;
    created_at: string;
}

// Dashboard Overview from OpenAPI OverviewResponse
export interface DashboardOverview {
    total_projects: number;
    active_projects: number;
    total_builds: number;
    successful_builds: number;
    failed_builds: number;
    success_rate: number;
    average_duration: number;
    last_updated: string;
}

// Project Statistics Response from OpenAPI spec
export interface ProjectStatisticsResponse {
    project_id: string;
    project_name: string;
    total_builds: number;
    successful_builds: number;
    failed_builds: number;
    success_rate: number;
    average_duration: number;
    recent_builds: {
        run_number: number;
        status: string;
        conclusion: string;
        created_at: string;
        duration: number;
    }[];
    build_trends: {
        date: string;
        successful: number;
        failed: number;
    }[];
}export interface Metrics {
    successRate: number; // percentage
    averageDuration: number; // in seconds
    deploymentFrequency: number; // number of deployments over a period
    buildsToday: number;
    buildsThisWeek: number;
    buildsThisMonth: number;
}

export interface ChartData {
    name: string;
    value: number;
    date?: string;
}

export interface Subscription {
    id: string;
    userId: string;
    projectId: string;
    preferences: NotificationPreference[];
    isActive: boolean;
    createdAt: string;
}

export interface NotificationPreference {
    type: 'success' | 'failure' | 'started' | 'cancelled';
    enabled: boolean;
}

export interface User {
    id: string;
    username: string;
    email?: string;
    telegramUserId?: string;
    role: 'admin' | 'user' | 'viewer';
    isActive: boolean;
    createdAt: string;
    lastActiveAt?: string;
}

export interface WebSocketMessage {
    type: 'project_update' | 'build_event' | 'notification';
    payload: any;
    timestamp: string;
}

export interface ApiResponse<T> {
    data: T;
    message: string;
    success: boolean;
    errors?: string[];
}

export interface PaginatedResponse<T> {
    data: T[];
    pagination: {
        page: number;
        limit: number;
        total: number;
        totalPages: number;
    };
}

export interface ApiError {
    message: string;
    code?: string;
    field?: string;
}

// Theme types
export interface ThemeConfig {
    mode: 'light' | 'dark';
    primaryColor: string;
    secondaryColor: string;
}

// Dashboard layout types
export interface DashboardStats {
    title: string;
    value: number | string;
    change?: number;
    changeType?: 'increase' | 'decrease';
    icon?: string;
    color?: 'primary' | 'secondary' | 'success' | 'warning' | 'error';
}

// Filter and search types
export interface ProjectFilter {
    status?: ProjectStatus[];
    dateRange?: {
        start: string;
        end: string;
    };
    search?: string;
}

export interface BuildFilter {
    status?: ProjectStatus[];
    projectIds?: string[];
    dateRange?: {
        start: string;
        end: string;
    };
    author?: string;
    branch?: string;
}