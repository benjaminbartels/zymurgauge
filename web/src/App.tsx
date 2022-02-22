import DashboardIcon from "@mui/icons-material/Dashboard";
import KitchenIcon from "@mui/icons-material/Kitchen";
import SettingsIcon from "@mui/icons-material/Settings";
import AppBar from "@mui/material/AppBar";
import Box from "@mui/material/Box";
import Drawer from "@mui/material/Drawer";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import { createTheme, ThemeProvider } from "@mui/material/styles";
import Toolbar from "@mui/material/Toolbar";
import Typography from "@mui/material/Typography";
import { Link, Outlet, useLocation } from "react-router-dom";

const theme = createTheme();
const drawerWidth = 240;

const getDisplayName = (pathname: string) => {
  switch (pathname.split("/")[1]) {
    case "chambers":
      return "Chambers";
    case "settings":
      return "Settings";
    default:
      return "Dashboard";
  }
};

function TopAppBar() {
  const location = useLocation();
  return (
    <AppBar
      position="fixed"
      sx={{
        width: `calc(100% - ${drawerWidth}px)`,
        ml: `${drawerWidth}px`,
      }}
    >
      <Toolbar>
        <Typography variant="h6" component="div" noWrap>
          {getDisplayName(location.pathname)}
        </Typography>
      </Toolbar>
    </AppBar>
  );
}

function NavigationDrawer() {
  return (
    <Drawer
      variant="permanent"
      anchor="left"
      sx={{
        width: drawerWidth,
        flexShrink: 0,
        "& .MuiDrawer-paper": {
          width: drawerWidth,
          boxSizing: "border-box",
        },
      }}
    >
      <List>
        <ListItem button component={Link} to="/">
          <ListItemIcon>
            <DashboardIcon />
          </ListItemIcon>
          <ListItemText primary="Dashboard" />
        </ListItem>
        <ListItem button component={Link} to="chambers">
          <ListItemIcon>
            <KitchenIcon />
          </ListItemIcon>
          <ListItemText primary="Chambers" />
        </ListItem>
        <ListItem button component={Link} to="settings">
          <ListItemIcon>
            <SettingsIcon />
          </ListItemIcon>
          <ListItemText primary="Settings" />
        </ListItem>
      </List>
    </Drawer>
  );
}

export default function App() {
  return (
    <ThemeProvider theme={theme}>
      <Box sx={{ display: "flex" }}>
        <TopAppBar />
        <NavigationDrawer />
        <Box component="main" sx={{ p: 3 }}>
          <Toolbar />
          <Outlet />
        </Box>
      </Box>
    </ThemeProvider>
  );
}
