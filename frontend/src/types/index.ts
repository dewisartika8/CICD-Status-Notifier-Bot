// src/types/index.ts

export interface Project {
    id: string;
    name: string;
    status: string;
    buildHistory: Build[];
}

export interface Build {
    id: string;
    projectId: string;
    status: string;
    duration: number; // in seconds
    timestamp: string; // ISO date string
}

export interface DashboardOverview {
    totalProjects: number;
    successfulBuilds: number;
    failedBuilds: number;
    averageBuildDuration: number; // in seconds
}

export interface Metrics {
    successRate: number; // percentage
    averageDuration: number; // in seconds
    deploymentFrequency: number; // number of deployments over a period
}

export interface Subscription {
    id: string;
    userId: string;
    projectId: string;
    preferences: string[]; // array of preferences
}

export interface ApiResponse<T> {
    data: T;
    message: string;
    success: boolean;
}