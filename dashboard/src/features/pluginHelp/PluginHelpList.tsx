import React from 'react';
import { useDispatch, useSelector } from 'react-redux'
import { makeStyles } from '@material-ui/core/styles';
import { PluginHelp } from './PluginHelp';
import { FetchAll } from 'app/reducers/pluginHelp';
import { RootState } from 'app/reducers';
import { useSnackbar } from 'notistack';

const useStyles = makeStyles((theme) => ({
  root: {
    margin: theme.spacing(1),
  },
}));

export function PluginHelpList() {
  const classes = useStyles();
  const models = useSelector((state: RootState) => (state.pluginHelp && state.pluginHelp.state === 'finished') ? state.pluginHelp.data : {});
  const dispatch = useDispatch();
  const { enqueueSnackbar } = useSnackbar();

  React.useEffect(() => {
    dispatch(FetchAll(enqueueSnackbar))
  }, [dispatch, enqueueSnackbar]);

  return (
    <div className={classes.root}>
      {Object.entries(models).map(([key, value]) => (
        <PluginHelp model={value} name={key} />
      ))}
    </div>
  )
};
