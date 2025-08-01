import React from 'react';
import { Box, Toolbar } from '@mui/material';
import { useAppSelector } from '@/store';
import Header from './Header';
import Sidebar from './Sidebar';

interface LayoutProps {
  children: React.ReactNode;
}

const DRAWER_WIDTH = 280;
const DRAWER_WIDTH_COLLAPSED = 80;

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const { sidebarOpen, sidebarCollapsed } = useAppSelector((state) => state.ui);

  const drawerWidth = sidebarCollapsed ? DRAWER_WIDTH_COLLAPSED : DRAWER_WIDTH;

  return (
    <Box sx={{ display: 'flex', minHeight: '100vh' }}>
      <Header />
      <Sidebar />
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          width: { 
            xs: '100%', 
            md: sidebarOpen ? `calc(100% - ${drawerWidth}px)` : '100%' 
          },
          ml: { 
            xs: 0, 
            md: sidebarOpen ? `${drawerWidth}px` : 0 
          },
          minHeight: '100vh',
          bgcolor: 'background.default',
          transition: (theme) =>
            theme.transitions.create(['width', 'margin'], {
              easing: theme.transitions.easing.sharp,
              duration: theme.transitions.duration.leavingScreen,
            }),
        }}
      >
        <Toolbar />
        <Box
          sx={{
            p: { xs: 2, sm: 3 },
            minHeight: 'calc(100vh - 64px)',
          }}
        >
          {children}
        </Box>
      </Box>
    </Box>
  );
};

export default Layout;