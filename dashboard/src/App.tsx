import React from 'react';
import './App.css';
import { HashRouter as Router } from 'react-router-dom';
import { createMuiTheme, ThemeProvider } from '@material-ui/core/styles';
import { CssBaseline } from '@material-ui/core';
import {
  Content,
  SideNav,
  ThemeDark,
  ThemeLight,
  TopBar,
  UseStyles
} from './containers';

export function App(props: any) {
  const classes = UseStyles();
  const [open, setOpen] = React.useState(false);
  const [light, setLight] = React.useState(true);

  const handleDrawerOpen = () => {
    setOpen(true);
  };

  const handleDrawerClose = () => {
    setOpen(false);
  };

  const handleToggleTheme = () => {
    if (light) {
      console.log("set theme 2");
      setLight(!light);
    } else {
      console.log("set theme 1");
      setLight(!light);
    }
  };

  React.useEffect(() => function () {
    if (props.onUnload) {
      props.onUnload();
    }
  }, [props]);

  // React.useEffect(() => {
  //   if (props.webSocketConnected) {
  //     props.fetchActions();
  //     props.fetchTenants();
  //   }
  // }, [props.webSocketConnected]);

  return (
    <ThemeProvider theme={light ? createMuiTheme(ThemeLight) : createMuiTheme(ThemeDark)}>
      <CssBaseline />
      <div className={classes.root}>
        <Router>
          <TopBar open={open} handleDrawerOpen={handleDrawerOpen} handleToggleTheme={handleToggleTheme} />
          <SideNav open={open} handleDrawerClose={handleDrawerClose} />
          <Content />
        </Router>
      </div>
    </ThemeProvider>
  );
}
