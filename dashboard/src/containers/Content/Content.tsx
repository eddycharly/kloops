import React from 'react';
import { Route, Switch } from 'react-router-dom';
import { makeStyles } from '@material-ui/core/styles';
import {
  Plugins,
  RepoConfig,
  RepoConfigs
} from '..';

const useStyles = makeStyles((theme) => ({
  content: {
    flexGrow: 1,
    height: `calc(100vh - 64px)`,
    marginTop: "64px",
    overflow: 'auto',
    display: 'flex',
    flexDirection: 'column',
    backgroundColor: "#f4f4f4",
  },
}));

function Content() {
  const classes = useStyles();

  return (
    <>
    <main className={classes.content}>
      <Switch>
        <Route path="/help/plugins" exact component={Plugins} />
        <Route path="/config/repositories" exact component={RepoConfigs} />
        <Route path="/config/repository/:name" component={RepoConfig} />
      </Switch>
    </main>
    </>
  );
}

export default Content;
