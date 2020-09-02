import React from 'react';

import { IterateOrders, IterateReviews } from './Iterations/Iterations'

function User({ props }) {
    return (
        <div className="card text-black bg-transparent m-3 mr-4">
            <div className="card-body">
                <p className="card-text">ID: {props.id}</p>
                <p className="card-text">Username: {props.username}</p>
                <p className="card-text">Email: {props.email}</p>

                <p className="card-text"><strong>Orders</strong></p>
                <IterateOrders props={props.orders} />

                <p className="card-text"><strong>Reviews</strong></p>
                <IterateReviews props={props.reviews} />

            </div>
        </div>
    )
}

export default User;
