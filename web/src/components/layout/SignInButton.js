import React, { Component } from 'react'
import SignInForm from './SignInForm'
import '../../App.css'
import styled from 'styled-components';
import Modal from 'react-modal';

import { connect } from 'react-redux';

class SignInButton extends Component {
    constructor(props) {
        super(props)
        this.state = {
            display: false,
        }
    }

    formClick = () => {
        this.setState({
            display: !this.state.display,
        });
    }

    render() {
        const HoverButton = styled.div`
        font-family: 'Righteous';
        background: none;
        border-radius: 10px;  
        &:hover {
            background-color: rgba(32, 134, 175, 0.438);
        }
        `
        let userIcon;
        let loggedIn = this.props.state.login.authenticated;

        if (loggedIn === true) {
            userIcon = <button className="btn border-none text-light" type="button" onClick={this.formClick}>
                            <i className="fa fa-user-circle fa-lg text-success p-2"></i>
                        </button>
        } else {
            userIcon = <button className="btn border-none text-light" type="button" onClick={this.formClick}>
                            <i className="fa fa-user-circle fa-lg text-danger p-2"></i>
                        </button>
        }

        return (
            <div>
                <HoverButton>
                    {userIcon}
                </HoverButton>

                <Modal className="bg-dark" isOpen={this.state.display ? true : false} style={{ overlay: { backgroundColor: 'rgba(0, 0, 0, 0.5)' } }} ariaHideApp={false}>
                    {this.state.display ? <SignInForm closeForm={this.formClick} /> : null}
                </Modal>

            </div>
        )
    }
}

const mapStateToProps = state => {
    return {
        state: state,
        authenticated: state.authenticated
    }
}

export default connect(mapStateToProps)(SignInButton)