import React from 'react';

// Subcomponents
import { IterateProducts, IterateReviews } from './Iterations/Iterations';

function Shop({ props }) {
    return (
        <div className="card text-black bg-transparent m-3 mr-4">
            <div className="card-body">
                <p className="card-text">ID: {props.id}</p>
                <p className="card-text">Name: {props.name}</p>

                <p className="card-text"><strong>Location</strong></p>
                <p className="card-text">Country: {props.location.country}</p>
                <p className="card-text">State: {props.location.state}</p>
                <p className="card-text">Zip code: {props.location.zip_code}</p>
                <p className="card-text">City: {props.location.city}</p>
                <p className="card-text">Address: {props.location.address}</p>

                <p className="card-text"><strong>Products</strong></p>
                <IterateProducts products={props.products} />

                <p className="card-text"><strong>Reviews</strong></p>
                <IterateReviews reviews={props.reviews} />
            </div>
        </div>
    )
}

export default Shop;