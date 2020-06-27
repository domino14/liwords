import React from 'react';
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import './App.scss';
import 'antd/dist/antd.css';
import useWebSocket from 'react-use-websocket';

import { Table } from './gameroom/table';
import { Lobby } from './lobby/lobby';
import { useStoreContext } from './store/store';

import { getSocketURI } from './socket/socket';
import { decodeToMsg } from './utils/protobuf';
import { onSocketMsg } from './store/socket_handlers';

type Props = {
  username: string;
};

const App = (props: Props) => {
  const store = useStoreContext();
  const socketUrl = getSocketURI(props.username);
  const { sendMessage } = useWebSocket(socketUrl, {
    onOpen: () => console.log('connected to socket'),
    // Will attempt to reconnect on all close events, such as server shutting down
    shouldReconnect: (closeEvent) => true,
    onMessage: (event: MessageEvent) =>
      decodeToMsg(event.data, onSocketMsg(store)),
  });

  return (
    <div className="App">
      <Router>
        <Switch>
          <Route path="/" exact>
            <Lobby username={props.username} sendSocketMsg={sendMessage} />
          </Route>
          <Route path="/game/:gameID">
            {/* Table meaning a game table */}
            <Table
              sendSocketMsg={sendMessage}
              username={props.username}
            />
          </Route>
        </Switch>
      </Router>
    </div>
  );
};

export default App;
