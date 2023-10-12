import * as React from "react";
import AppBar from "@mui/material/AppBar";
import Box from "@mui/material/Box";
import CssBaseline from "@mui/material/CssBaseline";
import Divider from "@mui/material/Divider";
import Drawer from "@mui/material/Drawer";
import IconButton from "@mui/material/IconButton";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemText from "@mui/material/ListItemText";
import MenuIcon from "@mui/icons-material/Menu";
import PersonIcon from "@mui/icons-material/Person";
import Toolbar from "@mui/material/Toolbar";
import Typography from "@mui/material/Typography";
import { Route, Routes, useParams, useNavigate } from "react-router-dom";
import { BoardPage, NamedMiniBoard } from "./Bingo";
import {
  Alert,
  Button,
  Link,
  Menu,
  MenuItem,
  Paper,
  ToggleButton,
  ToggleButtonGroup,
} from "@mui/material";
import { useStore } from "./state";
import { AuthStatus, Login } from "./Login";
import { Users } from "./Users";
import { FLBingoUserBoard, LoginState, UserPermission } from "./generated";
import { DataGrid } from "@mui/x-data-grid";

const drawerWidth = 240;

export default function Page() {
  const init = useStore((state) => state.init);
  const [mobileOpen, setMobileOpen] = React.useState(false);
  const auth = useStore((state) => state.auth);
  const joinBoard = useStore((state) => state.joinBoard);
  const overview = useStore((state) => state.overview);

  React.useEffect(() => void init(), []);
  const handleDrawerToggle = () => {
    setMobileOpen(!mobileOpen);
  };

  const drawer = (
    <div>
      <Toolbar>
        <Box width="100%" fontSize={25} textAlign="center">
          Flötenleague
        </Box>
      </Toolbar>
      <Divider />
      <List>
        {overview.bingos.map((bingo) => (
          <ListItem key={bingo.id} disablePadding>
            <ListItemButton href={`/bingo/${bingo.id}/overview`}>
              <ListItemText primary={bingo.name} />
            </ListItemButton>
            {auth.permission === UserPermission.User ||
            auth.permission === UserPermission.Moderator ? (
              bingo.boards.some((board) => board.userId === auth.id) ? (
                <ListItemButton href={`/bingo/${bingo.id}/board/${auth.id}`}>
                  <ListItemText primary="My Board" />
                </ListItemButton>
              ) : (
                <ListItemButton onClick={() => joinBoard(bingo.id)}>
                  <ListItemText primary="Join" />
                </ListItemButton>
              )
            ) : undefined}
          </ListItem>
        ))}
      </List>
      <Divider />
      {auth.permission === UserPermission.Moderator ? (
        <List>
          <ListItem disablePadding>
            <ListItemButton href="/users">
              <ListItemText primary="User Management" />
            </ListItemButton>
          </ListItem>
          <ListItem disablePadding>
            <ListItemButton href="/review">
              <ListItemText primary="Review" />
            </ListItemButton>
          </ListItem>
          <Divider />
        </List>
      ) : undefined}
      <Typography padding={2}>
        This product isn't affiliated with or endorsed by Grinding Gear Games in
        any way.
      </Typography>
    </div>
  );

  return (
    <Box sx={{ display: "flex" }}>
      <CssBaseline />
      <AppBar
        position="fixed"
        sx={{
          width: { sm: `calc(100% - ${drawerWidth}px)` },
          ml: { sm: `${drawerWidth}px` },
        }}
      >
        <Toolbar>
          <IconButton
            color="inherit"
            aria-label="open drawer"
            edge="start"
            onClick={handleDrawerToggle}
            sx={{ mr: 2, display: { sm: "none" } }}
          >
            <MenuIcon />
          </IconButton>
          <Typography variant="h6" noWrap component="div" flexGrow={1}>
            Bingo league
          </Typography>
          <AuthNav />
        </Toolbar>
      </AppBar>
      <Box
        component="nav"
        sx={{ width: { sm: drawerWidth }, flexShrink: { sm: 0 } }}
        aria-label="mailbox folders"
      >
        <Drawer
          variant="temporary"
          open={mobileOpen}
          onClose={handleDrawerToggle}
          ModalProps={{
            keepMounted: true, // Better open performance on mobile.
          }}
          sx={{
            display: { xs: "block", sm: "none" },
            "& .MuiDrawer-paper": {
              boxSizing: "border-box",
              width: drawerWidth,
            },
          }}
        >
          {drawer}
        </Drawer>
        <Drawer
          variant="permanent"
          sx={{
            display: { xs: "none", sm: "block" },
            "& .MuiDrawer-paper": {
              boxSizing: "border-box",
              width: drawerWidth,
            },
          }}
          open
        >
          {drawer}
        </Drawer>
      </Box>
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          p: 3,
          width: { sm: `calc(100% - ${drawerWidth}px)` },
        }}
      >
        <Toolbar />
        <NewUserBanner auth={auth} />
        <Routes>
          <Route path="/auth/:type" element={<AuthStatus />}></Route>
          <Route path="/bingo/:id/overview" element={<BingoOverview />}></Route>
          <Route
            path="/bingo/:id/board/:userId"
            element={<BoardPage />}
          ></Route>
          <Route path="/internal/login" element={<Login />}></Route>
          <Route path="/users" element={<Users />}></Route>
        </Routes>
      </Box>
    </Box>
  );
}

const NewUserBanner = ({ auth }: { auth: LoginState }) => {
  if (!auth.loggedIn) {
    return null;
  }

  switch (auth.permission) {
    case UserPermission.Banned:
      return (
        <Box paddingBottom={3}>
          <Alert severity="error">
            You are banned and unable to participate on this website.
          </Alert>
        </Box>
      );
    case UserPermission.Unverified:
      return (
        <Box paddingBottom={3}>
          <Alert severity="warning">
            Your account needs to be verified by a moderator. Contact someone on
            Discord.
          </Alert>
        </Box>
      );
    default:
      return null;
  }
};

const AuthNav = () => {
  const auth = useStore((state) => state.auth);
  const logout = useStore((state) => state.logout);
  const poeLogin = useStore((state) => state.poeLogin);
  const navigate = useNavigate();
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);
  if (!auth.loggedIn) {
    return (
      <Button
        color="inherit"
        onClick={() =>
          poeLogin().then((url) => {
            if (url === "internal") {
              navigate("/internal/login");
            } else {
              window.location.href = url;
            }
          })
        }
      >
        Login
      </Button>
    );
  }

  return (
    <>
      <Button
        color="inherit"
        onClick={(e) => setAnchorEl(e.currentTarget)}
        startIcon={<PersonIcon />}
      >
        {auth.username}
      </Button>
      <Menu
        id="basic-menu"
        anchorEl={anchorEl}
        open={!!anchorEl}
        onClose={() => setAnchorEl(null)}
      >
        <MenuItem onClick={logout}>Logout</MenuItem>
      </Menu>
    </>
  );
};

enum ViewType {
  Grid = "Grid",
  List = "List",
  Faq = "Faq",
}

const BingoOverview = () => {
  const { id } = useParams();
  const [viewType, setViewType] = React.useState(ViewType.Grid);
  const bingo = useStore((state) => state.overview).bingos.find(
    (bingo) => bingo.id === parseInt(id ?? "0")
  );
  if (bingo == null) {
    return <>not found</>;
  }
  const inner = (() => {
    switch (viewType) {
      case ViewType.Grid:
        return (
          <Box display="flex" flexWrap="wrap" gap={3}>
            {bingo.boards.map((board) => (
              <Link
                key={board.userId}
                color="inherit"
                href={`/bingo/${bingo.id}/board/${board.userId}`}
              >
                <NamedMiniBoard
                  name={board.username}
                  score={board.score}
                  columns={bingo.size}
                  fields={bingo.fields}
                  userFields={board.fields}
                />
              </Link>
            ))}
          </Box>
        );
      case ViewType.List:
        return (
          <Box display="flex" justifyContent="center">
            <Box maxWidth={300} width="100%" height={500}>
              <DataGrid<FLBingoUserBoard>
                rows={bingo.boards}
                columns={[
                  { field: "username", headerName: "Name", width: 200 },
                  { field: "score", headerName: "Score", width: 50 },
                ]}
              ></DataGrid>
            </Box>
          </Box>
        );
      case ViewType.Faq:
        return (
          <Box display="flex" justifyContent="center">
            <Box maxWidth={700} width="100%" height={500}>
              <Paper>
                <Box padding={2}>
                  <Typography>
                    1. Jeder Char muss mit dem Tag BNG_ beginnen. Zum Beispiel:
                    BNG_Vishous
                  </Typography>
                  <Typography>
                    2. Bei jedem eingereichten Screenshot muss mindestens der
                    League-Name zu sehen sein! (Overlay Map aktivieren)
                  </Typography>
                  <Typography>
                    3. Wenn 2. mal nicht gehen sollte muss das Character Fenster
                    zu sehen sein (C)
                  </Typography>
                  <Typography>
                    4. Wenn ein Objective completet ist und auch gewertet wurde,
                    düfen die Items benutzt/gelöscht werden. Bsp. Stack Decks
                    öffnen.
                  </Typography>
                  <Typography>
                    5. Alle Challenges müssen im Discord eingereicht werden und
                    zählen erst als Completed, wenn sie von uns geprüft wurden
                    und auf der Website sichtbar sind!
                  </Typography>
                </Box>
              </Paper>
            </Box>
          </Box>
        );
    }
  })();
  return (
    <>
      <Box display="flex" justifyContent="center" marginBottom={3}>
        <ToggleButtonGroup
          color="primary"
          value={viewType}
          exclusive
          onChange={(_, v) => {
            setViewType(v);
          }}
          aria-label="Platform"
        >
          <ToggleButton value={ViewType.Grid}>Grid</ToggleButton>
          <ToggleButton value={ViewType.List}>List</ToggleButton>
          <ToggleButton value={ViewType.Faq}>FAQ</ToggleButton>
        </ToggleButtonGroup>
      </Box>
      {inner}
    </>
  );
};
