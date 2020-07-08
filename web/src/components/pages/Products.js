import React, { Component } from 'react'
import axios from 'axios'

// Subcomponents
import Product from '../subcomponents/Product'

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
                    <Product props={product} />
                ))}
            </div>
        )
    }
}

export default Products;