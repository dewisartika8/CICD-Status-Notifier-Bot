import React, { useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Box,
  Typography,
  Paper,
  Grid,
  Card,
  CardContent,
  List,
  ListItem,
  ListItemText,
  ListItemIcon,
  Chip,
  Button,
  IconButton,
  Skeleton,
  Alert,
  Divider,
  LinearProgress,
} from '@mui/material';
import {
  ArrowBack as ArrowBackIcon,
  Build as BuildIcon,
  CheckCircle as SuccessIcon,
  Error as ErrorIcon,
  Schedule as ScheduleIcon,
  Refresh as RefreshIcon,
  OpenInNew as OpenInNewIcon,
} from '@mui/icons-material';
import { useQuery } from 'react-query';
import { useAppDispatch, useAppSelector } from '@/store';
import { fetchProjectById } from '@/store/slices/projectSlice';
import { buildApi } from '@/services/api';

const ProjectDetails: React.FC = () => {
  const { projectId } = useParams<{ projectId: string }>();
  const navigate = useNavigate();
  const dispatch = useAppDispatch();
  const { currentProject, status } = useAppSelector((state) => state.projects);

  // Fetch project details
  useEffect(() => {
    if (projectId) {
      dispatch(fetchProjectById(projectId));
    }
  }, [dispatch, projectId]);

  // Fetch build history
  const { data: buildHistory, isLoading: buildsLoading, refetch: refetchBuilds } = useQuery(
    ['builds', projectId],
    () => buildApi.getBuildHistory(projectId!),
    {
      enabled: !!projectId,
      refetchInterval: 30000,
    }
  );

  const handleBack = () => {
    navigate('/projects');
  };

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

  if (status === 'loading') {
    return (
      <Box>
        <Skeleton variant="rectangular" height={200} sx={{ mb: 3 }} />
        <Grid container spacing={3}>
          <Grid item xs={12} md={8}>
            <Skeleton variant="rectangular" height={400} />
          </Grid>
          <Grid item xs={12} md={4}>
            <Skeleton variant="rectangular" height={400} />
          </Grid>
        </Grid>
      </Box>
    );
  }

  if (!currentProject) {
    return (
      <Box>
        <Alert severity="error" sx={{ mb: 2 }}>
          Project not found or failed to load project details.
        </Alert>
        <Button startIcon={<ArrowBackIcon />} onClick={handleBack}>
          Back to Projects
        </Button>
      </Box>
    );
  }

  return (
    <Box>
      {/* Header */}
      <Box sx={{ mb: 4 }}>
        <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
          <IconButton onClick={handleBack} sx={{ mr: 1 }}>
            <ArrowBackIcon />
          </IconButton>
          <Typography variant="h4" component="h1" fontWeight="bold">
            {currentProject.name}
          </Typography>
          <Chip
            label={currentProject.status}
            color={getStatusColor(currentProject.status) as any}
            sx={{ ml: 2 }}
          />
        </Box>
        <Typography variant="body1" color="text.secondary">
          {currentProject.description || 'No description available'}
        </Typography>
      </Box>

      <Grid container spacing={3}>
        {/* Project Overview */}
        <Grid item xs={12} md={4}>
          <Paper sx={{ p: 3, mb: 3 }}>
            <Typography variant="h6" gutterBottom fontWeight="bold">
              Project Information
            </Typography>
            <List dense>
              <ListItem sx={{ px: 0 }}>
                <ListItemText
                  primary="Repository"
                  secondary={
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <Typography variant="body2">{currentProject.repository_url}</Typography>
                      <IconButton size="small">
                        <OpenInNewIcon fontSize="small" />
                      </IconButton>
                    </Box>
                  }
                />
              </ListItem>
              <Divider />
              <ListItem sx={{ px: 0 }}>
                <ListItemText
                  primary="Status"
                  secondary={
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      {getStatusIcon(currentProject.status)}
                      <Typography variant="body2">{currentProject.status}</Typography>
                    </Box>
                  }
                />
              </ListItem>
              <Divider />
              <ListItem sx={{ px: 0 }}>
                <ListItemText
                  primary="Last Build"
                  secondary={
                    currentProject.lastBuild?.timestamp
                      ? new Date(currentProject.lastBuild.timestamp).toLocaleString()
                      : 'Never'
                  }
                />
              </ListItem>
              <Divider />
              <ListItem sx={{ px: 0 }}>
                <ListItemText
                  primary="Created"
                  secondary={new Date(currentProject.created_at).toLocaleDateString()}
                />
              </ListItem>
            </List>
          </Paper>

          {/* Quick Stats */}
          <Paper sx={{ p: 3 }}>
            <Typography variant="h6" gutterBottom fontWeight="bold">
              Build Statistics
            </Typography>
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <Card variant="outlined">
                  <CardContent sx={{ textAlign: 'center', py: 2 }}>
                    <Typography variant="h4" color="success.main" fontWeight="bold">
                      {currentProject.buildHistory?.filter(b => b.status === 'success').length || 0}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      Success
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>
              <Grid item xs={6}>
                <Card variant="outlined">
                  <CardContent sx={{ textAlign: 'center', py: 2 }}>
                    <Typography variant="h4" color="error.main" fontWeight="bold">
                      {currentProject.buildHistory?.filter(b => b.status === 'failed').length || 0}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      Failed
                    </Typography>
                  </CardContent>
                </Card>
              </Grid>
            </Grid>
          </Paper>
        </Grid>

        {/* Build History */}
        <Grid item xs={12} md={8}>
          <Paper sx={{ p: 3 }}>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
              <Typography variant="h6" fontWeight="bold">
                Build History
              </Typography>
              <Button
                startIcon={<RefreshIcon />}
                onClick={() => refetchBuilds()}
                disabled={buildsLoading}
              >
                Refresh
              </Button>
            </Box>

            {buildsLoading ? (
              <Box>
                {[...Array(5)].map((_, index) => (
                  <Box key={index} sx={{ mb: 2 }}>
                    <Skeleton variant="rectangular" height={80} />
                  </Box>
                ))}
              </Box>
            ) : buildHistory?.data.data && buildHistory.data.data.length > 0 ? (
              <List>
                {buildHistory.data.data.map((build) => (
                  <ListItem
                    key={build.id}
                    divider
                    sx={{
                      border: '1px solid',
                      borderColor: 'divider',
                      borderRadius: 1,
                      mb: 1,
                      flexDirection: 'column',
                      alignItems: 'flex-start',
                    }}
                  >
                    <Box sx={{ display: 'flex', alignItems: 'center', width: '100%', mb: 1 }}>
                      <ListItemIcon sx={{ minWidth: 32 }}>
                        {getStatusIcon(build.status)}
                      </ListItemIcon>
                      <ListItemText
                        primary={
                          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                            <Typography variant="subtitle1" fontWeight="medium">
                              Build #{build.buildNumber || build.id.substring(0, 8)}
                            </Typography>
                            <Chip
                              label={build.status}
                              size="small"
                              color={getStatusColor(build.status) as any}
                              variant="outlined"
                            />
                          </Box>
                        }
                        secondary={
                          <Typography variant="body2" color="text.secondary">
                            {new Date(build.timestamp).toLocaleString()}
                            {build.duration && ` • ${build.duration}s`}
                            {build.author && ` • by ${build.author}`}
                          </Typography>
                        }
                      />
                    </Box>
                    {build.commitMessage && (
                      <Typography
                        variant="body2"
                        color="text.secondary"
                        sx={{ ml: 4, fontStyle: 'italic' }}
                      >
                        "{build.commitMessage}"
                      </Typography>
                    )}
                    {build.status === 'building' && (
                      <LinearProgress sx={{ width: '100%', mt: 1, height: 4, borderRadius: 2 }} />
                    )}
                  </ListItem>
                ))}
              </List>
            ) : (
              <Alert severity="info">
                No build history available for this project.
              </Alert>
            )}
          </Paper>
        </Grid>
      </Grid>
    </Box>
  );
};

export default ProjectDetails;
