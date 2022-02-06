import {
  Button,
  FormControl,
  Grid,
  InputLabel,
  MenuItem,
  Paper,
  Select,
  TextField,
  Typography,
} from "@mui/material";
import React, { useState } from "react";
import { Controller, useForm } from "react-hook-form";
import { useParams } from "react-router-dom";
import BatchService from "../services/BatchService";
import ChamberService from "../services/ChamberService";
import { BatchSummary } from "../types/Batch";
import { Chamber, DeviceConfig } from "../types/Chamber";

export default function EditChamberView() {
  const { handleSubmit, control } = useForm();
  const [chamber, setChamber] = useState<Chamber>();
  const [batchSummaries, setBatchSummaries] = useState<BatchSummary[]>();
  const [beerThermometerType, setBeerThermometerType] = useState("");
  const [beerThermometerID, setBeerThermometerID] = useState("");
  const [auxiliaryThermometerType, setAuxiliaryThermometerType] = useState("");
  const [auxiliaryThermometeID, setAuxiliaryThermometerID] = useState("");
  const [externalThermometerType, setExternalThermometerType] = useState("");
  const [externalThermometeID, setExternalThermometerID] = useState("");
  const [chillerGPIO, setChillerGPIO] = useState<String>();
  const [heaterGPIO, setHeaterGPIO] = useState("");

  const params = useParams();

  React.useEffect(() => {
    if (params.chamberId != null) {
      BatchService.getAllSummaries()
        .then((response: any) => {
          setBatchSummaries(response.data);
        })
        .catch((e: Error) => {
          console.log(e);
        });

      ChamberService.get(params.chamberId)
        .then((response: any) => {
          console.log("Chamber:", response.data);

          response.data.deviceConfigs.forEach((deviceConfig: DeviceConfig) => {
            // console.log(deviceConfig);

            deviceConfig.roles.forEach((role: string) => {
              // console.log(role);

              switch (role) {
                case "beerThermometer":
                  console.log(
                    "setting setBeerThermometerType:",
                    deviceConfig.type
                  );
                  setBeerThermometerType(deviceConfig.type);
                  setBeerThermometerID(deviceConfig.id);
                  break;
                case "auxiliaryThermometer":
                  setAuxiliaryThermometerType(deviceConfig.type);
                  setAuxiliaryThermometerID(deviceConfig.id);
                  break;
                case "externalThermometer":
                  setExternalThermometerType(deviceConfig.type);
                  setExternalThermometerID(deviceConfig.id);
                  break;
                case "chiller":
                  setChillerGPIO(deviceConfig.id);
                  break;
                case "heater":
                  setHeaterGPIO(deviceConfig.id);
                  break;

                default:
                  break;
              }
            });
          });
          setChamber(response.data);
        })

        .catch((e: Error) => {
          console.log(e);
        });
    }
  }, [params.chamberId]);

  const onSubmit = (data: any) => {
    console.log(data);
  };

  const getGPIOs = () => {
    var gpios: string[] = [];

    for (let i = 0; i <= 25; i++) {
      gpios.push("GPIO" + i);
    }

    return gpios;
  };

  return chamber != null ? (
    <Paper elevation={2} sx={{ maxWidth: 600 }}>
      <form onSubmit={handleSubmit(onSubmit)}>
        <Grid
          p={2}
          container
          justifyContent="flex-start"
          alignItems="flex-start"
          spacing={2}
        >
          <Grid item xs={12}>
            <Controller
              name="name"
              control={control}
              defaultValue={chamber?.name}
              rules={{ required: "Name required" }}
              render={({
                field: { onChange, value },
                fieldState: { error },
              }) => (
                <TextField
                  fullWidth
                  label="Name"
                  type="text"
                  value={value}
                  onChange={onChange}
                  error={!!error}
                  helperText={error ? error.message : null}
                />
              )}
            />
          </Grid>
          <Grid item xs={12} md={4}>
            <Controller
              name="chillerKp"
              control={control}
              defaultValue={chamber?.chillerKp}
              rules={{ required: "Chiller Kp required" }}
              render={({
                field: { onChange, value },
                fieldState: { error },
              }) => (
                <TextField
                  label="Chiller Kp"
                  type="number"
                  value={value}
                  onChange={onChange}
                  error={!!error}
                  helperText={error ? error.message : null}
                />
              )}
            />
          </Grid>
          <Grid item xs={12} md={4}>
            <Controller
              name="chillerKi"
              control={control}
              defaultValue={chamber?.chillerKi}
              rules={{ required: "Chiller Ki required" }}
              render={({
                field: { onChange, value },
                fieldState: { error },
              }) => (
                <TextField
                  label="Chiller Ki"
                  type="number"
                  value={value}
                  onChange={onChange}
                  error={!!error}
                  helperText={error ? error.message : null}
                />
              )}
            />
          </Grid>
          <Grid item xs={12} md={4}>
            <Controller
              name="chillerKd"
              control={control}
              defaultValue={chamber?.chillerKd}
              rules={{ required: "Chiller Kd required" }}
              render={({
                field: { onChange, value },
                fieldState: { error },
              }) => (
                <TextField
                  label="Chiller Kd"
                  type="number"
                  value={value}
                  onChange={onChange}
                  error={!!error}
                  helperText={error ? error.message : null}
                />
              )}
            />
          </Grid>
          <Grid item xs={12} md={4}>
            <Controller
              name="heaterKp"
              control={control}
              defaultValue={chamber?.heaterKp}
              rules={{ required: "Heater Kp required" }}
              render={({
                field: { onChange, value },
                fieldState: { error },
              }) => (
                <TextField
                  label="Heater Kp"
                  type="number"
                  value={value}
                  onChange={onChange}
                  error={!!error}
                  helperText={error ? error.message : null}
                />
              )}
            />
          </Grid>
          <Grid item xs={12} md={4}>
            <Controller
              name="heaterKi"
              control={control}
              defaultValue={chamber?.heaterKi}
              rules={{ required: "Heater Ki required" }}
              render={({
                field: { onChange, value },
                fieldState: { error },
              }) => (
                <TextField
                  label="Heater Ki"
                  type="number"
                  value={value}
                  onChange={onChange}
                  error={!!error}
                  helperText={error ? error.message : null}
                />
              )}
            />
          </Grid>
          <Grid item xs={12} md={4}>
            <Controller
              name="heaterKd"
              control={control}
              defaultValue={chamber?.heaterKd}
              rules={{ required: "Heater Kd required" }}
              render={({
                field: { onChange, value },
                fieldState: { error },
              }) => (
                <TextField
                  label="Heater Kd"
                  type="number"
                  value={value}
                  onChange={onChange}
                  error={!!error}
                  helperText={error ? error.message : null}
                />
              )}
            />
          </Grid>
          <Grid item xs={12} md={6}>
            <Controller
              name="chillerGPIO"
              control={control}
              defaultValue={chillerGPIO}
              render={({ field: { onChange, value } }) => (
                <FormControl fullWidth>
                  <InputLabel>Chiller GPIO</InputLabel>
                  <Select
                    label="Chiller GPIO"
                    value={value}
                    onChange={onChange}
                  >
                    {getGPIOs().map((option) => (
                      <MenuItem key={option} value={option}>
                        {option}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              )}
            />
          </Grid>
          <Grid item xs={12} md={6}>
            <Controller
              name="heaterGPIO"
              control={control}
              defaultValue={heaterGPIO}
              render={({ field: { onChange, value } }) => (
                <FormControl fullWidth>
                  <InputLabel>Heater GPIO</InputLabel>
                  <Select label="Heater GPIO" value={value} onChange={onChange}>
                    {getGPIOs().map((option) => (
                      <MenuItem key={option} value={option}>
                        {option}
                      </MenuItem>
                    ))}
                  </Select>
                </FormControl>
              )}
            />
          </Grid>
          <Grid item xs={12} md={6}>
            <Controller
              name="beerThermometerType"
              control={control}
              defaultValue={beerThermometerType}
              render={({ field: { onChange, value } }) => (
                <FormControl fullWidth>
                  <InputLabel>Beer Thermometer Type</InputLabel>
                  <Select
                    label="Beer Thermometer Type"
                    value={value}
                    onChange={onChange}
                  >
                    <MenuItem value="ds18b20">DS18B20</MenuItem>
                    <MenuItem value="tilt">Tilt</MenuItem>
                  </Select>
                </FormControl>
              )}
            />
          </Grid>

          <Grid item xs={12}>
            <Controller
              name="currentBatch"
              control={control}
              defaultValue={chamber?.currentBatch.id}
              render={({ field: { onChange, value } }) => (
                <FormControl fullWidth>
                  <InputLabel>Batch</InputLabel>
                  <Select label="Batch" value={value} onChange={onChange}>
                    {batchSummaries != null &&
                      batchSummaries.map((option) => (
                        <MenuItem key={option.id} value={option.id}>
                          {option.name} #{option.number} - {option.recipeName}
                        </MenuItem>
                      ))}
                  </Select>
                </FormControl>
              )}
            />
          </Grid>
          <Grid item xs={12} md={12}>
            <Button type="submit" variant="contained">
              Save
            </Button>
          </Grid>
        </Grid>
      </form>
    </Paper>
  ) : (
    <Typography>Chamber not found</Typography>
  );
}
