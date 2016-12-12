import React, { PropTypes } from 'react';
import { Link, IndexLink } from 'react-router';
import {observer} from 'mobx-react';
import Loading from './Loading';

const AppTodo = ({children, uiState: {appIsInSync}}) =>
  <div>
    {appIsInSync || <Loading />}
    <IndexLink to="/">Todos</IndexLink>
    {' | '}
    <Link to="/new-todo">New Todo</Link>
    {' | '}
    <Link to="/about">About</Link>
    <br/>
    {children}
  </div>;


AppTodo.propTypes = {
  children: PropTypes.element,
  uiState: PropTypes.object
};

export default observer(AppTodo);
