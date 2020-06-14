import React, { Component } from 'react'
import axios from 'axios'

class Products extends Component {
    state = {
        products: []
    }

    async componentDidMount() {
        const res = await axios.get('http://localhost:4000/products')
        this.setState({
            products: res.data
        })
    }

    render() {
        return (
            <div className="container">
                {this.state.products.map(product => (
                    <div className="card text-black bg-transparent mb-3 mr-4" key={product.ID}>
                        <div className="card-body">
                            <p className="card-text">ID: {product.ID}</p>
                            <p className="card-text">Brand: {product.Brand}</p>
                            <p className="card-text">Category: {product.Category}</p>
                            <p className="card-text">Type: {product.Type}</p>
                            <p className="card-text">Description: {product.Description}</p>
                            <p className="card-text">Weight: {product.Weight}</p>
                            <p className="card-text">Price: {product.Price}</p>
                        </div>
                    </div>
                ))}
            </div>
        )
    }
}

export default Products;