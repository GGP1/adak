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
            <div className="container mt-4">
                {this.state.users.map(user => (
                    <div className="card text-black bg-transparent mb-3 mr-4" key={user.ID}>
                        <div className="card-body">
                            <p className="card-text"><strong>ID:</strong> {user.ID}</p>
                            <p className="card-text"><strong>Name:</strong> {user.name}</p>
                            <p className="card-text"><strong>Email:</strong> {user.email}</p>
                            <p className="card-text"><strong>Cart ID:</strong> {user.cart_id}</p>
                            <p className="card-text"><strong>Orders:</strong> {user.orders}</p>
                            <p className="card-text"><strong>Reviews:</strong> {user.reviews}</p>
                        </div>
                    </div>
                ))}
            </div>
        )
    }
}

export default Users;