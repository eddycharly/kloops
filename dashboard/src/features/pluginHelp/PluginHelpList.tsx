import React from 'react';
import { useDispatch, useSelector } from 'react-redux'
import { makeStyles } from '@material-ui/core/styles';
import { PluginHelp } from './PluginHelp';
import { FetchAll } from './Slice';
import { RootState } from 'reducers'

const useStyles = makeStyles((theme) => ({
  root: {
    margin: theme.spacing(1),
  },
}));

export function PluginHelpList() {
  const classes = useStyles();
  const models = useSelector((state: RootState) => (state.pluginHelp && state.pluginHelp.state === 'finished') ? state.pluginHelp.data : {});
  const dispatch = useDispatch();

  React.useEffect(() => {
    dispatch(FetchAll())
  }, [dispatch]);

  return (
    <div className={classes.root}>
      {Object.entries(models).map(([key, value]) => (
        <PluginHelp model={value} name={key} />
      ))}
    </div>
  )
};
