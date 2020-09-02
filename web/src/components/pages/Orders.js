import React, { useEffect, useState } from 'react'
import axios from 'axios'

// Subcomponents
import { IterateOrders } from '../subcomponents/Iterations/Iterations'

function Orders() {
    const [orders, setOrders] = useState([])

    useEffect(() => {
        function getOrders() {
            const res = await axios.get('http://localhost:4000/orders')
            setOrders(res.data)
        }
        getOrders();
    }, [])

    return (
        <div className="container">
            <IterateOrders products={orders} />
        </div>
    )
}

export default Orders;