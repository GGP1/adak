import React, { Component } from 'react'
import axios from 'axios'

export default class Reviews extends Component {
    state = {
        reviews: []
    }

    async componentDidMount() {
        const res = await axios.get("http://localhost:4000/reviews")
        this.setState({
            reviews: res.data
        })
    }

    render() {
        return (
            <div className="container">
                {this.state.reviews.map(review => (
                    <div className="card text-black bg-transparent mb-3 mr-4" key={review.ID}>
                            <div className="card-body">
                                <p className="card-text">ID: {review.ID}</p>
                                <p className="card-text">Stars: {review.Stars}</p>
                                <p className="card-text">Comment: {review.Comment}</p>
                                <p className="card-text">UserID: {review.UserID}</p>
                                <p className="card-text">ProductID: {review.ProductID}</p>
                                <p className="card-text">ShopID: {review.ShopID}</p>
                    </div>
                </div>
            ))}
            </div>
        )
    }
}
