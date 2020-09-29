import React from 'react';
import { useDispatch, useSelector } from 'react-redux'
import { UseStyles } from 'containers/Utils';
import { FetchAll } from 'app/reducers/pluginHelp';
import { RootState } from 'app/reducers';
import { useSnackbar } from 'notistack';
import {
  Paper,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow
} from '@material-ui/core';

export function CommandHelp() {
  const classes = UseStyles();
  const models = useSelector((state: RootState) => (state.pluginHelp && state.pluginHelp.state === 'finished') ? state.pluginHelp.data : {});
  const dispatch = useDispatch();
  const { enqueueSnackbar } = useSnackbar();

  React.useEffect(() => {
    dispatch(FetchAll(enqueueSnackbar))
  }, [dispatch, enqueueSnackbar]);

  return (
    <Paper className={classes.paper}>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>Command</TableCell>
            <TableCell>Example</TableCell>
            <TableCell>Description</TableCell>
            <TableCell>Who Can Use</TableCell>
            <TableCell>Plugin</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {Object.entries(models).map(([key, value]) => (
            <>
              {value.commands && value.commands.map(cmd => (
                <TableRow>
                  <TableCell>{cmd.usage}</TableCell>
                  <TableCell>
                    <ul>
                      {cmd.examples.map(y => (
                        <li>{y}</li>
                      ))}
                    </ul>
                  </TableCell>
                  <TableCell>{cmd.description}</TableCell>
                  <TableCell>{cmd.whoCanUse}</TableCell>
                  <TableCell>{key}</TableCell>
                </TableRow>
              ))}
            </>
          )
          )}
        </TableBody>
      </Table>
    </Paper>
  );
}
