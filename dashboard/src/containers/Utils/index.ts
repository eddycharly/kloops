import { makeStyles } from '@material-ui/core/styles';

export const UseStyles = makeStyles((theme) => ({
    content: {
        flexGrow: 1,
        height: `calc(100vh - 64px)`,
        marginTop: "64px",
        overflow: 'auto',
        display: 'flex',
        flexDirection: 'column',
        backgroundColor: "#f4f4f4",
    },
    paper: {
        margin: theme.spacing(1),
        padding: theme.spacing(1),
        height: "100%",
    },
    home: {
        display: "flex",
        flexDirection: "column",
        justifyContent: "center"
    },
    fab: {
        position: 'absolute',
        bottom: theme.spacing(4),
        right: theme.spacing(4),
    },
}));
