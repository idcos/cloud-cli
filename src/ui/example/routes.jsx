import React from 'react';
import { Route, IndexRoute } from 'react-router';

import AppTodo from './components/AppTodo';
import IndexPage from './components/IndexPage';
import NewTodoPage from './components/NewTodoPage';
import AboutPage from './components/AboutPage';
import NotFoundPage from './components/NotFoundPage';

import uiState from './store/uiState';
import todoStore from './store/todoStore';

const appTodo = ({children}) =>
  <AppTodo uiState={uiState} state={todoStore}>
    {children}
  </AppTodo>;

const indexPage = () =>
  <IndexPage state={todoStore} />;

export default (
  <Route path="/"
      component={appTodo}>
    <IndexRoute component={indexPage}/>
    <Route path="new-todo" component={NewTodoPage}/>
    <Route path="about" component={AboutPage}/>
    <Route path="*" component={NotFoundPage}/>
  </Route>
);
