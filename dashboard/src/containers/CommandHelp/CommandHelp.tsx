import React from 'react';
import { UseStyles } from '..';
import {
  Paper,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow
} from '@material-ui/core';
import {
  getPluginHelp
} from '../../api/PluginHelp';
import {
  PluginHelp
} from '../../models';

export function CommandHelp() {
  const classes = UseStyles();
  const [items, setItems] = React.useState<{ [name: string]: PluginHelp } | null>(null);

  React.useEffect(() => {
    getPluginHelp().then(result => setItems(result));
  }, []);

  if (!items) {
    return null;
  }

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
          {Object.entries(items).map(([key, value]) => (
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
