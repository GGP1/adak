import React from 'react';

import OrderCart from './OrderCart';
import { IterateOrderProducts } from './Iterations/Iterations';

function Order({ props }) {
    return (
        <div className="card text-black bg-transparent m-3 mr-4">
            <div className="card-body">
                <p className="card-text">ID: {props.id}</p>
                <p className="card-text">User ID: {props.user_id}</p>
                <p className="card-text">Currency: {props.currency}</p>
                <p className="card-text">Address: {props.brand}</p>
                <p className="card-text">City: {props.category}</p>
                <p className="card-text">State: {props.type}</p>
                <p className="card-text">Zip code: {props.description}</p>
                <p className="card-text">Country: {props.weight}</p>
                <p className="card-text">Status: {props.discount}</p>
                <p className="card-text">Ordered at: {props.taxes}</p>
                <p className="card-text">Delivery date: {props.subtotal}</p>

                <p className="card-text"><strong>Cart</strong></p>
                <OrderCart props={props.cart}/>

                <p className="card-text"><strong>Products</strong></p>
                <IterateOrderProducts products={props.products} />
            </div>
        </div>
    )
}

export default Order;