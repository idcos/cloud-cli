import React from 'react';
import {observer} from 'mobx';

const Todo = ({todo: {id, completed, task}}) =>
  <h2>{id} - {completed} - {task}</h2>

export default observer(Todo)
