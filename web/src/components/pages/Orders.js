import React, { useEffect, useState } from 'react';
import axios from 'axios';
import cookies from 'js-cookie';

import NotFound from './NotFound';

// Subcomponents
import { IterateOrders } from '../subcomponents/Iterations/Iterations'

function Orders() {
    const [orders, setOrders] = useState([])
    const [status, setStatus] = useState()
    
    let aid = cookies.get("AID")
    const headers = {
        'AID': aid
    }

    useEffect(() => {
        async function getOrders() {
            const res = await axios.get(`http://localhost:4000/orders`, { headers: headers})
            setOrders(res.data)
            setStatus(res.status)
        }
        getOrders();
    }, [headers])

    if (aid === undefined || aid === "" || status === 404) {
        return (
            <NotFound />
        )
    } else {
        return (
            <div className="container">
                <IterateOrders products={orders} />
            </div>
        )
    }
}

export default Orders;