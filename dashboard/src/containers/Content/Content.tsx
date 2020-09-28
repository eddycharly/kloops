import React from 'react';
import { UseStyles } from '..';
import { Route, Switch } from 'react-router-dom';
import {
  Home,
  PluginConfig,
  PluginConfigs,
  RepoConfig,
  RepoConfigs
} from '..';
import { CommandHelp } from 'features/commandHelp/CommandHelp';
import { PluginHelpList } from 'features/pluginHelp/PluginHelpList';

export function Content() {
  const classes = UseStyles();
  return (
    <>
      <main className={classes.content}>
        <Switch>
          <Route path="/" exact component={Home} />
          <Route path="/help/commands" exact component={CommandHelp} />
          <Route path="/help/plugins" exact component={PluginHelpList} />
          <Route path="/config/repositories" exact component={RepoConfigs} />
          <Route path="/config/repository/:name" component={RepoConfig} />
          <Route path="/config/plugins" exact component={PluginConfigs} />
          <Route path="/config/plugin/:name" component={PluginConfig} />
        </Switch>
      </main>
    </>
  );
}
