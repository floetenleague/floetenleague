import {
  Alert,
  Box,
  Button,
  Paper,
  TextField,
  Typography,
} from "@mui/material";
import React from "react";
import { useParams } from "react-router-dom";
import { useStore } from "./state";

export const Login = () => {
  const [username, setUsername] = React.useState("");
  const [password, setPassword] = React.useState("");
  const login = useStore((state) => state.login);
  return (
    <Box display="flex" justifyContent="center">
      <Box maxWidth={500}>
        <Paper>
          <Box padding={2}>
            <Typography variant="h4">Internal Login</Typography>
            <TextField
              variant="outlined"
              fullWidth
              value={username}
              placeholder="Username"
              onChange={(e) => setUsername(e.target.value)}
            />
            <TextField
              variant="outlined"
              placeholder="Password"
              type="password"
              fullWidth
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
            <Button
              fullWidth
              variant="outlined"
              onClick={() => {
                login(username, password);
              }}
            >
              Login
            </Button>
          </Box>
        </Paper>
      </Box>
    </Box>
  );
};

export const AuthStatus = () => {
  const { type } = useParams();
  switch (type) {
    case "denied":
      return (
        <Box paddingBottom={3}>
          <Alert severity="error">Login failed. You've denied access.</Alert>
        </Box>
      );
    case "error":
      return (
        <Box paddingBottom={3}>
          <Alert severity="error">
            Something went wrong while logging in, try again later.
          </Alert>
        </Box>
      );
    case "success":
      return (
        <Box paddingBottom={3}>
          <Alert severity="success">Successfully logged in!</Alert>
        </Box>
      );
    default:
      return <>Yikes</>;
  }
};
