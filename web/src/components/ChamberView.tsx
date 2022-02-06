import Button from "@mui/material/Button";
import Card from "@mui/material/Card";
import CardActions from "@mui/material/CardActions";
import Grid from "@mui/material/Grid";
import Typography from "@mui/material/Typography";
import React, { useState } from "react";
import { Link, useParams } from "react-router-dom";
import ChamberService from "../services/ChamberService";
import { Chamber } from "../types/Chamber";

export default function ChamberView() {
  const [chamber, setChamber] = useState<Chamber>();
  const [currentFermentationStep, setCurrentFermentationStep] = useState("");
  const params = useParams();

  React.useEffect(() => {
    if (params.chamberId != null) {
      ChamberService.get(params.chamberId)
        .then((response: any) => {
          setChamber(response.data);
          setCurrentFermentationStep(response.data.currentFermentationStep);
        })
        .catch((e: Error) => {
          console.log(e);
        });
    }
  }, [params.chamberId]);

  function startFermentation(name: string) {
    if (params.chamberId != null) {
      ChamberService.startFermentation(params.chamberId, name);

      setCurrentFermentationStep(name);
    }
  }

  function stopFermentation() {
    if (params.chamberId != null) {
      ChamberService.stopFermentation(params.chamberId);

      setCurrentFermentationStep("");
    }
  }

  return chamber != null ? (
    <Card>
      <Typography gutterBottom>{chamber.name}</Typography>
      <Typography>Gravity: {chamber.readings.hydrometerGravity} SG</Typography>
      <Typography>
        Beer Temperature: {chamber.readings.beerTemperature} °C
      </Typography>
      <Typography>
        Auxiliary Temperature: {chamber.readings.auxiliaryTemperature} °C
      </Typography>
      <Typography>
        External Temperature: {chamber.readings.externalTemperature} °C
      </Typography>
      <Typography>STEP ==== {chamber.currentFermentationStep}</Typography>
      <Typography>{chamber.currentBatch.name}</Typography>
      {chamber.currentBatch != null &&
        chamber.currentBatch.fermentation != null &&
        chamber.currentBatch.fermentation.steps != null && (
          <Grid container columns={5}>
            {chamber.currentBatch.fermentation.steps.map((step: any) => (
              <>
                <Grid item xs={1}>
                  <Typography>{step.type}</Typography>
                </Grid>
                <Grid item xs={1}>
                  <Typography>{step.stepTemperature}°C</Typography>
                </Grid>
                <Grid item xs={1}>
                  <Typography>{step.stepTime} Days</Typography>
                </Grid>
                <Grid item xs={1}>
                  <Button
                    onClick={() => startFermentation(step.type)}
                    disabled={currentFermentationStep === step.type}
                  >
                    Start Fermentation
                  </Button>
                </Grid>
                <Grid item xs={1}>
                  <Button
                    onClick={stopFermentation}
                    disabled={currentFermentationStep !== step.type}
                  >
                    Stop Fermentation
                  </Button>
                </Grid>
              </>
            ))}
          </Grid>
        )}

      <CardActions>
        <Button component={Link} to="edit">
          Edit
        </Button>
      </CardActions>
    </Card>
  ) : (
    <Typography>Chamber not found</Typography>
  );
}
