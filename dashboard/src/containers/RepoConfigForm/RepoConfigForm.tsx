import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  FormControl,
  MenuItem,
  Select,
  Switch,
  TextField,
  Typography
} from '@material-ui/core';
import {
  createRepoConfig
} from '../../api/repoconfigs';

const useDataBind = (initial: any) => {
  const [value, setVal] = React.useState(initial)
  const onChange = (e: any) => setVal(e.target.value)
  return { value, onChange }
}

const plugins = [
  "branchcleaner",
  "cat",
  "dog",
  "shrug",
  "welcome",
  "yuks",
];

const useStyles = makeStyles((theme) => ({
  content: {
    flexDirection: 'row',
  },
  pluginSwitch: {
    justifyContent: 'flex-end',
  },
}));

function RepoConfigForm(props: any) {
  const classes = useStyles();

  const botName = useDataBind("kloops-bot");
  const provider = useDataBind("Gitea");
  const owner = useDataBind(null);
  const repo = useDataBind(null);

  const handleOk = () => {
    createRepoConfig({
      metadata: {
        name: `${provider.value}-${owner.value}-${repo.value}`.toLowerCase(),
      },
      spec: {
        botName: botName.value,
        autoMerge: {
          batchSizeLimit: 12,
          mergeType: "squash",
          labels: [
            "lgtm",
            "approve",
          ],
          missingLabels: [
            "blabla",
          ],
          reviewApprovedRequired: true,
        },
        pluginConfig: {
          ref: "toto",
          plugins: [
            "branchcleaner",
            "cat",
            "dog",
            "shrug",
            "welcome",
            "yuks",
          ]
        },
        [provider.value]: {
          owner: owner.value,
          repo: repo.value,
          hmacToken: {
            valueFrom: {
              secretKeyRef: {
                name: "bla",
                key: "hmac",
              }
            }
          },
          token: {
            valueFrom: {
              secretKeyRef: {
                name: "bla",
                key: "token",
              }
            }
          }
        }
      },
    });
    props.onClose();
  }

  return (
    <Dialog onClose={props.onClose} open={props.open}>
      <DialogTitle>New repo config</DialogTitle>
      <DialogContent>
        <TextField label="Bot name" fullWidth {...botName} />
        <Select fullWidth {...provider} >
          <MenuItem value="GitHub">GitHub</MenuItem>
          <MenuItem value="Gitea">Gitea</MenuItem>
        </Select>
        <TextField label="Organization" fullWidth {...owner} />
        <TextField label="Repository" fullWidth {...repo} />
        {plugins.map((name) => (
          <FormControl fullWidth className={classes.content}>
            <Typography>{name}</Typography>
            <Switch className={classes.pluginSwitch} />
          </FormControl>
        ))}
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose} color="primary">
          Cancel
        </Button>
        <Button onClick={handleOk} color="primary">
          Create
        </Button>
      </DialogActions>
    </Dialog>
  );
}

export default RepoConfigForm;
