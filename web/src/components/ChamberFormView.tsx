import {
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

  // TODO: Use Redux to store Batches and Thermometers

  // Load Batches on load
  React.useEffect(() => {
    BatchService.getAllSummaries()
      .then((response: any) => {
        console.debug("batches: ", response.data);
        setBatchSummaries(response.data);
      })
      .catch((e: any) => {
        console.error("Get Batch Summaries Error:", e); // TODO: handle errors
      });
  }, []);

  // Load Thermometers on load
  React.useEffect(() => {
    ThermometerService.getAll()
      .then((response: any) => {
        console.debug("thermometers: ", response.data);
        setThermometers(response.data);
      })
      .catch((e: any) => {
        console.error("Get Thermometers Error:", e);
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
          console.error("Get Batch Detail Error:", e);
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
          console.error("Get Chamber Error:", e);
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
      chillerKp: +data.chillerKp,
      chillerKi: +data.chillerKi,
      chillerKd: +data.chillerKd,
      heaterKp: +data.heaterKp,
      heaterKi: +data.heaterKi,
      heaterKd: +data.heaterKd,
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
        console.error("Save Error:", e);
      });
  };

  return (
    <>
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
                  <Grid item xs={12} md={4}>
                    <Controller
                      name="chillerKp"
                      control={control}
                      defaultValue={chamber?.chillerKp || ""}
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
                      defaultValue={chamber?.chillerKi || ""}
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
                      defaultValue={chamber?.chillerKd || ""}
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
                      defaultValue={chamber?.heaterKp || ""}
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
                      defaultValue={chamber?.heaterKi || ""}
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
                      defaultValue={chamber?.heaterKd || ""}
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
                                  {option.name} #{option.number} -{" "}
                                  {option.recipeName}
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
                                    Temperature (°C)
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
                                    <TableCell>{step.type}</TableCell>
                                    <TableCell align="right">
                                      {step.stepTemperature}
                                    </TableCell>
                                    <TableCell align="right">
                                      {step.stepTime}
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
  var items: string[] = []; // TODO: fix "any's"

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
