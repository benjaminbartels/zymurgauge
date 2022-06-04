import {
  Alert,
  Button,
  Card,
  CardActions,
  CardContent,
  FormControl,
  FormControlLabel,
  FormLabel,
  Grid,
  Radio,
  RadioGroup,
  Stack,
  TextField,
} from "@mui/material";
import React, { useState } from "react";
import { Controller, useForm } from "react-hook-form";
import SettingsService from "../services/settings-service";
import { AppSettings } from "../types/Settings";

export default function SettingsFormView() {
  const { handleSubmit, control } = useForm();
  const [settings, setSettings] = useState<AppSettings>();
  const [errorMessage, setErrorMessage] = useState<String>();

  React.useEffect(() => {
    SettingsService.get()
      .then((response: any) => {
        setSettings(response.data);
      })
      .catch((e: Error) => {
        setErrorMessage("Could not get Settings: " + e);
      });
  }, []);

  const onSubmit = (data: any) => {
    let settings: AppSettings = {
      temperatureUnits: data.temperatureUnits,
      authSecret: data.authSecret,
      brewfatherApiUserId: data.brewfatherApiUserId,
      brewfatherApiKey: data.brewfatherApiKey,
      brewfatherLogUrl: data.brewfatherLogUrl,
      influxDbUrl: data.influxDbUrl,
      influxDbReadToken: data.influxDbReadToken,
      statsDAddress: data.statsDAddress,
    };

    SettingsService.save(settings)
      .then((response: any) => {
        console.debug("Settings saved: ", response.data);
      })
      .catch((e: any) => {
        setErrorMessage("Could not save Settings: " + e);
      });

    localStorage.setItem("temperatureUnits", data.temperatureUnits);
  };

  return (
    <>
      {errorMessage != null && <Alert severity="error">{errorMessage}</Alert>}
      {settings != null && (
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
                  <FormControl component="fieldset">
                    <FormLabel component="legend">Temperature Units</FormLabel>
                    <Controller
                      name="temperatureUnits"
                      control={control}
                      defaultValue={settings?.temperatureUnits || ""}
                      render={({ field }) => (
                        <RadioGroup {...field}>
                          <FormControlLabel
                            value="Fahrenheit"
                            control={<Radio />}
                            label="°F"
                          />
                          <FormControlLabel
                            value="Celsius"
                            control={<Radio />}
                            label="°C"
                          />
                        </RadioGroup>
                      )}
                    />
                  </FormControl>
                </Grid>
                <Grid item xs={12}>
                  <Controller
                    name="authSecret"
                    control={control}
                    defaultValue={settings?.authSecret || ""}
                    rules={{ required: "Authorization Secret required" }}
                    render={({
                      field: { onChange, value },
                      fieldState: { error },
                    }) => (
                      <TextField
                        fullWidth
                        label="Authorization Secret"
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
                    name="brewfatherApiUserId"
                    control={control}
                    defaultValue={settings?.brewfatherApiUserId || ""}
                    rules={{ required: "Brewfather API User ID required" }}
                    render={({
                      field: { onChange, value },
                      fieldState: { error },
                    }) => (
                      <TextField
                        fullWidth
                        label="Brewfather API User ID"
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
                    name="brewfatherApiKey"
                    control={control}
                    defaultValue={settings?.brewfatherApiKey || ""}
                    rules={{ required: "Brewfather API Key required" }}
                    render={({
                      field: { onChange, value },
                      fieldState: { error },
                    }) => (
                      <TextField
                        fullWidth
                        label="Brewfather API Key"
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
                    name="brewfatherLogUrl"
                    control={control}
                    defaultValue={settings?.brewfatherLogUrl || ""}
                    render={({
                      field: { onChange, value },
                      fieldState: { error },
                    }) => (
                      <TextField
                        fullWidth
                        label="Brewfather Log URL"
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
                    name="influxDbUrl"
                    control={control}
                    defaultValue={settings?.influxDbUrl || ""}
                    render={({
                      field: { onChange, value },
                      fieldState: { error },
                    }) => (
                      <TextField
                        fullWidth
                        label="InfluxDB URL"
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
                    name="influxDbReadToken"
                    control={control}
                    defaultValue={settings?.influxDbReadToken || ""}
                    render={({
                      field: { onChange, value },
                      fieldState: { error },
                    }) => (
                      <TextField
                        fullWidth
                        label="InfluxDB Token"
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
                    name="statsDAddress"
                    control={control}
                    defaultValue={settings?.statsDAddress || ""}
                    render={({
                      field: { onChange, value },
                      fieldState: { error },
                    }) => (
                      <TextField
                        fullWidth
                        label="StatsD Address (telegraf)"
                        type="text"
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
      )}
    </>
  );
}
