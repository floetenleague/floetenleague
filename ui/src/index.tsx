import React from "react";
import ReactDOM from "react-dom/client";
import "./index.css";
import App from "./App";
import {
  createTheme,
  LinkProps,
  ListItemButtonProps,
  ThemeProvider,
} from "@mui/material";
import {
  HashRouter,
  Link as RouterLink,
  LinkProps as RouterLinkProps,
} from "react-router-dom";
import { SnackbarProvider } from "notistack";

const root = ReactDOM.createRoot(
  document.getElementById("root") as HTMLElement
);

export const LinkBehavior = React.forwardRef<
  HTMLAnchorElement,
  Omit<RouterLinkProps, "to"> & { href: RouterLinkProps["to"] }
>((props, ref) => {
  const { href, ...other } = props;
  // Map href (MUI) -> to (react-router)
  return <RouterLink ref={ref} to={href} {...other} />;
});

const myTheme = createTheme({
  components: {
    MuiLink: {
      defaultProps: {
        component: LinkBehavior,
      } as LinkProps,
    },
    MuiButtonBase: {
      defaultProps: {
        LinkComponent: LinkBehavior,
      },
    },
    MuiListItemButton: {
      defaultProps: {
        component: LinkBehavior,
      } as ListItemButtonProps,
    },
  },
  palette: {
    background: {
      default: "#282828",
      paper: "#32302f",
    },
    text: {
      primary: "#fbf1d4",
    },
    primary: {
      main: "#a89984",
    },
    secondary: {
      main: "#f44336",
    },
    mode: "dark",
  },
});

root.render(
  <React.StrictMode>
      <SnackbarProvider maxSnack={3}>
    <HashRouter>
        <ThemeProvider theme={myTheme}>
          <App />
        </ThemeProvider>
    </HashRouter>
      </SnackbarProvider>
  </React.StrictMode>
);
