# KLoops Dahsboard frontend

The KLoops Dashboard frontend is a modern React application.

It is written in Typescript and builds on top of the following main packages:
- [Typescript](https://www.typescriptlang.org/)
- [React](https://reactjs.org/)
- [Redux](https://redux.js.org/)
- [React Redux](https://react-redux.js.org/)
- [Redux Toolkit](https://redux-toolkit.js.org/)
- [material-ui](https://material-ui.com/)
- [reconnecting-websocket](https://github.com/pladaria/reconnecting-websocket)

The frontend communicates with the backend through a [rest API](#rest-api) and a [websocket](#web-socket).

## Rest API

TODO

## Web socket

The websocket is available at the `/resources` endpoint. This is used to update the frontend in realtime when some resources are created, updated or deleted.

Basically, the backend sends json messages to the frontend through the websocket in the form of React actions.

Those actions are read on the frontend side and dispatched to the Redux store.