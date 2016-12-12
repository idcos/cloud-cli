import {
  observable,
  action
} from 'mobx';

import { ajax } from 'jquery';

class Tansport {
  @observable pending = 0;

  @action success = (result) => {
    this.pending -= 1;
    return result;
  };

  @action error = (reason) => {
    this.pending -= 1;
    return reason;
  };

  @action doAjax = (promise) => {
    this.pending += 1;
    return promise.then(this.success, this.error);
  };

  @action fetchTodos() {
    // return this
    //   .doAjax(getJSON('/api/todos'))
    //   .then(result => result.data);

    return new Promise((resolve) => {
      resolve([]);
    });
  }

  @action saveTodo(json) {
    return this
      .doAjax(ajax({
        type: 'POST',
        url: '/api/todos',
        data: { todo: json }
      }))
      .then(result => result.data);
  }
}

const transport = new Tansport();

export default transport;
