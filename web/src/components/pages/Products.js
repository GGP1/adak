import React, { useState, useEffect } from 'react'
import axios from 'axios'

import { IterateProducts } from '../subcomponents/Iterations/Iterations'

function Products() {
    const [products, setProducts] = useState([])

    useEffect(() => {
        async function getProducts() {
            const res = await axios.get('http://localhost:4000/products')
            setProducts(res.data)
        }
        getProducts();
    }, [])

    return (
        <div className="container">
            <IterateProducts products={products} />
        </div>
    )
}

export default Products;