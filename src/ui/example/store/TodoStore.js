import {
  observable,
  action
} from 'mobx';

import Todo from './models/Todo';
import transport from './transport';

class TodoStore {
  transport;

  @observable todos = [];

  constructor(transport) {
    this.transport = transport;
    this.loadTodos();
  }

  @action pushTodos = (todos) => {
    todos.forEach(todo =>
      // !! new Todo(...) 必须在@action标记的函数里面运行
      this.todos.push(new Todo(todo.id, todo.completed, todo.task))
    );
  };


  @action clearTodos = () => {
    this.todos.clear();
  };

  @action removeTodo = (todo) => {
    this.todos.splice(this.todos.indexOf(todo), 1);
  }

  loadTodos() {
    this.clearTodos();
    this.transport.fetchTodos()
      .then(todos => this.pushTodos(todos));
  }

  createTodo(completed, task) {
    this.transport.saveTodo({completed, task})
      .then(todo => this.pushTodos([todo]));
  }
}

const store = new TodoStore(transport);

export default store;
