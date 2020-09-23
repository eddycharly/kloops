import React from 'react';
import { Route, Switch } from 'react-router-dom';
import { makeStyles } from '@material-ui/core/styles';
import {
  PluginConfig,
  PluginConfigs,
  PluginsHelp,
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

export function Content() {
  const classes = useStyles();

  return (
    <>
    <main className={classes.content}>
      <Switch>
        <Route path="/help/plugins" exact component={PluginsHelp} />
        <Route path="/config/repositories" exact component={RepoConfigs} />
        <Route path="/config/repository/:name" component={RepoConfig} />
        <Route path="/config/plugins" exact component={PluginConfigs} />
        <Route path="/config/plugin/:name" component={PluginConfig} />
      </Switch>
    </main>
    </>
  );
}
