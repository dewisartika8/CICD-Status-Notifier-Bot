import axios, { AxiosResponse } from 'axios';
import { 
  Project, 
  ProjectResponse,
  ListProjectResponse,
  CreateProjectRequest,
  UpdateProjectRequest,
  ProjectStatusUpdateRequest,
  ProjectStatisticsResponse,
  BuildEvent, 
  DashboardOverview, 
  Metrics, 
  ProjectStatus,
  ApiResponse,
  PaginatedResponse,
  ProjectFilter,
  BuildFilter,
  TelegramSubscriptionResponse,
  TelegramSubscriptionListResponse,
  CreateTelegramSubscriptionRequest,
  UpdateTelegramSubscriptionRequest,
  WebhookEventResponse
} from '@/types';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';

// Create axios instance with default config
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor for adding auth token
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('authToken');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor for handling errors
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('authToken');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Dashboard API - aligned with OpenAPI specification
export const dashboardApi = {
  // GET /api/v1/dashboard/overview
  getOverview: (): Promise<AxiosResponse<{ message: string; data: DashboardOverview }>> => 
    apiClient.get('/api/v1/dashboard/overview'),
  
  // GET /health
  getSystemHealth: (): Promise<AxiosResponse<{ status: string; database: string; timestamp: string }>> =>
    apiClient.get('/health'),
};

// Projects API - aligned with OpenAPI specification
export const projectApi = {
  // GET /api/v1/projects
  getProjects: (filter?: ProjectFilter): Promise<AxiosResponse<ListProjectResponse>> => 
    apiClient.get('/api/v1/projects', { params: filter }),
  
  // GET /api/v1/projects/{id}
  getProjectById: (projectId: string): Promise<AxiosResponse<{ message: string; data: ProjectResponse }>> =>
    apiClient.get(`/api/v1/projects/${projectId}`),
  
  // GET /api/v1/projects/{id}
  getProjectDetails: (projectId: string): Promise<AxiosResponse<{ message: string; data: ProjectResponse }>> => 
    apiClient.get(`/api/v1/projects/${projectId}`),
  
  // GET /api/v1/projects/{id}/status
  getProjectStatus: (projectId: string): Promise<AxiosResponse<{ message: string; data: { status: ProjectStatus } }>> =>
    apiClient.get(`/api/v1/projects/${projectId}/status`),
  
  // GET /api/v1/projects/{id}/statistics
  getProjectStatistics: (projectId: string): Promise<AxiosResponse<{ message: string; data: ProjectStatisticsResponse }>> =>
    apiClient.get(`/api/v1/projects/${projectId}/statistics`),
  
  // POST /api/v1/projects
  createProject: (project: CreateProjectRequest): Promise<AxiosResponse<{ message: string; data: ProjectResponse }>> =>
    apiClient.post('/api/v1/projects', project),
  
  // PUT /api/v1/projects/{id}
  updateProject: (projectId: string, project: UpdateProjectRequest): Promise<AxiosResponse<{ message: string; data: ProjectResponse }>> =>
    apiClient.put(`/api/v1/projects/${projectId}`, project),
  
  // PATCH /api/v1/projects/{id}/status
  updateProjectStatus: (projectId: string, status: ProjectStatusUpdateRequest): Promise<AxiosResponse<{ message: string; data: ProjectResponse }>> =>
    apiClient.patch(`/api/v1/projects/${projectId}/status`, status),
  
  // DELETE /api/v1/projects/{id}
  deleteProject: (projectId: string): Promise<AxiosResponse<{ message: string }>> =>
    apiClient.delete(`/api/v1/projects/${projectId}`),
};

// Webhook Events API - aligned with OpenAPI specification
export const webhookApi = {
  // GET /api/v1/webhooks/events/{projectId}
  getWebhookEventsByProject: (
    projectId: string, 
    params?: { limit?: number; offset?: number }
  ): Promise<AxiosResponse<{ message: string; data: WebhookEventResponse[] }>> =>
    apiClient.get(`/api/v1/webhooks/events/${projectId}`, { params }),
  
  // GET /api/v1/webhooks/events/{projectId}/{eventId} 
  getBuildById: (projectId: string, buildId: string): Promise<AxiosResponse<{ message: string; data: WebhookEventResponse }>> =>
    apiClient.get(`/api/v1/webhooks/events/${projectId}/${buildId}`),
};

// Telegram Bot API - aligned with OpenAPI specification  
export const telegramApi = {
  // GET /api/v1/telegram/subscriptions
  getSubscriptions: (params?: { page?: number; limit?: number }): Promise<AxiosResponse<TelegramSubscriptionListResponse>> =>
    apiClient.get('/api/v1/telegram/subscriptions', { params }),
  
  // POST /api/v1/telegram/subscriptions
  createSubscription: (subscription: CreateTelegramSubscriptionRequest): Promise<AxiosResponse<{ message: string; data: TelegramSubscriptionResponse }>> =>
    apiClient.post('/api/v1/telegram/subscriptions', subscription),
  
  // GET /api/v1/telegram/subscriptions/{id}
  getSubscriptionById: (id: string): Promise<AxiosResponse<{ message: string; data: TelegramSubscriptionResponse }>> =>
    apiClient.get(`/api/v1/telegram/subscriptions/${id}`),
  
  // PUT /api/v1/telegram/subscriptions/{id}
  updateSubscription: (id: string, subscription: UpdateTelegramSubscriptionRequest): Promise<AxiosResponse<{ message: string; data: TelegramSubscriptionResponse }>> =>
    apiClient.put(`/api/v1/telegram/subscriptions/${id}`, subscription),
  
  // DELETE /api/v1/telegram/subscriptions/{id}
  deleteSubscription: (id: string): Promise<AxiosResponse<{ message: string }>> =>
    apiClient.delete(`/api/v1/telegram/subscriptions/${id}`),
  
  // GET /api/v1/telegram/projects/{projectId}/subscriptions
  getProjectSubscriptions: (projectId: string, params?: { page?: number; limit?: number }): Promise<AxiosResponse<TelegramSubscriptionListResponse>> =>
    apiClient.get(`/api/v1/telegram/projects/${projectId}/subscriptions`, { params }),
  
  // POST /api/v1/telegram/webhook/set
  setWebhook: (webhook_url: string): Promise<AxiosResponse<{ message: string; data: any }>> =>
    apiClient.post('/api/v1/telegram/webhook/set', { webhook_url }),
  
  // DELETE /api/v1/telegram/webhook
  deleteWebhook: (): Promise<AxiosResponse<{ message: string }>> =>
    apiClient.delete('/api/v1/telegram/webhook'),
};

// Build Events API (derived from webhook events)
export const buildApi = {
  // Get build history via webhook events
  getBuildHistory: (
    projectId: string, 
    params?: { limit?: number; offset?: number }
  ): Promise<AxiosResponse<ApiResponse<BuildEvent[]>>> =>
    apiClient.get(`/api/v1/webhooks/events/${projectId}?event_type=workflow_run`, { params }),
  
  // Get specific build event
  getBuildById: (projectId: string, buildId: string): Promise<AxiosResponse<ApiResponse<BuildEvent>>> =>
    apiClient.get(`/api/v1/webhooks/events/${projectId}/${buildId}`),
};

// Metrics API  
export const metricsApi = {
  getAverageBuildDuration: (timeRange?: string): Promise<AxiosResponse<ApiResponse<any>>> =>
    apiClient.get('/metrics/average-build-duration', { params: { timeRange } }),
  
  getBuildTrends: (timeRange?: string): Promise<AxiosResponse<ApiResponse<any>>> =>
    apiClient.get('/metrics/build-trends', { params: { timeRange } }),
  
  getProjectMetrics: (projectId: string): Promise<AxiosResponse<ApiResponse<Metrics>>> =>
    apiClient.get(`/metrics/projects/${projectId}`),
  
  getDeploymentFrequency: (timeRange?: string): Promise<AxiosResponse<ApiResponse<any>>> =>
    apiClient.get('/metrics/deployment-frequency', { params: { timeRange } }),
};

// Notifications API
export const notificationApi = {
  getSubscriptions: (): Promise<AxiosResponse<ApiResponse<any[]>>> =>
    apiClient.get('/notifications/subscriptions'),
  
  subscribe: (projectId: string, preferences: any): Promise<AxiosResponse<ApiResponse<any>>> =>
    apiClient.post('/notifications/subscribe', { projectId, preferences }),
  
  unsubscribe: (subscriptionId: string): Promise<AxiosResponse<ApiResponse<void>>> =>
    apiClient.delete(`/notifications/subscriptions/${subscriptionId}`),
  
  updatePreferences: (subscriptionId: string, preferences: any): Promise<AxiosResponse<ApiResponse<any>>> =>
    apiClient.put(`/notifications/subscriptions/${subscriptionId}`, { preferences }),
};

// Health check API
export const healthApi = {
  checkHealth: (): Promise<AxiosResponse<ApiResponse<any>>> =>
    apiClient.get('/health'),
  
  checkDatabaseHealth: (): Promise<AxiosResponse<ApiResponse<any>>> =>
    apiClient.get('/health/database'),
  
  checkServicesHealth: (): Promise<AxiosResponse<ApiResponse<any>>> =>
    apiClient.get('/health/services'),
};

// Export axios instance for custom requests
export { apiClient };

// Legacy exports for backward compatibility
export const getDashboardOverview = () => dashboardApi.getOverview();
export const getProjects = () => projectApi.getProjects();
export const getProjectDetails = (projectId: string) => projectApi.getProjectDetails(projectId);
export const getProjectStatus = (projectId: string) => projectApi.getProjectStatus(projectId);
export const getAverageBuildDuration = () => metricsApi.getAverageBuildDuration();
export const getBuildTrends = () => metricsApi.getBuildTrends();

// Additional legacy aliases for components
export const fetchDashboardOverview = getDashboardOverview;
export const fetchProjects = getProjects;
export const fetchProjectDetails = getProjectDetails;