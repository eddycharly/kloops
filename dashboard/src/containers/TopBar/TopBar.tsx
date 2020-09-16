import React from 'react';
// import { connect } from 'react-redux';
import clsx from 'clsx';
import { makeStyles } from '@material-ui/core/styles';
import {
  AppBar,
  IconButton,
  LinearProgress,
  Toolbar,
  Typography
} from '@material-ui/core';

import {
  Menu as MenuIcon
} from '@material-ui/icons';

// import { isFetching } from '../../reducers';

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
    marginRight: 36,
  },
  menuButtonHidden: {
    display: 'none',
  },
}));

function TopBar(props: any) {
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
        <Typography variant="h6" noWrap>
          KLOOPS
      </Typography>
      </Toolbar>
      {props.isFetching && (
        <LinearProgress />
      )}
    </AppBar>
  );
}

// function mapStateToProps(state: any) {
//   return {
//     isFetching: isFetching(state)
//   };
// }

// export default connect(mapStateToProps, null)(TopBar);
export default TopBar;
