import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { ThemeConfig } from '@/types';

export interface UIState {
  theme: ThemeConfig;
  sidebarOpen: boolean;
  sidebarCollapsed: boolean;
  loading: {
    dashboard: boolean;
    projects: boolean;
    builds: boolean;
  };
  errors: {
    dashboard: string | null;
    projects: string | null;
    builds: string | null;
  };
  selectedProjectId: string | null;
  websocketConnected: boolean;
  lastWebsocketMessage: string | null;
  refreshInterval: number; // in seconds
  autoRefresh: boolean;
}

const initialState: UIState = {
  theme: {
    mode: 'light',
    primaryColor: '#1976d2',
    secondaryColor: '#dc004e',
  },
  sidebarOpen: true,
  sidebarCollapsed: false,
  loading: {
    dashboard: false,
    projects: false,
    builds: false,
  },
  errors: {
    dashboard: null,
    projects: null,
    builds: null,
  },
  selectedProjectId: null,
  websocketConnected: false,
  lastWebsocketMessage: null,
  refreshInterval: 30,
  autoRefresh: true,
};

const uiSlice = createSlice({
  name: 'ui',
  initialState,
  reducers: {
    toggleTheme: (state) => {
      state.theme.mode = state.theme.mode === 'light' ? 'dark' : 'light';
    },
    setTheme: (state, action: PayloadAction<ThemeConfig>) => {
      state.theme = action.payload;
    },
    toggleSidebar: (state) => {
      state.sidebarOpen = !state.sidebarOpen;
    },
    setSidebarOpen: (state, action: PayloadAction<boolean>) => {
      state.sidebarOpen = action.payload;
    },
    toggleSidebarCollapsed: (state) => {
      state.sidebarCollapsed = !state.sidebarCollapsed;
    },
    setSidebarCollapsed: (state, action: PayloadAction<boolean>) => {
      state.sidebarCollapsed = action.payload;
    },
    setLoading: (state, action: PayloadAction<{ key: keyof UIState['loading']; value: boolean }>) => {
      state.loading[action.payload.key] = action.payload.value;
    },
    setError: (state, action: PayloadAction<{ key: keyof UIState['errors']; value: string | null }>) => {
      state.errors[action.payload.key] = action.payload.value;
    },
    clearError: (state, action: PayloadAction<keyof UIState['errors']>) => {
      state.errors[action.payload] = null;
    },
    clearAllErrors: (state) => {
      state.errors = {
        dashboard: null,
        projects: null,
        builds: null,
      };
    },
    setSelectedProjectId: (state, action: PayloadAction<string | null>) => {
      state.selectedProjectId = action.payload;
    },
    setWebsocketConnected: (state, action: PayloadAction<boolean>) => {
      state.websocketConnected = action.payload;
    },
    setLastWebsocketMessage: (state, action: PayloadAction<string | null>) => {
      state.lastWebsocketMessage = action.payload;
    },
    setRefreshInterval: (state, action: PayloadAction<number>) => {
      state.refreshInterval = action.payload;
    },
    setAutoRefresh: (state, action: PayloadAction<boolean>) => {
      state.autoRefresh = action.payload;
    },
  },
});

export const {
  toggleTheme,
  setTheme,
  toggleSidebar,
  setSidebarOpen,
  toggleSidebarCollapsed,
  setSidebarCollapsed,
  setLoading,
  setError,
  clearError,
  clearAllErrors,
  setSelectedProjectId,
  setWebsocketConnected,
  setLastWebsocketMessage,
  setRefreshInterval,
  setAutoRefresh,
} = uiSlice.actions;

export default uiSlice.reducer;
