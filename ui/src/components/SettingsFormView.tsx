import {
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
import { Settings } from "../types/Settings";

export default function SettingsView() {
  const { handleSubmit, control } = useForm();
  const [settings, setSettings] = useState<Settings>();

  React.useEffect(() => {
    SettingsService.get()
      .then((response: any) => {
        console.log("settings:", response.data);
        setSettings(response.data);
      })
      .catch((e: Error) => {
        console.log(e);
      });
  }, []);

  const onSubmit = (data: any) => {
    console.log("!!!!! ", data.temperatureUnits);
    let settings: Settings = {
      brewfatherApiUserId: data.brewfatherApiUserId,
      brewfatherApiKey: data.brewfatherApiKey,
      brewfatherLogUrl: data.brewfatherLogUrl,
      temperatureUnits: data.temperatureUnits,
    };

    SettingsService.save(settings)
      .then((response: any) => {
        console.debug("Settings saved: ", response.data);
      })
      .catch((e: any) => {
        console.error("Save Error:", e);
      });
  };

  return (
    <>
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
