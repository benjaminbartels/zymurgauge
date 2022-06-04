// import DashboardIcon from "@mui/icons-material/Dashboard";
import KeyIcon from "@mui/icons-material/Key";
import KitchenIcon from "@mui/icons-material/Kitchen";
import LogoutIcon from "@mui/icons-material/Logout";
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
import React from "react";
import { Link, Outlet, useLocation, useNavigate } from "react-router-dom";
import AuthService from "./services/auth-service";
import SettingsService from "./services/settings-service";

const theme = createTheme();
const drawerWidth = 240;

const getDisplayName = (pathname: string) => {
  switch (pathname.split("/")[1]) {
    case "chambers":
      return "Chambers";
    case "settings":
      return "Settings";
    case "updateLogin":
      return "Update Login Credentials";
    case "login":
      return "Login";
    case "":
      return "Chambers";
    default:
      return "";
    // return "Dashboard";
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
  const navigate = useNavigate();
  const token = AuthService.getToken();

  const onLogoutClicked = (data: any) => {
    AuthService.logout();
    navigate(`/login`);
  };

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
      {token && (
        <List>
          {/* <ListItem button component={Link} to="/">
                <ListItemIcon>
                  <DashboardIcon />
                </ListItemIcon>
                <ListItemText primary="Dashboard" />
              </ListItem> */}
          {/* <ListItem button component={Link} to="chambers"> */}
          <ListItem button component={Link} to="/">
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
          <ListItem button component={Link} to="updateLogin">
            <ListItemIcon>
              <KeyIcon />
            </ListItemIcon>
            <ListItemText primary="Update Login" />
          </ListItem>
          <ListItem button onClick={onLogoutClicked}>
            <ListItemIcon>
              <LogoutIcon />
            </ListItemIcon>
            <ListItemText primary="Logout" />
          </ListItem>
        </List>
      )}
    </Drawer>
  );
}

export default function App() {
  React.useEffect(() => {
    if (localStorage.getItem("temperatureUnits") === null) {
      SettingsService.get()
        .then((response: any) => {
          localStorage.setItem(
            "temperatureUnits",
            response.data.temperatureUnits
          );
        })
        .catch((e: Error) => {
          console.log("Could not get Settings: " + e);
        });
    }
  }, []);

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
