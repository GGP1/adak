import React, { useState, useEffect } from 'react';
import axios from 'axios';

import { IterateOrders, IterateReviews } from '../subcomponents/Iterations/Iterations';

function Users() {
    const [users, setUsers] = useState([])

    useEffect(() => {
        async function getUsers() {
            const res = await axios.get('http://localhost:4000/users');
            setUsers(res.data);
        }
        getUsers();
    }, [])

    return (
        <div className="container mt-4">
            {users.map(user => (
                <div className="card text-black bg-transparent mb-3 mr-4" key={user.ID}>
                    <div className="card-body">
                        <p className="card-text"><strong>ID:</strong> {user.ID}</p>
                        <p className="card-text"><strong>Username:</strong> {user.username}</p>
                        <p className="card-text"><strong>Email:</strong> {user.email}</p>
                        <p className="card-text"><strong>Cart ID:</strong> {user.cart_id}</p>
                        <p className="card-text"><strong>Orders</strong></p>
                        <IterateOrders orders={user.orders} />
                        <p className="card-text"><strong>Reviews</strong></p>
                        <IterateReviews reviews={user.reviews} />
                    </div>
                </div>
            ))}
        </div>
    )

}

export default Users;