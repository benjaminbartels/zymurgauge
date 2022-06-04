import {
  Alert,
  Button,
  Card,
  CardActionArea,
  CardContent,
  Grid,
  Typography,
} from "@mui/material";
import React, { useState } from "react";
import { Link } from "react-router-dom";
import ChamberService from "../services/chamber-service";
import { Chamber } from "../types/Chamber";

export default function Chambers() {
  const [chambers, setChambers] = useState<Chamber[]>([]);
  const [temperatureUnitsLabel, SetTemperatureUnitsLabel] = useState("°C");
  const [errorMessage, setErrorMessage] = useState<String>();

  React.useEffect(() => {
    if (localStorage.getItem("temperatureUnits") === "Fahrenheit") {
      SetTemperatureUnitsLabel("°F");
    }
  }, []);

  React.useEffect(() => {
    ChamberService.getAll()
      .then((response: any) => {
        setChambers(response.data);
      })
      .catch((e: Error) => {
        setErrorMessage("Could not get Chambers: " + e);
      });
  }, []);
  // TODO: Use redux to store Chambers
  return (
    <>
      {errorMessage != null && <Alert severity="error">{errorMessage}</Alert>}
      <Grid container spacing={5}>
        {chambers.map((chamber: Chamber) => (
          <Grid item key={chamber.id} xs={12} sm={12} md={6} lg={3} xl={3}>
            <Card sx={{ minWidth: 400 }}>
              <CardActionArea component={Link} to={chamber!.id!}>
                <CardContent>
                  <Typography gutterBottom>
                    <b>{chamber.name}</b>
                  </Typography>
                  <Typography gutterBottom>
                    {chamber?.currentBatch?.recipe.name !== undefined
                      ? chamber?.currentBatch?.recipe.name
                      : "No Recipe"}
                  </Typography>
                  <Grid container>
                    <Grid item xs={9}>
                      <Typography variant="body2" noWrap>
                        Beer Temperature:
                      </Typography>
                    </Grid>
                    <Grid item xs={3}>
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
                </CardContent>
              </CardActionArea>
            </Card>
          </Grid>
        ))}
        <Grid item xs={12}>
          <Button component={Link} to="new" variant="contained">
            New
          </Button>
        </Grid>
      </Grid>
    </>
  );
}

const convertDisplayTemperature = (temperature: number) => {
  if (localStorage.getItem("temperatureUnits") === "Fahrenheit") {
    return (temperature * 9) / 5 + 32;
  }

  return temperature;
};
