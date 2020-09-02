import React from 'react';


function OrderProduct({ props }) {
    return (
        <div className="card text-black bg-transparent m-3 mr-4">
            <div className="card-body">
                <p className="card-text">ID: {props.product_id}</p>
                <p className="card-text">Quantity: {props.quantity}</p>
                <p className="card-text">Brand: {props.brand}</p>
                <p className="card-text">Category: {props.category}</p>
                <p className="card-text">Type: {props.type}</p>
                <p className="card-text">Description: {props.description}</p>
                <p className="card-text">Weight: {props.weight}</p>
                <p className="card-text">Discount: {props.discount}</p>
                <p className="card-text">Taxes: {props.taxes}</p>
                <p className="card-text">Subtotal: {props.subtotal}</p>
                <p className="card-text">Total: {props.total}</p>
            </div>
        </div>
    )
}

export default OrderProduct;