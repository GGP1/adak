import React from 'react';

function OrderCart({ props }) {
    return (
        <div className="card text-black bg-transparent m-3 mr-4">
            <div className="card-body">
                <p className="card-text">Order ID: {props.order_id}</p>
                <p className="card-text">Counter: {props.counter}</p>
                <p className="card-text">Weight: {props.weight}</p>
                <p className="card-text">Discount: {props.discount}</p>
                <p className="card-text">Taxes: {props.taxes}</p>
                <p className="card-text">Subtotal: {props.subtotal}</p>
                <p className="card-text">Total: {props.total}</p>
            </div>
        </div>
    )
}

export default OrderCart;