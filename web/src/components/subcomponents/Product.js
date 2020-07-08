import React from 'react';

function Product({props}) {
    return(
            <div className="card text-black bg-transparent mb-3 mr-4" key={props.product.ID}>
                <div className="card-body">
                    <p className="card-text">ID: {props.product.ID}</p>
                    <p className="card-text">Brand: {props.product.Brand}</p>
                    <p className="card-text">Category: {props.product.Category}</p>
                    <p className="card-text">Type: {props.product.Type}</p>
                    <p className="card-text">Description: {props.product.Description}</p>
                    <p className="card-text">Weight: {props.product.Weight}</p>
                    <p className="card-text">Taxes: {props.product.Taxes}</p>
                    <p className="card-text">Discount: {props.product.Discount}</p>
                    <p className="card-text">Subtotal: {props.product.Subtotal}</p>
                    <p className="card-text">Total: {props.product.Total}</p>
                    <p className="card-text">Reviews: {props.product.Reviews}</p>
                </div>
            </div>
    )
}

export default Product;