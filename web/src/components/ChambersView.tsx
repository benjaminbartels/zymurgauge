import {
  Button,
  Card,
  CardActionArea,
  CardContent,
  Grid,
  Typography,
} from "@mui/material";
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

  // TODO: Use redux to store
  return (
    <>
      <Grid container spacing={5}>
        {chambers.map((chamber: Chamber) => (
          <Grid item key={chamber.id} xs={12} sm={12} md={6} lg={3} xl={3}>
            <Card>
              <CardActionArea component={Link} to={chamber!.id!}>
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
                    <Grid item xs={3}>
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
