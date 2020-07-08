import React, { Component } from 'react'
import '../../App.css'

// Redux
import { connect } from 'react-redux';
import PropTypes from "prop-types";
import { logoutUser } from '../../redux/actions';

class DeployableNavbar extends Component {

    // Logout function
    logout = () => {
        this.props.logoutUser()
    }
    
    render() {
        const navStyle1 = {
            position: 'absolute',
            top: 85,
            right: -1,
            height: '198px',
            width: '135px',
            zIndex: 1001
        };

        const navStyle2 = {
            position: 'absolute',
            top: 92,
            right: -1,
            height: '198px',
            width: '135px',
            zIndex: 1001
        };

        let loggedIn = this.props.state.login.authenticated;

        if (loggedIn === true) {
            return (
                <div className="container collapse rounded" id="dep-navbar" data-toggle="collapse" style={navStyle1}>
                    <ul className="p-0">

                        <li className=" list-unstyled" id="toggleLi" data-toggle="collapse" data-target="#navbar-collapse.in">
                            <a className="navbar brand text-white" id="dnav-a" href="/" style={{ fontFamily: 'Open Sans', width: '100%' }}>
                                Home
                            </a>
                        </li>

                        <li className="list-unstyled" id="toggleLi" data-toggle="collapse" data-target="#navbar-collapse.in">
                            <a className="navbar brand text-white" id="dnav-a" href="/categories" style={{ fontFamily: 'Open Sans', width: '100%' }}>
                                Categories
                            </a>
                        </li>

                        <li className="list-unstyled" id="toggleLi" data-toggle="collapse" data-target="#navbar-collapse.in">
                            <a className="navbar brand text-success" id="dnav-a" href="/profile" style={{ fontFamily: 'Open Sans', width: '100%' }}>
                                Profile
                            </a>
                        </li>

                        <li className="list-unstyled" id="toggleLi" data-toggle="collapse" data-target="#navbar-collapse.in">
                            <button className="navbar brand text-danger" onClick={this.logout} style={{ fontFamily: 'Open Sans', width: '100%', background: 'none' }}>
                                Logout
                            </button>
                        </li>

                    </ul>
                </div>
            )
        } else {
            return (
                <div className="container collapse rounded" id="dep-navbar" data-toggle="collapse" style={navStyle2}>
                    <ul className="p-0">

                        <li className=" list-unstyled" id="toggleLi" data-toggle="collapse" data-target="#navbar-collapse.in">
                            <a className="navbar brand text-white" id="dnav-a" href="/" style={{ fontFamily: 'Open Sans', width: '100%' }}>
                                Home
                            </a>
                        </li>

                        <li className="list-unstyled" id="toggleLi" data-toggle="collapse" data-target="#navbar-collapse.in">
                            <a className="navbar brand text-white" id="dnav-a" href="/tasks" style={{ fontFamily: 'Open Sans', width: '100%' }}>
                                Browse tasks
                            </a>
                        </li>

                        <li className="list-unstyled" id="toggleLi" data-toggle="collapse" data-target="#navbar-collapse.in">
                            <a className="navbar brand text-white" id="dnav-a" href="/categories" style={{ fontFamily: 'Open Sans', width: '100%' }}>
                                Categories
                            </a>
                        </li>

                        <li className="list-unstyled" id="toggleLi" data-toggle="collapse" data-target="#navbar-collapse.in">
                            <a className="navbar brand text-danger" id="dnav-a" href="/users/login" style={{ fontFamily: 'Open Sans', width: '100%' }}>
                                Sign In
                            </a>
                        </li>

                        <li className="list-unstyled" id="toggleLi" data-toggle="collapse" data-target="#navbar-collapse.in">
                            <a className="navbar brand text-primary" id="dnav-a" href="/users/create" style={{ fontFamily: 'Open Sans', width: '100%' }}>
                                Sign Up
                            </a>
                        </li>

                    </ul>
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
DeployableNavbar.propTypes = {
    logoutUser: PropTypes.func.isRequired
};

export default connect(mapStateToProps, { logoutUser })(DeployableNavbar)