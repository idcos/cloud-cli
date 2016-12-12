import React from 'react';
import {observer} from 'mobx-react';

const Index = ({state: {todos}}) =>
  <div>
    <h1>List Todos </h1>
    <table>
      <tbody>
        {todos.map((todo) => (
          <tr key={todo.id}>
            <td>{todo.id}</td>
            <td>{todo.task}</td>
            <td>{todo.completed}</td>
          </tr>
        ))}
      </tbody>
    </table>
  </div>;



export default observer(Index);
