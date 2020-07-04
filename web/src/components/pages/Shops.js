import React, { Component } from 'react'
import axios from 'axios'

export default class Shops extends Component {
    state = {
        shops: []
    }

    async componentDidMount(){
        const res = await axios.get('http://localhost:4000/shops')
        this.setState({
            shops: res.data
        }) 
    }

    render() {
        return (
            <div className="container">
            {this.state.shops.map(shop => (
                <div className="card text-black bg-transparent mb-3 mr-4" key={shop.ID}>
                        <div className="card-body">
                            <p className="card-text">ID: {shop.ID}</p>
                            <p className="card-text">Name: {shop.Name}</p>
                            <p className="card-text">Country: {shop.Location.Country}</p>
                            <p className="card-text">City: {shop.Location.City}</p>
                            <p className="card-text">Address: {shop.Location.Address}</p>
                            <p className="card-text">Reviews: {shop.Reviews}</p>
                            <p className="card-text">Products: {shop.Products}</p>
                        </div>
                    </div>
            ))}
            </div>
        )
    }
}
