import React from 'react';
import { Route, Switch } from 'react-router-dom';
import { makeStyles } from '@material-ui/core/styles';
import {
  Toolbar
} from '@material-ui/core';
import {
  Plugins
} from '..';

const useStyles = makeStyles((theme) => ({
  content: {
    flexGrow: 1,
    height: '100vh',
    overflow: 'auto',
    backgroundColor: "#f4f4f4",
  },
  paper: {
    height: "100%",
    width: "100%",
  },
}));

function Content() {
  const classes = useStyles();

  return (
    <main className={classes.content}>
      <Toolbar />
      <Switch>
        <Route path="/help/plugins" exact component={Plugins} />
        {/* <Route path="/actions" exact component={Actions} />
        <Route path="/actions/:actionId" component={Action} />
        <Route path="/activity" exact component={Activity} />
        <Route path="/tenants" exact component={Tenants} />
        <Route path="/tenants/:tenantName" component={Tenant} />
        <Route path="/history" component={History} /> */}
      </Switch>
    </main>
  );
}

export default Content;
