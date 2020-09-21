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
    Delete as DeleteIcon,
    FileCopy as FileCopyIcon,
    RssFeed as RssFeedIcon

} from '@material-ui/icons';
import {
    listPluginConfigs
} from '../../api/PluginConfig';
import {
    PluginConfig
} from '../../models';

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

export function PluginConfigs() {
    const classes = useStyles();
    const [items, setItems] = React.useState<Array<PluginConfig>>([]);

    React.useEffect(() => {
        listPluginConfigs().then(result => setItems(result));
    }, []);

    return (
        <>
            <Paper className={classes.paper}>
                <Table aria-label="simple table">
                    <TableHead>
                        <TableRow>
                            <TableCell>Name</TableCell>
                            <TableCell>Namespace</TableCell>
                            <TableCell>Age</TableCell>
                            <TableCell>Actions</TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {items.map((row) => {
                            return (
                                <TableRow key={row.name}>
                                    <TableCell component="th" scope="row">
                                        <Link to={`/config/plugin/${row.name}`}>{row.name}</Link>
                                    </TableCell>
                                    <TableCell component="th" scope="row">{row.namespace}</TableCell>
                                    <TableCell><Moment fromNow>{row.creationTimestamp}</Moment></TableCell>
                                    <TableCell>
                                        <IconButton size="small">
                                            <DeleteIcon />
                                        </IconButton>
                                        <IconButton size="small">
                                            <FileCopyIcon />
                                        </IconButton>
                                        <IconButton size="small">
                                            <RssFeedIcon />
                                        </IconButton>
                                    </TableCell>
                                </TableRow>
                            )
                        })}
                    </TableBody>
                </Table>
                <Fab color="secondary" aria-label="add" className={classes.fab}>
                    <AddIcon />
                </Fab>
            </Paper>
        </>
    );
}
