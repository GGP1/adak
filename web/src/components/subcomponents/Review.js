import React from 'react';

function Review({ props }) {
    return (
        <div className="card text-black bg-transparent m-3 mr-4">
            <div className="card-body">
                <p className="card-text">ID: {props.stars}</p>
                <p className="card-text">Brand: {props.name}</p>
            </div>
        </div>
    )
}

export default Review;
