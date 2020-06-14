import React, { Component } from 'react'
import axios from 'axios'

class Users extends Component {
    state = {
        users: []
    }

    async componentDidMount () {
        const res = await axios.get('http://localhost:4000/users');
        this.setState({
            users: res.data
        })
    }

    render() {
        return (
            <div className="container">
                {this.state.users.map(user => (
                    <div className="card text-black bg-transparent mb-3 mr-4" key={user.ID}>
                        <div className="card-body">
                            <p className="card-text">ID: {user.ID}</p>
                            <p className="card-text">Name: {user.Firstname} {user.Lastname}</p>
                            <p className="card-text">Email: {user.Email}</p>
                        </div>
                    </div>
                ))}
            </div>
        )
    }
}

export default Users;