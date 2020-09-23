import React from 'react';
import clsx from 'clsx';
import { makeStyles } from '@material-ui/core/styles';
import {
  AppBar,
  IconButton,
  LinearProgress,
  Link,
  Toolbar,
  Typography
} from '@material-ui/core';
import {
  Menu as MenuIcon
} from '@material-ui/icons';

const drawerWidth = 240;

const useStyles = makeStyles((theme) => ({
  appBar: {
    zIndex: theme.zIndex.drawer + 1,
    transition: theme.transitions.create(['width', 'margin'], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
  },
  appBarShift: {
    marginLeft: drawerWidth,
    width: `calc(100% - ${drawerWidth}px)`,
    transition: theme.transitions.create(['width', 'margin'], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.enteringScreen,
    }),
  },
  menuButton: {
    marginRight: 30,
  },
  menuButtonHidden: {
    display: 'none',
  },
  logo: {
    maxHeight: 48,
    marginRight: 30,
  },
}));

export function TopBar(props: any) {
  const classes = useStyles();

  return (
    <AppBar position="fixed" className={clsx(classes.appBar, props.open && classes.appBarShift)}>
      <Toolbar>
        <IconButton
          edge="start"
          color="inherit"
          aria-label="open drawer"
          onClick={props.handleDrawerOpen}
          className={clsx(classes.menuButton, props.open && classes.menuButtonHidden)}
        >
          <MenuIcon />
        </IconButton>
        <Link href="#/">
          <img src="logo.png" alt="logo" className={classes.logo} />
        </Link>
        <Typography variant="h6" noWrap><Link href="#/" color="inherit" underline="none">KLOOPS</Link></Typography>
      </Toolbar>
      {props.isFetching && (
        <LinearProgress />
      )}
    </AppBar>
  );
}
