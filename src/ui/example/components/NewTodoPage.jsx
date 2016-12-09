import React from 'react';
import {Row, Col} from 'antd';
import NewTodoForm from './NewTodoForm';

class NewTodoPage extends React.Component {
  static propTypes = {
    route: React.PropTypes.object
  }

  static contextTypes = {
    router: React.PropTypes.object
  }

  componentDidMount() {
    this.context.router.setRouteLeaveHook(this.props.route, this.routerWillLeave);
  }

  routerWillLeave() {
    return '未保存, 确定要离开吗?';
  }

  render() {
    return (
      <Row type="flex" justify="center" align="top">
        <Col span={4}>
          <NewTodoForm/>
        </Col>
      </Row>
    );
  }
}

export default NewTodoPage;
