import React from 'react';
import './App.css';
import { HashRouter as Router } from 'react-router-dom';
import { makeStyles } from '@material-ui/core/styles';
import { CssBaseline } from '@material-ui/core';
import {
  Content,
  SideNav,
  TopBar
} from './containers';

const useStyles = makeStyles((theme) => ({
  root: {
    display: 'flex',
  }
}));

function App(props: any) {
  const classes = useStyles();
  const [open, setOpen] = React.useState(false);

  const handleDrawerOpen = () => {
    setOpen(true);
  };

  const handleDrawerClose = () => {
    setOpen(false);
  };

  React.useEffect(() => function () {
    if (props.onUnload) {
      props.onUnload();
    }
  }, []);

  // React.useEffect(() => {
  //   if (props.webSocketConnected) {
  //     props.fetchActions();
  //     props.fetchTenants();
  //   }
  // }, [props.webSocketConnected]);

  return (
    <div className={classes.root}>
      <CssBaseline />
      <Router>
        <TopBar open={open} handleDrawerOpen={handleDrawerOpen} />
        <SideNav open={open} handleDrawerClose={handleDrawerClose} />
        <Content />
      </Router>
    </div>
  );
}

export default App;
