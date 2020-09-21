import React from 'react';
import { Link } from 'react-router-dom';
import { makeStyles } from '@material-ui/core/styles';
import {
  Drawer,
  IconButton,
  Divider,
  List,
  ListItem,
  ListItemText,
  ListSubheader,
  Toolbar
} from '@material-ui/core';
import {
  ChevronLeft as ChevronLeftIcon
} from '@material-ui/icons';

const drawerWidth = 240;

const useStyles = makeStyles((theme) => ({
  drawer: {
    width: drawerWidth,
    flexShrink: 0,
  },
  drawerPaper: {
    width: drawerWidth,
  },
  drawerContainer: {
    overflow: 'auto',
  },
  toolbarIcon: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'flex-end',
    padding: '0 8px',
    ...theme.mixins.toolbar,
  },
}));

function SideNav(props: any) {
  const classes = useStyles();

  return (
    <Drawer
      className={classes.drawer}
      variant="temporary"
      onEscapeKeyDown={props.handleDrawerClose}
      onBackdropClick={props.handleDrawerClose}
      open={props.open}
      classes={{
        paper: classes.drawerPaper,
      }}
    >
      <Toolbar className={classes.toolbarIcon}>
        <IconButton onClick={props.handleDrawerClose}>
          <ChevronLeftIcon />
        </IconButton>
      </Toolbar>
      <Divider />
      <div className={classes.drawerContainer}>
        <List>
          <ListSubheader>Configuration</ListSubheader>
          <Divider />
          <ListItem button component={Link} to="/config/kloops" onClick={props.handleDrawerClose}>
            <ListItemText primary="Kloops" />
          </ListItem>
          <ListItem button component={Link} to="/config/repositories" onClick={props.handleDrawerClose}>
            <ListItemText primary="Repositories" />
          </ListItem>
          <ListItem button component={Link} to="/config/plugins" onClick={props.handleDrawerClose}>
            <ListItemText primary="Plugins" />
          </ListItem>
          <Divider />
          <ListSubheader>Jobs</ListSubheader>
          <Divider />
          <ListItem button component={Link} to="/help/plugins" onClick={props.handleDrawerClose}>
            <ListItemText primary="Jobs" />
          </ListItem>
          <Divider />
          <ListSubheader>Help</ListSubheader>
          <Divider />
          <ListItem button component={Link} to="/help/plugins" onClick={props.handleDrawerClose}>
            <ListItemText primary="Plugins" />
          </ListItem>
          <ListItem button component={Link} to="/help/about" onClick={props.handleDrawerClose}>
            <ListItemText primary="About" />
          </ListItem>
        </List>
      </div>
    </Drawer>
  );
}

export default SideNav;
