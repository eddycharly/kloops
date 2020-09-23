import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Typography
} from '@material-ui/core';
import {
  getPluginHelp
} from '../../api/PluginHelp';
import {
  PluginHelp
} from '../../models';
import {
  ExpandMore as ExpandMoreIcon
} from '@material-ui/icons';

const useStyles = makeStyles((theme) => ({
  root: {
    margin: theme.spacing(1),
  },
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

export function PluginsHelp() {
  const classes = useStyles();
  const [items, setItems] = React.useState<{ [name: string]: PluginHelp } | null>(null);

  React.useEffect(() => {
    getPluginHelp().then(result => setItems(result));
  }, []);

  if (!items) {
    return null;
  }

  const getShortHelp = (pluginHelp: PluginHelp)=>{
    if (pluginHelp.shortDescription){
      return pluginHelp.shortDescription;
    }
    if (pluginHelp.description){
      return pluginHelp.description.split(".")[0];
    }
    return "No description";
  }

  return (
    <div className={classes.root}>
      {Object.entries(items).map(([key, value]) => (
        <Accordion>
          <AccordionSummary expandIcon={<ExpandMoreIcon />}>
            <Typography variant="h5" color="primary" className={classes.heading}>{key}</Typography>
            <Typography className={classes.secondaryHeading}>{getShortHelp(value)}</Typography>
          </AccordionSummary>
          <AccordionDetails className={classes.details}>
            <Typography variant="h5">Description</Typography>
            {value.description}
            {value.excludedProviders && (
              <>
                <Typography variant="h5">Excluded providers</Typography>
                {value.excludedProviders}
              </>
            )}
            {value.events && (
              <>
                <Typography variant="h5">Events</Typography>
                {value.events}
              </>
            )}
            {value.commands && value.commands.map(cmd => (
              <Typography>{cmd.usage}</Typography>
            ))}
          </AccordionDetails>
        </Accordion>
      ))}
    </div>
  );
}
