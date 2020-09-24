import React from 'react';
import clsx from 'clsx';
import { UseStyles } from '..';
import {
  List,
  ListItem,
  ListItemText,
  Paper,
  Typography
} from '@material-ui/core';
import { Link } from 'react-router-dom';

export function Home() {
  const classes = UseStyles();
  return (
    <Paper variant="outlined" square className={clsx(classes.paper, classes.home)}>
      <Typography variant="h1" align="center" color="primary">Welcome to KLoops !</Typography>
      <Typography variant="h4" align="center" color="secondary">Start using KLoops by :</Typography>
      <List>
        <ListItem button component={Link} to="/config/repositories">
          <ListItemText primary="Managing your repositories" style={{ textAlign: "center" }} />
        </ListItem>
        <ListItem button component={Link} to="/config/plugins">
          <ListItemText primary="Managing your plugins configurations" style={{ textAlign: "center" }} />
        </ListItem>
        <ListItem button component={Link} to="/help/plugins">
          <ListItemText primary="Working with jobs" style={{ textAlign: "center" }} />
        </ListItem>
        <ListItem button component={Link} to="/help/plugins">
          <ListItemText primary="Browsing plugins help" style={{ textAlign: "center" }} />
        </ListItem>
        <ListItem button component={Link} to="/help/about">
          <ListItemText primary="Or landing in the about page" style={{ textAlign: "center" }} />
        </ListItem>
      </List>
    </Paper>
  );
}
