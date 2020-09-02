import React from 'react';

import Order from '../Order'
import OrderProduct from '../OrderProduct';
import Product from '../Product';
import Review from '../Review';
import Shop from '../Shop';

export function IterateOrders(props) {
    const o = props.orders;

    if (o !== undefined) {
        let orders = o.map(order => {
            return <Order props={order} key={order.id} />
        })
        return orders
    }

    return <p className="card-text">No orders found</p>
}


export function IterateOrderProducts(props) {
    const p = props.products;

    if (p !== undefined) {
        let products = p.map(product => {
            return <OrderProduct props={product} key={product.id} />
        })
        return products
    }

    return <p className="card-text">No order products found</p>
}

export function IterateProducts(props) {
    const p = props.products;

    if (p !== undefined) {
        let products = p.map(product => {
            return <Product props={product} key={product.id} />
        })
        return products
    }

    return <p className="card-text">No products found</p>
}

export function IterateReviews(props) {
    const r = props.reviews;

    if (r !== undefined) {
        let reviews = r.map(review => {
            return <Review props={review} key={review.id} />
        })
        return reviews
    }

    return <p className="card-text">No reviews found</p>
}

export function IterateShops(props) {
    const s = props.shops;

    if (s !== undefined) {
        let shops = s.map(shop => {
            return <Shop props={shop} key={shop.id} />
        })
        return shops
    }

    return <p className="card-text">No shops found</p>
}