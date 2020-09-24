import React from 'react';
import { UseStyles } from '..';
import { Route, Switch } from 'react-router-dom';
import {
  Home,
  PluginConfig,
  PluginConfigs,
  PluginsHelp,
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
