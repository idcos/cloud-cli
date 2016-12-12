import React from 'react';
import {observer} from 'mobx-react';

@observer
class User extends React.Component {
  render() {
    const {user: {fullName, changeName, remove}} = this.props;

    return (
      <div>
        <h3>{fullName}</h3>
        <button onClick={changeName}>Change</button>
        <button onClick={remove}>Delete</button>
      </div>
  );
  }
}

export default User;
