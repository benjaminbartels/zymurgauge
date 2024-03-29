import {
  Alert,
  Button,
  Card,
  CardActions,
  CardContent,
  Grid,
  TextField,
} from "@mui/material";
import { useState } from "react";
import { Controller, useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import AuthService from "../services/auth-service";

export default function LoginFormView() {
  const { handleSubmit, control } = useForm();
  const navigate = useNavigate();
  const [errorMessage, setErrorMessage] = useState<String>();

  const onSubmit = async (data: any) => {
    try {
      await AuthService.login(data.username, data.password);

      navigate(`/chambers`);
    } catch (e: any) {
      setErrorMessage("Could not login: " + e);
    }
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
                  rules={{ required: "Username required" }}
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
                  rules={{ required: "Password required" }}
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
            </Grid>
          </CardContent>
          <CardActions>
            <Button type="submit" variant="contained">
              Login
            </Button>
          </CardActions>
        </form>
      </Card>
    </>
  );
}
