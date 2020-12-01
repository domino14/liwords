import { useCallback, useEffect, useMemo, useRef } from 'react';
import axios from 'axios';
import jwt from 'jsonwebtoken';
import useWebSocket from 'react-use-websocket';
import { useLocation } from 'react-router-dom';
import { message } from 'antd';
import { useMountedState } from '../utils/mounted';
import { useLoginStateStoreContext } from '../store/store';
import {
  useOnSocketMsg,
  ReverseMessageType,
  enableShowSocket,
  parseMsgs,
} from '../store/socket_handlers';
import { decodeToMsg } from '../utils/protobuf';
import { toAPIUrl } from '../api/api';
import { ActionType } from '../actions/actions';

const getSocketURI = (): string => {
  const loc = window.location;
  let protocol;
  if (loc.protocol === 'https:') {
    protocol = 'wss:';
  } else {
    protocol = 'ws:';
  }
  const host = window.RUNTIME_CONFIGURATION.socketEndpoint || loc.host;

  return `${protocol}//${host}/ws`;
};

// this only depends on protocol and host, will never change as we navigate SPA.
const socketUrl = getSocketURI();

type TokenResponse = {
  token: string;
  cid: string;
  app_version: string;
};

type DecodedToken = {
  unn: string;
  uid: string;
  a: boolean; // authed
};

// Returning undefined from useEffect is fine, but some linters dislike it.
const doNothing = () => {};

export const LiwordsSocket = (props: {
  disconnect: false; // statically deprecate this.
  resetSocket: () => void;
  setValues: (_: {
    sendMessage: (msg: Uint8Array) => void;
    justDisconnected: boolean;
  }) => void;
}): null => {
  const isMountedRef = useRef(true);
  useEffect(() => () => void (isMountedRef.current = false), []);
  const { useState } = useMountedState();

  const { disconnect, resetSocket, setValues } = props;
  void resetSocket; // for later use
  const onSocketMsg = useOnSocketMsg();

  const loginStateStore = useLoginStateStoreContext();
  const location = useLocation();
  const { pathname } = location;

  // const [socketToken, setSocketToken] = useState('');
  const [justDisconnected, setJustDisconnected] = useState(false);

  // Source-of-truth must be local, not the store.
  const [isConnectedToSocket, setIsConnectedToSocket] = useState(false);
  const { dispatchLoginState } = loginStateStore;
  const getFullSocketUrlAsync = useCallback(async () => {
    console.log('About to request token');
    // Unfortunately this function must return a valid url.
    const failUrl = `${socketUrl}?${new URLSearchParams({
      path: pathname,
    })}`;

    try {
      const resp = await axios.post<TokenResponse>(
        toAPIUrl('user_service.AuthenticationService', 'GetSocketToken'),
        {},
        { withCredentials: true }
      );
      // Important: resetSocket does not resetStore, be very careful to avoid
      // dispatching stuffs from a decommissioned socket after axios returns.
      if (!isMountedRef.current) return failUrl;

      const socketToken = resp.data.token;
      const { cid, app_version } = resp.data;

      const ret = `${socketUrl}?${new URLSearchParams({
        token: socketToken,
        path: pathname,
        cid,
      })}`;

      const decoded = jwt.decode(socketToken) as DecodedToken;
      dispatchLoginState({
        actionType: ActionType.SetAuthentication,
        payload: {
          username: decoded.unn,
          userID: decoded.uid,
          loggedIn: decoded.a,
          connID: cid,
          path: pathname,
        },
      });
      if (!isMountedRef.current) return failUrl;
      console.log('Got token, setting state, and will try to connect...');
      if (window.RUNTIME_CONFIGURATION.appVersion !== app_version) {
        console.log(
          'app version mismatch',
          'local',
          window.RUNTIME_CONFIGURATION.appVersion,
          'remote',
          app_version
        );

        // bring back when we fix circleci sed
        /*
          message.warning(
            'Woogles has been updated. Please refresh this page at your leisure.',
            0
          ); */
      }

      return ret;
    } catch (e) {
      if (e.response) {
        window.console.log(e.response);
      }
      return failUrl;
    }
  }, [dispatchLoginState, pathname]);

  useEffect(() => {
    if (isConnectedToSocket) {
      console.log('connected to socket');
      dispatchLoginState({
        actionType: ActionType.SetConnectedToSocket,
        payload: true,
      });
      message.destroy('connecting-socket');
      setJustDisconnected(false);
      return () => {
        if (isMountedRef.current) {
          console.log('disconnected from socket :(');
        } else {
          // Yes, the smiley matters!
          console.log('disconnected from socket :)');
        }
        // Special case: useEffect cleanups seem to be run in forward order,
        // but resetSocket does not imply resetStore, and it is important that
        // we inform loginStateStore of the unmount.
        dispatchLoginState({
          actionType: ActionType.SetConnectedToSocket,
          payload: false,
        });
        setJustDisconnected(true);
      };
    }
    return doNothing;
  }, [dispatchLoginState, isConnectedToSocket]);

  useEffect(() => {
    if (isConnectedToSocket) {
      return doNothing;
    }
    const t = setTimeout(() => {
      message.warning({
        content: 'Connecting to server...',
        duration: 0,
        key: 'connecting-socket',
      });
    }, 2000);
    return () => {
      clearTimeout(t);
    };
  }, [isConnectedToSocket]);

  const { sendMessage: originalSendMessage } = useWebSocket(
    getFullSocketUrlAsync,
    {
      onOpen: () => {
        setIsConnectedToSocket(true);
      },
      onClose: () => {
        setIsConnectedToSocket(false);
      },
      reconnectAttempts: Infinity,
      reconnectInterval: 1000,
      retryOnError: true,
      shouldReconnect: (closeEvent) => true,
      onMessage: (event: MessageEvent) => decodeToMsg(event.data, onSocketMsg),
    },
    !disconnect
  );

  const sendMessage = useMemo(() => {
    if (!enableShowSocket) return originalSendMessage;

    return (msg: Uint8Array) => {
      const msgs = parseMsgs(msg);

      msgs.forEach((msg) => {
        const { msgType, parsedMsg } = msg;

        console.log(
          '%csent',
          'background: cyan',
          ReverseMessageType[msgType] ?? msgType,
          parsedMsg.toObject(),
          performance.now()
        );
      });

      return originalSendMessage(msg);
    };
  }, [originalSendMessage]);

  const ret = useMemo(() => ({ sendMessage, justDisconnected }), [
    sendMessage,
    justDisconnected,
  ]);
  useEffect(() => {
    setValues(ret);
  }, [setValues, ret]);

  return null;
};
