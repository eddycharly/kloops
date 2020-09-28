import React from 'react';
import { useDispatch, useSelector } from 'react-redux'
import { UseStyles } from 'containers/Utils';
import { FetchAll } from './Slice';
import { RootState } from 'reducers'
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
  const models = useSelector((state: RootState) => state.commandHelp.state === 'finished' ? state.commandHelp.data : {});
  const dispatch = useDispatch();

  React.useEffect(() => {
    dispatch(FetchAll())
  }, []);

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
