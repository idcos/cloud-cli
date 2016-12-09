import {
  observable,
  computed
} from 'mobx';

class Todo {
  id = null;

  @observable completed = false;
  @observable task = "";

  constructor(id, completed, task) {
    this.id = id;
    this.completed = completed;
    this.task = task;
  }

  @computed get asJson() {
    return {
      id: this.id,
      completed: this.completed,
      task: this.task,
      authorId: this.author ? this.author.id : null
    };
  }
}


export default Todo;
