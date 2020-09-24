import React from 'react';
import clsx from 'clsx';
import { UseStyles } from '..';
import {
  AppBar,
  IconButton,
  LinearProgress,
  Link,
  Toolbar,
  Typography
} from '@material-ui/core';
import {
  Brightness4 as ThemeIcon,
  GitHub as GitHubIcon,
  Menu as MenuIcon,
} from '@material-ui/icons';

export function TopBar(props: any) {
  const classes = UseStyles();
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
        <div className={classes.grow} />
        <IconButton color="inherit" onClick={props.handleToggleTheme}>
          <ThemeIcon />
        </IconButton>
        <IconButton color="inherit" target="_blank" href="https://github.com/eddycharly/kloops/">
          <GitHubIcon />
        </IconButton>
      </Toolbar>
      {props.isFetching && (
        <LinearProgress />
      )}
    </AppBar>
  );
}
