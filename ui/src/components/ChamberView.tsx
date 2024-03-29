import { ClientOptions, InfluxDB } from "@influxdata/influxdb-client-browser";
import {
  Alert,
  Button,
  Card,
  CardActions,
  CardContent,
  Container,
  Grid,
  Paper,
  Stack,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
} from "@mui/material";
import type { ChartData } from "chart.js";
import {
  Chart as ChartJS,
  ChartOptions,
  Legend,
  LinearScale,
  LineElement,
  PointElement,
  TimeScale,
  Tooltip,
} from "chart.js";
import "chartjs-adapter-date-fns";
import Annotations from "chartjs-plugin-annotation";
import { useEffect, useState } from "react";
import { Line } from "react-chartjs-2";
import { Link, useNavigate, useParams } from "react-router-dom";
import ChamberService from "../services/chamber-service";
import SettingsService from "../services/settings-service";
import { Chamber } from "../types/Chamber";

ChartJS.register(
  LinearScale,
  PointElement,
  LineElement,
  TimeScale,
  Annotations,
  Tooltip,
  Legend
);

export default function ChamberView() {
  const params = useParams();
  const navigate = useNavigate();
  const [graphData, setGraphData] = useState<ChartData<any, any, any>>();
  const [graphOptions, setGraphOptions] = useState<ChartOptions<any>>();
  const [chamber, setChamber] = useState<Chamber>();
  const [currentFermentationStep, setCurrentFermentationStep] = useState("");
  const [temperatureUnitsLabel, SetTemperatureUnitsLabel] = useState("°C");
  const [errorMessage, setErrorMessage] = useState<String>();

  useEffect(() => {
    if (localStorage.getItem("temperatureUnits") === "Fahrenheit") {
      SetTemperatureUnitsLabel("°F");
    }
  }, []);

  useEffect(() => {
    var influxDbUrl: string;
    var influxDbReadToken: string;

    SettingsService.get()
      .then((response: any) => {
        influxDbUrl = response.data.influxDbUrl;
        influxDbReadToken = response.data.influxDbReadToken;
      })
      .catch((e: Error) => {
        setErrorMessage("Could not get Settings: " + e);
      });

    if (params.chamberId != null) {
      ChamberService.get(params.chamberId)
        .then((response: any) => {
          setChamber(response.data);
          setCurrentFermentationStep(response.data.currentFermentationStep);

          const originalGravity = parseFloat(
            response.data.currentBatch.recipe.originalGravity
          );
          const finalGravity = parseFloat(
            response.data.currentBatch.recipe.finalGravity
          );

          const influxQuery = async (url: string) => {
            const beerTemperatureData: { x: any; y: any }[] = [];
            const auxiliaryTemperatureData: { x: any; y: any }[] = [];
            const externalTemperatureData: { x: any; y: any }[] = [];
            const hydrometerGravityData: { x: any; y: any }[] = [];

            const chamberName = response.data.name.replace(" ", "");

            let query =
              `from(bucket: "telegraf")
              |> range(start: -12h)
              |> filter(fn: (r) => r._measurement == "zymurgauge_` +
              chamberName +
              `")
              |> sample(n:2, pos: 0)`;

            const clientOptions: ClientOptions = {
              url: url,
              token: influxDbReadToken,
            };

            const queryApi = await new InfluxDB(clientOptions).getQueryApi(
              "zymurgauge"
            );

            await queryApi.queryRows(query, {
              next(row, tableMeta) {
                const o = tableMeta.toObject(row);
                switch (o._field) {
                  case "beer_temperature": {
                    beerTemperatureData.push({ x: o._time, y: o._value });
                    break;
                  }
                  case "auxiliary_temperature": {
                    auxiliaryTemperatureData.push({ x: o._time, y: o._value });
                    break;
                  }
                  case "external_temperature": {
                    externalTemperatureData.push({ x: o._time, y: o._value });
                    break;
                  }
                  case "hydrometer_gravity": {
                    hydrometerGravityData.push({ x: o._time, y: o._value });
                    break;
                  }
                }
              },

              complete() {
                const options: ChartOptions<any> = {
                  responsive: true,
                  interaction: {
                    mode: "index" as const,
                    intersect: false,
                  },
                  plugins: {
                    annotation: {
                      annotations: {
                        og: {
                          type: "line",
                          yScaleID: "y1",
                          yMin: originalGravity,
                          yMax: originalGravity,
                          borderDash: [2, 2],
                          label: {
                            backgroundColor: "rgba(0, 0, 0, 0.0)",
                            color: "rgba(0, 0, 0)",
                            enabled: true,
                            content: "OG",
                            position: "start",
                            font: { weight: "normal" },
                            yAdjust: -10,
                          },
                        },
                        fg: {
                          type: "line",
                          yScaleID: "y1",
                          yMin: finalGravity,
                          yMax: finalGravity,
                          borderDash: [2, 2],
                          label: {
                            backgroundColor: "rgba(0, 0, 0, 0.0)",
                            color: "rgba(0, 0, 0)",
                            enabled: true,
                            content: "FG",
                            position: "start",
                            font: { weight: "normal" },
                            yAdjust: 10,
                          },
                        },
                      },
                    },
                  },
                  scales: {
                    x: {
                      type: "time",
                      title: {
                        display: true,
                        text: "Date",
                      },
                    },
                    y: {
                      title: {
                        display: true,
                        text: "Temperature " + temperatureUnitsLabel,
                      },
                      display: true,
                      position: "left" as const,
                      suggestedMax: 25,
                      suggestedMin: 19,
                    },
                    y1: {
                      title: {
                        display: true,
                        text: "Specific Gravity",
                      },
                      type: "linear" as const,
                      display: true,
                      position: "right" as const,
                      suggestedMax: originalGravity + 0.001,
                      suggestedMin: finalGravity - 0.001,
                    },
                  },
                };

                const data: ChartData<any, any, any> = {
                  datasets: [],
                };

                if (beerTemperatureData.length > 0) {
                  data.datasets.push({
                    label: "Beer Temperature",
                    data: beerTemperatureData,
                    yAxisID: "y",
                    xAxisID: "x",
                    pointRadius: 0,
                    borderColor: "rgb(255, 0, 0)",
                    backgroundColor: "rgba(255, 0, 0, 0.5)",
                  });
                }

                if (auxiliaryTemperatureData.length > 0) {
                  data.datasets.push({
                    label: "Auxiliary Temperature",
                    data: auxiliaryTemperatureData,
                    yAxisID: "y",
                    xAxisID: "x",
                    pointRadius: 0,
                    borderColor: "rgb(0, 0, 255)",
                    backgroundColor: "rgba(0, 0, 255, 0.5)",
                  });
                }

                if (externalTemperatureData.length > 0) {
                  data.datasets.push({
                    label: "External Temperature",
                    data: externalTemperatureData,
                    yAxisID: "y",
                    xAxisID: "x",
                    pointRadius: 0,
                    borderColor: "rgb(0, 255, 0)",
                    backgroundColor: "rgba(0, 255, 0, 0.5)",
                  });
                }

                if (hydrometerGravityData.length > 0) {
                  data.datasets.push({
                    label: "Gravity",
                    data: hydrometerGravityData,
                    yAxisID: "y1",
                    xAxisID: "x",
                    pointRadius: 0,
                    borderColor: "rgb( 255, 0, 255)",
                    backgroundColor: "rgba( 255, 0, 255, 0.5)",
                  });
                }

                setGraphOptions(options);
                setGraphData(data);
              },
              error(error) {
                console.log("query failed- ", error);
              },
            });
          };

          if (influxDbUrl != null && influxDbUrl !== "") {
            influxQuery(influxDbUrl);
          }
        })
        .catch((e: Error) => {
          setErrorMessage("Could not get Chamber: " + e);
        });
    }
  }, [params.chamberId, temperatureUnitsLabel]);

  function startFermentation(name: string) {
    if (params.chamberId != null) {
      ChamberService.startFermentation(params.chamberId, name);

      // TODO: handler error

      setCurrentFermentationStep(name);
    }
  }

  function stopFermentation() {
    if (params.chamberId != null) {
      ChamberService.stopFermentation(params.chamberId);

      // TODO: handler error

      setCurrentFermentationStep("");
    }
  }

  function remove(id: string) {
    ChamberService.delete(id)
      .then(() => {
        navigate(`../../../chambers`);
      })

      .catch((e: Error) => {
        setErrorMessage("Could not delete Chamber: " + e);
      });
  }
  return (
    <>
      {errorMessage != null && <Alert severity="error">{errorMessage}</Alert>}
      {chamber != null && (
        <Card sx={{ maxWidth: 800 }}>
          <CardContent>
            <Typography gutterBottom>
              <b>{chamber.name}</b>
            </Typography>
            <Typography gutterBottom>
              {chamber?.currentBatch?.recipe?.name !== undefined
                ? chamber?.currentBatch?.recipe?.name
                : "No Recipe"}
            </Typography>
            <Grid container>
              <Grid item xs={9}>
                <Typography variant="body2" noWrap>
                  Beer Temperature:
                </Typography>
              </Grid>
              <Grid item xs={3} alignContent="right">
                <Typography align="right" noWrap>
                  {chamber?.readings?.beerTemperature &&
                    convertDisplayTemperature(
                      chamber?.readings?.beerTemperature
                    )}{" "}
                  {temperatureUnitsLabel}
                </Typography>
              </Grid>
              <Grid item xs={9}>
                <Typography variant="body2" noWrap>
                  Auxiliary Temperature:
                </Typography>
              </Grid>
              <Grid item xs={3}>
                <Typography align="right" noWrap>
                  {chamber?.readings?.auxiliaryTemperature &&
                    convertDisplayTemperature(
                      chamber?.readings?.auxiliaryTemperature
                    )}{" "}
                  {temperatureUnitsLabel}
                </Typography>
              </Grid>
              <Grid item xs={9}>
                <Typography variant="body2" noWrap>
                  External Temperature:
                </Typography>
              </Grid>
              <Grid item xs={3}>
                <Typography align="right" noWrap>
                  {chamber?.readings?.externalTemperature &&
                    convertDisplayTemperature(
                      chamber?.readings?.externalTemperature
                    )}{" "}
                  {temperatureUnitsLabel}
                </Typography>
              </Grid>
              <Grid item xs={9}>
                <Typography variant="body2" noWrap>
                  Gravity:
                </Typography>
              </Grid>
              <Grid item xs={3}>
                <Typography align="right" noWrap>
                  {chamber?.readings?.hydrometerGravity} SG
                </Typography>
              </Grid>
            </Grid>
            {chamber.currentBatch != null &&
              chamber.currentBatch.recipe.fermentation != null &&
              chamber.currentBatch.recipe.fermentation.steps != null && (
                <Grid container>
                  <Grid item xs={12}>
                    <Typography align="center" noWrap variant="h6">
                      Fermentation Profile:{" "}
                      {chamber.currentBatch?.recipe.fermentation?.name}
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
                            <TableCell></TableCell>
                          </TableRow>
                        </TableHead>
                        <TableBody>
                          {chamber.currentBatch.recipe.fermentation.steps.map(
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
                                  {convertDisplayTemperature(step.temperature)}
                                </TableCell>
                                <TableCell align="right">
                                  {step.duration}
                                </TableCell>
                                <TableCell align="right">
                                  <Stack
                                    direction="row"
                                    spacing={1}
                                    justifyContent="flex-end"
                                  >
                                    <Button
                                      variant="contained"
                                      size="small"
                                      onClick={() =>
                                        startFermentation(step.name)
                                      }
                                      disabled={
                                        currentFermentationStep === step.name
                                      }
                                    >
                                      Start
                                    </Button>
                                    <Button
                                      variant="contained"
                                      size="small"
                                      onClick={stopFermentation}
                                      disabled={
                                        currentFermentationStep !== step.name
                                      }
                                    >
                                      Stop
                                    </Button>
                                  </Stack>
                                </TableCell>
                              </TableRow>
                            )
                          )}
                        </TableBody>
                      </Table>
                    </TableContainer>
                  </Grid>
                  {graphData != null && (
                    <Grid item xs={12}>
                      <Container
                        sx={{
                          height: 400,
                        }}
                      >
                        <Line options={graphOptions} data={graphData} />
                      </Container>
                    </Grid>
                  )}
                </Grid>
              )}
          </CardContent>
          <CardActions>
            <Stack direction="row" spacing={1}>
              <Button component={Link} to="edit" variant="contained">
                Edit
              </Button>
              <Button variant="contained" onClick={() => remove(chamber.id!)}>
                Delete
              </Button>
            </Stack>
          </CardActions>
        </Card>
      )}
    </>
  );
}

const convertDisplayTemperature = (temperature: number) => {
  if (localStorage.getItem("temperatureUnits") === "Fahrenheit") {
    return (temperature * 9) / 5 + 32;
  }

  return temperature;
};
