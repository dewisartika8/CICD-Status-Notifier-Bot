import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { projectApi } from '@/services/api';
import { Project, ProjectStatus, BuildEvent } from '@/types';

export interface ProjectState {
  projects: Project[];
  currentProject: Project | null;
  status: 'idle' | 'loading' | 'succeeded' | 'failed';
  error: string | null;
  lastUpdated: string | null;
}

const initialState: ProjectState = {
  projects: [],
  currentProject: null,
  status: 'idle',
  error: null,
  lastUpdated: null,
};

// Async thunks
export const fetchProjects = createAsyncThunk(
  'projects/fetchProjects',
  async (_, { rejectWithValue }) => {
    try {
      const response = await projectApi.getProjects();
      return response.data;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.message || 'Failed to fetch projects');
    }
  }
);

export const fetchProjectById = createAsyncThunk(
  'projects/fetchProjectById',
  async (projectId: string, { rejectWithValue }) => {
    try {
      const response = await projectApi.getProjectById(projectId);
      return response.data;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.message || 'Failed to fetch project');
    }
  }
);

export const updateProjectStatus = createAsyncThunk(
  'projects/updateProjectStatus',
  async ({ projectId, status }: { projectId: string; status: ProjectStatus }, { rejectWithValue }) => {
    try {
      const response = await projectApi.updateProjectStatus(projectId, { status });
      return response.data;
    } catch (error: any) {
      return rejectWithValue(error.response?.data?.message || 'Failed to update project status');
    }
  }
);

const projectSlice = createSlice({
  name: 'projects',
  initialState,
  reducers: {
    clearError: (state) => {
      state.error = null;
    },
    setCurrentProject: (state, action: PayloadAction<Project | null>) => {
      state.currentProject = action.payload;
    },
    updateProjectBuildEvent: (state, action: PayloadAction<{ projectId: string; buildEvent: BuildEvent }>) => {
      const { projectId, buildEvent } = action.payload;
      const project = state.projects.find(p => p.id === projectId);
      if (project) {
        project.lastBuild = buildEvent;
        // Only update status if it's a valid ProjectStatus
        if (['active', 'inactive', 'archived'].includes(buildEvent.status as any)) {
          project.status = buildEvent.status as ProjectStatus;
        }
        project.updated_at = new Date().toISOString();
      }
      if (state.currentProject?.id === projectId) {
        state.currentProject.lastBuild = buildEvent;
        // Only update status if it's a valid ProjectStatus
        if (state.currentProject && ['active', 'inactive', 'archived'].includes(buildEvent.status as any)) {
          state.currentProject.status = buildEvent.status as ProjectStatus;
        }
        state.currentProject.updated_at = new Date().toISOString();
      }
      state.lastUpdated = new Date().toISOString();
    },
    addProject: (state, action: PayloadAction<Project>) => {
      state.projects.push(action.payload);
      state.lastUpdated = new Date().toISOString();
    },
    removeProject: (state, action: PayloadAction<string>) => {
      state.projects = state.projects.filter(p => p.id !== action.payload);
      if (state.currentProject?.id === action.payload) {
        state.currentProject = null;
      }
      state.lastUpdated = new Date().toISOString();
    },
  },
  extraReducers: (builder) => {
    builder
      // Fetch projects
      .addCase(fetchProjects.pending, (state) => {
        state.status = 'loading';
        state.error = null;
      })
      .addCase(fetchProjects.fulfilled, (state, action) => {
        state.status = 'succeeded';
        const projects = action.payload.projects || action.payload;
        state.projects = projects.map((p: any) => ({
          ...p,
          buildHistory: p.buildHistory || [],
          isActive: p.status === 'active'
        }));
        state.lastUpdated = new Date().toISOString();
      })
      .addCase(fetchProjects.rejected, (state, action) => {
        state.status = 'failed';
        state.error = action.payload as string;
      })
      // Fetch project by ID
      .addCase(fetchProjectById.pending, (state) => {
        state.status = 'loading';
        state.error = null;
      })
      .addCase(fetchProjectById.fulfilled, (state, action) => {
        state.status = 'succeeded';
        const projectData = action.payload.data || action.payload;
        const mappedProject = {
          ...projectData,
          buildHistory: [],
          isActive: projectData.status === 'active'
        };
        state.currentProject = mappedProject;
        // Update project in list if exists
        const index = state.projects.findIndex(p => p.id === projectData.id);
        if (index !== -1) {
          const mappedProject = {
            ...projectData,
            buildHistory: state.projects[index]?.buildHistory || [],
            isActive: projectData.status === 'active'
          };
          state.projects[index] = mappedProject;
        }
        state.lastUpdated = new Date().toISOString();
      })
      .addCase(fetchProjectById.rejected, (state, action) => {
        state.status = 'failed';
        state.error = action.payload as string;
      })
      // Update project status
      .addCase(updateProjectStatus.fulfilled, (state, action) => {
        const updatedProject = action.payload.data || action.payload;
        const index = state.projects.findIndex(p => p.id === updatedProject.id);
        if (index !== -1) {
          const mappedUpdatedProject = {
            ...updatedProject,
            buildHistory: state.projects[index]?.buildHistory || [],
            isActive: updatedProject.status === 'active'
          };
          state.projects[index] = mappedUpdatedProject;
        }
        if (state.currentProject?.id === updatedProject.id) {
          const mappedCurrentProject = {
            ...updatedProject,
            buildHistory: state.currentProject?.buildHistory || [],
            isActive: updatedProject.status === 'active'
          };
          state.currentProject = mappedCurrentProject;
        }
        state.lastUpdated = new Date().toISOString();
      });
  },
});

export const {
  clearError,
  setCurrentProject,
  updateProjectBuildEvent,
  addProject,
  removeProject,
} = projectSlice.actions;

export default projectSlice.reducer;
