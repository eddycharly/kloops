import { makeStyles, ThemeOptions } from '@material-ui/core/styles';

export const ThemeDark: ThemeOptions = {
    palette: {
        type: "dark",
    }
};

export const ThemeLight: ThemeOptions = {
    palette: {
        type: "light",
    }
};

const drawerWidth = 240;

export const UseStyles = makeStyles((theme) => ({
    root: {
        display: 'flex',
    },
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
    grow: {
        flexGrow: 1,
    },
    fab: {
        position: 'absolute',
        bottom: theme.spacing(4),
        right: theme.spacing(4),
    },
    appBar: {
        zIndex: theme.zIndex.drawer + 1,
        transition: theme.transitions.create(['width', 'margin'], {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.leavingScreen,
        }),
    },
    appBarShift: {
        marginLeft: drawerWidth,
        width: `calc(100% - ${drawerWidth}px)`,
        transition: theme.transitions.create(['width', 'margin'], {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.enteringScreen,
        }),
    },
    menuButton: {
        marginRight: 30,
    },
    menuButtonHidden: {
        display: 'none',
    },
    logo: {
        maxHeight: 48,
        marginRight: 30,
    },
}));
