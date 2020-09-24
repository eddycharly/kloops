import React from 'react';
import { Link } from 'react-router-dom';
import { makeStyles } from '@material-ui/core/styles';
import Moment from 'react-moment';
import {
  Fab,
  IconButton,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow
} from '@material-ui/core';
import {
  Add as AddIcon,
  RssFeed as RssFeedIcon
} from '@material-ui/icons';
import {
  createHook,
  listRepoConfigs
} from '../../api';
import {
  RepoConfigForm
} from '..';

const useStyles = makeStyles((theme) => ({
  paper: {
    margin: theme.spacing(2),
    padding: theme.spacing(2),
    height: "100%",
  },
  fab: {
    position: 'absolute',
    bottom: theme.spacing(4),
    right: theme.spacing(4),
  },
}));

function RepoConfigs() {
  const classes = useStyles();
  const [items, setItems] = React.useState<any[]>([]);
  const [open, setOpen] = React.useState(false);

  React.useEffect(() => {
    listRepoConfigs().then(result => setItems(result));
  }, []);

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

  const handleClose = () => {
    setOpen(false);
  }

  const onHook = (name: string) => {
    createHook(name);
  }

  return (
    <>
      <RepoConfigForm open={open} onClose={handleClose} />
      <Paper className={classes.paper} square>
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
              <TableCell>Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {items && items.map((row) => {
              const infos = scmInfos(row);
              return (
                <TableRow key={row.metadata.name}>
                  <TableCell>
                    <Link to={`/config/repository/${row.metadata.name}`}>{row.metadata.name}</Link>
                  </TableCell>
                  <TableCell>{row.metadata.namespace}</TableCell>
                  <TableCell>{row.spec.botName}</TableCell>
                  <TableCell>{infos.provider}</TableCell>
                  <TableCell>{infos.organization}</TableCell>
                  <TableCell>{infos.repository}</TableCell>
                  <TableCell><Moment fromNow>{row.metadata.creationTimestamp}</Moment></TableCell>
                  <TableCell>
                    <IconButton onClick={() => onHook(row.metadata.name)} size="small">
                      <RssFeedIcon />
                    </IconButton>
                  </TableCell>
                </TableRow>
              )
            })}
          </TableBody>
        </Table>
        <Fab color="secondary" aria-label="add" className={classes.fab} onClick={() => setOpen(true)}>
          <AddIcon />
        </Fab>
      </Paper>
    </>
  );
}

export default RepoConfigs;
