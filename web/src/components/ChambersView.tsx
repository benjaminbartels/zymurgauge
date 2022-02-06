import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Card from "@mui/material/Card";
import CardActionArea from "@mui/material/CardActionArea";
import Grid from "@mui/material/Grid";
import Typography from "@mui/material/Typography";
import React, { useState } from "react";
import { Link } from "react-router-dom";
import ChamberService from "../services/ChamberService";
import { Chamber } from "../types/Chamber";

export default function Chambers() {
  const [chambers, setChambers] = useState<Chamber[]>([]);

  React.useEffect(() => {
    ChamberService.getAll()
      .then((response: any) => {
        setChambers(response.data);
      })
      .catch((e: Error) => {
        console.log(e);
      });
  }, []);

  return (
    <Box>
      <Grid container>
        <Grid item>
          {chambers.map((chamber: Chamber) => (
            <Card key={chamber.id}>
              <CardActionArea component={Link} to={chamber.id}>
                <Typography gutterBottom>{chamber.name}</Typography>
                <Typography>
                  Gravity: {chamber.readings.hydrometerGravity} SG
                </Typography>
                <Typography>
                  Beer Temperature: {chamber.readings.beerTemperature} °C
                </Typography>
                <Typography>
                  Auxiliary Temperature: {chamber.readings.auxiliaryTemperature}{" "}
                  °C
                </Typography>
                <Typography>
                  External Temperature: {chamber.readings.externalTemperature}{" "}
                  °C
                </Typography>
                <Typography>{chamber.currentBatch.name}</Typography>
              </CardActionArea>
            </Card>
          ))}
        </Grid>
      </Grid>
      <Button component={Link} to="new">
        New
      </Button>
    </Box>
  );
}
