import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Moment from 'react-moment';
import {
  Divider,
  GridList,
  GridListTile,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography
} from '@material-ui/core';
import {
  listRepoConfigs
} from '../../api/repoconfigs';

const useStyles = makeStyles((theme) => ({
  paper: {
    margin: "20px",
    padding: "20px",
    height: "100%",
    // width: "100%",
  },
}));

function RepoConfigs() {
  const classes = useStyles();
  const [items, setItems] = React.useState<any[]>([]);

  React.useEffect(() => {
    listRepoConfigs().then(result => setItems(result));
  }, []);

  const scmInfos = (repo: any) => {
    if (repo.spec.gitea) return {
      provider:  "Gitea",
      organization: repo.spec.gitea.owner,
      repository: repo.spec.gitea.repo,
    };
    if (repo.spec.gitHub) return {
      provider:  "GitHub",
      organization: repo.spec.gitHub.owner,
      repository: repo.spec.gitHub.repo,
    };
    return {
      provider: "Unknown",
      organization: "Unknown",
      repository: "Unknown",
    };
  };
  return (
    <Paper className={classes.paper}>
      <Table aria-label="simple table">
        <TableHead>
          <TableRow>
            <TableCell>Name</TableCell>
            <TableCell>Namespace</TableCell>
            <TableCell>Bot name</TableCell>
            <TableCell>Scm provider</TableCell>
            <TableCell>Organization</TableCell>
            <TableCell>Repository</TableCell>
            <TableCell>Age</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {items.map((row) => {
            const infos = scmInfos(row);
            return (
              <TableRow key={row.metadata.name}>
                <TableCell component="th" scope="row">{row.metadata.name}</TableCell>
                <TableCell component="th" scope="row">{row.metadata.namespace}</TableCell>
                <TableCell component="th" scope="row">{row.spec.botName}</TableCell>
                <TableCell component="th" scope="row">{infos.provider}</TableCell>
                <TableCell component="th" scope="row">{infos.organization}</TableCell>
                <TableCell component="th" scope="row">{infos.repository}</TableCell>
                <TableCell><Moment fromNow>{row.metadata.creationTimestamp}</Moment></TableCell>
              </TableRow>
            )
          })}
        </TableBody>
      </Table>
    </Paper>
  );
}

export default RepoConfigs;
