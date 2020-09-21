import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Moment from 'react-moment';
import {
  Chip,
  Fab,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography, Toolbar
} from '@material-ui/core';
import {
  Delete as DeleteIcon
} from '@material-ui/icons';

import {
  getRepoConfig
} from '../../api';

const useStyles = makeStyles((theme) => ({
  paper: {
    margin: theme.spacing(2),
    padding: theme.spacing(2),
  },
  fab: {
    position: 'absolute',
    bottom: theme.spacing(3),
    right: theme.spacing(3),
  },
}));

function RepoConfig(props: any) {
  const classes = useStyles();
  const [item, setItem] = React.useState<any>(null);

  React.useEffect(() => {
    getRepoConfig(props.match.params.name).then(result => setItem(result));
  }, []);

  if (item == null) {
    return null;
  }

  const scmInfos = (repo: any) => {
    if (repo.spec.gitea) return {
      provider: "Gitea",
      organization: repo.spec.gitea.owner,
      repository: repo.spec.gitea.repo,
    };
    if (repo.spec.gitHub) return {
      provider: "GitHub",
      organization: repo.spec.gitHub.owner,
      repository: repo.spec.gitHub.repo,
    };
    return {
      provider: "Unknown",
      organization: "Unknown",
      repository: "Unknown",
    };
  };

  const infos = scmInfos(item);

  const rows = [
    ["Name", item.metadata.name],
    ["Namespace", item.metadata.namespace],
    ["Bot name", item.spec.botName],
    ["Scm provider", infos.provider],
    ["Organization", infos.organization],
    ["Repository", infos.repository],
    ["Age", (<Moment fromNow>{item.metadata.creationTimestamp}</Moment>)],
  ];

  return (
    <Paper className={classes.paper}>
      <Toolbar>
        <Typography>Spec</Typography>
      </Toolbar>
      <Table>
        <TableBody>
          {rows.map((row) => (
            <TableRow>
              <TableCell>{row[0]}</TableCell>
              <TableCell>{row[1]}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
      <Toolbar>
        <Typography>Automerge</Typography>
      </Toolbar>
      <Table>
        <TableBody>
        <TableRow>
              <TableCell>Batch size limit</TableCell>
              <TableCell>{item.spec.autoMerge.batchSizeLimit}</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>Merge type</TableCell>
              <TableCell>{item.spec.autoMerge.mergeType}</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>Labels</TableCell>
              <TableCell>{item.spec.autoMerge.labels.map((label: string) => (
                <Chip label={label} color="primary" variant="outlined" />
              ))}</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>Missing labels</TableCell>
              <TableCell>{item.spec.autoMerge.missingLabels.map((label: string) => (
                <Chip label={label} color="secondary" variant="outlined" />
              ))}</TableCell>
            </TableRow>
            <TableRow>
              <TableCell>Review approved required</TableCell>
              <TableCell>{`${item.spec.autoMerge.reviewApprovedRequired}`}</TableCell>
            </TableRow>
        </TableBody>
      </Table>
      <Fab color="secondary" aria-label="add" className={classes.fab}>
        <DeleteIcon />
      </Fab>
    </Paper>
  );
}

export default RepoConfig;
