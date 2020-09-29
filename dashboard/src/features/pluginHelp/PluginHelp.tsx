import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import { ExpandMore } from '@material-ui/icons';
import { PluginHelp as Model } from 'models';
import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Grid,
  List,
  ListItemText,
  ListItem,
  Typography
} from '@material-ui/core';

const useStyles = makeStyles((theme) => ({
  heading: {
    flexBasis: '33.33%',
    flexShrink: 0,
  },
  secondaryHeading: {
    fontSize: theme.typography.pxToRem(15),
    color: theme.palette.text.secondary,
  },
  details: {
    overflow: 'auto',
    display: 'flex',
    flexDirection: 'column',
  },
}));

function getShortHelp(model: Model) {
  if (model.shortDescription) {
    return model.shortDescription;
  }
  if (model.description) {
    return model.description.split(".")[0];
  }
  return "No description";
}

function getName(prefix: string | undefined, name: string) {
  if (prefix) {
    return `/[${prefix}]${name}`;
  }
  return `/${name}`;
}

interface props {
  model: Model,
  name: string,
};

export function PluginHelp({ model, name }: props) {
  const classes = useStyles();
  return (
    <Accordion>
      <AccordionSummary expandIcon={<ExpandMore />}>
        <Typography variant="h5" color="primary" className={classes.heading}>{name}</Typography>
        <Typography className={classes.secondaryHeading}>{getShortHelp(model)}</Typography>
      </AccordionSummary>
      <AccordionDetails className={classes.details}>
        <Grid container spacing={2}>
          <Grid item xs>
            <Typography variant="h6" color="primary">Description</Typography>
            <Typography>{model.description}</Typography>
          </Grid>
          {model.excludedProviders && (
            <Grid item xs>
              <Typography variant="h6" color="primary">Excluded providers</Typography>
              {model.excludedProviders}
            </Grid>
          )}
          {model.events && (
            <Grid item xs>
              <Typography variant="h6" color="primary">Events</Typography>
              <List>
                {model.events.map(x => (
                  <ListItem><ListItemText>{x}</ListItemText></ListItem>
                ))}
              </List>
            </Grid>
          )}
          {model.commands && model.commands.map(cmd => (
            <Grid item xs>
              <Typography variant="h6" color="primary">Command</Typography>
              <Typography color="secondary">
                <ul>
                  {cmd.names.map(name => (<li>{getName(cmd.prefix, name)}</li>))}
                </ul>
              </Typography>
              <Typography>{cmd.description}</Typography>
              {/* <Typography>Examples</Typography>
              <ul>
                {x.examples.map(y => (
                  <li>{y}</li>
                ))}
              </ul> */}
            </Grid>
          ))}
        </Grid>
      </AccordionDetails>
    </Accordion>
  );
};
