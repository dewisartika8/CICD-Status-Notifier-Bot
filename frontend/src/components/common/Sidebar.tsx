import React from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import {
  Drawer,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Typography,
  Divider,
  Box,
  Collapse,
  IconButton,
  Tooltip,
} from '@mui/material';
import {
  Dashboard as DashboardIcon,
  Folder as ProjectsIcon,
  Analytics as AnalyticsIcon,
  Settings as SettingsIcon,
  ExpandLess,
  ExpandMore,
  Build as BuildIcon,
  Notifications as NotificationsIcon,
  ChevronLeft as ChevronLeftIcon,
  ChevronRight as ChevronRightIcon,
} from '@mui/icons-material';
import { useAppDispatch, useAppSelector } from '@/store';
import { setSidebarOpen, toggleSidebarCollapsed } from '@/store/slices/uiSlice';

interface MenuItem {
  id: string;
  label: string;
  path: string;
  icon: React.ReactElement;
  children?: MenuItem[];
}

const menuItems: MenuItem[] = [
  {
    id: 'dashboard',
    label: 'Dashboard',
    path: '/',
    icon: <DashboardIcon />,
  },
  {
    id: 'projects',
    label: 'Projects',
    path: '/projects',
    icon: <ProjectsIcon />,
    children: [
      {
        id: 'all-projects',
        label: 'All Projects',
        path: '/projects',
        icon: <ProjectsIcon />,
      },
      {
        id: 'builds',
        label: 'Builds',
        path: '/projects/builds',
        icon: <BuildIcon />,
      },
    ],
  },
  {
    id: 'analytics',
    label: 'Analytics',
    path: '/analytics',
    icon: <AnalyticsIcon />,
  },
  {
    id: 'notifications',
    label: 'Notifications',
    path: '/notifications',
    icon: <NotificationsIcon />,
  },
  {
    id: 'settings',
    label: 'Settings',
    path: '/settings',
    icon: <SettingsIcon />,
  },
];

const DRAWER_WIDTH = 280;
const DRAWER_WIDTH_COLLAPSED = 80;

const Sidebar: React.FC = () => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const location = useLocation();
  const { sidebarOpen, sidebarCollapsed } = useAppSelector((state) => state.ui);

  const [expandedItems, setExpandedItems] = React.useState<string[]>(['projects']);

  const handleItemClick = (path: string) => {
    navigate(path);
  };

  const handleExpandClick = (itemId: string) => {
    setExpandedItems(prev =>
      prev.includes(itemId)
        ? prev.filter(id => id !== itemId)
        : [...prev, itemId]
    );
  };

  const handleDrawerClose = () => {
    dispatch(setSidebarOpen(false));
  };

  const handleCollapseToggle = () => {
    dispatch(toggleSidebarCollapsed());
  };

  const isActive = (path: string) => {
    if (path === '/') {
      return location.pathname === '/';
    }
    return location.pathname.startsWith(path);
  };

  const drawerWidth = sidebarCollapsed ? DRAWER_WIDTH_COLLAPSED : DRAWER_WIDTH;

  const renderMenuItem = (item: MenuItem, depth = 0) => {
    const hasChildren = item.children && item.children.length > 0;
    const isExpanded = expandedItems.includes(item.id);
    const active = isActive(item.path);

    return (
      <React.Fragment key={item.id}>
        <ListItem disablePadding sx={{ display: 'block' }}>
          <ListItemButton
            onClick={() => {
              if (hasChildren) {
                handleExpandClick(item.id);
              } else {
                handleItemClick(item.path);
              }
            }}
            sx={{
              minHeight: 48,
              justifyContent: sidebarCollapsed ? 'center' : 'initial',
              px: 2.5,
              ml: depth * 2,
              bgcolor: active ? 'action.selected' : 'transparent',
              '&:hover': {
                bgcolor: active ? 'action.selected' : 'action.hover',
              },
            }}
          >
            <ListItemIcon
              sx={{
                minWidth: 0,
                mr: sidebarCollapsed ? 0 : 3,
                justifyContent: 'center',
                color: active ? 'primary.main' : 'inherit',
              }}
            >
              {sidebarCollapsed ? (
                <Tooltip title={item.label} placement="right">
                  {item.icon}
                </Tooltip>
              ) : (
                item.icon
              )}
            </ListItemIcon>
            {!sidebarCollapsed && (
              <>
                <ListItemText
                  primary={item.label}
                  sx={{
                    opacity: 1,
                    color: active ? 'primary.main' : 'inherit',
                    fontWeight: active ? 600 : 400,
                  }}
                />
                {hasChildren && (isExpanded ? <ExpandLess /> : <ExpandMore />)}
              </>
            )}
          </ListItemButton>
        </ListItem>
        {hasChildren && !sidebarCollapsed && (
          <Collapse in={isExpanded} timeout="auto" unmountOnExit>
            <List component="div" disablePadding>
              {item.children?.map(child => renderMenuItem(child, depth + 1))}
            </List>
          </Collapse>
        )}
      </React.Fragment>
    );
  };

  return (
    <Drawer
      variant="permanent"
      open={sidebarOpen}
      sx={{
        width: sidebarOpen ? drawerWidth : 0,
        flexShrink: 0,
        whiteSpace: 'nowrap',
        boxSizing: 'border-box',
        ...(sidebarOpen && {
          '& .MuiDrawer-paper': {
            width: drawerWidth,
            transition: (theme) =>
              theme.transitions.create('width', {
                easing: theme.transitions.easing.sharp,
                duration: theme.transitions.duration.enteringScreen,
              }),
            overflowX: 'hidden',
          },
        }),
        ...(!sidebarOpen && {
          '& .MuiDrawer-paper': {
            width: 0,
            transition: (theme) =>
              theme.transitions.create('width', {
                easing: theme.transitions.easing.sharp,
                duration: theme.transitions.duration.leavingScreen,
              }),
            overflowX: 'hidden',
          },
        }),
      }}
    >
      <Box
        sx={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'space-between',
          padding: (theme) => theme.spacing(1),
          ...theme => theme.mixins.toolbar,
        }}
      >
        {!sidebarCollapsed && (
          <Typography variant="h6" noWrap component="div" sx={{ ml: 2, fontWeight: 'bold' }}>
            Menu
          </Typography>
        )}
        <IconButton onClick={handleCollapseToggle}>
          {sidebarCollapsed ? <ChevronRightIcon /> : <ChevronLeftIcon />}
        </IconButton>
      </Box>
      <Divider />
      <List>
        {menuItems.map(item => renderMenuItem(item))}
      </List>
    </Drawer>
  );
};

export default Sidebar;