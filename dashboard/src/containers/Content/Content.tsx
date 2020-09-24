import React from 'react';
import { UseStyles } from '..';
import { Route, Switch } from 'react-router-dom';
import {
  CommandHelp,
  Home,
  PluginConfig,
  PluginConfigs,
  PluginHelp,
  RepoConfig,
  RepoConfigs
} from '..';

export function Content() {
  const classes = UseStyles();
  return (
    <>
      <main className={classes.content}>
        <Switch>
          <Route path="/" exact component={Home} />
          <Route path="/help/commands" exact component={CommandHelp} />
          <Route path="/help/plugins" exact component={PluginHelp} />
          <Route path="/config/repositories" exact component={RepoConfigs} />
          <Route path="/config/repository/:name" component={RepoConfig} />
          <Route path="/config/plugins" exact component={PluginConfigs} />
          <Route path="/config/plugin/:name" component={PluginConfig} />
        </Switch>
      </main>
    </>
  );
}
