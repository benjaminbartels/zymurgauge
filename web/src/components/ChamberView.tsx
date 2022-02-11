import {
  Button,
  Card,
  CardActions,
  CardContent,
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
import React, { useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import ChamberService from "../services/chamber-service";
import { Chamber } from "../types/Chamber";

export default function ChamberView() {
  const params = useParams();
  const navigate = useNavigate();
  const [chamber, setChamber] = useState<Chamber>();
  const [currentFermentationStep, setCurrentFermentationStep] = useState("");

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
        console.log(e);
      });
  }

  return chamber != null ? (
    <Card sx={{ maxWidth: 600 }}>
      <CardContent>
        <Typography gutterBottom>
          <b>{chamber.name}</b>
        </Typography>
        <Typography gutterBottom>
          {chamber?.currentBatch?.recipeName !== undefined
            ? chamber?.currentBatch?.recipeName
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
              {chamber?.readings?.beerTemperature} °C
            </Typography>
          </Grid>
          <Grid item xs={9}>
            <Typography variant="body2" noWrap>
              Auxiliary Temperature:
            </Typography>
          </Grid>
          <Grid item xs={3}>
            <Typography align="right" noWrap>
              {chamber?.readings?.auxiliaryTemperature} °C
            </Typography>
          </Grid>
          <Grid item xs={9}>
            <Typography variant="body2" noWrap>
              External Temperature:
            </Typography>
          </Grid>
          <Grid item xs={3}>
            <Typography align="right" noWrap>
              {chamber?.readings?.externalTemperature} °C
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
          chamber.currentBatch.fermentation != null &&
          chamber.currentBatch.fermentation.steps != null && (
            <Grid container>
              <Grid item xs={12}>
                <Typography align="center" noWrap variant="h6">
                  Fermentation Profile:{" "}
                  {chamber.currentBatch?.fermentation?.name}
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
                          <Typography noWrap>Temperature (°C)</Typography>
                        </TableCell>
                        <TableCell align="right">
                          <Typography noWrap>Time (days)</Typography>
                        </TableCell>
                        <TableCell></TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {chamber.currentBatch.fermentation.steps.map(
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
                            <TableCell align="right">{step.stepTime}</TableCell>
                            <TableCell align="right">
                              <Stack
                                direction="row"
                                spacing={1}
                                justifyContent="flex-end"
                              >
                                <Button
                                  variant="contained"
                                  size="small"
                                  onClick={() => startFermentation(step.type)}
                                  disabled={
                                    currentFermentationStep === step.type
                                  }
                                >
                                  Start
                                </Button>
                                <Button
                                  variant="contained"
                                  size="small"
                                  onClick={stopFermentation}
                                  disabled={
                                    currentFermentationStep !== step.type
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
  ) : (
    <Typography>Chamber not found</Typography>
  );
}
