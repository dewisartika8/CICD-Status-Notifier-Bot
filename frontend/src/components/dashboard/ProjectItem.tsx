import React from 'react';
import { Box, Card, CardContent, Typography, Chip, IconButton } from '@mui/material';
import { PlayArrow, Stop, Refresh } from '@mui/icons-material';
import StatusBadge from '../projects/StatusBadge';
import { Project } from '../../types';

interface ProjectItemProps {
  project: Project;
}

const ProjectItem: React.FC<ProjectItemProps> = ({ project }) => {
  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'success':
      case 'completed':
        return 'success';
      case 'failed':
      case 'error':
        return 'error';
      case 'running':
      case 'in_progress':
        return 'warning';
      default:
        return 'default';
    }
  };

  return (
    <Card 
      sx={{ 
        mb: 2, 
        '&:hover': { 
          boxShadow: 3,
          transform: 'translateY(-2px)',
          transition: 'all 0.2s ease-in-out'
        } 
      }}
    >
      <CardContent>
        <Box display="flex" justifyContent="space-between" alignItems="flex-start">
          <Box flex={1}>
            <Typography variant="h6" component="h3" gutterBottom>
              {project.name}
            </Typography>
            <Typography variant="body2" color="text.secondary" gutterBottom>
              {project.description || 'No description available'}
            </Typography>
            <Box display="flex" alignItems="center" gap={1} mt={1}>
              <StatusBadge status={project.status} />
              <Chip
                label={`Build #${project.lastBuildNumber || 'N/A'}`}
                size="small"
                variant="outlined"
              />
                            {project.repository_url && (
                <p className="text-sm text-gray-600 mb-2">
                  {project.repository_url}
                </p>
              )}
            </Box>
          </Box>
          
          <Box display="flex" gap={1}>
            <IconButton 
              size="small" 
              color="primary"
              title="Start Build"
            >
              <PlayArrow />
            </IconButton>
            <IconButton 
              size="small" 
              color="secondary"
              title="Stop Build"
            >
              <Stop />
            </IconButton>
            <IconButton 
              size="small"
              title="Refresh Status"
            >
              <Refresh />
            </IconButton>
          </Box>
        </Box>
        
        {project.lastUpdated && (
          <Typography variant="caption" color="text.secondary" display="block" mt={1}>
            Last updated: {new Date(project.lastUpdated).toLocaleString()}
          </Typography>
        )}
      </CardContent>
    </Card>
  );
};

export default ProjectItem;
