import React from 'react';
import {Form, Input, Button, Checkbox} from 'antd';
const FormItem = Form.Item;

class NewTodoForm extends React.Component {
    handleSubmit = (e) => {
        e.preventDefault();
        this.props.form.validateFields((errors, values) => {
            if (errors) {
                console.log("Error!", errors, values);
                return;
            }
            console.log('Submit!!!');
            console.log(values);
        });
    }

    taskExists(rule, value, callback) {
        if (!value) {
            callback();
        } else {
            setTimeout(() => {
                if (value === 'hello') {
                    callback([new Error('抱歉，Task名称不能是hello。')]);
                } else {
                    callback();
                }
            }, 800);
        }
    }

    render() {
        const {getFieldProps, getFieldError, isFieldValidating} = this.props.form;

        const taskProps = getFieldProps('task', {
            rules: [
                {required: true, min: 5, message: 'Task名至少为 5 个字符'},
                {validator: this.taskExists},
            ],
        });
        return (
            <Form>
                <FormItem label="任务"
                          hasFeedback
                          help={isFieldValidating('task') ? '校验中...' : (getFieldError('task') || []).join(', ')}
                >
                    <Input placeholder="请输入任务"
                           {...taskProps}
                    />
                </FormItem>
                <FormItem>
                    <Checkbox {...getFieldProps('finished')}>完成</Checkbox>
                </FormItem>
                <Button type="primary" htmlType="submit" onClick={this.handleSubmit}>添加</Button>
            </Form>
        );
    }
}


export default Form.create()(NewTodoForm);
