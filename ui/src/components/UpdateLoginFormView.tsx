import {
  Alert,
  Button,
  Card,
  CardActions,
  CardContent,
  Grid,
  Stack,
  TextField,
} from "@mui/material";
import { useState } from "react";
import { Controller, useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import AuthService from "../services/auth-service";
import { Credentials } from "../types/Auth";

export default function UpdateLoginFormView() {
  const { handleSubmit, control } = useForm();
  const navigate = useNavigate();
  const [errorMessage, setErrorMessage] = useState<String>();

  const onSubmit = (data: any) => {
    if (data.password !== data.confirmPassword) {
      setErrorMessage("Passwords do not match");
      return;
    }

    let credentials: Credentials = {
      username: data.username,
      password: data.password,
    };

    AuthService.save(credentials)
      .then((response: any) => {
        console.debug("Credentials saved: ", response.data);
        AuthService.logout();
        navigate(`/login`);
      })
      .catch((e: any) => {
        setErrorMessage("Could not save Settings: " + e);
      });
  };

  return (
    <>
      {errorMessage != null && <Alert severity="error">{errorMessage}</Alert>}
      <Card sx={{ maxWidth: 600 }}>
        <form onSubmit={handleSubmit(onSubmit)}>
          <CardContent>
            <Grid
              container
              justifyContent="flex-start"
              alignItems="flex-start"
              spacing={2}
            >
              <Grid item xs={12}>
                <Controller
                  name="username"
                  control={control}
                  defaultValue={localStorage.getItem("username")}
                  rules={{ required: "Username is required" }}
                  render={({
                    field: { onChange, value },
                    fieldState: { error },
                  }) => (
                    <TextField
                      fullWidth
                      label="Username"
                      type="text"
                      value={value}
                      onChange={onChange}
                      error={!!error}
                      helperText={error ? error.message : null}
                    />
                  )}
                />
              </Grid>
              <Grid item xs={12}>
                <Controller
                  name="password"
                  control={control}
                  // defaultValue={settings?.authSecret || ""}
                  rules={{ required: "Password is required" }}
                  render={({
                    field: { onChange, value },
                    fieldState: { error },
                  }) => (
                    <TextField
                      fullWidth
                      label="Password"
                      type="password"
                      value={value}
                      onChange={onChange}
                      error={!!error}
                      helperText={error ? error.message : null}
                    />
                  )}
                />
              </Grid>
              <Grid item xs={12}>
                <Controller
                  name="confirmPassword"
                  control={control}
                  // defaultValue={settings?.authSecret || ""}
                  rules={{ required: "Password is required" }}
                  render={({
                    field: { onChange, value },
                    fieldState: { error },
                  }) => (
                    <TextField
                      fullWidth
                      label="Confirm Password"
                      type="password"
                      value={value}
                      onChange={onChange}
                      error={!!error}
                      helperText={error ? error.message : null}
                    />
                  )}
                />
              </Grid>
            </Grid>
          </CardContent>
          <CardActions>
            <Stack direction="row" spacing={1}>
              <Button type="submit" variant="contained">
                Save
              </Button>
            </Stack>
          </CardActions>
        </form>
      </Card>
    </>
  );
}
