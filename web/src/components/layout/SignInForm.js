import React, { Component } from 'react';
import { Button, Form } from 'react-bootstrap';
import axios from 'axios';
import styled from 'styled-components';

// Redux
import PropTypes from "prop-types";
import { connect } from 'react-redux';
import { loginUser, logoutUser } from '../../redux/actions';

class SignInForm extends Component {

  _isMounted = false;

  constructor(props) {
    super(props)
    this.state = {
      email: '',
      password: ''
    }
  }

  componentDidMount = async () => {
    // Check if component is mounted
    this._isMounted = true;
    // Load all the users from the database
    const res = await axios.get('http://localhost:4000/users');
    this.setState({
      users: {
        email: res.data.map(user => user.email),
        password: res.data.map(user => user.password)
      }
    });
  }

  componentWillUnmount = () => {
    this._isMounted = false;
  }

  // Load input data to the component state
  onInputChange = (e) => {
    this.setState({
      [e.target.name]: e.target.value
    });
  }

  // Login function
  signIn = async () => {
    // User input data
    const user = {
      email: this.state.email,
      password: this.state.password,
    };
    // Sending user object to the server
    const res = await axios.post('http://localhost:4000/login', user);

    // Check if the login is correct
    if (res.status === 200) {
      this.props.loginUser(user)
      this.props.closeForm()
    } else {
      console.log("Invalid email or password.")
    }
  }

  // Logout function
  logout = () => {
    if (this._isMounted)
      this.props.logoutUser()
      this.props.closeForm()
  }

  render() {

    const containerStyle = {
      minHeight: '60px',
      margin: '0 auto',
      width: '120px',
      bottom: 0,
      left: -230,
      top: '25vh',
      right: 0,
      zIndex: 1003
    }

    const HoverText = styled.div`
      padding: 0;
      margin: 0;
      &:hover {
        color: black
      }
      `

    let loggedIn = this.props.state.login.authenticated;

    if (loggedIn === false) {
      return (
        <div id="form" className="container position-fixed" style={containerStyle}>
          <div className="position-absolute rounded p-3"
            style={{ 'backgroundColor': 'rgba(100, 22, 42, 0.95)', 'width': '20em' }}>
            <Form onSubmit={this.signIn}>
              <button type="button" className="close ml-2 mb-1" onClick={this.props.closeForm}>
                <span className="text-white" aria-hidden="true">×</span>
              </button>
              <h6 className="dropdown-header text-center text-white">Sign in</h6>
              <div className="dropdown-divider" style={{ 'width': '50%', 'marginLeft': '25%' }}></div>

              <Form.Group className="mt-3 text-white" controlId="formUsername">
                <Form.Label>Username</Form.Label>
                <input
                  className="form-control text-white"
                  style={{ 'background': 'transparent' }}
                  type="string"
                  name="email"
                  placeholder="Enter your username"
                  value={this.state.email}
                  onChange={this.onInputChange}
                  required
                />
              </Form.Group>

              <Form.Group className="text-white" controlId="formPassword">
                <Form.Label>Password</Form.Label>
                <input
                  className="form-control text-white"
                  style={{ 'background': 'transparent' }}
                  type="password"
                  name="password"
                  placeholder="Password"
                  value={this.state.password}
                  onChange={this.onInputChange} />
                <Form.Text className="text-white">Keep your password safe.</Form.Text>
              </Form.Group>

              {/*<Form.Group className="text-white" controlId="formCheckbox">
              <Form.Check type="checkbox" label="Check me out" />
            </Form.Group>*/}

              <Button className="float-none btn-block rounded-pill btn-light" variant="primary" type="submit">
                Sign in
            </Button>

              <div className="dropdown-divider bg-dark"></div>
              <div id="questions">
                <a className="dropdown-item text-center text-white" href="/users/add">
                  <HoverText>Not on joblib yet? Register</HoverText>
                </a>
                <a className="dropdown-item text-center text-white" href="/recover">
                  <HoverText>Forgot password?</HoverText>
                </a>
              </div>
            </Form>
          </div>
        </div>
      )
    } else { // If user is logged in, return this modal
      return (
        <div id="form" className="container position-fixed" style={containerStyle}>
          <div className="position-absolute rounded p-3"
            style={{ 'backgroundColor': 'rgba(100, 22, 42, 0.95)', 'width': '20em' }}>
            <Form onSubmit={this.logout}>
              <button type="button" className="close ml-2 mb-1" onClick={this.props.closeForm}>
                <span className="text-white" aria-hidden="true">×</span>
              </button>
              <h6 className="dropdown-header text-center text-white">{this.props.state.login.token.username}</h6>
              <div className="dropdown-divider" style={{ 'width': '50%', 'marginLeft': '25%' }}></div>
              <Button className="float-none btn-block btn-light mt-4 mb-2" variant="primary" type="submit">
                Logout
            </Button>
            </Form>
          </div>
        </div>
      )
    }
  }
}

const mapStateToProps = state => {
  return {
    state: state,
    authenticated: state.authenticated
  }
}

SignInForm.propTypes = {
  loginUser: PropTypes.func.isRequired,
  logoutUser: PropTypes.func.isRequired
};

export default connect(mapStateToProps, { loginUser, logoutUser })(SignInForm)