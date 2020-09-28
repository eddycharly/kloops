import React from 'react';
import clsx from 'clsx';
import { UseStyles } from '..';
import { Brightness4, GitHub, Menu } from '@material-ui/icons';
import {
  AppBar,
  IconButton,
  LinearProgress,
  Link,
  Toolbar,
  Typography
} from '@material-ui/core';

interface props {
  open: boolean;
  handleDrawerOpen: () => void;
  handleToggleTheme: () => void;
};

export function TopBar({ handleDrawerOpen, handleToggleTheme, open }: props) {
  const classes = UseStyles();
  const isFetching = false;
  return (
    <AppBar position="fixed" className={clsx(classes.appBar, open && classes.appBarShift)}>
      <Toolbar>
        <IconButton
          edge="start"
          color="inherit"
          aria-label="open drawer"
          onClick={handleDrawerOpen}
          className={clsx(classes.menuButton, open && classes.menuButtonHidden)}
        >
          <Menu />
        </IconButton>
        <Link href="#/">
          <img src="logo.png" alt="logo" className={classes.logo} />
        </Link>
        <Typography variant="h6" noWrap><Link href="#/" color="inherit" underline="none">KLOOPS</Link></Typography>
        <div className={classes.grow} />
        <IconButton color="inherit" onClick={handleToggleTheme}>
          <Brightness4 />
        </IconButton>
        <IconButton color="inherit" target="_blank" href="https://github.com/eddycharly/kloops/">
          <GitHub />
        </IconButton>
      </Toolbar>
      {isFetching && (
        <LinearProgress />
      )}
    </AppBar>
  );
}
