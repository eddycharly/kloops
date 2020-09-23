import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import {
  Paper,
  Typography
} from '@material-ui/core';

const useStyles = makeStyles((theme) => ({
  paper: {
    margin: theme.spacing(1),
    padding: theme.spacing(1),
    height: "100%",
    display: "flex",
    flexDirection: "column",
    justifyContent: "center"
  },
}));

export function Home() {
  const classes = useStyles();

  return (
    <Paper variant="outlined" square className={classes.paper}>
      <Typography variant="h1" align="center" color="primary">Welcome to KLoops !</Typography>
    </Paper>
  );
}
