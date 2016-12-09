/* eslint-disable import/default */

import './favicon.ico';
import './styles/styles.less';

import React from 'react';
import { Router, browserHistory } from 'react-router';
import routes from './routes';
import { render } from 'react-dom';
import mobx from 'mobx';

import todoStore from './store/todoStore';
import Todo from './store/models/Todo';

mobx.useStrict(true);

render(
    <Router history={browserHistory} routes={routes} />,
    document.getElementById('app')
);

window.todoStore = todoStore;
window.toJS = mobx.toJS;
window.Todo = Todo;
