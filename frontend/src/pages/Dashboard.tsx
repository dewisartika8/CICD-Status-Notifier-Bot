import React, { useEffect } from 'react';
import {
  Box,
  Grid,
  Typography,
  Paper,
  Card,
  CardContent,
  List,
  ListItem,
  ListItemText,
  ListItemIcon,
  Chip,
  LinearProgress,
  Alert,
  Skeleton,
} from '@mui/material';
import {
  TrendingUp as TrendingUpIcon,
  TrendingDown as TrendingDownIcon,
  CheckCircle as SuccessIcon,
  Error as ErrorIcon,
  Build as BuildIcon,
  Schedule as ScheduleIcon,
} from '@mui/icons-material';
import { useQuery } from 'react-query';
import { useAppDispatch, useAppSelector } from '@/store';
import { fetchProjects } from '@/store/slices/projectSlice';
import { dashboardApi, projectApi } from '@/services/api';
import { DashboardStats } from '@/types';

// Mock data untuk development
const mockStats: DashboardStats[] = [
  {
    title: 'Total Projects',
    value: 12,
    change: 2,
    changeType: 'increase',
    icon: 'folder',
    color: 'primary',
  },
  {
    title: 'Successful Builds',
    value: '89%',
    change: 5,
    changeType: 'increase',
    icon: 'success',
    color: 'success',
  },
  {
    title: 'Failed Builds',
    value: 8,
    change: -3,
    changeType: 'decrease',
    icon: 'error',
    color: 'error',
  },
  {
    title: 'Avg Build Time',
    value: '2m 34s',
    change: -15,
    changeType: 'decrease',
    icon: 'timer',
    color: 'secondary',
  },
];

const Dashboard: React.FC = () => {
  const dispatch = useAppDispatch();
  const { projects, status: projectsStatus, error } = useAppSelector((state) => state.projects);

  // Fetch dashboard overview
  const { data: overview, isLoading: overviewLoading, error: overviewError } = useQuery(
    'dashboard-overview',
    () => dashboardApi.getOverview(),
    {
      refetchInterval: 30000, // Refresh every 30 seconds
      retry: 2,
    }
  );

  // Fetch projects if not already loaded
  useEffect(() => {
    if (projectsStatus === 'idle') {
      dispatch(fetchProjects());
    }
  }, [dispatch, projectsStatus]);

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'success':
        return 'success';
      case 'failed':
        return 'error';
      case 'building':
        return 'warning';
      case 'pending':
        return 'default';
      default:
        return 'default';
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'success':
        return <SuccessIcon color="success" />;
      case 'failed':
        return <ErrorIcon color="error" />;
      case 'building':
        return <BuildIcon color="warning" />;
      case 'pending':
        return <ScheduleIcon color="disabled" />;
      default:
        return <BuildIcon color="disabled" />;
    }
  };

  if (overviewError || error) {
    return (
      <Box>
        <Alert severity="error" sx={{ mb: 2 }}>
          Failed to load dashboard data. Please try refreshing the page.
        </Alert>
      </Box>
    );
  }

  return (
    <Box>
      {/* Page Header */}
      <Box sx={{ mb: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom fontWeight="bold">
          Dashboard Overview
        </Typography>
        <Typography variant="body1" color="text.secondary">
          Monitor your CI/CD pipelines and project status
        </Typography>
      </Box>

      {/* Statistics Cards */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        {mockStats.map((stat, index) => (
          <Grid item xs={12} sm={6} md={3} key={index}>
            <Card
              sx={{
                height: '100%',
                background: (theme) =>
                  theme.palette.mode === 'dark'
                    ? 'linear-gradient(135deg, rgba(255,255,255,0.1) 0%, rgba(255,255,255,0.05) 100%)'
                    : 'linear-gradient(135deg, rgba(0,0,0,0.02) 0%, rgba(0,0,0,0.01) 100%)',
                border: '1px solid',
                borderColor: 'divider',
              }}
            >
              <CardContent>
                <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                  <Box>
                    <Typography color="text.secondary" gutterBottom variant="overline">
                      {stat.title}
                    </Typography>
                    <Typography variant="h4" component="div" fontWeight="bold">
                      {overviewLoading ? <Skeleton width={60} /> : stat.value}
                    </Typography>
                    {stat.change && (
                      <Box sx={{ display: 'flex', alignItems: 'center', mt: 1 }}>
                        {stat.changeType === 'increase' ? (
                          <TrendingUpIcon color="success" fontSize="small" />
                        ) : (
                          <TrendingDownIcon color="error" fontSize="small" />
                        )}
                        <Typography
                          variant="body2"
                          color={stat.changeType === 'increase' ? 'success.main' : 'error.main'}
                          sx={{ ml: 0.5 }}
                        >
                          {Math.abs(stat.change)}%
                        </Typography>
                      </Box>
                    )}
                  </Box>
                </Box>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>

      <Grid container spacing={3}>
        {/* Project Status List */}
        <Grid item xs={12} md={8}>
          <Paper sx={{ p: 3, height: 'fit-content' }}>
            <Typography variant="h6" gutterBottom fontWeight="bold">
              Recent Projects
            </Typography>
            {projectsStatus === 'loading' ? (
              <Box>
                {[...Array(5)].map((_, index) => (
                  <Box key={index} sx={{ mb: 2 }}>
                    <Skeleton variant="rectangular" height={60} />
                  </Box>
                ))}
              </Box>
            ) : projects.length > 0 ? (
              <List>
                {projects.slice(0, 8).map((project) => (
                  <ListItem
                    key={project.id}
                    divider
                    sx={{
                      border: '1px solid',
                      borderColor: 'divider',
                      borderRadius: 1,
                      mb: 1,
                      '&:hover': {
                        bgcolor: 'action.hover',
                      },
                    }}
                  >
                    <ListItemIcon>{getStatusIcon(project.status)}</ListItemIcon>
                    <ListItemText
                      primary={
                        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                          <Typography variant="subtitle1" fontWeight="medium">
                            {project.name}
                          </Typography>
                          <Chip
                            label={project.status}
                            size="small"
                            color={getStatusColor(project.status) as any}
                            variant="outlined"
                          />
                        </Box>
                      }
                      secondary={
                        <Box sx={{ mt: 1 }}>
                          <Typography variant="body2" color="text.secondary">
                            Last build: {project.lastBuild?.timestamp 
                              ? new Date(project.lastBuild.timestamp).toLocaleString()
                              : 'Never'
                            }
                          </Typography>
                          {project.buildHistory?.some(b => b.status === 'building') && (
                            <LinearProgress sx={{ mt: 1, height: 4, borderRadius: 2 }} />
                          )}
                        </Box>
                      }
                    />
                  </ListItem>
                ))}
              </List>
            ) : (
              <Alert severity="info">
                No projects found. Add your first project to get started!
              </Alert>
            )}
          </Paper>
        </Grid>

        {/* Activity Timeline */}
        <Grid item xs={12} md={4}>
          <Paper sx={{ p: 3, height: 'fit-content' }}>
            <Typography variant="h6" gutterBottom fontWeight="bold">
              Recent Activity
            </Typography>
            <List dense>
              {[...Array(6)].map((_, index) => (
                <ListItem key={index} sx={{ px: 0 }}>
                  <ListItemIcon sx={{ minWidth: 32 }}>
                    <Box
                      sx={{
                        width: 8,
                        height: 8,
                        borderRadius: '50%',
                        bgcolor: index % 2 === 0 ? 'success.main' : 'error.main',
                      }}
                    />
                  </ListItemIcon>
                  <ListItemText
                    primary={
                      <Typography variant="body2">
                        {index % 2 === 0 ? 'Build succeeded' : 'Build failed'} for Project {index + 1}
                      </Typography>
                    }
                    secondary={
                      <Typography variant="caption" color="text.secondary">
                        {index + 1} hours ago
                      </Typography>
                    }
                  />
                </ListItem>
              ))}
            </List>
          </Paper>
        </Grid>
      </Grid>
    </Box>
  );
};

export default Dashboard;