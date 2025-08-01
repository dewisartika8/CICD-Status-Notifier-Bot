import React from 'react';
import {
  AppBar,
  Toolbar,
  Typography,
  IconButton,
  Badge,
  Menu,
  MenuItem,
  Switch,
  FormControlLabel,
  Box,
  Tooltip,
} from '@mui/material';
import {
  Menu as MenuIcon,
  Notifications as NotificationsIcon,
  Brightness4 as DarkModeIcon,
  Brightness7 as LightModeIcon,
  Settings as SettingsIcon,
  Refresh as RefreshIcon,
} from '@mui/icons-material';
import { useAppDispatch, useAppSelector } from '@/store';
import { toggleSidebar, toggleTheme, setAutoRefresh } from '@/store/slices/uiSlice';

const Header: React.FC = () => {
  const dispatch = useAppDispatch();
  const { theme, sidebarOpen, autoRefresh, websocketConnected } = useAppSelector((state) => state.ui);
  const { unreadCount } = useAppSelector((state) => state.notifications);

  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);

  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const handleSidebarToggle = () => {
    dispatch(toggleSidebar());
  };

  const handleThemeToggle = () => {
    dispatch(toggleTheme());
  };

  const handleAutoRefreshToggle = () => {
    dispatch(setAutoRefresh(!autoRefresh));
  };

  return (
    <AppBar
      position="fixed"
      sx={{
        zIndex: (theme) => theme.zIndex.drawer + 1,
        background: (theme) => 
          theme.palette.mode === 'dark' 
            ? 'linear-gradient(45deg, #1e3a8a 30%, #3b82f6 90%)'
            : 'linear-gradient(45deg, #2196F3 30%, #21CBF3 90%)',
      }}
    >
      <Toolbar>
        <IconButton
          color="inherit"
          aria-label="toggle sidebar"
          onClick={handleSidebarToggle}
          edge="start"
          sx={{ mr: 2 }}
        >
          <MenuIcon />
        </IconButton>

        <Typography 
          variant="h6" 
          noWrap component="div" 
          sx={{ 
            flexGrow: 1,
            fontWeight: 'bold',
            letterSpacing: '0.5px',
          }}
        >
          CI/CD Status Notifier
        </Typography>

        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          {/* WebSocket Connection Indicator */}
          <Box
            sx={{
              width: 8,
              height: 8,
              borderRadius: '50%',
              bgcolor: websocketConnected ? 'success.main' : 'error.main',
              mr: 1,
            }}
          />

          {/* Auto Refresh Toggle */}
          <Tooltip title="Toggle Auto Refresh">
            <FormControlLabel
              control={
                <Switch
                  checked={autoRefresh}
                  onChange={handleAutoRefreshToggle}
                  size="small"
                  color="default"
                />
              }
              label="Auto"
              sx={{ 
                color: 'inherit', 
                mr: 1,
                '& .MuiFormControlLabel-label': {
                  fontSize: '0.75rem',
                },
              }}
            />
          </Tooltip>

          {/* Refresh Button */}
          <Tooltip title="Refresh Data">
            <IconButton color="inherit" size="small">
              <RefreshIcon />
            </IconButton>
          </Tooltip>

          {/* Notifications */}
          <Tooltip title="Notifications">
            <IconButton color="inherit">
              <Badge badgeContent={unreadCount} color="error">
                <NotificationsIcon />
              </Badge>
            </IconButton>
          </Tooltip>

          {/* Theme Toggle */}
          <Tooltip title={`Switch to ${theme.mode === 'dark' ? 'light' : 'dark'} mode`}>
            <IconButton color="inherit" onClick={handleThemeToggle}>
              {theme.mode === 'dark' ? <LightModeIcon /> : <DarkModeIcon />}
            </IconButton>
          </Tooltip>

          {/* Settings Menu */}
          <Tooltip title="Settings">
            <IconButton
              color="inherit"
              onClick={handleMenuOpen}
              size="large"
            >
              <SettingsIcon />
            </IconButton>
          </Tooltip>

          <Menu
            anchorEl={anchorEl}
            open={Boolean(anchorEl)}
            onClose={handleMenuClose}
            transformOrigin={{ horizontal: 'right', vertical: 'top' }}
            anchorOrigin={{ horizontal: 'right', vertical: 'bottom' }}
          >
            <MenuItem onClick={handleMenuClose}>Profile</MenuItem>
            <MenuItem onClick={handleMenuClose}>Settings</MenuItem>
            <MenuItem onClick={handleMenuClose}>Logout</MenuItem>
          </Menu>
        </Box>
      </Toolbar>
    </AppBar>
  );
};

export default Header;