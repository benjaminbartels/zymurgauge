import {
  Alert,
  Button,
  Card,
  CardActions,
  CardContent,
  FormControl,
  Grid,
  InputLabel,
  MenuItem,
  Paper,
  Select,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TextField,
  Typography,
} from "@mui/material";
import React, { useState } from "react";
import { Controller, useForm } from "react-hook-form";
import { useNavigate, useParams } from "react-router-dom";
import BatchService from "../services/batch-service";
import ChamberService from "../services/chamber-service";
import ThermometerService from "../services/thermometer-service";
import { BatchDetail, BatchSummary } from "../types/Batch";
import { Chamber } from "../types/Chamber";

export default function ChamberFormView() {
  const params = useParams();
  const navigate = useNavigate();
  const { handleSubmit, control, watch } = useForm();
  const [thermometers, setThermometers] = useState<String[]>();
  const [currentBatchId, setCurrentBatchId] = useState();
  const [batchSummaries, setBatchSummaries] = useState<BatchSummary[]>();
  const [batchDetail, setBatchDetail] = useState<BatchDetail>();
  const [chamber, setChamber] = useState<Chamber>();
  const [temperatureUnitsLabel, SetTemperatureUnitsLabel] = useState("°C");
  const [errorMessage, setErrorMessage] = useState<String>();

  // TODO: Use Redux to store Batches and Thermometers

  React.useEffect(() => {
    if (localStorage.getItem("temperatureUnits") === "Fahrenheit") {
      SetTemperatureUnitsLabel("°F");
    }
  }, []);

  // Load Batches on load
  React.useEffect(() => {
    BatchService.getAllSummaries()
      .then((response: any) => {
        setBatchSummaries(response.data);
      })
      .catch((e: any) => {
        var arr: BatchSummary[] = [];
        setBatchSummaries(arr);
        setErrorMessage("Could not get Batches: " + e);
      });
  }, []);

  // Load Thermometers on load
  React.useEffect(() => {
    ThermometerService.getAll()
      .then((response: any) => {
        setThermometers(response.data);
      })
      .catch((e: any) => {
        setErrorMessage("Could not get Thermometers: " + e);
      });
  }, []);

  React.useEffect(() => {
    if (currentBatchId != null && currentBatchId !== "") {
      BatchService.getDetail(currentBatchId)
        .then((response: any) => {
          console.debug("batch: ", response.data);
          setBatchDetail(response.data);
        })
        .catch((e: Error) => {
          setErrorMessage("Could not get Batch: " + e);
        });
    } else {
      setBatchDetail(undefined);
    }
  }, [currentBatchId]);

  // Load chamber when batches and theremoeters done
  React.useEffect(() => {
    if (
      params.chamberId != null &&
      thermometers != null &&
      batchSummaries != null
    ) {
      ChamberService.get(params.chamberId)
        .then((response: any) => {
          console.debug("chamber: ", response.data);
          setChamber(response.data);
          setCurrentBatchId(response.data.currentBatch?.id);
        })
        .catch((e: any) => {
          setErrorMessage("Could not get Chamber: " + e);
        });
    }
  }, [params.chamberId, thermometers, batchSummaries]);

  const onSubmit = (data: any) => {
    let chamber: Chamber = {
      id: params.chamberId,
      name: data.name,
      deviceConfig: {
        chillerGpio: data.chillerGpio,
        heaterGpio: data.heaterGpio,
        beerThermometerType: data.beerThermometerType,
        beerThermometerId: data.beerThermometerId,
        auxiliaryThermometerType: data.auxiliaryThermometerType,
        auxiliaryThermometerId: data.auxiliaryThermometerId,
        externalThermometerType: data.externalThermometerType,
        externalThermometerId: data.externalThermometerId,
        hydrometerType: data.hydrometerType,
        hydrometerId: data.hydrometerId,
      },
      chillingDifferential: +data.chillingDifferential,
      heatingDifferential: +data.heatingDifferential,
      currentFermentationStep: "",
      currentBatch: batchDetail !== undefined ? batchDetail : undefined,
      readings: null, // TODO: solve this...
    };

    ChamberService.save(chamber)
      .then((response: any) => {
        console.debug("Chamber saved: ", response.data);
        navigate(`../../../chambers/${response.data.id}`);
      })
      .catch((e: any) => {
        setErrorMessage("Could not save Chamber: " + e);
      });
  };

  return (
    <>
      {errorMessage != null && <Alert severity="error">{errorMessage}</Alert>}
      {thermometers != null &&
        batchSummaries != null &&
        (chamber != null || params.chamberId == null) && (
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
                      name="name"
                      control={control}
                      defaultValue={chamber?.name || ""}
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
                  <Grid item xs={12} md={6}>
                    <Controller
                      name="chillingDifferential"
                      control={control}
                      defaultValue={chamber?.chillingDifferential || ""}
                      rules={{ required: "Chilling Differential required" }}
                      render={({
                        field: { onChange, value },
                        fieldState: { error },
                      }) => (
                        <TextField
                          label={`Chilling Differential (${temperatureUnitsLabel})`}
                          type="number"
                          inputProps={{ step: ".1" }}
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
                      name="heatingDifferential"
                      control={control}
                      defaultValue={chamber?.heatingDifferential || ""}
                      rules={{ required: "Heating Differential required" }}
                      render={({
                        field: { onChange, value },
                        fieldState: { error },
                      }) => (
                        <TextField
                          label={`Heating Differential (${temperatureUnitsLabel})`}
                          type="number"
                          inputProps={{ step: ".1" }}
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
                      name="chillerGpio"
                      control={control}
                      defaultValue={chamber?.deviceConfig.chillerGpio || ""}
                      render={({ field: { onChange, value } }) => (
                        <FormControl fullWidth>
                          <InputLabel>Chiller Gpio</InputLabel>
                          <Select
                            label="Chiller Gpio"
                            value={value}
                            onChange={onChange}
                          >
                            {getGpioItems()}
                          </Select>
                        </FormControl>
                      )}
                    />
                  </Grid>
                  <Grid item xs={12} md={6}>
                    <Controller
                      name="heaterGpio"
                      control={control}
                      defaultValue={chamber?.deviceConfig.heaterGpio || ""}
                      render={({ field: { onChange, value } }) => (
                        <FormControl fullWidth>
                          <InputLabel>Heater Gpio</InputLabel>
                          <Select
                            label="Heater Gpio"
                            value={value}
                            onChange={onChange}
                          >
                            {getGpioItems()}
                          </Select>
                        </FormControl>
                      )}
                    />
                  </Grid>
                  <Grid item xs={12} md={6}>
                    <Controller
                      name="beerThermometerType"
                      control={control}
                      defaultValue={
                        chamber?.deviceConfig.beerThermometerType || ""
                      }
                      render={({ field: { onChange, value } }) => (
                        <FormControl fullWidth>
                          <InputLabel>Beer Thermometer Type</InputLabel>
                          <Select
                            label="Beer Thermometer Type"
                            value={value}
                            onChange={onChange}
                          >
                            {getThermometerTypeItems()}
                          </Select>
                        </FormControl>
                      )}
                    />
                  </Grid>
                  <Grid item xs={12} md={6}>
                    <Controller
                      name="beerThermometerId"
                      control={control}
                      defaultValue={
                        chamber?.deviceConfig.beerThermometerId || ""
                      }
                      render={({ field: { onChange, value } }) => (
                        <FormControl fullWidth>
                          <InputLabel>Beer Thermometer ID</InputLabel>
                          <Select
                            label="Beer Thermometer ID"
                            value={value}
                            onChange={onChange}
                          >
                            {watch("beerThermometerType") === "tilt" &&
                              getTiltColorItems()}
                            {watch("beerThermometerType") === "ds18b20" &&
                              getThermometerIDItems(thermometers)}
                          </Select>
                        </FormControl>
                      )}
                    />
                  </Grid>
                  <Grid item xs={12} md={6}>
                    <Controller
                      name="auxiliaryThermometerType"
                      control={control}
                      defaultValue={
                        chamber?.deviceConfig.auxiliaryThermometerType || ""
                      }
                      render={({ field: { onChange, value } }) => (
                        <FormControl fullWidth>
                          <InputLabel>Auxiliary Thermometer Type</InputLabel>
                          <Select
                            label="Auxiliary Thermometer Type"
                            value={value}
                            onChange={onChange}
                          >
                            {getThermometerTypeItems()}
                          </Select>
                        </FormControl>
                      )}
                    />
                  </Grid>
                  <Grid item xs={12} md={6}>
                    <Controller
                      name="auxiliaryThermometerId"
                      control={control}
                      defaultValue={
                        chamber?.deviceConfig.auxiliaryThermometerId || ""
                      }
                      render={({ field: { onChange, value } }) => (
                        <FormControl fullWidth>
                          <InputLabel>Auxiliary Thermometer ID</InputLabel>
                          <Select
                            label="Auxiliary Thermometer ID"
                            value={value}
                            onChange={onChange}
                          >
                            {watch("auxiliaryThermometerType") === "tilt" &&
                              getTiltColorItems()}
                            {watch("auxiliaryThermometerType") === "ds18b20" &&
                              getThermometerIDItems(thermometers)}
                          </Select>
                        </FormControl>
                      )}
                    />
                  </Grid>
                  <Grid item xs={12} md={6}>
                    <Controller
                      name="externalThermometerType"
                      control={control}
                      defaultValue={
                        chamber?.deviceConfig.externalThermometerType || ""
                      }
                      render={({ field: { onChange, value } }) => (
                        <FormControl fullWidth>
                          <InputLabel>External Thermometer Type</InputLabel>
                          <Select
                            label="External Thermometer Type"
                            value={value}
                            onChange={onChange}
                          >
                            {getThermometerTypeItems()}
                          </Select>
                        </FormControl>
                      )}
                    />
                  </Grid>
                  <Grid item xs={12} md={6}>
                    <Controller
                      name="externalThermometerId"
                      control={control}
                      defaultValue={
                        chamber?.deviceConfig.externalThermometerId || ""
                      }
                      render={({ field: { onChange, value } }) => (
                        <FormControl fullWidth>
                          <InputLabel>External Thermometer ID</InputLabel>
                          <Select
                            label="External Thermometer ID"
                            value={value}
                            onChange={onChange}
                          >
                            {watch("externalThermometerType") === "tilt" &&
                              getTiltColorItems()}
                            {watch("externalThermometerType") === "ds18b20" &&
                              getThermometerIDItems(thermometers)}
                          </Select>
                        </FormControl>
                      )}
                    />
                  </Grid>
                  <Grid item xs={12} md={6}>
                    <Controller
                      name="hydrometerType"
                      control={control}
                      defaultValue={chamber?.deviceConfig.hydrometerType || ""}
                      render={({ field: { onChange, value } }) => (
                        <FormControl fullWidth>
                          <InputLabel>Hydrometer Type</InputLabel>
                          <Select
                            label="Hydrometer Type"
                            value={value}
                            onChange={onChange}
                          >
                            {getHydrometerTypeItems()}
                          </Select>
                        </FormControl>
                      )}
                    />
                  </Grid>
                  <Grid item xs={12} md={6}>
                    <Controller
                      name="hydrometerId"
                      control={control}
                      defaultValue={chamber?.deviceConfig.hydrometerId || ""}
                      render={({ field: { onChange, value } }) => (
                        <FormControl fullWidth>
                          <InputLabel>Hydrometer ID</InputLabel>
                          <Select
                            label="Hydrometer ID"
                            value={value}
                            onChange={onChange}
                          >
                            {getTiltColorItems()}
                          </Select>
                        </FormControl>
                      )}
                    />
                  </Grid>
                  <Grid item xs={12}>
                    <Controller
                      name="currentBatchId"
                      control={control}
                      defaultValue={chamber?.currentBatch?.id || ""}
                      render={({ field: { onChange, value } }) => (
                        <FormControl fullWidth>
                          <InputLabel>Batch</InputLabel>
                          <Select
                            label="Batch"
                            value={value}
                            onChange={(event: any, child: any) => {
                              onChange(event, child);
                              setCurrentBatchId(event.target.value);
                            }}
                          >
                            <MenuItem key="" value="">
                              None
                            </MenuItem>
                            {batchSummaries != null &&
                              batchSummaries.map((option) => (
                                <MenuItem key={option.id} value={option.id}>
                                  Batch #{option.number} - {option.recipeName}
                                </MenuItem>
                              ))}
                          </Select>
                        </FormControl>
                      )}
                    />
                  </Grid>
                  {batchDetail != null && (
                    <>
                      <Grid item xs={12}>
                        <Typography align="center" noWrap variant="h6">
                          Fermentation Profile:{" "}
                          {batchDetail?.recipe?.fermentation?.name}
                        </Typography>
                      </Grid>
                      <Grid item xs={12}>
                        <TableContainer component={Paper}>
                          <Table>
                            <TableHead>
                              <TableRow>
                                <TableCell>
                                  <Typography>Type</Typography>
                                </TableCell>
                                <TableCell align="right">
                                  <Typography noWrap>
                                    Temperature ({temperatureUnitsLabel})
                                  </Typography>
                                </TableCell>
                                <TableCell align="right">
                                  <Typography noWrap>Time (days)</Typography>
                                </TableCell>
                              </TableRow>
                            </TableHead>
                            <TableBody>
                              {batchDetail?.recipe?.fermentation?.steps?.map(
                                (step, index) => (
                                  <TableRow
                                    key={index}
                                    sx={{
                                      "&:last-child td, &:last-child th": {
                                        border: 0,
                                      },
                                    }}
                                  >
                                    <TableCell>{step.name}</TableCell>
                                    <TableCell align="right">
                                      {convertDisplayTemperature(
                                        step.temperature
                                      )}
                                    </TableCell>
                                    <TableCell align="right">
                                      {step.duration}
                                    </TableCell>
                                  </TableRow>
                                )
                              )}
                            </TableBody>
                          </Table>
                        </TableContainer>
                      </Grid>
                    </>
                  )}
                </Grid>
              </CardContent>
              <CardActions>
                <Button type="submit" variant="contained">
                  Save
                </Button>
              </CardActions>
            </form>
          </Card>
        )}
    </>
  );
}

const getGpioItems = () => {
  var items: string[] = [];

  for (let i = 0; i <= 25; i++) {
    items.push(i.toString());
  }

  return items.map((option) => (
    <MenuItem key={option} value={option}>
      {option}
    </MenuItem>
  ));
};

const getThermometerTypeItems = () => {
  return [
    <MenuItem key="" value="">
      None
    </MenuItem>,
    <MenuItem key="ds18b20" value="ds18b20">
      DS18B20
    </MenuItem>,
    <MenuItem key="tilt" value="tilt">
      Tilt
    </MenuItem>,
  ];
};

const getHydrometerTypeItems = () => {
  return [
    <MenuItem key="" value="">
      None
    </MenuItem>,
    <MenuItem key="tilt" value="tilt">
      Tilt
    </MenuItem>,
  ];
};

const getThermometerIDItems = (thermometers: any) => {
  return thermometers.map((thermometer: any) => (
    <MenuItem key={thermometer} value={thermometer}>
      {thermometer}
    </MenuItem>
  ));
};

const getTiltColorItems = () => {
  return [
    <MenuItem key="red" value="red">
      Red
    </MenuItem>,
    <MenuItem key="green" value="green">
      Green
    </MenuItem>,
    <MenuItem key="black" value="black">
      Black
    </MenuItem>,
    <MenuItem key="purple" value="purple">
      Purple
    </MenuItem>,
    <MenuItem key="orange" value="orange">
      Orange
    </MenuItem>,
    <MenuItem key="blue" value="blue">
      Blue
    </MenuItem>,
    <MenuItem key="yellow" value="yellow">
      Yellow
    </MenuItem>,
    <MenuItem key="pink" value="pink">
      Pink
    </MenuItem>,
  ];
};

const convertDisplayTemperature = (temperature: number) => {
  if (localStorage.getItem("temperatureUnits") === "Fahrenheit") {
    return (temperature * 9) / 5 + 32;
  }

  return temperature;
};
