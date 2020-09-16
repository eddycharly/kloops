import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import {
  Divider,
  GridList,
  GridListTile,
  Paper,
  Typography
} from '@material-ui/core';
import {
  getRepoConfigs
} from '../../api/repoconfigs';

const useStyles = makeStyles((theme) => ({
  content: {
    flexGrow: 1,
    height: '100vh',
    overflow: 'auto',
    backgroundColor: "#f4f4f4",
  },
  paper: {
    height: "100%",
    width: "100%",
  },
}));

function RepoConfigs() {
  const classes = useStyles();
  const [items, setItems] = React.useState<any[]>([]);

  React.useEffect(() => {
    getRepoConfigs().then(result => setItems(result))
  }, [])


  return (
    <>
      {items.map(x => (
        <Typography>{x.metadata.name}</Typography>
      ))}
    </>
  );
}

export default RepoConfigs;
